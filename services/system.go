package services

import (
	"taskManager/types"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

func GetSystemStats() (types.SystemStats, error) {
	cpuPercent, err := cpu.Percent(time.Second, false)
	if err != nil {
		return types.SystemStats{}, err
	}

	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return types.SystemStats{}, err
	}

	cpuCount, err := cpu.Counts(true)
	if err != nil {
		return types.SystemStats{}, err
	}

	return types.SystemStats{
		CPUPercent: cpuPercent[0],
		MemoryInfo: memInfo,
		CPUCount:   cpuCount,
	}, nil
}
