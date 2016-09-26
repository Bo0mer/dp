package main

import (
	"bytes"
	"encoding/json"
)

// Formatter formats bytes that represent text in different encodings.
type Formatter interface {
	Format([]byte) ([]byte, error)
}

type nopformatter struct{}

func (n nopformatter) Format(src []byte) ([]byte, error) { return src, nil }

// NopFormatter returns formatter that does nothing.
func NopFormatter() Formatter {
	return nopformatter{}
}

// JSONFormatter formats JSON.
type JSONFormatter struct {
	// Prefix is appendend to each new line of the formatted string.
	Prefix string
	// Indent is the symbol to use for indentation.
	Indent string
}

// Format returns src pretty formatted.
func (f *JSONFormatter) Format(src []byte) ([]byte, error) {
	if len(src) == 0 {
		return []byte{}, nil
	}
	var buf bytes.Buffer
	err := json.Indent(&buf, src, f.Prefix, f.Indent)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
