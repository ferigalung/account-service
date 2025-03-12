package logger

import (
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v2/log"
)

type LogData struct {
	Time    string      `json:"time"`
	Level   string      `json:"level"`
	Message string      `json:"message"`
	Context interface{} `json:"context,omitempty"`
}

func Log(level string, message string, ctx interface{}) {
	data := LogData{
		Time:    time.Now().Format(time.RFC3339),
		Level:   level,
		Message: message,
		Context: ctx,
	}
	jsonLog, _ := json.Marshal(data)

	switch level {
	case "info":
		log.Info(string(jsonLog))
	case "warn":
		log.Warn(string(jsonLog))
	case "error":
		log.Error(string(jsonLog))
	case "fatal":
		log.Fatal(string(jsonLog))
	default:
		log.Info(string(jsonLog))
	}
}
