package lib

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"google.golang.org/grpc/status"
)

var (
	red    = color.New(color.FgRed).FprintlnFunc()
	yellow = color.New(color.FgYellow).FprintlnFunc()
)

func Crash(a ...any) {
	if len(a) == 1 {
		if _, ok := a[0].(error); ok {
			if s, ok := status.FromError(a[0].(error)); ok {
				red(os.Stderr, fmt.Sprintf("[%s] %s", s.Code().String(), s.Message()))
			}
		}
	} else if len(a) == 0 {
		red(os.Stderr, "unknown error!")
	} else {
		red(os.Stderr, a...)
	}
	os.Exit(1)
}

func Warn(a ...any) {
	yellow(os.Stderr, a...)
}
