package core

import "github.com/BurntSushi/toml"

// Config represents the Goed configuration data.
type Config struct {
	SyntaxHighlighting bool
	Theme              string // ie: theme1.toml
	MaxCmdBufferLines  int    // Max # of lines to keep in buffer when running a command
	GuiFont            string // full path to a monospace TTF font
	GuiFontSize        int
	GuiFontDpi         int
	MinViewWidth       int // preferred minimum view width (in characters)
	LineWidthIndicator int // line width indicator (ie: 80 cols)
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
	if conf.GuiFontSize == 0 {
		conf.GuiFontSize = 10
	}
	if conf.GuiFontDpi == 0 {
		conf.GuiFontDpi = 96
	}
	if conf.MinViewWidth == 0 {
		conf.MinViewWidth = 80
	}
	if conf.LineWidthIndicator == 0 {
		conf.LineWidthIndicator = 80
	}
	return conf
}
