package services

import (
	"taskManager/types"

	"github.com/shirou/gopsutil/process"
)

func GetProcessInfo() ([]types.ProcessInfo, error) {
	processes, err := process.Processes()
	if err != nil {
		return nil, err
	}

	var processInfos []types.ProcessInfo
	for _, p := range processes {
		name, err := p.Name()
		if err != nil {
			continue
		}

		cpuPercent, err := p.CPUPercent()
		if err != nil {
			cpuPercent = 0
		}

		memPercent, err := p.MemoryPercent()
		if err != nil {
			memPercent = 0
		}

		status, err := p.Status()
		if err != nil {
			status = "unknown"
		}

		createTime, err := p.CreateTime()
		if err != nil {
			createTime = 0
		}

		processInfos = append(processInfos, types.ProcessInfo{
			PID:        p.Pid,
			Name:       name,
			CPUPercent: cpuPercent,
			MemPercent: memPercent,
			Status:     status,
			CreateTime: createTime,
		})
	}

	return processInfos, nil
}
