package songtools

import (
	"io"
	"strings"
)

// FormatReader represents the ability to read a SongSet.
type FormatReader interface {
	Read(io.Reader) (*Song, error)
}

// FormatWriter represents the ability to write a SongSet.
type FormatWriter interface {
	Write(io.Writer, *Song) error
}

// Format represents a named ability to read and write a SongSet.
type Format struct {
	Name   string
	Reader FormatReader
	Writer FormatWriter
}

var registeredFormats = []*Format{}

// Formats returns all the registered formats.
func Formats() []*Format {
	return registeredFormats
}

// FormatNames returns the names of the registered formats.
func FormatNames() []string {
	names := []string{}
	for _, f := range registeredFormats {
		names = append(names, f.Name)
	}

	return names
}

// FormatByName returns a registered format by name. It compares names case-insensitively.
func FormatByName(name string) (*Format, bool) {

	name = strings.ToLower(name)
	for _, f := range registeredFormats {
		if strings.ToLower(f.Name) == name {
			return f, true
		}
	}

	return nil, false
}

// RegisterFormat registers a format.
func RegisterFormat(f *Format) {
	registeredFormats = append(registeredFormats, f)
}
