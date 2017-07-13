package core

import (
	"io/ioutil"
	"net/rpc"
	"os"
	"path"
	"strings"
	"time"
)

func CleanupDotGoed() {
	cleanStinckySocks()
	cleanOldBuffers()
}

// socks that wheree not cleanup by previous (unclean) shutdown
func cleanStinckySocks() {
	dir := path.Join(GoedHome(), "instances")
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}
	for _, fi := range files {
		if strings.HasSuffix(fi.Name(), ".sock") {
			f := path.Join(dir, fi.Name())
			c, err := rpc.DialHTTP("unix", f)
			if err != nil {
				// no longer working
				os.Remove(f)
			} else {
				c.Close()
			}
		}
	}
}

// cleanup buffers over 30 days old
func cleanOldBuffers() {
	tooOld := time.Now().Add(-30 * 24 * time.Hour)
	dir := path.Join(GoedHome(), "buffers")
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}
	for _, fi := range files {
		if fi.ModTime().Before(tooOld) {
			os.Remove(path.Join(dir, fi.Name()))
		}
	}
}
