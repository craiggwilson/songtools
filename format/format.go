package format

import (
	"io"
	"strings"

	"github.com/songtools/songtools"
)

const (
	// Default is the default format to use when none is specified or inferred.
	Default = "chordsOverLyrics"
)

// Reader represents the ability to read a SongSet.
type Reader interface {
	Read(io.Reader) (*songtools.Song, error)
}

// Writer represents the ability to write a SongSet.
type Writer interface {
	Write(io.Writer, *songtools.Song) error
}

// Formats is a slice of Formats.
type Formats []*Format

// Filter uses the filter to filter out certain formats.
func (fs Formats) Filter(filter func(*Format) bool) Formats {
	formats := Formats{}
	for _, f := range fs {
		if filter(f) {
			formats = append(formats, f)
		}
	}

	return formats
}

var registeredFormats = Formats{}

// Format represents a named ability to read and write a SongSet.
type Format struct {
	Name       string
	Reader     Reader
	Writer     Writer
	Extensions []string
}

// CanRead indicates whether the format can be used to read a song.
func (f *Format) CanRead() bool {
	return f.Reader != nil
}

// CanWrite indicates whether the format can be used to write a song.
func (f *Format) CanWrite() bool {
	return f.Writer != nil
}

// RegisteredFormats returns all the registered formats.
func RegisteredFormats() Formats {
	return registeredFormats
}

// RegisteredNames returns the names of the registered formats that can read and write.
func RegisteredNames() []string {
	return FilteredNames(func(f *Format) bool {
		return true
	})
}

// FilteredNames returns the names of the registered formats that match the criteria.
func FilteredNames(filter func(*Format) bool) []string {
	names := []string{}
	formats := registeredFormats.Filter(filter)
	for _, f := range formats {
		names = append(names, f.Name)
	}

	return names
}

// ByName returns a registered format by name. It compares names case-insensitively.
func ByName(name string) (*Format, bool) {

	name = strings.ToLower(name)
	for _, f := range registeredFormats {
		if strings.ToLower(f.Name) == name {
			return f, true
		}
	}

	return nil, false
}

// Register registers a format.
func Register(f *Format) {
	registeredFormats = append(registeredFormats, f)
}
