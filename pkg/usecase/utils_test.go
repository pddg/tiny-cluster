package usecase_test

import (
	"testing"
	"time"

	"github.com/pddg/tiny-cluster/pkg/models"
	"github.com/pddg/tiny-cluster/pkg/usecase"
)

var (
	machineFixtures = []*models.Machine{
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

func Test_MachineQueryMatch(t *testing.T) {
	testCases := map[string]struct {
		target *models.Machine
		query  *usecase.MachineQuery
		expect bool
	}{
		"match by name": {
			target: machineFixtures[0],
			query:  &usecase.MachineQuery{"name": machineFixtures[0].Name},
			expect: true,
		},
		"match by ipv4 addr": {
			target: machineFixtures[0],
			query:  &usecase.MachineQuery{"ipv4": machineFixtures[0].IPv4Addr},
			expect: true,
		},
		"match by mac": {
			target: machineFixtures[0],
			query:  &usecase.MachineQuery{"mac": machineFixtures[0].MAC},
			expect: true,
		},
		"match by name and addr and mac": {
			target: machineFixtures[0],
			query: &usecase.MachineQuery{
				"mac":  machineFixtures[0].MAC,
				"name": machineFixtures[0].Name,
				"ipv4": machineFixtures[0].IPv4Addr,
			},
			expect: true,
		},
		"only name is matched (and)": {
			target: machineFixtures[0],
			query: &usecase.MachineQuery{
				"mac":  "not found",
				"name": machineFixtures[0].Name,
				"ipv4": "not found",
			},
			expect: false,
		},
		"only name is matched (or)": {
			target: machineFixtures[0],
			query: &usecase.MachineQuery{
				"mac":  "not found",
				"name": machineFixtures[0].Name,
				"ipv4": "not found",
				"and":  "false",
			},
			expect: true,
		},
		"do not match": {
			target: machineFixtures[0],
			query: &usecase.MachineQuery{
				"mac":  "don't match",
				"name": "don't match",
				"ipv4": "don't match",
			},
			expect: false,
		},
	}
	for tn, tc := range testCases {
		t.Run(tn, func(t *testing.T) {
			t.Parallel()
			actual := tc.query.Match(tc.target)
			if actual != tc.expect {
				t.Errorf("Matching result is invalid. Expected: %v, Actual: %v", tc.expect, actual)
			}
		})
	}
}
