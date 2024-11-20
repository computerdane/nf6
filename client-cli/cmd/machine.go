package cmd

import (
	"context"
	"fmt"

	"github.com/computerdane/nf6/lib"
	"github.com/computerdane/nf6/nf6"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var (
	addrIpv6 string
	wgPubKey string
)

func init() {
	machineCmd.AddCommand(machineAddCmd)
	machineCmd.AddCommand(machineLsCmd)
	machineCmd.AddCommand(machineRenameCmd)

	rootCmd.AddCommand(machineCmd)

	machineAddCmd.PersistentFlags().StringVar(&addrIpv6, "addr-ipv6", "", "ipv6 address of the machine")
	machineAddCmd.PersistentFlags().StringVar(&wgPubKey, "wg-public-key", "", "wireguard public key for the machine")
}

var machineCmd = &cobra.Command{
	Use:   "machine",
	Short: "Manage your machines",
}

var machineAddCmd = &cobra.Command{
	Use:    "add [name]",
	Short:  "Add a machine",
	Args:   cobra.ExactArgs(1),
	PreRun: RequireSecureClient,
	Run: func(cmd *cobra.Command, args []string) {
		hostName := args[0]
		if err := lib.ValidateHostName(hostName); err != nil {
			Crash(err)
		}
		if err := lib.PromptOrValidate(&addrIpv6, promptui.Prompt{
			Label:    "IPv6 Address",
			Validate: lib.ValidateIpv6Address,
		}); err != nil {
			Crash(err)
		}
		if err := lib.PromptOrValidate(&wgPubKey, promptui.Prompt{
			Label:    "WireGuard Public Key",
			Validate: lib.ValidateWireguardKey,
		}); err != nil {
			Crash(err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		reply, err := clientSecure.AddMachine(ctx, &nf6.AddMachineRequest{HostName: hostName, WgPublicKey: wgPubKey, AddrIpv6: addrIpv6})
		if err != nil {
			Crash(err)
		}
		if !reply.GetSuccess() {
			Crash()
		}
	},
}

var machineLsCmd = &cobra.Command{
	Use:    "ls",
	Short:  "List your machines",
	PreRun: RequireSecureClient,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		reply, err := clientSecure.ListMachines(ctx, &nf6.ListMachinesRequest{})
		if err != nil {
			Crash(err)
		}
		for _, machineName := range reply.Names {
			fmt.Println(machineName)
		}
	},
}

var machineRenameCmd = &cobra.Command{
	Use:    "rename [oldName] [newName]",
	Short:  "Rename a machine",
	Args:   cobra.ExactArgs(2),
	PreRun: RequireSecureClient,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		reply, err := clientSecure.RenameMachine(ctx, &nf6.RenameMachineRequest{OldName: args[0], NewName: args[1]})
		if err != nil {
			Crash(err)
		}
		if !reply.GetSuccess() {
			Crash()
		}
	},
}
