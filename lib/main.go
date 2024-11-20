package lib

import (
	"net"
	"regexp"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ValidateHostName(name string) error {
	if match, _ := regexp.MatchString(`^[A-Za-z0-9]([A-Za-z0-9_-]{0,61}[A-Za-z0-9])?$`, name); match {
		return nil
	}
	return status.Error(codes.InvalidArgument, "Invalid host name. Must be alphanumeric with optional underscores in the middle.")
}

func ValidateIpv6Address(addr string) error {
	ip := net.ParseIP(addr)
	if ip == nil || ip.To4() != nil {
		return status.Error(codes.InvalidArgument, "Invalid IPv6 address")
	}
	return nil
}

func ValidateRepoName(name string) error {
	if matchValid, _ := regexp.MatchString(`^[A-Za-z0-9][A-Za-z0-9\-_]+[A-Za-z0-9]$`, name); matchValid {
		if matchInvalid, _ := regexp.MatchString(`^.*(\-\-|__|\-_|_\-).*$`, name); !matchInvalid {
			return nil
		}
	}
	return status.Error(codes.InvalidArgument, "Repo name must only contain characters A-Z, a-z, 0-9, -, and _. Repo name must not start or end with - or _. Repo name must not have two or more consecutive - and/or _.")
}

func ValidateWireguardKey(key string) error {
	_, err := wgtypes.ParseKey(key)
	if err != nil {
		return status.Error(codes.InvalidArgument, "Invalid WireGuard key")
	}
	return nil
}

func PromptOrValidate(value *string, prompt promptui.Prompt) error {
	if *value == "" {
		result, err := prompt.Run()
		if err != nil {
			return err
		}
		*value = result
	} else {
		if err := prompt.Validate(*value); err != nil {
			return err
		}
	}
	return nil
}

type Option struct {
	P         any
	Name      string
	Shorthand string
	Value     any
	Usage     string
}

var options = []Option{}

func AddOption(cmd *cobra.Command, o Option) {
	switch o.Value.(type) {
	case string:
		if o.Shorthand == "" {
			cmd.PersistentFlags().StringVar(o.P.(*string), o.Name, o.Value.(string), o.Usage)
		} else {
			cmd.PersistentFlags().StringVarP(o.P.(*string), o.Name, o.Shorthand, o.Value.(string), o.Usage)
		}
	case int:
		if o.Shorthand == "" {
			cmd.PersistentFlags().IntVar(o.P.(*int), o.Name, o.Value.(int), o.Usage)
		} else {
			cmd.PersistentFlags().IntVarP(o.P.(*int), o.Name, o.Shorthand, o.Value.(int), o.Usage)
		}
	case bool:
		if o.Shorthand == "" {
			cmd.PersistentFlags().BoolVar(o.P.(*bool), o.Name, o.Value.(bool), o.Usage)
		} else {
			cmd.PersistentFlags().BoolVarP(o.P.(*bool), o.Name, o.Shorthand, o.Value.(bool), o.Usage)
		}
	case time.Duration:
		if o.Shorthand == "" {
			cmd.PersistentFlags().DurationVar(o.P.(*time.Duration), o.Name, o.Value.(time.Duration), o.Usage)
		} else {
			cmd.PersistentFlags().DurationVarP(o.P.(*time.Duration), o.Name, o.Shorthand, o.Value.(time.Duration), o.Usage)
		}
	default:
		return
	}
	viper.BindPFlag(o.Name, cmd.PersistentFlags().Lookup(o.Name))
	options = append(options, o)
}

func LoadOptions() {
	for _, o := range options {
		switch o.Value.(type) {
		case string:
			*o.P.(*string) = viper.GetString(o.Name)
		case int:
			*o.P.(*int) = viper.GetInt(o.Name)
		case bool:
			*o.P.(*bool) = viper.GetBool(o.Name)
		case time.Duration:
			*o.P.(*time.Duration) = viper.GetDuration(o.Name)
		}
	}
}
