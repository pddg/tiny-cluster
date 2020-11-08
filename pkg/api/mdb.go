package api

import (
	"context"

	"github.com/pddg/tiny-cluster/pkg/api/pb"
	"github.com/pddg/tiny-cluster/pkg/models"
	"github.com/pddg/tiny-cluster/pkg/usecase"
)

type machineDatabaseServer struct {
	pb.UnimplementedMachineDatabaseServer
	machineUsecase usecase.MachineUsecase
}

func (mdb *machineDatabaseServer) GetMachines(ctx context.Context, req *pb.GetMachinesRequest) (*pb.GetMachinesResponse, error) {
	var queries usecase.MachineQuery
	for _, item := range req.GetQueries() {
		queries[item.GetKey()] = item.GetValue()
	}
	var resp *pb.GetMachinesResponse
	machines, err := mdb.machineUsecase.GetMachineByQuery(ctx, &queries)
	if err != nil {
		return resp, err
	}
	for _, machine := range machines {
		m := &pb.Machine{
			Name:         machine.Name,
			Ipv4Addr:     machine.IPv4Addr,
			Mac:          machine.MAC,
			DeployedDate: machine.DeployedDate,
			Spec: &pb.MachineSpec{
				Memory: machine.Spec.Memory,
				Disk:   machine.Spec.Disk,
				Core:   machine.Spec.Core,
			},
		}
		resp.Machines = append(resp.Machines, m)
	}
	return resp, nil
}

func (mdb *machineDatabaseServer) RegisterOrUpdateMachine(ctx context.Context, req *pb.RegisterOrUpdateMachineRequest) (*pb.RegisterOrUpdateMachineResponse, error) {
	machine := &models.Machine{
		Name:     req.Machine.Name,
		IPv4Addr: req.Machine.Ipv4Addr,
		MAC:      req.Machine.Mac,
		Spec: models.MachineSpec{
			Core:   req.Machine.GetSpec().Core,
			Memory: req.Machine.GetSpec().Memory,
			Disk:   req.Machine.GetSpec().Disk,
		},
	}
	var resp *pb.RegisterOrUpdateMachineResponse
	err := mdb.machineUsecase.RegisterOrUpdateMachine(ctx, machine)
	if err != nil {
		resp.Success = false
		resp.Message = err.Error()
		return resp, err
	}
	resp.Success = true
	resp.Message = "ok"
	return resp, nil
}

func (mdb *machineDatabaseServer) DeleteMachine(ctx context.Context, req *pb.DeleteMachineRequest) (*pb.DeleteMachineResponse, error) {
	machine := &models.Machine{
		Name:     req.Machine.Name,
		IPv4Addr: req.Machine.Ipv4Addr,
		MAC:      req.Machine.Mac,
		Spec: models.MachineSpec{
			Core:   req.Machine.GetSpec().Core,
			Memory: req.Machine.GetSpec().Memory,
			Disk:   req.Machine.GetSpec().Disk,
		},
	}
	var resp *pb.DeleteMachineResponse
	err := mdb.machineUsecase.DeleteMachine(ctx, machine)
	if err != nil {
		resp.Success = false
		resp.Message = err.Error()
		return resp, err
	}
	resp.Success = true
	resp.Message = "ok"
	return resp, nil
}

// NewMachineDatabaseServer returns actual implementation of MachineDatabaseServer.
func NewMachineDatabaseServer(machineUsecase usecase.MachineUsecase) pb.MachineDatabaseServer {
	return &machineDatabaseServer{
		machineUsecase: machineUsecase,
	}
}
