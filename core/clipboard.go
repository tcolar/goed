package core

import "github.com/atotto/clipboard"

var mockCbText string = ""

func ClipboardRead() (string, error) {
	if Testing {
		return mockCbText, nil
	} else {
		return clipboard.ReadAll()
	}
}

func ClipboardWrite(s string) error {
	if Testing {
		mockCbText = s
		return nil
	} else {
		return clipboard.WriteAll(s)
	}
}
