package core

import (
	"bytes"
	"io"
	"strings"

	"github.com/tcolar/goed/assert"
	. "gopkg.in/check.v1"
)

func (s *CoreSuite) TestDropCrLf(t *C) {
	tr := DropCrLfTransformer{}
	src := []byte("foo\nbar")
	dst := make([]byte, 100)
	d, _, err := tr.Transform(dst, src, true)
	if err != nil {
		panic(err)
	}
	assert.DeepEq(t, dst[:d], []byte("foo\nbar"))

	src = []byte("aaa\r\nbbb\r\nccc")
	dst = make([]byte, 100)
	d, _, err = tr.Transform(dst, src, true)
	if err != nil {
		panic(err)
	}
	assert.DeepEq(t, dst[:d], []byte("aaa\nbbb\nccc"))

	src = []byte("\r\n\r\naaa\r\n")
	dst = make([]byte, 100)
	d, _, err = tr.Transform(dst, src, true)
	if err != nil {
		panic(err)
	}
	assert.DeepEq(t, dst[:d], []byte("\n\naaa\n"))
}

func (s *CoreSuite) TestAddCrLf(t *C) {
	tr := AddCrLfTransformer{}
	src := []byte("foo\nbar")
	dst := make([]byte, 100)
	d, _, err := tr.Transform(dst, src, true)
	if err != nil {
		panic(err)
	}
	assert.DeepEq(t, dst[:d], []byte("foo\r\nbar"))

	src = []byte("aaa\nbbb\nccc")
	dst = make([]byte, 100)
	d, _, err = tr.Transform(dst, src, true)
	if err != nil {
		panic(err)
	}
	assert.DeepEq(t, dst[:d], []byte("aaa\r\nbbb\r\nccc"))

	src = []byte("\na\r\nb\rc\n\n")
	dst = make([]byte, 100)
	d, _, err = tr.Transform(dst, src, true)
	if err != nil {
		panic(err)
	}
	assert.DeepEq(t, dst[:d], []byte("\r\na\r\nb\rc\r\n\r\n"))
}

func (s *CoreSuite) TestCrLfEncoding(t *C) {
	es := "\r\nFoo\r\nBar\r\nBazz\rBuzz\r\n\r\nZzz"
	ds := "\nFoo\nBar\nBazz\rBuzz\n\nZzz"
	enc := CrLfEncoding{}.NewEncoder()
	dec := CrLfEncoding{}.NewDecoder()

	// test decoding
	r := strings.NewReader(es)
	w := bytes.NewBuffer(make([]byte, 0, 3))
	io.Copy(w, dec.Reader(r))
	assert.DeepEq(t, w.Bytes(), []byte(ds))

	// test encoding
	r = strings.NewReader(ds)
	w = bytes.NewBuffer(make([]byte, 0, 3))
	io.Copy(enc.Writer(w), r)
	assert.DeepEq(t, w.Bytes(), []byte(es))

	// test encoding with ErrShortDst
	b := make([]byte, 5000, 5000)
	b[0] = 'A'
	b[4999] = 'Z'
	for i := 1; i != 4999; i++ {
		b[i] = '\n'
	}
	r2 := bytes.NewReader(b)
	w = bytes.NewBuffer(make([]byte, 0, 3))
	_, err := io.Copy(enc.Writer(w), r2)
	data := w.Bytes()
	assert.Nil(t, err)
	assert.Eq(t, len(data), 9998)
	for i := 0; i != 9998; i++ {
		c := int32(data[i])
		switch i {
		case 0:
			assert.Eq(t, c, 'A')
		case 9997:
			assert.Eq(t, c, 'Z')
		default:
			if i%2 == 0 {
				assert.Eq(t, c, '\n')
			} else {
				assert.Eq(t, c, '\r')
			}
		}
	}
}
