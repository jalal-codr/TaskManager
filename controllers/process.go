package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strconv"
	"taskManager/types"

	"github.com/shirou/gopsutil/process"
)

func StartProcessHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req types.StartProcessRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create and start the command
	cmd := exec.Command(req.Command, req.Args...)
	if err := cmd.Start(); err != nil {
		response := types.ProcessResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to start process: %v", err),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	response := types.ProcessResponse{
		Success: true,
		Message: "Process started successfully",
		PID:     cmd.Process.Pid,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func StopProcessHandler(w http.ResponseWriter, r *http.Request) {
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
		response := types.ProcessResponse{
			Success: false,
			Message: fmt.Sprintf("Process not found: %v", err),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	if err := p.Kill(); err != nil {
		response := types.ProcessResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to stop process: %v", err),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	response := types.ProcessResponse{
		Success: true,
		Message: fmt.Sprintf("Process %d stopped successfully", pid),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
