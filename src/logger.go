package cryptella

import (
	"fmt"
	"log"
	"os"
	"time"
)

type Logger struct {
	file *os.File
}

func NewLogger() (*Logger, error) {

	path := "/Users/buildup/Projects/cryptella/logs"
	filename := fmt.Sprintf("%s%sLogCryptella%s.txt", path, string(os.PathSeparator), time.Now().Format("02012006"))

	// Ensure the directory exists
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to create directory: %v", err)
	}

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return &Logger{
		file: file,
	}, nil
}

func (l *Logger) Log(message string) {

	logEntry := fmt.Sprintf("%s: %s\n", time.Now().Format("02/01/2006 15:04:05"), message)
	print(logEntry)
	if _, err := l.file.WriteString(logEntry); err != nil {
		log.Println("Failed to write log entry:", err)
	}
}

func (l *Logger) Logf(format string, a ...any) {
	l.Log(fmt.Sprintf(format, a...))
	log.Printf(format, a...)
}

func (l *Logger) Close() {
	if err := l.file.Close(); err != nil {
		log.Println("Failed to close log file:", err)
	}
}
