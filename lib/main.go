package lib

import (
	"regexp"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var validRepoRegex = `^[A-Za-z0-9][A-Za-z0-9\-_]+[A-Za-z0-9]$`
var invalidRepoRegex = `^.*(\-\-|__|\-_|_\-).*$`

func ValidateRepoName(name string) (valid bool, err error) {
	if matchValid, _ := regexp.MatchString(validRepoRegex, name); matchValid {
		if matchInvalid, _ := regexp.MatchString(invalidRepoRegex, name); !matchInvalid {
			return true, nil
		}
	}
	return false, status.Error(codes.InvalidArgument, "Repo name must only contain characters A-Z, a-z, 0-9, -, and _. Repo name must not start or end with - or _. Repo name must not have two or more consecutive - and/or _.")
}

type StringOption struct {
	P     *string
	Name  string
	Value string
	Usage string
}

func AddStringOptions(cmd *cobra.Command, options []StringOption) {
	for _, o := range options {
		cmd.PersistentFlags().StringVar(o.P, o.Name, o.Value, o.Usage)
		viper.BindPFlag(o.Name, cmd.PersistentFlags().Lookup(o.Name))
		*o.P = viper.GetString(o.Name)
	}
}
func LoadStringOptions(cmd *cobra.Command, options []StringOption) {
	for _, o := range options {
		*o.P = viper.GetString(o.Name)
	}
}

type IntOption struct {
	P     *int
	Name  string
	Value int
	Usage string
}

func AddIntOptions(cmd *cobra.Command, options []IntOption) {
	for _, o := range options {
		cmd.PersistentFlags().IntVar(o.P, o.Name, o.Value, o.Usage)
		viper.BindPFlag(o.Name, cmd.PersistentFlags().Lookup(o.Name))
		*o.P = viper.GetInt(o.Name)
	}
}
func LoadIntOptions(cmd *cobra.Command, options []IntOption) {
	for _, o := range options {
		*o.P = viper.GetInt(o.Name)
	}
}

type DurationOption struct {
	P     *time.Duration
	Name  string
	Value time.Duration
	Usage string
}

func AddDurationOptions(cmd *cobra.Command, options []DurationOption) {
	for _, o := range options {
		cmd.PersistentFlags().DurationVar(o.P, o.Name, o.Value, o.Usage)
		viper.BindPFlag(o.Name, cmd.PersistentFlags().Lookup(o.Name))
		*o.P = viper.GetDuration(o.Name)
	}
}
func LoadDurationOptions(cmd *cobra.Command, options []DurationOption) {
	for _, o := range options {
		*o.P = viper.GetDuration(o.Name)
	}
}
