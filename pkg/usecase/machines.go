package usecase

import (
	"context"

	"github.com/pddg/tiny-cluster/pkg/models"
	"github.com/pddg/tiny-cluster/pkg/repositories"
)

type MachineUsecase interface {
	GetAllMachines(ctx context.Context) ([]*models.Machine, error)
	GetMachineByName(ctx context.Context, name string) (*models.Machine, error)
	GetMachineByQuery(ctx context.Context, query map[string]string) ([]*models.Machine, error)
}

type machineUseCaseImpl struct {
	repo repositories.MachineRepository
}

func (m *machineUseCaseImpl) GetAllMachines(ctx context.Context) ([]*models.Machine, error) {
	return m.repo.GetMachines(ctx)
}

func (m *machineUseCaseImpl) GetMachineByQuery(ctx context.Context, query map[string]string) ([]*models.Machine, error) {
	machines, err := m.repo.GetMachines(ctx)
	if err != nil {
		return nil, err
	}
	var matchedMachines []*models.Machine
	for _, machine := range machines {
		if mac := query["mac"]; len(mac) > 0 {
			if machine.MAC == mac {
				matchedMachines = 
			}
		}
	}
}
