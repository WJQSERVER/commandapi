package forward

import (
	"CommandApi/logger"
	"encoding/json"
	"net/http"
	"os/exec"
	"time"
)

type CommandRequest struct {
	Command string `json:"command"`
}

type CommandResponse struct {
	Stdout     string `json:"stdout"`
	Stderr     string `json:"stderr"`
	ReturnCode int    `json:"returncode"`
	Error      string `json:"error,omitempty"`
}

// 日志模块
var (
	logw       = logger.Logw
	logInfo    = logger.LogInfo
	LogWarning = logger.LogWarning
	logError   = logger.LogError
)

func ExecuteCommand(w http.ResponseWriter, r *http.Request) {
	var req CommandRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if req.Command == "" {
		http.Error(w, "No command provided", http.StatusBadRequest)
		return
	}

	// 记录请求信息
	logInfo("Received command: %s, Time: %s, User-Agent: %s, Method: %s, RemoteAddr: %s\n", req.Command, time.Now().Format(time.RFC3339), r.UserAgent(), r.Method, r.RemoteAddr)

	cmd := exec.Command("sh", "-c", req.Command)
	stdout, err := cmd.CombinedOutput()
	returnCode := cmd.ProcessState.ExitCode()

	response := CommandResponse{
		Stdout:     string(stdout),
		Stderr:     "",
		ReturnCode: returnCode,
	}

	if err != nil {
		response.Error = err.Error()
		response.Stderr = string(stdout)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
