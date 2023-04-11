package main

import (
	"io"
	"log"
	"os"
)

// Thing is a very silly thing to illustrate how test structure
// can indicate what is or is not important.
type Thing struct {
	Msg string

	// These are different ways to implement some of the same functionality.
	// They are useful in different ways, but also signal what's important,
	// and that can show up in tests.

	// LogFile is the least flexible, definitely gonna write to filesystem.
	LogFile string
	// Logger has more config options, but still a specific type.
	Logger *log.Logger
	// Writer can be anything that implements the interface.
	// generally, an interface may be a good spot to mock stuff,
	// if its implementation details are unimportant to your tests.
	Writer io.Writer
}

// NewThing is the default happy-path constructor that most things will use.
// This one just writes to stdout.
func NewThing(msg string) *Thing {
	return &Thing{
		Msg:    msg,
		Writer: os.Stdout,
	}
}

// WriteLogFile makes a file and writes to it, always.
func (t *Thing) WriteLogFile() error {
	f, err := os.Create(t.LogFile)
	if err != nil {
		return err
	}
	_, err = f.Write([]byte(t.Msg))
	return err
}

// WriteLogger uses the logger, however it's configured.
func (t *Thing) WriteLogger() {
	t.Logger.Println(t.Msg)
}

// WriteWriter does any kind of Write() in the whole wide world.
func (t *Thing) WriteWriter() error {
	_, err := t.Writer.Write([]byte(t.Msg))
	return err
}
