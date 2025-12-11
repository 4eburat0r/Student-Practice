package logger

import (
    "log"
    "os"
)

func NewLogger(env string) *log.Logger {
    logger := log.New(os.Stdout, "[AUTH-SERVICE] ", log.LstdFlags|log.Lshortfile)
    
    if env == "production" {
        logger.SetFlags(log.LstdFlags)
    }
    
    return logger
}
