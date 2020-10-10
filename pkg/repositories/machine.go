package repositories

import (
    "context"

    "github.com/pddg/tiny-cluster/pkg/models"
)

// MachineRepository is a repository about Machine.
type MachineRepository interface {
    // GetMachines returns all machines.
    // This returns empty list and no error if no machines were found.
    GetMachines(ctx context.Context) ([]*models.Machine, error)
    // RegisterMachine creates a record of the machine.
    // This returns error when the item has been created.
    RegisterMachine(ctx context.Context, machine *models.Machine) error
    // DeleteMachine deletes the record of the machine.
    // This returns error when the item does not exist.
    DeleteMachine(ctx context.Context, machine *models.Machine) error
}

