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
			fmt.Fprintln(os.Stderr, "must specify at least one file")
			os.Exit(1)
		}

		if list && write {
			fmt.Fprintln(os.Stderr, "cannot specify both list and write")
			os.Exit(2)
		}

		if format != "" && format != "plain" {
			fmt.Fprintln(os.Stderr, fmt.Sprintf("%q is not a valid format", format))
			os.Exit(3)
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
				fmt.Fprintln(os.Stderr, fmt.Sprintf("path %q is invalid: %v", p, err))
				os.Exit(3)
			}

			for _, path := range paths {
				err := formatSong(path, opt)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
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
			fmt.Fprintln(os.Stdout, path)
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
