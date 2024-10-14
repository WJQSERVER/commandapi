package main

import (
	"CommandApi/config"
	"CommandApi/forward"
	"CommandApi/logger"
	"flag"
	"fmt"
	"log"
	"net/http"
)

var (
	cfg        *config.Config
	configfile = "/data/go/config/config.yaml"
)

// 日志模块
var (
	logw       = logger.Logw
	logInfo    = logger.LogInfo
	LogWarning = logger.LogWarning
	logError   = logger.LogError
)

func ReadFlag() {
	cfgfile := flag.String("cfg", configfile, "config file path")
	configfile = *cfgfile
}

func loadConfig() {
	var err error
	// 初始化配置
	cfg, err = config.LoadConfig(configfile)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	fmt.Printf("Loaded config: %v\n", cfg)
}

func setupLogger() {
	// 初始化日志模块
	err := logger.Init(cfg.Log.LogFilePath, cfg.Log.MaxLogSize) // 传递日志文件路径
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	logInfo("Logger initialized")
	logInfo("Init Completed")
}

func init() {
	ReadFlag()
	loadConfig()
	setupLogger()
}

func main() {
	defer logger.Close() // 确保在退出时关闭日志文件
	http.HandleFunc("/execute", forward.ExecuteCommand)
	logInfo("Server starting on port %d", cfg.Server.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", cfg.Server.Port), nil)
}
