package logger

import (
	"context"
	"fmt"
	"log"
)

func Info(ctx context.Context, v ...any) {
	args := make([]any, 0, 1+len(v))
	args = append(append(args, "[info]"), v...)
	_ = log.Output(2, fmt.Sprintln(args...))
}

func Err(ctx context.Context, v ...any) {
	args := make([]any, 0, 1+len(v))
	args = append(append(args, "[error] "), v...)
	_ = log.Output(2, fmt.Sprintln(args...))
}
