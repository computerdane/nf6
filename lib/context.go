package lib

import (
	"context"
	"time"
)

var timeout = 5 * time.Second

func SetTimeout(t time.Duration) {
	timeout = t
}

func Context() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}
