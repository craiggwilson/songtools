package cmd

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/songtools/songtools"
	"github.com/songtools/songtools/format"
	_ "github.com/songtools/songtools/format/chordsOverLyrics" // formats are registered in the init functions.
	_ "github.com/songtools/songtools/format/html"             // formats are registered in the init functions.
	_ "github.com/songtools/songtools/format/pdf"              // formats are registered in the init functions.
)

// Version is the version of the applications.
const Version = "0.1"

// FormatFinder is a function definition for finding a format.
type FormatFinder func(string, string, *bytes.Buffer) (*format.Format, error)

// FindReadWriteFormat gets a format that can read and write for the given path, buffer, and format.
func FindReadWriteFormat(name, path string, buffer *bytes.Buffer) (*format.Format, error) {
	f, err := findFormat(name, path, buffer)
	if err != nil {
		return nil, err
	}

	if f.Reader == nil {
		return nil, fmt.Errorf("the format %q cannot be used for reading", f.Name)
	}

	if f.Writer == nil {
		return nil, fmt.Errorf("the format %q cannot be used for writer", f.Name)
	}

	return f, nil
}

// FindReadFormat gets a format that can read for the given path, buffer, and format.
func FindReadFormat(name, path string, buffer *bytes.Buffer) (*format.Format, error) {
	f, err := findFormat(name, path, buffer)
	if err != nil {
		return nil, err
	}

	if f.Reader == nil {
		return nil, fmt.Errorf("the format %q cannot be used for reading", f.Name)
	}

	return f, nil
}

func findFormat(name, path string, buffer *bytes.Buffer) (*format.Format, error) {
	var f *format.Format
	if name == "" {
		// TODO: try to detect format.
		f, _ = format.ByName(format.Default)
	} else {
		var ok bool
		if f, ok = format.ByName(name); !ok {
			return nil, fmt.Errorf("unable to find format %q", name)
		}
	}

	return f, nil
}

// SetSongTitleIfNecessary sets the song title to the path if the song doesn't already have a title.
func SetSongTitleIfNecessary(path string, s *songtools.Song) {
	if path != "" && s.Title == "" {
		_, filename := filepath.Split(path)
		ext := filepath.Ext(path)
		if ext != "" && strings.HasSuffix(filename, ext) {
			filename = strings.TrimSuffix(filename, ext)
		}
		s.Title = filename
	}
}
