package lib

import (
	"os"
	"path"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var IsDevShell = os.Getenv("IS_DEV_SHELL") != ""

type Option struct {
	P         any
	Name      string
	Shorthand string
	Value     any
	Usage     string
}

var (
	configPath = ""
	options    = []*Option{}
)

func SetConfigPath(path string) {
	configPath = path
	viper.SetConfigFile(configPath)
}

func SetHomeConfigPath(dir string) {
	home, err := os.UserHomeDir()
	if err != nil {
		Crash(err)
	}
	SetConfigPath(home + "/.config/" + dir + "/config.yaml")
}

func SetSystemConfigPath(dir string) {
	SetConfigPath("/var/lib/" + dir + "/config/config.yaml")
}

func SaveConfig() {
	configDir := path.Dir(configPath)
	if err := os.MkdirAll(configDir, os.ModePerm); err != nil {
		Warn("failed to make config directory: ", err)
	}
	if _, err := os.OpenFile(configPath, os.O_CREATE|os.O_WRONLY, 0600); err != nil {
		Warn("failed to open config file: ", err)
	}
	if err := viper.WriteConfig(); err != nil {
		Warn("failed to write config: ", err)
	}
}

func InitConfig(save bool) {
	if _, err := os.Stat(configPath); err != nil {
		SaveConfig()
	}
	if err := viper.ReadInConfig(); err == nil {
		LoadOptions()
	}
	if save {
		SaveConfig()
	}
}

func AddOption(cmd *cobra.Command, o *Option) {
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
