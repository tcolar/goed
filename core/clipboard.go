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

func ClipboardWrite(s string) {
	if Testing {
		mockCbText = s
	} else {
		clipboard.WriteAll(s)
	}
}
