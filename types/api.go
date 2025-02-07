package types

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
