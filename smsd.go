package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"./cfg"
	"go.uber.org/zap"
)

func main() {
	configPath := flag.String("c", "./settings.toml", "Config file location")
	rawJSON := []byte(`{
      "level": "debug",
      "encoding": "json",
      "outputPaths": ["./logs/smsd.log"],
      "errorOutputPaths": ["stderr"],
      "encoderConfig": {
        "messageKey": "message",
        "levelKey": "level",
        "levelEncoder": "lowercase"
      }
    }`)

	var cfg_log zap.Config
	if err := json.Unmarshal(rawJSON, &cfg_log); err != nil {
		panic(fmt.Errorf("Error read json config for logger, err: %s", err.Error()))
	}
	logger, err := cfg_log.Build()
	if err != nil {
		panic(fmt.Errorf("Error building logger, err: %s", err.Error()))
	}
	defer logger.Sync()

	flag.Parse()
	cfg := cfg.New(configPath, logger)

	listener, err := net.Listen("tcp", cfg.Host+":"+cfg.Port)
	if err != nil {
		logger.Error("Fail in start daemon", zap.String("err", err.Error()), zap.Any("time", time.Now().Format(time.RFC3339)))
		os.Exit(1)
	}
	logger.Info("Smsd is running on", zap.String("port", cfg.Port), zap.Any("time", time.Now().Format(time.RFC3339)))
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Info("Failed accepting a connection request", zap.String("err", err.Error()), zap.Any("time", time.Now().Format(time.RFC3339)))
		}
		go handleRequest(conn, logger)
	}
}

func handleRequest(conn net.Conn, logger *zap.Logger) {
	message, err := bufio.NewReader(conn).ReadString('\n')
	defer conn.Close()
	if err == io.EOF {
		logger.Error("Reached EOF - close this connection", zap.String("err", err.Error()), zap.String("host", conn.RemoteAddr().String()), zap.Any("time", time.Now().Format(time.RFC3339)))
		return
	}
	conn.Write([]byte(message))
}
