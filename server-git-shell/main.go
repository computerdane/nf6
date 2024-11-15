package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func parseGitCommand(args []string) error {
	if args[0] == "git-receive-pack" || args[0] == "git-upload-pack" || args[0] == "git-upload-archive" {
		for i, arg := range args {
			// ignore command name
			if i == 0 {
				continue
			}
			// ignore optional arguments
			if strings.HasPrefix(arg, "--") {
				continue
			}
			// get repo name without outer quotes
			name := arg[1 : len(arg)-1]
			// substitute path corresponding to repo name
			for _, path := range os.Args[1:] {
				if strings.HasSuffix(path, name) {
					args[i] = "'" + path + "'"
					return nil
				}
			}
			return errors.New("repo not found")
		}
	}
	return errors.New("invalid git command")
}

func main() {
	cmdStr := os.Getenv("SSH_ORIGINAL_COMMAND")
	args := strings.Split(cmdStr, " ")

	if err := parseGitCommand(args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cmd := exec.Command("git-shell", "-c", strings.Join(args, " "))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
