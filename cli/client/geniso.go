package client

import (
	"github.com/computerdane/nf6/iso"
	"github.com/computerdane/nf6/lib"
	"github.com/spf13/cobra"
)

var (
	genisoAccountSshPubKey string
	genisoHostAddr6        string
	genisoHostSystem       string
	genisoHostWgPrivKey    string
	genisoVipWgEndpoint    string
	genisoVipWgPubKey      string
)

func init() {
	genisoCmd.Flags().StringVar(&genisoAccountSshPubKey, "account-ssh-pub-key", "", "your SSH public key")
	genisoCmd.Flags().StringVar(&genisoHostAddr6, "host-addr6", "", "host's IPv6 address")
	genisoCmd.Flags().StringVar(&genisoHostSystem, "host-system", "", "host's system type (Nix system name)")
	genisoCmd.Flags().StringVar(&genisoHostWgPrivKey, "host-wg-priv-key", "", "host's WireGuard private key")
	genisoCmd.Flags().StringVar(&genisoVipWgEndpoint, "vip-wg-endpoint", "auto", "VIP WireGuard endpoint (if auto, will grab from server)")
	genisoCmd.Flags().StringVar(&genisoVipWgPubKey, "vip-wg-pub-key", "auto", "VIP WireGuard public key (if auto, will grab from server)")
}

var genisoCmd = &cobra.Command{
	Use:   "geniso",
	Short: "Generate a new Nix install ISO pre-configured for Nf6",
	Run: func(cmd *cobra.Command, args []string) {
		if genisoVipWgEndpoint == "auto" || genisoVipWgPubKey == "auto" {
			ConnectPublic(cmd, args)
			ctx, cancel := lib.Context()
			defer cancel()
			reply, err := apiPublic.GetIpv6Info(ctx, nil)
			if err != nil {
				lib.Crash(err)
			}
			if genisoVipWgEndpoint == "auto" {
				genisoVipWgEndpoint = reply.GetVipWgEndpoint()
			}
			if genisoVipWgPubKey == "auto" {
				genisoVipWgPubKey = reply.GetVipWgPubKey()
			}
		}
		isoPath, err := iso.Generate("/tmp/nf6-geniso", &iso.Config{
			AccountSshPubKey: genisoAccountSshPubKey,
			HostAddr6:        genisoHostAddr6,
			HostSystem:       genisoHostSystem,
			HostWgPrivKey:    genisoHostWgPrivKey,
			VipWgEndpoint:    genisoVipWgEndpoint,
			VipWgPubKey:      genisoVipWgPubKey,
		})
		if err != nil {
			lib.Crash(err)
		}
		lib.Output(map[string]string{"isoPath": isoPath})
	},
}
