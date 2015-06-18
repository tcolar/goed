package core

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"runtime"
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

func StringToRunes(s string) [][]rune {
	b := []byte(s)
	lines := bytes.Split(b, []byte("\n"))
	runes := [][]rune{}
	for i, l := range lines {
		if i != len(lines)-1 ||
			(len(l) != 0 || strings.HasSuffix(s, "\n")) {
			runes = append(runes, bytes.Runes(l))
		}
	}
	return runes
}

// RunesToString returns a rune section as a srting.
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

func InitHome() {
	usr, err := user.Current()
	t := ""
	if Testing {
		t = "_test"
	}
	if err != nil {
		fmt.Printf("Error : %s \n", err.Error())
		Home = "goed"
	} else if runtime.GOOS == "windows" { // meh
		Home = path.Join(usr.HomeDir, fmt.Sprintf("goed%s", t))
	} else {
		Home = path.Join(usr.HomeDir, fmt.Sprintf(".goed%s", t))
	}
	os.MkdirAll(Home, 0777)
	os.MkdirAll(path.Join(Home, "buffers"), 0777)
	// TODO : Update config if new version ??
	ioutil.WriteFile(path.Join(Home, "Version.txt"), []byte(Version), 644)
}
