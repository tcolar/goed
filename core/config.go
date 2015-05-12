package core

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/BurntSushi/toml"
)

type Config struct {
	SyntaxHighlighting bool
	Theme              string // ie: theme1.toml
}

func LoadConfig(file string) *Config {
	conf := LoadDefaultConfig()
	loc := path.Join(Home, file)
	if _, err := toml.DecodeFile(loc, &conf); err != nil {
		panic(err)
	}
	return conf
}

func LoadDefaultConfig() *Config {
	loc := path.Join(Home, "config.toml")
	// If the config does not exist yet(first start ?), create it
	if _, err := os.Stat(loc); os.IsNotExist(err) {
		os.MkdirAll(Home, 0755)
		err := ioutil.WriteFile(loc, []byte(defaultConfig), 0755)
		if err != nil {
			panic(err)
		}
	}
	var conf *Config
	if _, err := toml.DecodeFile(loc, &conf); err != nil {
		panic(err)
	}
	return conf
}

var defaultConfig = `
SyntaxHighlighting=true
Theme="default.toml"
`
