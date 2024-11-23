package client

import (
	"github.com/computerdane/nf6/iso"
	"github.com/computerdane/nf6/lib"
	"github.com/spf13/cobra"
)

var (
	genisoHostAddr6      string
	genisoServerAddr6    string
	genisoServerWgPubKey string
	genisoSshPubKey      string
	genisoSystem         string
	genisoWgPrivKey      string
)

func init() {
	genisoCmd.Flags().StringVar(&genisoHostAddr6, "host-addr6", "", "host's IPv6 address")
	genisoCmd.Flags().StringVar(&genisoServerAddr6, "server-addr6", "auto", "server's IPv6 address (if auto, will grab from server)")
	genisoCmd.Flags().StringVar(&genisoServerWgPubKey, "server-wg-pub-key", "auto", "server's WireGuard public key (if auto, will grab from server)")
	genisoCmd.Flags().StringVar(&genisoSshPubKey, "ssh-pub-key", "", "host's SSH public key")
	genisoCmd.Flags().StringVar(&genisoSystem, "system", "", "host's system type (Nix system name)")
	genisoCmd.Flags().StringVar(&genisoWgPrivKey, "wg-priv-key", "", "host's WireGuard private key")
}

var genisoCmd = &cobra.Command{
	Use:   "geniso",
	Short: "Generate a new Nix install ISO pre-configured for Nf6",
	Run: func(cmd *cobra.Command, args []string) {
		if genisoServerAddr6 == "auto" || genisoServerWgPubKey == "auto" {
			ConnectPublic(cmd, args)
			ctx, cancel := lib.Context()
			defer cancel()
			reply, err := apiPublic.GetIpv6Info(ctx, nil)
			if err != nil {
				lib.Crash(err)
			}
			if genisoServerAddr6 == "auto" {
				genisoServerAddr6 = reply.GetWgServerAddr6()
			}
			if genisoServerWgPubKey == "auto" {
				genisoServerWgPubKey = reply.GetWgServerWgPubKey()
			}
		}
		if err := iso.Generate("/tmp/nf6-geniso", &iso.Config{
			HostAddr:       genisoHostAddr6,
			ServerAddr:     genisoServerAddr6,
			ServerWgPubKey: genisoServerWgPubKey,
			SshPubKey:      genisoSshPubKey,
			System:         genisoSystem,
			WgPrivKey:      genisoWgPrivKey,
		}); err != nil {
			lib.Crash(err)
		}
	},
}
