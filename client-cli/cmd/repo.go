package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/computerdane/nf6/nf6"
	"github.com/spf13/cobra"
)

func init() {
	repoCmd.AddCommand(repoCloneCmd)
	repoCmd.AddCommand(repoCreateCmd)
	repoCmd.AddCommand(repoLsCmd)
	repoCmd.AddCommand(repoRenameCmd)

	rootCmd.AddCommand(repoCmd)
}

var repoCmd = &cobra.Command{
	Use:   "repo",
	Short: "Manage your repos",
}

// git clone -c core.sshCommand="ssh -i ~/.local/share/nf6/ssh/id_ed25519" git@nf6.sh:main2

var repoCloneCmd = &cobra.Command{
	Use:    "clone [name] [gitCloneArgs]",
	Short:  "Clone a repo",
	Args:   cobra.MinimumNArgs(1),
	PreRun: requireSecureClient,
	Run: func(cmd *cobra.Command, args []string) {
		sshCommand := fmt.Sprintf(`ssh -i "%s"`, sshPrivKeyPath)
		repoUrl := fmt.Sprintf("git@%s:%s", gitHost, args[0])

		gitArgs := []string{"clone", "-c", "core.sshCommand=" + sshCommand, repoUrl}
		gitArgs = append(gitArgs, args[1:]...)

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
	Args:   cobra.ExactArgs(1),
	PreRun: requireSecureClient,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		reply, err := clientSecure.CreateRepo(ctx, &nf6.CreateRepoRequest{Name: args[0]})
		if err != nil {
			crash(err)
		}
		if !reply.GetSuccess() {
			crash()
		}
	},
}

var repoLsCmd = &cobra.Command{
	Use:    "ls",
	Short:  "List your repos",
	PreRun: requireSecureClient,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		reply, err := clientSecure.ListRepos(ctx, &nf6.ListReposRequest{})
		if err != nil {
			crash(err)
		}
		for _, repoName := range reply.Names {
			fmt.Println(repoName)
		}
	},
}

var repoRenameCmd = &cobra.Command{
	Use:    "rename [oldName] [newName]",
	Short:  "Rename a repo",
	Args:   cobra.ExactArgs(2),
	PreRun: requireSecureClient,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		reply, err := clientSecure.RenameRepo(ctx, &nf6.RenameRepoRequest{OldName: args[0], NewName: args[1]})
		if err != nil {
			crash(err)
		}
		if !reply.GetSuccess() {
			crash()
		}
	},
}
