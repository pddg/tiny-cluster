package usecase_test

import (
	"testing"
	"time"

	"github.com/pddg/tiny-cluster/pkg/models"
	"github.com/pddg/tiny-cluster/pkg/usecase"
)

var (
	fixtures = []*models.Machine{
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
	}
)

func Test_MachineQueryMatch(t *testing.T) {
	testCases := map[string]struct {
		target *models.Machine
		query  *usecase.MachineQuery
		expect bool
	}{
		"match by name": {
			target: fixtures[0],
			query:  &usecase.MachineQuery{"name": fixtures[0].Name},
			expect: true,
		},
		"match by ipv4 addr": {
			target: fixtures[0],
			query:  &usecase.MachineQuery{"ipv4": fixtures[0].IPv4Addr},
			expect: true,
		},
		"match by mac": {
			target: fixtures[0],
			query:  &usecase.MachineQuery{"mac": fixtures[0].MAC},
			expect: true,
		},
		"match by name and addr and mac": {
			target: fixtures[0],
			query: &usecase.MachineQuery{
				"mac":  fixtures[0].MAC,
				"name": fixtures[0].Name,
				"ipv4": fixtures[0].IPv4Addr,
			},
			expect: true,
		},
		"do not match": {
			target: fixtures[0],
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
