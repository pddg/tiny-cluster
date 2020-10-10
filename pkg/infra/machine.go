package infra

import (
	"context"
	"encoding/json"
	"path"
	"time"

	"github.com/pddg/tiny-cluster/pkg/models"
	repo "github.com/pddg/tiny-cluster/pkg/repositories"
	"go.etcd.io/etcd/clientv3"
)

var machinePrefix = path.Join(BASE_PREFIX, "machines/v1")

type machineRepoImpl struct {
	*baseRepoImpl
}

func (m *machineRepoImpl) getKey(machine *models.Machine) string {
	return path.Join(machinePrefix, machine.ID)
}

func (m *machineRepoImpl) GetMachines(ctx context.Context) ([]*models.Machine, error) {
	var machines []*models.Machine
	client, err := m.newClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()
	resp, err := client.Get(ctx, machinePrefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	for _, kv := range resp.Kvs {
		m := new(models.Machine)
		if err := json.Unmarshal(kv.Value, m); err != nil {
			return nil, err
		}
		machines = append(machines, m)
	}
	return machines, nil
}

func (m *machineRepoImpl) RegisterMachine(ctx context.Context, machine *models.Machine) error {
	client, err := m.newClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()
	key := m.getKey(machine)
	valueByte, err := json.Marshal(machine)
	if err != nil {
		return err
	}
	value := string(valueByte)
	if err := doCreate(ctx, client, key, value); err != nil {
		return err
	}
	return nil
}

func (m *machineRepoImpl) DeleteMachine(ctx context.Context, machine *models.Machine) error {
	client, err := m.newClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()
	key := m.getKey(machine)
	return doDelete(ctx, client, key)
}

func NewMachineRepository(endpoints []string, timeout int) repo.MachineRepository {
	return &machineRepoImpl{
		baseRepoImpl: &baseRepoImpl{
			config: &clientv3.Config{
				Endpoints:   endpoints,
				DialTimeout: time.Duration(timeout),
			},
		},
	}
}
