package core

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

// UpdateResource creates or updates the bundled resources into GOED_HOME
func UpdateResources() {
	version, err := Asset("res/resources_version.txt")
	if err != nil {
		panic(err)
	}
	f := path.Join(Home, "resources_version.txt")
	curVersion, err := ioutil.ReadFile(f)
	if err != nil || string(curVersion) != string(version) {
		// new version, copy resources
		for _, nm := range AssetNames() {
			parts := strings.Split(nm, string(os.PathSeparator))
			target := path.Join(Home, path.Join(parts[1:]...))
			asset, _ := Asset(nm)
			os.MkdirAll(path.Dir(target), 0750)
			err := ioutil.WriteFile(target, asset, 0750)
			if err != nil {
				panic(err)
			}
			log.Printf("Copying %s to %s\n", nm, target)
		}
		loc := path.Join(Home, "config.toml")
		// If no custom config file yet, create one
		if _, err := os.Stat(loc); os.IsNotExist(err) {
			err := CopyFile(path.Join(Home, "standard", "config.toml"), loc)
			if err != nil {
				panic(err)
			}
		}
	}
}

// GetResource finds a GOED resource either from
// - GOED_HOME/<path>
// or - GOED/HOME/standard/<path>
func FindResource(relPath string) (absPath string) {
	absPath = path.Join(Home, relPath)
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return path.Join(Home, "standard", relPath)
	}
	return absPath
}
