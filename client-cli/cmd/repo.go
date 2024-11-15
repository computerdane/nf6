package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"slices"

	"github.com/computerdane/nf6/nf6"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	repoCmd.AddCommand(repoCloneCmd)
	repoCmd.AddCommand(repoCreateCmd)
	repoCmd.AddCommand(repoLsCmd)
	repoCmd.AddCommand(repoRenameCmd)
	repoCmd.AddCommand(repoSetDefaultCmd)

	rootCmd.AddCommand(repoCmd)
}

func repoNameOrDefault(args []string) string {
	if len(args) == 0 {
		return defaultRepo
	} else {
		return args[0]
	}
}

var repoCmd = &cobra.Command{
	Use:   "repo",
	Short: "Manage your repos",
}

var repoCloneCmd = &cobra.Command{
	Use:    "clone [name] [gitCloneArgs]",
	Short:  "Clone a repo",
	PreRun: RequireSecureClient,
	Run: func(cmd *cobra.Command, args []string) {
		sshCommand := fmt.Sprintf(`ssh -i "%s"`, sshPrivKeyPath)
		repoUrl := fmt.Sprintf("git@%s:%s", gitHost, repoNameOrDefault(args))

		gitArgs := []string{"clone", "-c", "core.sshCommand=" + sshCommand, repoUrl}
		if len(args) > 1 {
			gitArgs = append(gitArgs, args[1:]...)
		}

		gitCmd := exec.Command("git", gitArgs...)
		gitCmd.Stdin = os.Stdin
		gitCmd.Stdout = os.Stdout
		gitCmd.Stderr = os.Stderr
		gitCmd.Run()
	},
}

var repoCreateCmd = &cobra.Command{
	Use:    "create [name]",
	Short:  "Create a repo",
	Args:   cobra.MaximumNArgs(1),
	PreRun: RequireSecureClient,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		reply, err := clientSecure.CreateRepo(ctx, &nf6.CreateRepoRequest{Name: repoNameOrDefault(args)})
		if err != nil {
			Crash(err)
		}
		if !reply.GetSuccess() {
			Crash()
		}
	},
}

var repoLsCmd = &cobra.Command{
	Use:    "ls",
	Short:  "List your repos",
	PreRun: RequireSecureClient,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		reply, err := clientSecure.ListRepos(ctx, &nf6.ListReposRequest{})
		if err != nil {
			Crash(err)
		}
		for _, repoName := range reply.Names {
			if repoName == defaultRepo {
				color.New(color.FgGreen, color.Bold).Print(repoName)
				fmt.Print(" (default)")
			} else {
				fmt.Print(repoName)
			}
			fmt.Println()
		}
	},
}

var repoRenameCmd = &cobra.Command{
	Use:    "rename [oldName] [newName]",
	Short:  "Rename a repo",
	Args:   cobra.ExactArgs(2),
	PreRun: RequireSecureClient,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		reply, err := clientSecure.RenameRepo(ctx, &nf6.RenameRepoRequest{OldName: args[0], NewName: args[1]})
		if err != nil {
			Crash(err)
		}
		if !reply.GetSuccess() {
			Crash()
		}
	},
}

var repoSetDefaultCmd = &cobra.Command{
	Use:    "set-default [name]",
	Short:  "Set the default repo",
	Args:   cobra.ExactArgs(1),
	PreRun: RequireSecureClient,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		reply, err := clientSecure.ListRepos(ctx, &nf6.ListReposRequest{})
		if err != nil {
			Crash(err)
		}

		if slices.Contains(reply.GetNames(), args[0]) {
			defaultRepo = args[0]
			viper.Set("defaultRepo", defaultRepo)
			if err := viper.WriteConfig(); err != nil {
				Crash(err)
			}
		}
	},
}
