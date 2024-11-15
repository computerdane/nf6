package cmd

import (
	"context"
	"os"

	"github.com/computerdane/nf6/nf6"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(whoamiCmd)
}

var whoamiCmd = &cobra.Command{
	Use:    "whoami",
	Short:  "Get info from the server about your user",
	PreRun: RequireSecureClient,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		reply, err := clientSecure.WhoAmI(ctx, &nf6.WhoAmIRequest{})
		if err != nil {
			Crash(err)
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendRow(table.Row{"Email", reply.GetEmail()})
		t.AppendSeparator()
		t.AppendRow(table.Row{"SSL Pubkey", reply.GetSslPublicKey()})
		t.AppendSeparator()
		t.AppendRow(table.Row{"SSH Pubkey", reply.GetSshPublicKey()})
		t.AppendSeparator()
		t.SetStyle(table.StyleRounded)
		t.Render()
	},
}
