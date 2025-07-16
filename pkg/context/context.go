package context

import (
	"context"
	"time"
)

func WithTimeout(timeout int) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
}
