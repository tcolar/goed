package core

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"unicode/utf8"
)

var LineSep = []byte{'\n'}

func CopyFile(from, to string) error {
	_, err := copyFile(from, to, false, nil)
	return err
}

func CopyFileSkipBom(from, to string) (bom []byte, err error) {
	return copyFile(from, to, true, nil)
}

func CopyFileWithBom(from, to string, bom []byte) error {
	_, err := copyFile(from, to, false, bom)
	return err
}

func copyFile(from, to string, skipBom bool, addBom []byte) (bom []byte, err error) {
	if skipBom {
		bomLen, _, _ := CheckBom(from)
		bom = make([]byte, bomLen)
	}

	in, err := os.Open(from)
	if err != nil {
		return nil, err
	}
	defer in.Close()
	out, err := os.Create(to)
	if err != nil {
		return nil, err
	}
	defer out.Close()

	if skipBom && len(bom) > 0 {
		// read the bom to skip it (seek doesn't seem to have any effect on copy)
		in.Read(bom)
	}

	if len(addBom) > 0 {
		// write bom
		_, err := out.Write(addBom)
		if err != nil {
			return nil, err
		}
	}

	// write content
	_, err = io.Copy(out, in)
	return bom, err
}

// Mv file moves a file by copy, then delete
// because os.Rename does not always work
func MvFile(from, to string) error {
	if err := CopyFile(from, to); err != nil {
		return err
	}
	return os.Remove(from)
}

// InitHome initializes the ~/.goed directory structure
func InitHome(id int64) {
	InstanceId = id
	Home = GoedHome()
	os.MkdirAll(Home, 0750)
	os.MkdirAll(path.Join(Home, "buffers"), 0750)
	os.MkdirAll(path.Join(Home, "logs"), 0750)
	os.MkdirAll(path.Join(Home, "instances"), 0750)
	ioutil.WriteFile(path.Join(Home, "Version.txt"), []byte(Version), 644)

	// RCP instance socket
	Socket = GoedSocket(id)

	// Terminal app
	Terminal = os.Getenv("SHELL")
	if len(Terminal) == 0 {
		Terminal = "/bin/bash"
	}

	// Custom log file
	f := path.Join(Home, "logs", fmt.Sprintf("%d.log", id))
	var err error
	LogFile, err = os.Create(f)
	if err != nil {
		panic(err)
	}
	log.SetOutput(LogFile)

	UpdateResources()
}

// Instances returns a list of known Goed Instances
func Instances() (ids []int64) {
	files, err := ioutil.ReadDir(path.Join(GoedHome(), "instances"))
	if err != nil {
		return ids
	}
	for _, f := range files {
		nm := f.Name()
		if strings.HasSuffix(nm, ".sock") {
			i, err := strconv.ParseInt(nm[:len(nm)-5], 10, 64)
			if err == nil {
				ids = append(ids, i)
			}
		}
	}
	// we want newer first, so revese the list
	for i, j := 0, len(ids)-1; i < j; {
		ids[i], ids[j] = ids[j], ids[i]
		i++
		j--
	}
	return ids
}

func GoedSocket(id int64) string {
	return path.Join(GoedHome(), "instances", fmt.Sprintf("%d.sock", id))
}

func GoedHome() string {
	usr, err := user.Current()
	t := ""
	home := "goed"
	if Testing {
		t = "_test"
	}
	if err != nil {
		log.Printf("Error : %s \n", err.Error())
		home = "goed"
	} else if runtime.GOOS == "windows" { // meh
		home = path.Join(usr.HomeDir, fmt.Sprintf("goed%s", t))
	} else {
		home = path.Join(usr.HomeDir, fmt.Sprintf(".goed%s", t))
	}
	return home
}

// LookupLocation will try to locate the given location
// if not found relative to dir, then try up the directory tree
// this works great to open GO import path for example
func LookupLocation(dir, loc string) (string, bool) {
	f := path.Join(dir, loc)
	stat, err := os.Stat(f)
	if err == nil {
		return f, stat.IsDir()
	}
	if dir == filepath.Dir(dir) { // at root, not found
		return loc, strings.HasSuffix(loc, string(filepath.Separator))
	}
	return LookupLocation(filepath.Dir(dir), loc)
}

func Cleanup() {
	LogFile.Close()
	info, err := LogFile.Stat()
	if err == nil && info.Size() == 0 {
		os.Remove(LogFile.Name())
	}
	os.Remove(Socket)
}

func EnvWith(custom []string) []string {
	env := os.Environ()
	for i, e := range env {
		key := strings.Split(e, "=")[0] + "="
		for _, c := range custom {
			if strings.HasPrefix(c, key) {
				env[i] = c
			}
		}
	}
	return env
}

func RunesLen(runes []rune) int {
	l := 0
	for _, r := range runes {
		l += utf8.RuneLen(r)
	}
	return l
}

func IsDir(loc string) bool {
	info, err := os.Stat(loc)
	if err == nil {
		return info.IsDir()
	}
	return strings.HasSuffix(loc, string(filepath.Separator))
}
