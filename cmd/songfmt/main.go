package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/codegangsta/cli"
	"github.com/songtools/songtools"
	"github.com/songtools/songtools/plaintext"
)

func main() {
	app := cli.NewApp()
	app.Name = "songfmt"
	app.Version = songtools.Version
	app.Author = "Craig Wilson"
	app.Usage = "Formats a song(s) according to a specified format."
	app.UsageText = fmt.Sprintf("%v [flags] path...", app.Name)
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "list, l",
			Usage: fmt.Sprintf("List files whose formatting differs from %v", app.Name),
		},
		cli.BoolFlag{
			Name:  "write, w",
			Usage: "Write result to (source) file instead of stdout",
		},
		cli.StringFlag{
			Name:  "format, f",
			Usage: "Specifies the format to be used. Defaults to the input format. Valid options are (plain).",
		},
	}
	app.Action = func(c *cli.Context) {
		list := c.Bool("list")
		write := c.Bool("write")
		format := c.String("format")
		args := c.Args()
		if len(args) == 0 {
			os.Stderr.WriteString("must specify at least one file\n")
			os.Exit(1)
		}

		if list && write {
			os.Stderr.WriteString("cannot specify both list and write\n")
			os.Exit(2)
		}

		if format != "" && format != "plain" {
			os.Stderr.WriteString(fmt.Sprintf("%q is not a valid format\n", format))
		}

		opt := &formatOptions{
			diffonly:  list,
			overwrite: write,
			parser:    plaintext.ParseSongSet,
			writer:    plaintext.WriteSongSet,
		}

		for _, p := range args {

			paths, err := filepath.Glob(p)
			if err != nil {
				os.Stderr.WriteString(fmt.Sprintf("path %q is invalid: %v", p, err))
				os.Exit(3)
			}

			for _, path := range paths {
				formatSong(path, opt)
			}
		}
	}

	app.Run(os.Args)
}

type formatOptions struct {
	diffonly  bool
	overwrite bool
	parser    songtools.SongSetParser
	writer    songtools.SongSetWriter
}

// formatSong returns true if the formatting differs and an error if there was an issue with formatting the song.
func formatSong(path string, opt *formatOptions) error {
	inputBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("unable to read %q: %v", path, err)
	}

	input := bytes.NewBuffer(inputBytes)
	set, err := opt.parser(input)
	if err != nil {
		return fmt.Errorf("unable to parse %q: %v", path, err)
	}

	if opt.diffonly {
		var output bytes.Buffer
		opt.writer(&output, set)

		if output.String() != string(inputBytes) {
			os.Stdout.WriteString(path + "\n")
		}
	} else {
		if opt.overwrite {
			var output bytes.Buffer
			opt.writer(&output, set)
			err = ioutil.WriteFile(path, output.Bytes(), 0644)
			if err != nil {
				return fmt.Errorf("unable to write %q: %v", path, err)
			}
		} else {
			opt.writer(os.Stdout, set)
		}
	}

	return nil
}
