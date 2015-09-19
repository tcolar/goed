package core

import (
	"bytes"
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
	in, err := os.Open(from)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(to)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}

// Mv file moves a file by copy, then delete
// because os.Rename does not always work
func MvFile(from, to string) error {
	if err := CopyFile(from, to); err != nil {
		return err
	}
	return os.Remove(from)
}

// CountLines does a quick (buffered) line(\n) count of a file.
func CountLines(r io.Reader) (int, error) {
	buf := make([]byte, 8192)
	count := 0

	for {
		c, err := r.Read(buf)
		if err != nil {
			if err == io.EOF {
				return count, nil
			}
			return count, err
		}
		count += bytes.Count(buf[:c], LineSep)
	}

	return count, nil
}

// StringToRunes transforms a string into a rune matrix.
func StringToRunes(s string) [][]rune {
	b := []byte(s)
	lines := bytes.Split(b, []byte("\n"))
	runes := [][]rune{}
	for i, l := range lines {
		if len(l) > 0 && l[len(l)-1] == '\r' {
			l = l[:len(l)-1]
		}
		if i != len(lines)-1 ||
			(len(l) != 0 || strings.HasSuffix(s, "\n")) {
			runes = append(runes, bytes.Runes(l))
		}
	}
	return runes
}

// RunesToString transforms a rune matrix as a string.
func RunesToString(runes [][]rune) string {
	r := []rune{}
	for i, line := range runes {
		if i != 0 && i != len(runes) {
			r = append(r, '\n')
		}
		r = append(r, line...)
	}
	return string(r)
}

// IsTextFile checks if a filke appears to be text or not(binary)
func IsTextFile(file string) bool {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return true // new file ?
	}
	buf := make([]byte, 1024)
	c, err := f.Read(buf)
	if err != nil {
		return true
	}
	return utf8.Valid(buf[:c])
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
	dir = filepath.Dir(dir)
	if strings.HasSuffix(dir, string(os.PathSeparator)) { //root
		return loc, true
	}
	return LookupLocation(dir, loc)
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
