package infra

import (
	"context"
	"encoding/json"
	"path"
	"reflect"
	"testing"
	"time"

	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
	"go.etcd.io/etcd/clientv3"
	"golang.org/x/xerrors"

	tcErr "github.com/pddg/tiny-cluster/pkg/errors"
	"github.com/pddg/tiny-cluster/pkg/models"
)

type machineFixtureImpl []*models.Machine

func (mf *machineFixtureImpl) prepare(ctx context.Context, t *testing.T, client *clientv3.Client) {
	t.Helper()
	for _, v := range *mf {
		valueByte, _ := json.Marshal(v)
		_, err := client.Put(ctx, path.Join(machinePrefix, v.MAC), string(valueByte))
		if err != nil {
			t.Errorf("Failed to put value due to %v", err)
		}
	}
}

func (mf *machineFixtureImpl) clean(ctx context.Context, t *testing.T, client *clientv3.Client) {
	for _, v := range *mf {
		_, err := client.Delete(ctx, path.Join(machinePrefix, v.MAC))
		if err != nil {
			switch {
			case xerrors.Is(err, rpctypes.ErrKeyNotFound):
				continue
			default:
				t.Errorf("Failed to put value due to %v", err)
			}
		}
	}
}

func (mf *machineFixtureImpl) toSlice() []*models.Machine {
	return *mf
}

var (
	machineFixtures = &machineFixtureImpl{
		{
			Name:         "machine1",
			MAC:          "mac1",
			IPv4Addr:     "19.168.0.2",
			DeployedDate: time.Date(2020, 10, 1, 0, 0, 0, 0, time.UTC).Unix(),
			Spec: models.MachineSpec{
				Core:   4,
				Memory: 2048,
				Disk:   128,
			},
		},
		{
			Name:         "machine2",
			MAC:          "mac2",
			IPv4Addr:     "19.168.1.2",
			DeployedDate: time.Date(2020, 10, 2, 0, 0, 0, 0, time.UTC).Unix(),
			Spec: models.MachineSpec{
				Core:   2,
				Memory: 1024,
				Disk:   64,
			},
		},
	}
)

func Test_machineRepoImpl_GetMachines(t *testing.T) {
	testCases := map[string]struct {
		fixtures  *machineFixtureImpl
		expect    []*models.Machine
		expectErr error
	}{
		"get all": {
			fixtures:  machineFixtures,
			expect:    machineFixtures.toSlice(),
			expectErr: nil,
		},
		"get nothing": {
			fixtures:  &machineFixtureImpl{},
			expect:    []*models.Machine(nil),
			expectErr: nil,
		},
	}
	ctx := context.Background()
	endpoints := getTestEndpoints(t)
	for tn, tc := range testCases {
		t.Run(tn, func(t *testing.T) {
			client := getTestClient(t)
			setUpTest(ctx, t, client, tc.fixtures)
			defer tearDownTest(ctx, t, client, tc.fixtures)
			r := NewMachineRepository(endpoints, 10)
			actual, actualErr := r.GetMachines(ctx)
			if !xerrors.Is(actualErr, tc.expectErr) {
				t.Errorf("Invalid error. Expect: %#v, Actual: %v", tc.expectErr, actualErr)
				return
			}
			if !reflect.DeepEqual(actual, tc.expect) {
				t.Errorf("Invalid response. Expect: %#v, Actual: %#v", tc.expect, actual)
			}
		})
	}
}

func Test_machineRepoImpl_RegisterMachine(t *testing.T) {
	testCases := map[string]struct {
		fixtures *machineFixtureImpl
		machine  *models.Machine
		expect   error
	}{
		"register normally": {
			fixtures: &machineFixtureImpl{},
			machine:  machineFixtures.toSlice()[0],
			expect:   nil,
		},
		"duplicate entry": {
			fixtures: machineFixtures,
			machine:  machineFixtures.toSlice()[0],
			expect:   tcErr.ErrAlreadyExists,
		},
		// This will be error at the usecase layer.
		"insufficient field": {
			fixtures: &machineFixtureImpl{},
			machine:  &models.Machine{Name: "machine"},
			expect:   nil,
		},
	}
	ctx := context.Background()
	endpoints := getTestEndpoints(t)
	r := NewMachineRepository(endpoints, 10)
	for tn, tc := range testCases {
		t.Run(tn, func(t *testing.T) {
			client := getTestClient(t)
			setUpTest(ctx, t, client, tc.fixtures)
			defer func() {
				tearDownTest(ctx, t, client, tc.fixtures)
				tearDownTest(ctx, t, client, &machineFixtureImpl{tc.machine})
			}()
			actual := r.RegisterMachine(ctx, tc.machine)
			if !xerrors.Is(actual, tc.expect) {
				t.Errorf("Invalid error. Expect: %v, Actual: %v", tc.expect, actual)
			}
		})
	}
}

func Test_machineRepoImpl_UpdateMachine(t *testing.T) {
	updatedMachine := machineFixtures.toSlice()[0]
	updatedMachine.Name = "updated"
	testCases := map[string]struct {
		fixtures *machineFixtureImpl
		machine  *models.Machine
		expect   error
	}{
		"update normally": {
			fixtures: machineFixtures,
			machine:  updatedMachine,
			expect:   nil,
		},
		// This will be error at the usecase layer.
		"insufficient field": {
			fixtures: machineFixtures,
			machine:  &models.Machine{MAC: updatedMachine.MAC, Name: "machine"},
			expect:   nil,
		},
		"update non exist item": {
			fixtures: &machineFixtureImpl{},
			machine:  updatedMachine,
			expect:   tcErr.ErrNotFound,
		},
	}
	ctx := context.Background()
	endpoints := getTestEndpoints(t)
	r := NewMachineRepository(endpoints, 10)
	for tn, tc := range testCases {
		t.Run(tn, func(t *testing.T) {
			client := getTestClient(t)
			setUpTest(ctx, t, client, tc.fixtures)
			defer tearDownTest(ctx, t, client, tc.fixtures)
			actual := r.UpdateMachine(ctx, tc.machine)
			if !xerrors.Is(actual, tc.expect) {
				t.Errorf("Invalid error. Expect: %v, Actual: %v", tc.expect, actual)
			}
		})
	}
}

func Test_machineRepoImpl_DeleteMachine(t *testing.T) {
	testCases := map[string]struct {
		fixtures *machineFixtureImpl
		machine  *models.Machine
		expect   error
	}{
		"delete normally": {
			fixtures: machineFixtures,
			machine:  machineFixtures.toSlice()[0],
			expect:   nil,
		},
		"delete non exist item": {
			fixtures: &machineFixtureImpl{},
			machine:  machineFixtures.toSlice()[0],
			expect:   tcErr.ErrNotFound,
		},
	}
	ctx := context.Background()
	endpoints := getTestEndpoints(t)
	r := NewMachineRepository(endpoints, 10)
	for tn, tc := range testCases {
		t.Run(tn, func(t *testing.T) {
			client := getTestClient(t)
			setUpTest(ctx, t, client, tc.fixtures)
			defer tearDownTest(ctx, t, client, tc.fixtures)
			actual := r.DeleteMachine(ctx, tc.machine)
			if !xerrors.Is(actual, tc.expect) {
				t.Errorf("Invalid error. Expect: %v, Actual: %v", tc.expect, actual)
			}
		})
	}
}
