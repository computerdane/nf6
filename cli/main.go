package main

import (
	"github.com/computerdane/nf6/cli/client"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "nf",
	Short: "Nf6 simplifies computer networking and deployment",
}

func init() {
	rootCmd.AddCommand(client.RegisterCmd)
}

func main() {
	rootCmd.Execute()
}
