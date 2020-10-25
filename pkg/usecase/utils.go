package usecase

import (
	"github.com/pddg/tiny-cluster/pkg/models"
)

// MachineQuery indicate that the query to filter the machine instance.
type MachineQuery map[string]string

func (q MachineQuery) Match(machine *models.Machine) bool {
	var match bool
	for k, v := range q {
		switch {
		case k == "mac":
			match = machine.MAC == v
		case k == "name":
			match = machine.Name == v
		case k == "ipv4":
			match = machine.IPv4Addr == v
		}
	}
	return match
}
