package main

import (
	"bytes"
	"io"
	"os"
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
	buf := make([]byte, 8196)
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
