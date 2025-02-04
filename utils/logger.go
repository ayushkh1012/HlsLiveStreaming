package utils

import (
    "log"
    "os"
)

type Logger struct {
    *log.Logger
}

func NewLogger() *Logger {
    return &Logger{
        Logger: log.New(os.Stdout, "[HLS-Server] ", log.LstdFlags|log.Lshortfile),
    }
}
