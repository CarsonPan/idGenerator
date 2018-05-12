package idGenerator

import (
	"encoding/binary"
	"encoding/hex"
	"sync"
	"time"

	"idGenerator/machine"
)

const (
	seqMask       uint64 = 4095
	machineIdMask uint64 = 127
	processIdMask uint64 = 7
)

var mac machine.Machine
var machineId uint64
var processId uint64
var staticSeqId uint64
var lastTimestamp int64
var minTimestamp int64
var mu sync.Mutex

func init() {
	mac = machine.CreateMachine()
	machineId = mac.GetMachineId()
	processId = mac.GetProcessId()
	if machineId > machineIdMask {
		panic("machineId 范围为0-127")
	}
	if processId > processIdMask {
		panic("processId 范围为0-7")
	}

	machineId = machineId << 57 >> 42
	processId = processId << 61 >> 49
	staticSeqId = 0
	lastTimestamp = time.Now().UnixNano() / int64(time.Millisecond)
	minTime, _ := time.Parse("2006-01-02 15:04:05", "2018-05-11 00:00:00")
	minTimestamp = minTime.UnixNano() / int64(time.Millisecond)
	mu = sync.Mutex{}
}

func GenerateNewStringId() string {
	id := GenerateNewId()
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(id))
	return hex.EncodeToString(buf)
}

func GenerateNewId() int64 {
	var timestamp int64
	var seqId uint64
	mu.Lock()
flag:
	timestamp = time.Now().UnixNano() / int64(time.Millisecond)
	if lastTimestamp < timestamp {
		lastTimestamp = timestamp
		staticSeqId = 0
	}
	if staticSeqId > seqMask {
		time.Sleep(1 * time.Millisecond)
		goto flag
	}
	seqId = staticSeqId
	staticSeqId++
	mu.Unlock()
	return pack(uint64(timestamp), seqId)
}

func pack(timestamp uint64, seqId uint64) int64 {
	//雪花算法
	timestamp = timestamp << 22
	return int64(timestamp | machineId | processId | seqId)

}
