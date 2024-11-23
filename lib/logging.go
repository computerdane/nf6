package lib

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"google.golang.org/grpc/status"
)

var OutputType = ""
var TableStyle = table.StyleDefault
var ShowIds = false

func init() {
	TableStyle.Options.DrawBorder = false
	TableStyle.Options.SeparateColumns = false
}

func SetOutputType(t string) {
	OutputType = t
}

var (
	red    = color.New(color.FgRed).FprintlnFunc()
	yellow = color.New(color.FgYellow).FprintlnFunc()
)

func toJson(a any) []byte {
	output, err := json.Marshal(a)
	if err != nil {
		Crash(err)
	}
	return output
}

func Header(s string) {
	yellow(os.Stdout, s)
}

func sortedJsonKeys(a any) (keys []string, m map[string]interface{}) {
	j := toJson(a)
	if err := json.Unmarshal(j, &m); err != nil {
		Crash(err)
	}
	keys = make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	for i, s := range keys {
		if s == "id" {
			keys[0], keys[i] = keys[i], keys[0]
			break
		}
	}
	return keys, m
}

func OutputAll(a any) {
	switch OutputType {
	case "table":
		keys, m := sortedJsonKeys(a)
		for _, k := range keys {
			yellow(os.Stdout, k)
			Output(m[k])
			fmt.Println()
		}
	default:
		Output(a)
	}
}

func Output(a any) {
	switch OutputType {
	case "json":
		fmt.Println(string(toJson(a)))
	case "table":
		keys, m := sortedJsonKeys(a)
		tbl := table.NewWriter()
		tbl.SetOutputMirror(os.Stdout)
		for _, k := range keys {
			if k == "id" && !ShowIds {
				continue
			}
			tbl.AppendRow(table.Row{k, m[k]})
		}
		tbl.SetStyle(TableStyle)
		tbl.Render()
	default:
		fmt.Println(a)
	}
}

func OutputStringList(a []string) {
	switch OutputType {
	case "json":
		fmt.Println(string(toJson(a)))
	case "table":
		sort.Strings(a)
		for _, s := range a {
			fmt.Println(s)
		}
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
