package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
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

func executeCommand(w http.ResponseWriter, r *http.Request) {
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
	log.Printf("Received command: %s, Time: %s, User-Agent: %s\n", req.Command, time.Now().Format(time.RFC3339), r.UserAgent())

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

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile) // 设置日志格式

	logFile, err := os.OpenFile("./log/run.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Cannot open log file: ", err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)

	http.HandleFunc("/execute", executeCommand)
	log.Println("Server starting on port 6329...")
	http.ListenAndServe(":6329", nil)
}
