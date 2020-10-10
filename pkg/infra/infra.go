package infra

import (
	"context"
	"fmt"

	"go.etcd.io/etcd/clientv3"
    "golang.org/x/xerrors"
)

const BASE_PREFIX = "/tiny-cluster"

var (
    ErrNotFound = xerrors.New("The item was not found")
    ErrAlreadyExists = xerrors.New("The item has already been exist")
)

type baseRepoImpl struct {
    config *clientv3.Config
}

func (r *baseRepoImpl) newClient(ctx context.Context) (*clientv3.Client, error) {
    timeoutCtx, cancel := context.WithTimeout(ctx, 30)
    defer cancel()
    errCh := make(chan error, 1)
    clientCh := make(chan *clientv3.Client)
    go func() {
        client, err := clientv3.New(*r.config)
        if err != nil {
            errCh <- err
        } else {
            clientCh <- client
        }
        close(clientCh)
        close(errCh)
    }()
    select {
    case <-timeoutCtx.Done():
        return nil, fmt.Errorf("Connection timed out")
    case err := <-errCh:
        return nil, err
    case client := <-clientCh:
        return client, nil
    }
}

func doGetWithRev(ctx context.Context, client *clientv3.Client, key string, opts ...clientv3.OpOption) (string, int64, error) {
    resp, err := client.Get(ctx, key, opts...)
    if err != nil {
        return "", 0, err
    }
    if resp.Count == 0 {
        return "", 0, ErrNotFound
    }
    return resp.Kvs[0].String(), resp.Kvs[0].ModRevision, nil
}

func doGet(ctx context.Context, client *clientv3.Client, key string, opts ...clientv3.OpOption) (string, error) {
    value, _, err := doGetWithRev(ctx, client, key, opts...)
    return value, err
}

func doCreate(ctx context.Context, client *clientv3.Client, key string, value string) error {
    doesNotExist := clientv3.Compare(clientv3.Version(key), "=", 0)
    create := clientv3.OpPut(key, value)
    createResp, err := client.Txn(ctx).
        If(doesNotExist).
        Then(create).
        Commit()
    if err != nil {
        return err
    }
    if !createResp.Succeeded {
        return ErrAlreadyExists
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
        return err
    }
    if !updateResp.Succeeded {
        return ErrNotFound
    }
    // the item has been updated
    txnResp := updateResp.Responses[0].GetResponseTxn()
    if !txnResp.Succeeded {
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
        return err
    }
    if !deleteResp.Succeeded {
        return ErrNotFound
    }
    return nil
}

