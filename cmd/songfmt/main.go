package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/codegangsta/cli"
	"github.com/songtools/songtools/cmd"
	"github.com/songtools/songtools/format"
)

func main() {
	app := cli.NewApp()
	app.Name = "songfmt"
	app.Version = cmd.Version
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
			Name:  "infmt",
			Usage: fmt.Sprintf("Specifies the format the song is currently in. Valid options are %v.", format.Names()),
		},
		cli.StringFlag{
			Name:  "outfmt",
			Usage: fmt.Sprintf("Specifies the format for output. Defaults to the same as infmt. Valid options are %v.", format.Names()),
		},
	}
	app.Action = func(c *cli.Context) {
		list := c.Bool("list")
		write := c.Bool("write")

		args := c.Args()
		if len(args) == 0 {
			fmt.Fprintln(os.Stderr, "must specify at least one file")
			os.Exit(1)
		}

		if list && write {
			fmt.Fprintln(os.Stderr, "cannot specify both list and write")
			os.Exit(2)
		}

		opt := &formatOptions{
			diffonly:  list,
			overwrite: write,
			infmt:     c.String("infmt"),
			outfmt:    c.String("outfmt"),
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
	infmt     string
	outfmt    string
}

func formatSong(path string, opt *formatOptions) error {
	inputBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("unable to read %q: %v", path, err)
	}
	input := bytes.NewBuffer(inputBytes)

	readFormat, err := cmd.FindReadFormat(opt.infmt, path, input)
	if err != nil {
		return fmt.Errorf("unable to find input format for %q: %v", path, err)
	}

	var writeFormat *format.Format
	if opt.outfmt == "" {
		if readFormat.Writer == nil {
			return fmt.Errorf("the input format %q is unable to be used for writing", readFormat.Name)
		}

		writeFormat = readFormat
	} else {
		var ok bool
		if writeFormat, ok = format.ByName(opt.outfmt); !ok {
			return fmt.Errorf("unable to find output format %q for %q", opt.outfmt, path)
		}
	}

	song, err := readFormat.Reader.Read(input)
	if err != nil {
		return fmt.Errorf("unable to parse %q: %v", path, err)
	}

	if opt.diffonly {
		var output bytes.Buffer
		writeFormat.Writer.Write(&output, song)

		if output.String() != string(inputBytes) {
			fmt.Fprintln(os.Stdout, path)
		}
	} else {
		if opt.overwrite {
			var output bytes.Buffer
			writeFormat.Writer.Write(&output, song)
			err = ioutil.WriteFile(path, output.Bytes(), 0644)
			if err != nil {
				return fmt.Errorf("unable to write %q: %v", path, err)
			}
		} else {
			writeFormat.Writer.Write(os.Stdout, song)
		}
	}

	return nil
}
