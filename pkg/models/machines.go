package models

// MachineSpec is a spec of the host.
type MachineSpec struct {
	// Core is a number of CPU core.
	Core int32 `json:"core"`
	// Memory is an amount of DRAM (MB).
	Memory int64 `json:"memory"`
	// Disk is an amount of local disk (GB).
	Disk int64 `json:"disk"`
}

// Machine is a information of phisical host
type Machine struct {
	// MAC is Media Access Control address of this host
	MAC string `json:"mac"`
	// Name is a friendly name.
	Name string `json:"name"`
	// IPv4Addr is a IPv4 address of this host.
	IPv4Addr string `json:"ipv4_addr"`
	// DeployedDate is a UNIX time of the date when this host is deployed.
	DeployedDate int64 `json:"deployed_date"`
	// Spec indicates the machine spec of the host.
	Spec MachineSpec `json:"spec"`
}
