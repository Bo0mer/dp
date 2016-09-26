package main

// Printer provides functionality for printing optionally formated lines.
type Printer interface {
	// Printf formats according to a format specifier and writes to the printer.
	// It returns the number of bytes written and any write error encountered.
	Printf(string, ...interface{}) (int, error)
	// Println formats using the default formats for its operands and writes to
	// the printer. It should return the number of bytes written and any write
	// error encountered.
	Println(...interface{}) (int, error)
}
