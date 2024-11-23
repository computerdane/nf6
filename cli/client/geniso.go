package client

import (
	"github.com/computerdane/nf6/iso"
	"github.com/computerdane/nf6/lib"
	"github.com/spf13/cobra"
)

var (
	genisoHostAddr6           string
	genisoServerGlobalPrefix6 string
	genisoWgServerEndpoint    string
	genisoWgServerWgPubKey    string
	genisoSshPubKey           string
	genisoSystem              string
	genisoWgPrivKey           string
)

func init() {
	genisoCmd.Flags().StringVar(&genisoHostAddr6, "host-addr6", "", "host's IPv6 address")
	genisoCmd.Flags().StringVar(&genisoWgServerEndpoint, "wg-server-endpoint", "auto", "WireGuard server's endpoint")
	genisoCmd.Flags().StringVar(&genisoServerGlobalPrefix6, "server-global-prefix6", "auto", "Server's global IPv6 prefix")
	genisoCmd.Flags().StringVar(&genisoWgServerWgPubKey, "wg-server-wg-pub-key", "auto", "WireGuard server's public key (if auto, will grab from server)")
	genisoCmd.Flags().StringVar(&genisoSshPubKey, "ssh-pub-key", "", "your SSH public key")
	genisoCmd.Flags().StringVar(&genisoSystem, "system", "", "host's system type (Nix system name)")
	genisoCmd.Flags().StringVar(&genisoWgPrivKey, "wg-priv-key", "", "host's WireGuard private key")
}

var genisoCmd = &cobra.Command{
	Use:   "geniso",
	Short: "Generate a new Nix install ISO pre-configured for Nf6",
	Run: func(cmd *cobra.Command, args []string) {
		if genisoWgServerEndpoint == "auto" || genisoWgServerWgPubKey == "auto" {
			ConnectPublic(cmd, args)
			ctx, cancel := lib.Context()
			defer cancel()
			reply, err := apiPublic.GetIpv6Info(ctx, nil)
			if err != nil {
				lib.Crash(err)
			}
			if genisoWgServerEndpoint == "auto" {
				genisoWgServerEndpoint = reply.GetWgServerEndpoint()
			}
			if genisoWgServerWgPubKey == "auto" {
				genisoWgServerWgPubKey = reply.GetWgServerWgPubKey()
			}
		}
		isoPath, err := iso.Generate("/tmp/nf6-geniso", &iso.Config{
			HostAddr:            genisoHostAddr6,
			ServerGlobalPrefix6: genisoServerGlobalPrefix6,
			ServerWgEndpoint:    genisoWgServerEndpoint,
			ServerWgPubKey:      genisoWgServerWgPubKey,
			SshPubKey:           genisoSshPubKey,
			System:              genisoSystem,
			WgPrivKey:           genisoWgPrivKey,
		})
		if err != nil {
			lib.Crash(err)
		}
		lib.Output(map[string]string{"isoPath": isoPath})
	},
}
