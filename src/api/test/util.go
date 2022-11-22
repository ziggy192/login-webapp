package test

import (
	"bitbucket.org/ziggy192/ng_lu/src/logger"
	"context"
	"testing"
)

func NewContext(t *testing.T) context.Context {
	return logger.SaveRequestID(context.Background(), t.Name())
}
