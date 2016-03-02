package main

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/songtools/songtools"
	"github.com/songtools/songtools/format"
)

const version = "0.1"

type formatFinder func(string, string, *bytes.Buffer) (*format.Format, error)

func findReadWriteFormat(name, path string, buffer *bytes.Buffer) (*format.Format, error) {
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

func findReadFormat(name, path string, buffer *bytes.Buffer) (*format.Format, error) {
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
func setSongTitleIfNecessary(path string, s *songtools.Song) {
	if path != "" && s.Title == "" {
		_, filename := filepath.Split(path)
		ext := filepath.Ext(path)
		if ext != "" && strings.HasSuffix(filename, ext) {
			filename = strings.TrimSuffix(filename, ext)
		}
		s.Title = filename
	}
}
