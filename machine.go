package main

import (
	"runtime"
	"sort"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/process"
)

const (
	machineTopProcessLimit   = 5
	machineNoProcesses       = 0
	machineFirstCPUReading   = 0
	machineCPUSampleDuration = 200 * time.Millisecond
)

type MachineStatus struct {
	Hostname          string           `json:"hostname"`
	Platform          string           `json:"platform"`
	PlatformVersion   string           `json:"platformVersion"`
	Architecture      string           `json:"architecture"`
	CPUPercent        *float64         `json:"cpuPercent"`
	MemoryUsedBytes   uint64           `json:"memoryUsedBytes"`
	MemoryUsedPercent *float64         `json:"memoryUsedPercent"`
	LargestProcesses  []MachineProcess `json:"largestProcesses"`
	Network           NetworkStatus    `json:"network"`
}

type MachineProcess struct {
	PID           int32   `json:"pid"`
	Name          string  `json:"name"`
	MemoryBytes   uint64  `json:"memoryBytes"`
	MemoryPercent float64 `json:"memoryPercent"`
}

func (a *App) GetMachineStatus() MachineStatus {
	checkedAt := time.Now()
	networkStatus := buildNetworkStatus(readNetworkLink(), probeInternet(networkHTTPClient, networkProbeURL))
	status := MachineStatus{
		Architecture: runtime.GOARCH,
		Network:      a.networkMonitor.record(networkStatus, checkedAt),
	}
	totalMemory := uint64(machineNoProcesses)

	if hostInfo, err := host.Info(); err == nil {
		status.Hostname = hostInfo.Hostname
		status.Platform = hostInfo.Platform
		status.PlatformVersion = hostInfo.PlatformVersion
	}

	if cpuPercent, err := cpu.Percent(machineCPUSampleDuration, false); err == nil && len(cpuPercent) > machineNoProcesses {
		status.CPUPercent = &cpuPercent[machineFirstCPUReading]
	}

	if memory, err := mem.VirtualMemory(); err == nil {
		totalMemory = memory.Total
		status.MemoryUsedBytes = memory.Used
		status.MemoryUsedPercent = &memory.UsedPercent
	}

	status.LargestProcesses = readLargestProcesses(totalMemory)
	return status
}

func readLargestProcesses(totalMemory uint64) []MachineProcess {
	processes, err := process.Processes()
	if err != nil {
		return []MachineProcess{}
	}

	largestProcesses := make([]MachineProcess, machineNoProcesses, len(processes))
	for _, runningProcess := range processes {
		name, nameError := runningProcess.Name()
		memory, memoryError := runningProcess.MemoryInfo()
		if nameError != nil || memoryError != nil {
			continue
		}

		memoryPercent := float64(machineNoProcesses)
		if totalMemory > machineNoProcesses {
			memoryPercent = float64(memory.RSS) / float64(totalMemory) * networkPercentageScale
		}
		largestProcesses = append(largestProcesses, MachineProcess{
			PID:           runningProcess.Pid,
			Name:          name,
			MemoryBytes:   memory.RSS,
			MemoryPercent: memoryPercent,
		})
	}

	sort.Slice(largestProcesses, func(leftIndex, rightIndex int) bool {
		return largestProcesses[leftIndex].MemoryBytes > largestProcesses[rightIndex].MemoryBytes
	})
	if len(largestProcesses) > machineTopProcessLimit {
		largestProcesses = largestProcesses[:machineTopProcessLimit]
	}
	return largestProcesses
}
