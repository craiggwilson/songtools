package cmd

import (
	"bytes"
	"fmt"

	"github.com/songtools/songtools"
	_ "github.com/songtools/songtools/plain" // formats are registered in the init functions.
)

// Version is the version of the applications.
const Version = "0.1"

// FindReadFormat gets a format that can read and write for the given path, buffer, and format.
func FindReadWriteFormat(name, path string, buffer *bytes.Buffer) (*songtools.Format, error) {
	format, err := findFormat(name, path, buffer)
	if err != nil {
		return nil, err
	}

	if format.Reader == nil {
		return nil, fmt.Errorf("the format %q cannot be used for reading", format.Name)
	}

	if format.Writer == nil {
		return nil, fmt.Errorf("the format %q cannot be used for writer", format.Name)
	}

	return format, nil
}

// FindReadFormat gets a format that can read for the given path, buffer, and format.
func FindReadFormat(name, path string, buffer *bytes.Buffer) (*songtools.Format, error) {
	format, err := findFormat(name, path, buffer)
	if err != nil {
		return nil, err
	}

	if format.Reader == nil {
		return nil, fmt.Errorf("the format %q cannot be used for reading", format.Name)
	}

	return format, nil
}

func findFormat(name, path string, buffer *bytes.Buffer) (*songtools.Format, error) {
	var format *songtools.Format
	if name == "" {
		// TODO: try to detect format.
		format, _ = songtools.FormatByName("plain")
	} else {
		var ok bool
		if format, ok = songtools.FormatByName(name); !ok {
			return nil, fmt.Errorf("unable to find format %q", name)
		}
	}

	return format, nil
}
