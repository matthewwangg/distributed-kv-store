package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type LogWriter struct {
	ID     string
	Addr   string
	Writer io.Writer
}

func (lw LogWriter) Write(p []byte) (int, error) {
	timestamp := time.Now().Format("15:04:05")

	msg := string(p)
	msg = strings.TrimRight(msg, "\r\n")

	line := fmt.Sprintf("[%s] [INFO] [%s] [%s] %s\n", timestamp, lw.ID, lw.Addr, msg)
	return lw.Writer.Write([]byte(line))
}

func SetupLogger(nodeID, peerAddr string) error {
	date := time.Now().Format("2006-01-02")
	dir := filepath.Join("logs", nodeID)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create logs directory: %w", err)
	}

	logFilePath := filepath.Join(dir, date+".log")
	f, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	log.SetFlags(0)
	log.SetPrefix("")
	log.SetOutput(LogWriter{ID: nodeID, Addr: peerAddr, Writer: f})

	return nil
}
