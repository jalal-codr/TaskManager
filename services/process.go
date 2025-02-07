package services

import (
	"strings"
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

// List of common system processes to exclude
var systemProcesses = []string{
	"systemd", "init", "svchost.exe", "explorer.exe", "wininit.exe", "services.exe",
	"launchd", "kernel_task", "dock", "windowserver",
}

// Check if the process is a third-party (non-system) service
func isThirdPartyProcess(p *process.Process) bool {
	name, err := p.Name()
	if err != nil {
		return false
	}

	// Exclude known system processes
	for _, sysProc := range systemProcesses {
		if strings.EqualFold(name, sysProc) {
			return false
		}
	}

	// Check process executable path
	exePath, err := p.Exe()
	if err == nil {
		// Exclude processes from system directories
		if strings.Contains(exePath, "/System/Library/") || // macOS system services
			strings.Contains(exePath, "/usr/lib/") || // Linux system services
			strings.Contains(exePath, "C:\\Windows\\System32") { // Windows system services
			return false
		}
	}

	return true
}

// GetProcessInfo returns only third-party (non-system) processes
func GetThirdPartyProcessInfo() ([]types.ProcessInfo, error) {
	processes, err := process.Processes()
	if err != nil {
		return nil, err
	}

	var processInfos []types.ProcessInfo
	for _, p := range processes {
		if !isThirdPartyProcess(p) {
			continue // Skip system processes
		}

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
