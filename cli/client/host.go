package client

import (
	"github.com/computerdane/nf6/lib"
	"github.com/computerdane/nf6/nf6"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var (
	hostName     string
	hostAddr6    string
	hostWgPubKey string
)

func init() {
	hostCreateCmd.Flags().StringVarP(&hostAddr6, "addr6", "a", "", "IPv6 address")
	hostCreateCmd.Flags().StringVarP(&hostWgPubKey, "wg-pub-key", "w", "", "WireGuard public key")

	hostEditCmd.Flags().StringVarP(&hostName, "name", "n", "", "host name")
	hostEditCmd.Flags().StringVarP(&hostAddr6, "addr6", "a", "", "IPv6 address")
	hostEditCmd.Flags().StringVarP(&hostWgPubKey, "wg-pub-key", "w", "", "WireGuard public key")

	hostCmd.AddCommand(hostCreateCmd)
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
	PreRun: Connect,
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
		// if err := lib.PromptOrValidate(&hostAddr6, &promptui.Prompt{
		// 	Label:    "IPv6 address",
		// 	Validate: lib.ValidateIpv6Address,
		// }); err != nil {
		// 	lib.Crash(err)
		// }
		if err := lib.PromptOrValidate(&hostWgPubKey, &promptui.Prompt{
			Label:    "WireGuard public key",
			Validate: lib.ValidateWireguardKey,
		}); err != nil {
			lib.Crash(err)
		}
		ctx, cancel := lib.Context()
		defer cancel()
		if _, err := api.CreateHost(ctx, &nf6.CreateHost_Request{Name: newName, Addr6: hostAddr6, WgPubKey: hostWgPubKey}); err != nil {
			lib.Crash(err)
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
		lib.OutputStringList(reply.GetNames())
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
