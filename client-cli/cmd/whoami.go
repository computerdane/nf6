package cmd

import (
	"context"
	"log"

	"github.com/computerdane/nf6/nf6"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(whoamiCmd)
}

var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Get info from the server about your user",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		reply, err := client.WhoAmI(ctx, &nf6.WhoAmIRequest{})
		if err != nil {
			log.Fatal(err)
		}
		log.Print(reply.GetSslPublicKey())
	},
}
