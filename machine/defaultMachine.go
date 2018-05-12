package machine

import (
	"crypto/md5"
	"os"
)

const (
	machineIdMask uint64 = 127
	processIdMask uint64 = 7
)

type DefaultMacthine struct {
}

func (this *DefaultMacthine) GetMachineId() uint64 {
	hostname, err := os.Hostname()
	if err != nil {
		return 0
	}
	hashCode, err := md5.New().Write([]byte(hostname))
	if err != nil {
		return 0
	}
	code := uint64(hashCode) & machineIdMask //只保留后7位 且 0-128
	return code

}

func (this *DefaultMacthine) GetProcessId() uint64 {
	return uint64(os.Getpid()) & processIdMask //只保留后3位 且0-7
}
