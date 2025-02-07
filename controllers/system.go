package controllers

import (
	"encoding/json"
	"net/http"
	"taskManager/services"
	"taskManager/types"
	"time"
)

func SystemHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	processes, err := services.GetProcessInfo()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	systemStats, err := services.GetSystemStats()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := types.APIResponse{
		Timestamp:   time.Now().Format(time.RFC3339),
		SystemStats: systemStats,
		Processes:   processes,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
