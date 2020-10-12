//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock
package usecase

import (
	"context"

	"github.com/pddg/tiny-cluster/pkg/models"
	"github.com/pddg/tiny-cluster/pkg/repositories"
)

type MachineUsecase interface {
	GetAllMachines(ctx context.Context) ([]*models.Machine, error)
	GetMachineByName(ctx context.Context, name string) (*models.Machine, error)
	GetMachineByQuery(ctx context.Context, query *MachineQuery) ([]*models.Machine, error)
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

func NewMachineUseCase(repo repositories.MachineRepository) MachineUsecase {
	return &machineUseCaseImpl{
		repo: repo,
	}
}
