package core

import "github.com/BurntSushi/toml"

// Config represents the Goed configuration data.
type Config struct {
	SyntaxHighlighting bool
	Theme              string // ie: theme1.toml
	MaxCmdBufferLines  int    // Max # of lines to keep in buffer when running a command
}

func LoadConfig(file string) *Config {
	conf := &Config{}
	loc := FindResource(file)
	if _, err := toml.DecodeFile(loc, conf); err != nil {
		panic(err)
	}
	if conf.MaxCmdBufferLines == 0 {
		conf.MaxCmdBufferLines = 10000
	}
	return conf
}
