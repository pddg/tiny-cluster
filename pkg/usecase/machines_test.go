package usecase_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"golang.org/x/xerrors"

	"github.com/pddg/tiny-cluster/pkg/models"
	"github.com/pddg/tiny-cluster/pkg/repositories/mock"
	"github.com/pddg/tiny-cluster/pkg/usecase"
)

func Test_machineUseCaseImpl_GetMachineByName(t *testing.T) {
	sampleErr := xerrors.Errorf("Sample error")
	testCases := map[string]struct {
		fixtures   []*models.Machine
		errFixture error
		name       string
		expect     *models.Machine
		expectErr  error
	}{
		"match": {
			fixtures:   machineFixtures,
			errFixture: nil,
			name:       "machine1",
			expect:     machineFixtures[0],
			expectErr:  nil,
		},
		"do not match": {
			fixtures:   machineFixtures,
			errFixture: nil,
			name:       "not found",
			expect:     nil,
			expectErr:  nil,
		},
		"empty": {
			fixtures:   []*models.Machine{},
			errFixture: nil,
			name:       "machine1",
			expect:     nil,
			expectErr:  nil,
		},
		"error": {
			fixtures:   machineFixtures,
			errFixture: sampleErr,
			name:       "machine1",
			expect:     nil,
			expectErr:  sampleErr,
		},
	}
	for tn, tc := range testCases {
		ctx := context.TODO()
		ctrl := gomock.NewController(t)
		t.Run(tn, func(t *testing.T) {
			t.Parallel()
			repoMock := mock.NewMockMachineRepository(ctrl)
			repoMock.EXPECT().GetMachines(ctx).Return(tc.fixtures, tc.errFixture)
			machineUseCase := usecase.NewMachineUseCase(repoMock)
			actual, err := machineUseCase.GetMachineByName(ctx, tc.name)
			if err != tc.expectErr {
				t.Errorf("Invalid error. Expected: %#v, Actual: %#v", tc.expectErr, err)
				return
			}
			if !reflect.DeepEqual(tc.expect, actual) {
				t.Errorf("Invalid response. Expected: %#v, Actual: %#v", tc.expect, actual)
			}
		})
	}
}

func Test_machineUseCaseImpl_GetAllMachines(t *testing.T) {
	sampleErr := xerrors.Errorf("Sample error")
	testCases := map[string]struct {
		fixtures   []*models.Machine
		errFixture error
		expect     []*models.Machine
		expectErr  error
	}{
		"return machines normally": {
			fixtures:   machineFixtures,
			errFixture: nil,
			expect:     machineFixtures,
			expectErr:  nil,
		},
		"empty": {
			fixtures:   []*models.Machine{},
			errFixture: nil,
			expect:     []*models.Machine{},
			expectErr:  nil,
		},
		"error": {
			fixtures:   nil,
			errFixture: sampleErr,
			expect:     nil,
			expectErr:  sampleErr,
		},
	}
	ctx := context.TODO()
	ctrl := gomock.NewController(t)
	for tn, tc := range testCases {
		t.Run(tn, func(t *testing.T) {
			t.Parallel()
			repoMock := mock.NewMockMachineRepository(ctrl)
			repoMock.EXPECT().GetMachines(ctx).Return(tc.fixtures, tc.errFixture)
			machineUseCase := usecase.NewMachineUseCase(repoMock)
			actual, err := machineUseCase.GetAllMachines(ctx)
			if err != tc.expectErr {
				t.Errorf("Invalid error. Expected: %#v, Actual: %#v", err, tc.expectErr)
				return
			}
			if !reflect.DeepEqual(actual, tc.expect) {
				t.Errorf("Invalid response. Expected: %#v, Actual: %#v", actual, tc.expect)
			}
		})
	}
}

func Test_machineUseCaseImpl_GetMachineByQuery(t *testing.T) {
	var emptyMachines []*models.Machine
	sampleErr := xerrors.Errorf("Sample error")
	testCases := map[string]struct {
		fixtures   []*models.Machine
		errFixture error
		query      *usecase.MachineQuery
		expect     []*models.Machine
		expectErr  error
	}{
		"query by name": {
			fixtures:   machineFixtures,
			errFixture: nil,
			query:      &usecase.MachineQuery{"name": machineFixtures[0].Name},
			expect:     machineFixtures[0:1],
			expectErr:  nil,
		},
		"query by ip addr": {
			fixtures:   machineFixtures,
			errFixture: nil,
			query:      &usecase.MachineQuery{"ipv4": machineFixtures[0].IPv4Addr},
			expect:     machineFixtures[0:1],
			expectErr:  nil,
		},
		"query by MAC": {
			fixtures:   machineFixtures,
			errFixture: nil,
			query:      &usecase.MachineQuery{"mac": machineFixtures[0].MAC},
			expect:     machineFixtures[0:1],
			expectErr:  nil,
		},
		"not found query": {
			fixtures:   machineFixtures,
			errFixture: nil,
			query:      &usecase.MachineQuery{"mac": "not found"},
			expect:     emptyMachines,
			expectErr:  nil,
		},
		"complex query (match)": {
			fixtures:   machineFixtures,
			errFixture: nil,
			query: &usecase.MachineQuery{
				"mac":  machineFixtures[0].MAC,
				"name": machineFixtures[0].Name,
				"ipv4": machineFixtures[0].IPv4Addr,
			},
			expect:    machineFixtures[0:1],
			expectErr: nil,
		},
		"complex query (not match)": {
			fixtures:   machineFixtures,
			errFixture: nil,
			query: &usecase.MachineQuery{
				"mac":  machineFixtures[0].MAC,
				"name": "not found",
				"ipv4": "not found",
			},
			expect:    machineFixtures[0:1],
			expectErr: nil,
		},
		"query empty machines": {
			fixtures:   []*models.Machine{},
			errFixture: nil,
			query:      &usecase.MachineQuery{"name": machineFixtures[0].Name},
			expect:     emptyMachines,
			expectErr:  nil,
		},
		"error": {
			fixtures:   nil,
			errFixture: sampleErr,
			expect:     nil,
			expectErr:  sampleErr,
		},
	}
	ctx := context.TODO()
	ctrl := gomock.NewController(t)
	for tn, tc := range testCases {
		t.Run(tn, func(t *testing.T) {
			t.Parallel()
			repoMock := mock.NewMockMachineRepository(ctrl)
			repoMock.EXPECT().GetMachines(ctx).Return(tc.fixtures, tc.errFixture)
			machineUseCase := usecase.NewMachineUseCase(repoMock)
			actual, err := machineUseCase.GetMachineByQuery(ctx, tc.query)
			if err != tc.expectErr {
				t.Errorf("Invalid error. Expected: %#v, Actual: %#v", err, tc.expectErr)
				return
			}
			if !reflect.DeepEqual(actual, tc.expect) {
				t.Errorf("Invalid response. Expected: %#v, Actual: %#v", actual, tc.expect)
			}
		})
	}
}

func Test_machineUseCaseImpl_RegisterOrUpdateMachine(t *testing.T) {
	sampleErr := xerrors.Errorf("Sample error")
	testCases := map[string]struct {
		fixtures   []*models.Machine
		errFixture error
		isUpdate bool
		machine    *models.Machine
		expect     error
	}{
		"register normally": {
			fixtures:   []*models.Machine{},
			errFixture: nil,
			isUpdate: false,
			machine:    machineFixtures[0],
			expect:     nil,
		},
		"update normally": {
			fixtures:   machineFixtures,
			errFixture: nil,
			isUpdate: true,
			machine:    machineFixtures[0],
			expect:     nil,
		},
		"error": {
			fixtures:   machineFixtures,
			errFixture: sampleErr,
			isUpdate: false,
			machine:    machineFixtures[0],
			expect:     sampleErr,
		},
	}
	ctx := context.TODO()
	ctrl := gomock.NewController(t)
	for tn, tc := range testCases {
		t.Run(tn, func(t *testing.T) {
			t.Parallel()
			repoMock := mock.NewMockMachineRepository(ctrl)
			repoMock.EXPECT().GetMachines(ctx).Return(tc.fixtures, tc.errFixture)
			if tc.isUpdate {
				repoMock.EXPECT().UpdateMachine(ctx, tc.machine).Return(tc.errFixture)
			} else {
				repoMock.EXPECT().RegisterMachine(ctx, tc.machine).Return(tc.errFixture)
			}
			machineUseCase := usecase.NewMachineUseCase(repoMock)
			actual := machineUseCase.RegisterOrUpdateMachine(ctx, tc.machine)
			if actual != tc.expect {
				t.Errorf("Invalid response. Expected: %#v, Actual: %#v", actual, tc.expect)
			}
		})
	}
}
