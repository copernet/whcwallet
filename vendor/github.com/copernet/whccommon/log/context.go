package log

import (
	"context"

	"github.com/satori/go.uuid"
)

var rootContext = context.Background()

func NewContext() context.Context {
	uid, _ := uuid.NewV4()
	return context.WithValue(rootContext, DefaultTraceLabel, uid.String())
}
