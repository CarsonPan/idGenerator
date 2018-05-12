package machine

type Machine interface {
	GetMachineId() uint64
	GetProcessId() uint64
}
