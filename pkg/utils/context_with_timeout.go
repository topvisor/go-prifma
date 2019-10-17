package utils

import (
	"context"
	"time"
)

func ContextWithTimeout(timeout time.Duration) (context.Context, context.CancelFunc) {
	ctx := context.Background()
	if timeout > 0 {
		return context.WithTimeout(ctx, timeout)
	} else {
		return context.WithCancel(ctx)
	}
}
