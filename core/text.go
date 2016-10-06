// CountLines does a quick (buffered) line(\n) count of a file.
package core

import (
	"bytes"
	"io"
	"os"
	"strings"
	"unicode/utf8"
)

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

// IsTextFile checks if a file appears to be text or not(binary)
// It's not always deterministic so it tries best guess.
func IsTextFile(file string) bool {
	// if it's a new/empty file, it can be a text file
	if stats, err := os.Stat(file); os.IsNotExist(err) || stats.Size() == 0 {
		return true
	}
	// if we find  BOM, Vey High odds it's a text file
	bomLen, _, _ := CheckBom(file)
	if bomLen > 0 {
		return true
	}
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return true
	}
	// does it only contain ut8 characters ? -> likely utf8
	buf := make([]byte, 1024)
	c, err := f.Read(buf)
	if err != nil {
		return true
	}
	if utf8.Valid(buf[:c]) {
		return true
	}
	// ok, so it's either ut16 without bom or binary
	// trying to determine
	countNewLinesLe, countNewLinesBe := 0, 0
	oddNulls, evenNulls := 0, 0
	for i := 1; i < c; i++ {
		if i%2 == 0 && buf[i] == 0 {
			evenNulls++
		}
		if i%2 != 1 && buf[i] == 0 {
			oddNulls++
		}
		if buf[i-1] == 0x0A && buf[i] == 0x00 {
			countNewLinesLe++
		} else if buf[i-1] == 0x0 && buf[i] == 0x0A {
			countNewLinesBe++
		}
	}
	if countNewLinesLe >= 4 {
		return true // likely utf16 LittleEndian text
	}
	if countNewLinesBe >= 4 {
		return true // likely utf16 BigEndian text
	}
	if oddNulls > c/2 {
		return true // likely utf16 LittleEndian (lots of little endian ascii bytes)
	}
	if evenNulls > c/2 {
		return true // likely utf16 BigEndian text (tots of gib endian ascii bytes)
	}
	// doesn't look like text
	if c < 1000 {
		return true // probably binary but file is small, not too risky to try it as a text file
	}
	return false // all else failed, assume binary
}

// Check if the file starts with a bom and if so return it's info
// https://en.wikipedia.org/wiki/Byte_order_mark
func CheckBom(from string) (bomLen, bytesPerChar int, e Endianness) {
	in, err := os.Open(from)
	if err != nil {
		return 0, 1, LittleEndian
	}
	defer in.Close()
	bom := make([]byte, 5)
	read, _ := in.Read(bom) // @5
	if read >= 5 {
		if bom[0] == 0x2B && bom[1] == 0x2F && bom[2] == 0x76 && bom[3] == 0x38 && bom[4] == 0x2D {
			return 5, 1, LittleEndian // UTF-7
		}
	}
	if read >= 4 {
		if bom[0] == 0x00 && bom[1] == 0x00 && bom[2] == 0xFE && bom[3] == 0xFF {
			return 4, 4, BigEndian // UTF-32 BE
		}
		if bom[0] == 0xFF && bom[1] == 0xFE && bom[2] == 0x00 && bom[3] == 0x00 {
			return 4, 4, LittleEndian // UTF-32 LE
		}
		if bom[0] == 0xDD && bom[1] == 0x73 && bom[2] == 0x66 && bom[3] == 0x73 {
			return 4, 1, LittleEndian // UTF-EBCDIC
		}
		if bom[0] == 0x84 && bom[1] == 0x31 && bom[2] == 0x95 && bom[3] == 0x33 {
			return 4, 2, LittleEndian // GB-18030
		}
		if bom[0] == 0x2B && bom[1] == 0x2F && bom[2] == 0x76 &&
			(bom[3] == 0x38 || bom[3] == 0x39 || bom[3] == 0x2B || bom[3] == 0x2F) {
			return 4, 1, LittleEndian // UTF-7
		}
	}
	if read >= 3 {
		if bom[0] == 0xEF && bom[1] == 0xBB && bom[2] == 0xBF {
			return 3, 1, LittleEndian // UTF-8
		}
		if bom[0] == 0xF7 && bom[1] == 0x64 && bom[2] == 0x4C {
			return 3, 1, LittleEndian // UTF-1
		}
		if bom[0] == 0x0E && bom[1] == 0xFE && bom[2] == 0xFF {
			return 3, 1, LittleEndian // SCSU
		}
		if bom[0] == 0xFB && bom[1] == 0xEE && bom[2] == 0x28 {
			return 3, 1, LittleEndian // BOCU-1
		}
	}
	if read >= 2 {
		if bom[0] == 0xFE && bom[1] == 0xFF {
			return 2, 2, BigEndian // UTF-16 BE
		}
		if bom[0] == 0xFF && bom[1] == 0xFE {
			return 2, 2, LittleEndian // UTF-16 LE
		}
	}
	return 0, 0, LittleEndian
}
