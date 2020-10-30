package usecase

import (
	"strings"

	"github.com/pddg/tiny-cluster/pkg/models"
)

// MachineQuery indicate that the query to filter the machine instance.
type MachineQuery map[string]string

func (q MachineQuery) _match(machine *models.Machine, and bool) bool {
	var match bool
	for k, v := range q {
		switch {
		case k == "mac":
			match = machine.MAC == v
		case k == "name":
			match = machine.Name == v
		case k == "ipv4":
			match = machine.IPv4Addr == v
		default:
			continue
		}
		if match {
			if !and {
				return true
			}
		} else {
			if and {
				break
			}
		}
	}
	return match
}

// Match returns true if all factors in the query matches the machine.
func (q MachineQuery) Match(machine *models.Machine) bool {
	queryMap := map[string]string(q)
	and, ok := queryMap["and"]
	if !ok {
		return q._match(machine, true)
	}
	return q._match(machine, strings.ToLower(and) == "true")
}
