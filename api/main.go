package main

import (
	"github.com/computerdane/nf6/api/server"
	"github.com/computerdane/nf6/lib"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "nf6-api",
	Short: "API server for Nf6",
}

func init() {
	server.Init(rootCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		lib.Crash(err)
	}
}
