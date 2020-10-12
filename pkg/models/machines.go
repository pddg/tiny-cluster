package models

// MachineSpec is a spec of the host.
type MachineSpec struct {
	// Core is a number of CPU core.
	Core int
	// Memory is an amount of DRAM (MB).
	Memory int
	// Disk is an amount of local disk (GB).
	Disk int
}

// Machine is a information of phisical host
type Machine struct {
	// MAC is Media Access Control address of this host
	MAC string
	// Name is a friendly name.
	Name string
	// IPv4Addr is a IPv4 address of this host.
	IPv4Addr string
	// DeployedDate is a UNIX time of the date when this host is deployed.
	DeployedDate int64
	// Spec indicates the machine spec of the host.
	Spec MachineSpec
}
