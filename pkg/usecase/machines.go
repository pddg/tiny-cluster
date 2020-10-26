//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock
package usecase

import (
	"context"

	"github.com/pddg/tiny-cluster/pkg/models"
	"github.com/pddg/tiny-cluster/pkg/repositories"
)

// MachineUsecase is the interface to manipulate the machine data.
type MachineUsecase interface {
	// GetAllMachines returns all machines.
	GetAllMachines(ctx context.Context) ([]*models.Machine, error)
	// GetMachineByName returns the machine whose name is matched with the given name.
	GetMachineByName(ctx context.Context, name string) (*models.Machine, error)
	// GetMachineByQuery returns the machine which is filtered by given query.
	GetMachineByQuery(ctx context.Context, query *MachineQuery) ([]*models.Machine, error)
	// RegisterMachine register the machine if it has not been registered.
	RegisterOrUpdateMachine(ctx context.Context, machine *models.Machine) error
}

type machineUseCaseImpl struct {
	repo repositories.MachineRepository
}

func (m *machineUseCaseImpl) GetAllMachines(ctx context.Context) ([]*models.Machine, error) {
	return m.repo.GetMachines(ctx)
}

func (m *machineUseCaseImpl) GetMachineByName(ctx context.Context, name string) (*models.Machine, error) {
	query := &MachineQuery{"name": name}
	machines, err := m.GetMachineByQuery(ctx, query)
	if err != nil {
		return nil, err
	}
	if len(machines) == 0 {
		return nil, nil
	}
	return machines[0], nil
}

func (m *machineUseCaseImpl) GetMachineByQuery(ctx context.Context, query *MachineQuery) ([]*models.Machine, error) {
	machines, err := m.repo.GetMachines(ctx)
	if err != nil {
		return nil, err
	}
	var matchedMachines []*models.Machine
	for _, machine := range machines {
		if query.Match(machine) {
			matchedMachines = append(matchedMachines, machine)
		}
	}
	return matchedMachines, nil
}

func (m *machineUseCaseImpl) RegisterOrUpdateMachine(ctx context.Context, machine *models.Machine) error {
	query := &MachineQuery{
		"name": machine.Name,
		"ipv4": machine.IPv4Addr,
		"mac":  machine.MAC,
		"and":  "false",
	}
	existsMachines, err := m.GetMachineByQuery(ctx, query)
	if err != nil {
		return err
	}
	if len(existsMachines) != 0 {
		return m.repo.UpdateMachine(ctx, machine)
	}
	return m.repo.RegisterMachine(ctx, machine)
}

func NewMachineUseCase(repo repositories.MachineRepository) MachineUsecase {
	return &machineUseCaseImpl{
		repo: repo,
	}
}
