package client

import (
	"github.com/computerdane/nf6/lib"
	"github.com/computerdane/nf6/nf6"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var (
	repoName string
)

func init() {
	repoCmd.PersistentFlags().StringVarP(&repoName, "name", "n", "", "repo name")

	repoCmd.AddCommand(repoCreateCmd)
	repoCmd.AddCommand(repoGetCmd)
	repoCmd.AddCommand(repoListCmd)
	repoCmd.AddCommand(repoEditCmd)
}

var repoCmd = &cobra.Command{
	Use:   "repo",
	Short: "Manage your repos",
}

var repoCreateCmd = &cobra.Command{
	Use:    "create",
	Short:  "Create a new repo",
	Args:   cobra.MaximumNArgs(1),
	PreRun: Connect,
	Run: func(cmd *cobra.Command, args []string) {
		newName := ""
		if len(args) > 0 {
			newName = args[0]
		}
		if err := lib.PromptOrValidate(&newName, &promptui.Prompt{
			Label:    "Name",
			Validate: lib.ValidateRepoName,
		}); err != nil {
			lib.Crash(err)
		}
		ctx, cancel := lib.Context()
		defer cancel()
		if _, err := api.CreateRepo(ctx, &nf6.CreateRepo_Request{Name: newName}); err != nil {
			lib.Crash(err)
		}
	},
}

var repoGetCmd = &cobra.Command{
	Use:    "get [name]",
	Short:  "Get info about a repo",
	Args:   cobra.ExactArgs(1),
	PreRun: Connect,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := lib.Context()
		defer cancel()
		reply, err := api.GetRepo(ctx, &nf6.GetRepo_Request{Name: args[0]})
		if err != nil {
			lib.Crash(err)
		}
		lib.Output(reply)
	},
}

var repoListCmd = &cobra.Command{
	Use:    "list",
	Short:  "List your repos",
	PreRun: Connect,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := lib.Context()
		defer cancel()
		reply, err := api.ListRepos(ctx, nil)
		if err != nil {
			lib.Crash(err)
		}
		lib.OutputStringList(reply.GetNames())
	},
}

var repoEditCmd = &cobra.Command{
	Use:    "edit [name]",
	Short:  "Edit a repo",
	Args:   cobra.ExactArgs(1),
	PreRun: Connect,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := lib.Context()
		defer cancel()
		reply, err := api.GetRepo(ctx, &nf6.GetRepo_Request{Name: args[0]})
		if err != nil {
			lib.Crash(err)
		}
		if repoName == "" {
			if err := lib.PromptOrValidate(&repoName, &promptui.Prompt{
				Label:    "Name",
				Default:  reply.GetName(),
				Validate: lib.ValidateRepoName,
			}); err != nil {
				lib.Crash(err)
			}
		}
		req := nf6.UpdateRepo_Request{Id: reply.GetId()}
		if repoName != "" {
			req.Name = &repoName
		}
		{
			ctx, cancel := lib.Context()
			defer cancel()
			if _, err := api.UpdateRepo(ctx, &req); err != nil {
				lib.Crash(err)
			}
		}
	},
}
