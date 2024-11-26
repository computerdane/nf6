package client

import (
	"github.com/computerdane/nf6/iso"
	"github.com/computerdane/nf6/lib"
	"github.com/computerdane/nf6/nf6"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

var (
	hostName            string
	hostAddr6           string
	hostWgPubKey        string
	hostListAll         bool
	hostCreateIso       bool
	hostCreateIsoSystem string
)

func init() {
	hostCreateCmd.Flags().StringVarP(&hostAddr6, "addr6", "a", "", "IPv6 address")
	hostCreateCmd.Flags().StringVarP(&hostWgPubKey, "wg-pub-key", "w", "", "WireGuard public key")
	hostCreateCmd.Flags().BoolVarP(&hostCreateIso, "iso", "i", false, "create a host and generate an install ISO")
	hostCreateCmd.Flags().StringVarP(&hostCreateIsoSystem, "iso-system", "s", "", "host system type for ISO")

	hostListCmd.Flags().BoolVarP(&hostListAll, "all", "a", false, "show all host info")

	hostEditCmd.Flags().StringVarP(&hostName, "name", "n", "", "host name")
	hostEditCmd.Flags().StringVarP(&hostAddr6, "addr6", "a", "", "IPv6 address")
	hostEditCmd.Flags().StringVarP(&hostWgPubKey, "wg-pub-key", "w", "", "WireGuard public key")

	hostCmd.AddCommand(hostCreateCmd)
	hostCmd.AddCommand(hostDeleteCmd)
	hostCmd.AddCommand(hostGetCmd)
	hostCmd.AddCommand(hostListCmd)
	hostCmd.AddCommand(hostEditCmd)
}

var hostCmd = &cobra.Command{
	Use:   "host",
	Short: "Manage your hosts",
}

var hostCreateCmd = &cobra.Command{
	Use:    "create [name]",
	Short:  "Create a new host",
	Args:   cobra.MaximumNArgs(1),
	PreRun: ConnectBoth,
	Run: func(cmd *cobra.Command, args []string) {
		newName := ""
		if len(args) > 0 {
			newName = args[0]
		}
		if err := lib.PromptOrValidate(&newName, &promptui.Prompt{
			Label:    "Name",
			Validate: lib.ValidateHostName,
		}); err != nil {
			lib.Crash(err)
		}
		if hostCreateIso {
			ctx, cancel := lib.Context()
			defer cancel()
			global, err := apiPublic.GetGlobal(ctx, nil)
			if err != nil {
				lib.Crash(err)
			}

			if hostCreateIsoSystem == "" {
				prompt := promptui.Select{
					Label: "System",
					Items: lib.ValidNixSystems,
				}
				_, result, err := prompt.Run()
				if err != nil {
					lib.Crash(err)
				}
				hostCreateIsoSystem = result
			}
			if err := lib.ValidateNixSystem(hostCreateIsoSystem); err != nil {
				lib.Crash(err)
			}

			wgPrivKey, err := wgtypes.GeneratePrivateKey()
			if err != nil {
				lib.Crash(err)
			}

			ctx, cancel = lib.Context()
			defer cancel()
			if _, err := api.CreateHost(ctx, &nf6.CreateHost_Request{Name: newName, WgPubKey: wgPrivKey.PublicKey().String()}); err != nil {
				lib.Crash(err)
			}

			ctx, cancel = lib.Context()
			defer cancel()
			hostInfo, err := api.GetHost(ctx, &nf6.GetHost_Request{Name: newName})
			if err != nil {
				lib.Crash(err)
			}

			ctx, cancel = lib.Context()
			defer cancel()
			accountInfo, err := api.GetAccount(ctx, nil)
			if err != nil {
				lib.Crash(err)
			}

			isoPath, err := iso.Generate("/tmp/nf6-iso-"+hostName, &iso.Config{
				AccountSshPubKey: accountInfo.GetSshPubKey(),
				HostAddr6:        hostInfo.GetAddr6(),
				HostSystem:       hostCreateIsoSystem,
				HostWgPrivKey:    wgPrivKey.String(),
				VipWgEndpoint:    global.GetVipWgEndpoint(),
				VipWgPubKey:      global.GetVipWgPubKey(),
			})
			if err != nil {
				lib.Crash(err)
			}
			lib.Output(map[string]string{"isoPath": isoPath})
		} else {
			if err := lib.PromptOrValidate(&hostWgPubKey, &promptui.Prompt{
				Label:    "WireGuard public key",
				Validate: lib.ValidateWireguardKey,
			}); err != nil {
				lib.Crash(err)
			}
			ctx, cancel := lib.Context()
			defer cancel()
			if _, err := api.CreateHost(ctx, &nf6.CreateHost_Request{Name: newName, Addr6: &hostAddr6, WgPubKey: hostWgPubKey}); err != nil {
				lib.Crash(err)
			}
		}
	},
}

var hostDeleteCmd = &cobra.Command{
	Use:    "delete [name]",
	Short:  "Delete a host",
	Args:   cobra.ExactArgs(1),
	PreRun: Connect,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := lib.Context()
		defer cancel()
		host, err := api.GetHost(ctx, &nf6.GetHost_Request{Name: args[0]})
		if err != nil {
			lib.Crash(err)
		}
		{
			ctx, cancel := lib.Context()
			defer cancel()
			if _, err := api.DeleteHost(ctx, &nf6.DeleteHost_Request{Id: host.GetId()}); err != nil {
				lib.Crash(err)
			}
		}
	},
}

var hostGetCmd = &cobra.Command{
	Use:    "get [name]",
	Short:  "Get info about a host",
	Args:   cobra.ExactArgs(1),
	PreRun: Connect,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := lib.Context()
		defer cancel()
		reply, err := api.GetHost(ctx, &nf6.GetHost_Request{Name: args[0]})
		if err != nil {
			lib.Crash(err)
		}
		lib.Output(reply)
	},
}

var hostListCmd = &cobra.Command{
	Use:    "list",
	Short:  "List your hosts",
	PreRun: Connect,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := lib.Context()
		defer cancel()
		reply, err := api.ListHosts(ctx, nil)
		if err != nil {
			lib.Crash(err)
		}
		if hostListAll {
			output := map[string]interface{}{}
			for _, name := range reply.GetNames() {
				ctx, cancel := lib.Context()
				defer cancel()
				reply, err := api.GetHost(ctx, &nf6.GetHost_Request{Name: name})
				if err != nil {
					lib.Crash(err)
				}
				output[name] = reply
			}
			lib.OutputAll(output)
		} else {
			lib.OutputStringList(reply.GetNames())
		}
	},
}

var hostEditCmd = &cobra.Command{
	Use:    "edit [name]",
	Short:  "Edit a host",
	Args:   cobra.ExactArgs(1),
	PreRun: Connect,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := lib.Context()
		defer cancel()
		reply, err := api.GetHost(ctx, &nf6.GetHost_Request{Name: args[0]})
		if err != nil {
			lib.Crash(err)
		}
		if hostName == "" && hostAddr6 == "" && hostWgPubKey == "" {
			if err := lib.PromptOrValidate(&hostName, &promptui.Prompt{
				Label:    "Name",
				Default:  reply.GetName(),
				Validate: lib.ValidateHostName,
			}); err != nil {
				lib.Crash(err)
			}
			if err := lib.PromptOrValidate(&hostAddr6, &promptui.Prompt{
				Label:    "IPv6 address",
				Default:  reply.GetAddr6(),
				Validate: lib.ValidateIpv6Address,
			}); err != nil {
				lib.Crash(err)
			}
			if err := lib.PromptOrValidate(&hostWgPubKey, &promptui.Prompt{
				Label:    "WireGuard public key",
				Default:  reply.GetWgPubKey(),
				Validate: lib.ValidateWireguardKey,
			}); err != nil {
				lib.Crash(err)
			}
		}
		req := nf6.UpdateHost_Request{Id: reply.GetId()}
		if hostName != "" {
			req.Name = &hostName
		}
		if hostAddr6 != "" {
			req.Addr6 = &hostAddr6
		}
		if hostWgPubKey != "" {
			req.WgPubKey = &hostWgPubKey
		}
		{
			ctx, cancel := lib.Context()
			defer cancel()
			if _, err := api.UpdateHost(ctx, &req); err != nil {
				lib.Crash(err)
			}
		}
	},
}
