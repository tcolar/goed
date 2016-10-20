// CountLines does a quick (buffered) line(\n) count of a file.
package core

import (
	"bytes"
	"io"
	"os"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/encoding/unicode/utf32"
	"golang.org/x/text/transform"
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

func HasWindowsNewLine(r io.Reader) bool {
	buf := make([]byte, 1000)
	c, _ := r.Read(buf)
	lf := bytes.Count(buf[:c], LineSep)
	crlf := bytes.Count(buf[:c], []byte{'\r', '\n'})
	if crlf > lf/2 {
		return true
	}
	return false
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

type TextInfo struct {
	Enc encoding.Encoding
}

// ReadTextInfo checks if a file appears to be text or not(binary)
// Returns nil if the file appears binary or some unsupported encoding.
func ReadTextInfo(file string, srcHasWindowsNewLines bool) *TextInfo {
	// if it's a new/empty file, it can be a UTF8 text file
	if stats, err := os.Stat(file); os.IsNotExist(err) || stats.Size() == 0 {
		return CrLfTextInfo(nil, srcHasWindowsNewLines)
	}
	// if starts with a BOM, Vey High odds it's a text file
	bomEnc := BomEncoding(file)
	if bomEnc != nil {
		return CrLfTextInfo(bomEnc, srcHasWindowsNewLines)
	}
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return CrLfTextInfo(nil, srcHasWindowsNewLines)
	}
	// does it only contain ut8 characters ? -> likely utf8
	buf := make([]byte, 1024)
	c, err := f.Read(buf)
	if err != nil {
		return CrLfTextInfo(unicode.UTF8, srcHasWindowsNewLines)
	}
	if utf8.Valid(buf[:c]) {
		return CrLfTextInfo(unicode.UTF8, srcHasWindowsNewLines)
	}
	// ok, so it's either utf16 without bom or binary (or some other unsuported encoding)
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
		return CrLfTextInfo(unicode.UTF16(unicode.LittleEndian, unicode.UseBOM), srcHasWindowsNewLines) // likely utf16 LittleEndian text
	}
	if countNewLinesBe >= 4 {
		return CrLfTextInfo(unicode.UTF16(unicode.BigEndian, unicode.UseBOM), srcHasWindowsNewLines) // likely utf16 BigEndian text
	}
	if oddNulls > c/2 {
		return CrLfTextInfo(unicode.UTF16(unicode.LittleEndian, unicode.UseBOM), srcHasWindowsNewLines) // likely utf16 LittleEndian (lots of little endian ascii bytes)
	}
	if evenNulls > c/2 {
		return CrLfTextInfo(unicode.UTF16(unicode.BigEndian, unicode.UseBOM), srcHasWindowsNewLines) // likely utf16 BigEndian text (tots of gib endian ascii bytes)
	}
	// doesn't look like text
	if c < 1000 {
		return CrLfTextInfo(unicode.UTF8, srcHasWindowsNewLines) // probably binary but file is small, not too risky to try it as a text file
	}
	return nil // all else failed, assume binary / unsupported
}

// return TextInfo with extra CrLf encoding/decoding if needed
func CrLfTextInfo(enc encoding.Encoding, srcHasWindowsNewLines bool) *TextInfo {
	if !srcHasWindowsNewLines {
		return &TextInfo{
			Enc: enc,
		}
	}
	// wrap with CRLF encoder/decoder
	return &TextInfo{
		Enc: &CrLfEncoding{
			ChainWith: enc,
		},
	}
}

// Check if the file starts with a bom and if so return the encoding
// Returns nil if no BOM or unsupported encoding
func BomEncoding(from string) encoding.Encoding {
	in, err := os.Open(from)
	if err != nil {
		return nil
	}
	defer in.Close()
	bom := make([]byte, 4)
	read, _ := in.Read(bom) // @4
	if read >= 4 {
		if bom[0] == 0x00 && bom[1] == 0x00 && bom[2] == 0xFE && bom[3] == 0xFF {
			return utf32.UTF32(utf32.BigEndian, utf32.UseBOM) // UTF-32 BE
		}
		if bom[0] == 0xFF && bom[1] == 0xFE && bom[2] == 0x00 && bom[3] == 0x00 {
			return utf32.UTF32(utf32.LittleEndian, utf32.UseBOM) // UTF-32 LE
		}
		if bom[0] == 0x84 && bom[1] == 0x31 && bom[2] == 0x95 && bom[3] == 0x33 {
			return simplifiedchinese.GB18030 // GB-18030
		}
	}
	if read >= 3 {
		if bom[0] == 0xEF && bom[1] == 0xBB && bom[2] == 0xBF {
			return unicode.UTF8 // UTF-8
		}
	}
	if read >= 2 {
		if bom[0] == 0xFE && bom[1] == 0xFF {
			return unicode.UTF16(unicode.BigEndian, unicode.UseBOM) // UTF-16 BE
		}
		if bom[0] == 0xFF && bom[1] == 0xFE {
			return unicode.UTF16(unicode.LittleEndian, unicode.UseBOM) // UTF-16 LE
		}
	}
	return nil
}

// CrLfEncoding encode / decodes '\r\n' to '\n'
type CrLfEncoding struct {
	ChainWith encoding.Encoding
}

var _ encoding.Encoding = (*CrLfEncoding)(nil)

func (c CrLfEncoding) NewDecoder() *encoding.Decoder {
	if c.ChainWith != nil {
		return &encoding.Decoder{
			Transformer: transform.Chain(c.ChainWith.NewDecoder(), DropCrLfTransformer{}),
		}
	}
	return &encoding.Decoder{
		Transformer: DropCrLfTransformer{},
	}
}

func (c CrLfEncoding) NewEncoder() *encoding.Encoder {
	if c.ChainWith != nil {
		return &encoding.Encoder{
			Transformer: transform.Chain(AddCrLfTransformer{}, c.ChainWith.NewEncoder()),
		}
	}
	return &encoding.Encoder{
		Transformer: AddCrLfTransformer{},
	}
}

// Drop Windows "\r\n" combos (in favor of plain "\n")
type DropCrLfTransformer struct{}

func (t DropCrLfTransformer) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	dstAt := 0
	srcAt := 0
	for i := 1; i < len(src); i++ {
		if src[i-1] == '\r' && src[i] == '\n' {
			copy(dst[dstAt:], src[srcAt:i-1])
			dstAt += i - srcAt - 1
			srcAt = i
		}
	}
	ln := len(src) - srcAt
	copy(dst[dstAt:], src[srcAt:srcAt+ln])
	srcAt += ln
	dstAt += ln
	return dstAt, srcAt, nil
}

func (t DropCrLfTransformer) Reset() {}

// Replace "\n" with "\r\n" windows combos
type AddCrLfTransformer struct{}

func (t AddCrLfTransformer) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	dstAt := 0
	srcAt := 0
	for i := 0; i < len(src); i++ {
		if src[i] == '\n' && (i == 0 || src[i-1] != '\r') {
			if dstAt+i-srcAt+2 >= len(dst) {
				return dstAt, srcAt, transform.ErrShortDst
			}
			copy(dst[dstAt:], src[srcAt:i])
			dstAt += i - srcAt
			srcAt = i + 1
			copy(dst[dstAt:], []byte{'\r', '\n'})
			dstAt += 2
		}
	}
	ln := len(src) - srcAt
	if dstAt+ln >= len(dst) {
		return dstAt, srcAt, transform.ErrShortDst
	}
	copy(dst[dstAt:], src[srcAt:srcAt+ln])
	srcAt += ln
	dstAt += ln
	return dstAt, srcAt, nil
}

func (t AddCrLfTransformer) Reset() {}
