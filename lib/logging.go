package lib

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"google.golang.org/grpc/status"
)

var OutputType = ""

func SetOutputType(t string) {
	OutputType = t
}

var (
	red    = color.New(color.FgRed).FprintlnFunc()
	yellow = color.New(color.FgYellow).FprintlnFunc()
)

func toJson(a any) []byte {
	output, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		Crash(err)
	}
	return output
}

func Output(a any) {
	switch OutputType {
	case "json":
		fmt.Println(string(toJson(a)))
	case "table":
		j := toJson(a)
		var m map[string]interface{}
		if err := json.Unmarshal(j, &m); err != nil {
			Crash(err)
		}
		tbl := table.NewWriter()
		tbl.SetOutputMirror(os.Stdout)
		for k, v := range m {
			tbl.AppendRow(table.Row{k, v})
			tbl.AppendSeparator()
		}
		tbl.SetStyle(table.StyleRounded)
		tbl.Render()
	default:
		fmt.Println(a)
	}
}

func Crash(a ...any) {
	if len(a) == 1 {
		if _, ok := a[0].(error); ok {
			if s, ok := status.FromError(a[0].(error)); ok {
				red(os.Stderr, fmt.Sprintf("[%s] %s", s.Code().String(), s.Message()))
				os.Exit(1)
			}
		}
	} else if len(a) == 0 {
		red(os.Stderr, "unknown error!")
		os.Exit(1)
	}
	red(os.Stderr, a...)
	os.Exit(1)
}

func Warn(a ...any) {
	yellow(os.Stderr, a...)
}
