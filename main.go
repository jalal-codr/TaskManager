package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/process"
)

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

type APIResponse struct {
	Timestamp   string        `json:"timestamp"`
	SystemStats SystemStats   `json:"system_stats"`
	Processes   []ProcessInfo `json:"processes"`
}

type StartProcessRequest struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
}

type ProcessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	PID     int    `json:"pid,omitempty"`
}

func getProcessInfo() ([]ProcessInfo, error) {
	processes, err := process.Processes()
	if err != nil {
		return nil, err
	}

	var processInfos []ProcessInfo
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

		processInfos = append(processInfos, ProcessInfo{
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

func getSystemStats() (SystemStats, error) {
	cpuPercent, err := cpu.Percent(time.Second, false)
	if err != nil {
		return SystemStats{}, err
	}

	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return SystemStats{}, err
	}

	cpuCount, err := cpu.Counts(true)
	if err != nil {
		return SystemStats{}, err
	}

	return SystemStats{
		CPUPercent: cpuPercent[0],
		MemoryInfo: memInfo,
		CPUCount:   cpuCount,
	}, nil
}

func systemHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	processes, err := getProcessInfo()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	systemStats, err := getSystemStats()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := APIResponse{
		Timestamp:   time.Now().Format(time.RFC3339),
		SystemStats: systemStats,
		Processes:   processes,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func startProcessHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req StartProcessRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create and start the command
	cmd := exec.Command(req.Command, req.Args...)
	if err := cmd.Start(); err != nil {
		response := ProcessResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to start process: %v", err),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	response := ProcessResponse{
		Success: true,
		Message: "Process started successfully",
		PID:     cmd.Process.Pid,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func stopProcessHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	pidStr := r.URL.Query().Get("pid")
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		http.Error(w, "Invalid PID", http.StatusBadRequest)
		return
	}

	p, err := process.NewProcess(int32(pid))
	if err != nil {
		response := ProcessResponse{
			Success: false,
			Message: fmt.Sprintf("Process not found: %v", err),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	if err := p.Kill(); err != nil {
		response := ProcessResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to stop process: %v", err),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	response := ProcessResponse{
		Success: true,
		Message: fmt.Sprintf("Process %d stopped successfully", pid),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	// Set up routes
	http.HandleFunc("/api/system", systemHandler)
	http.HandleFunc("/api/process/start", startProcessHandler)
	http.HandleFunc("/api/process/stop", stopProcessHandler)

	port := ":8080"
	fmt.Printf("Server starting on port %s...\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
