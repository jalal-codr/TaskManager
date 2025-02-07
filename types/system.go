package types

import "github.com/shirou/gopsutil/mem"

type ProcessInfo struct {
	PID        int32   `json:"pid"`
	Name       string  `json:"name"`
	CPUPercent float64 `json:"cpu_percent"`
	MemPercent float32 `json:"memory_percent"`
	Status     string  `json:"status"`
	CreateTime int64   `json:"create_time"`
}

type SystemStats struct {
	CPUPercent float64                `json:"cpu_percent"`
	MemoryInfo *mem.VirtualMemoryStat `json:"memory_info"`
	CPUCount   int                    `json:"cpu_count"`
}
