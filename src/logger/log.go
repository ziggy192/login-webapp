package logger

import (
	"fmt"
	"log"
)

func Info(v ...any) {
	args := make([]any, 0, 1+len(v))
	args = append(append(args, "[info]"), v...)
	_ = log.Output(2, fmt.Sprintln(args...))
}

func Err(v ...any) {
	args := make([]any, 0, 1+len(v))
	args = append(append(args, "[error] "), v...)
	_ = log.Output(2, fmt.Sprintln(args...))
}
