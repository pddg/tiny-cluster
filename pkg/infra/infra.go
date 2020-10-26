package infra

import (
	"context"
	"time"

	"go.etcd.io/etcd/clientv3"
	"golang.org/x/xerrors"

	tcErr "github.com/pddg/tiny-cluster/pkg/errors"
)

const BasePrefix = "/tiny-cluster"

type baseRepoImpl struct {
	config *clientv3.Config
}

func (r *baseRepoImpl) newClient(ctx context.Context) (*clientv3.Client, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	errCh := make(chan error, 1)
	clientCh := make(chan *clientv3.Client)
	go func() {
		client, err := clientv3.New(*r.config)
		if err != nil {
			errCh <- xerrors.Errorf("Failed to create etcd client %w:", err)
		} else {
			clientCh <- client
		}
		close(clientCh)
		close(errCh)
	}()
	select {
	case <-timeoutCtx.Done():
		return nil, tcErr.ErrTimedOut
	case err := <-errCh:
		return nil, xerrors.Errorf("Could not create etcd client %w:", err)
	case client := <-clientCh:
		return client, nil
	}
}

func doGetWithRev(ctx context.Context, client *clientv3.Client, key string, opts ...clientv3.OpOption) ([]byte, int64, error) {
	var value []byte
	resp, err := client.Get(ctx, key, opts...)
	if err != nil {
		return value, 0, xerrors.Errorf("Failed to get the key ('%s') %w:", key, err)
	}
	if resp.Count == 0 {
		return value, 0, tcErr.ErrNotFound
	}
	value = resp.Kvs[0].Value
	return value, resp.Kvs[0].ModRevision, nil
}

func doGet(ctx context.Context, client *clientv3.Client, key string, opts ...clientv3.OpOption) ([]byte, error) {
	value, _, err := doGetWithRev(ctx, client, key, opts...)
	return value, err
}

func doGetAll(ctx context.Context, client *clientv3.Client, key string, opts ...clientv3.OpOption) ([][]byte, error) {
	var values [][]byte
	opts = append(opts, clientv3.WithPrefix())
	resp, err := client.Get(ctx, key, opts...)
	if err != nil {
		return values, xerrors.Errorf("Failed to get the values whose key starts with '%s' %w:", key, err)
	}
	for _, kv := range resp.Kvs {
		if kv.Version == 0 {
			continue
		}
		values = append(values, kv.Value)
	}
	return values, nil
}

func doCreate(ctx context.Context, client *clientv3.Client, key string, value string) error {
	doesNotExist := clientv3.Compare(clientv3.Version(key), "=", 0)
	create := clientv3.OpPut(key, value)
	createResp, err := client.Txn(ctx).
		If(doesNotExist).
		Then(create).
		Commit()
	if err != nil {
		return xerrors.Errorf("etcd client operation error %w:", err)
	}
	if !createResp.Succeeded {
		return tcErr.ErrAlreadyExists
	}
	return nil
}

func doUpdate(ctx context.Context, client *clientv3.Client, rev int64, key string, value string) error {
RETRY:
	exists := clientv3.Compare(clientv3.Version(key), ">", 0)
	isNotUpdated := []clientv3.Cmp{clientv3.Compare(clientv3.ModRevision(key), "=", rev)}
	updateOp := []clientv3.Op{clientv3.OpPut(key, value)}
	updateItem := clientv3.OpTxn(isNotUpdated, updateOp, nil)
	updateResp, err := client.Txn(ctx).
		If(exists).
		Then(updateItem).
		Commit()
	if err != nil {
		return xerrors.Errorf("etcd client operation error %w:", err)
	}
	if !updateResp.Succeeded {
		return tcErr.ErrNotFound
	}
	// the item has been updated
	txnResp := updateResp.Responses[0].GetResponseTxn()
	if !txnResp.Succeeded {
		_, latestRev, err := doGetWithRev(ctx, client, key)
		if err != nil {
			return err
		}
		rev = latestRev
		goto RETRY
	}
	return nil
}

func doDelete(ctx context.Context, client *clientv3.Client, key string) error {
	exists := clientv3.Compare(clientv3.Version(key), ">", 0)
	deleteItem := clientv3.OpDelete(key)
	deleteResp, err := client.Txn(ctx).
		If(exists).
		Then(deleteItem).
		Commit()
	if err != nil {
		return xerrors.Errorf("etcd client operation error %w:", err)
	}
	if !deleteResp.Succeeded {
		return tcErr.ErrNotFound
	}
	return nil
}
