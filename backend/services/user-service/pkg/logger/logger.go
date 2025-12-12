package logger

import "log"

func Info(v ...interface{})  { log.Println(v...) }
func Error(v ...interface{}) { log.Println(v...) }
