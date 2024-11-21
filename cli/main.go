package main

import (
	"github.com/computerdane/nf6/cli/client"
	"github.com/computerdane/nf6/lib"
	"github.com/spf13/cobra"
)

var (
	configPath string
	saveConfig bool
)

var rootCmd = &cobra.Command{
	Use:   "nf",
	Short: "Nf6 simplifies computer networking and deployment",
}

func init() {
	client.Init(rootCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		lib.Crash(err)
	}
}
