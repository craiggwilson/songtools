package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/codegangsta/cli"
	"github.com/songtools/songtools"
	"github.com/songtools/songtools/plaintext"
)

func main() {
	app := cli.NewApp()
	app.Name = "songtranspose"
	app.Version = songtools.Version
	app.Author = "Craig Wilson"
	app.Usage = "Transposes a song."
	app.UsageText = fmt.Sprintf("%v [flags] path", app.Name)
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "inkey",
			Usage: "The current key of the song",
		},
		cli.StringFlag{
			Name:  "outkey",
			Usage: "The desired key of the song",
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
		write := c.Bool("write")
		format := c.String("format")
		inkey := c.String("inkey")
		outkey := c.String("outkey")

		args := c.Args()
		if len(args) == 0 {
			fmt.Fprintln(os.Stderr, "must specify a file to transpose")
			os.Exit(1)
		}

		if format != "" && format != "plain" {
			fmt.Fprintln(os.Stderr, fmt.Sprintf("%q is not a valid format\n", format))
			os.Exit(2)
		}

		opt := &transposeOptions{
			overwrite: write,
			parser:    plaintext.ParseSongSet,
			writer:    plaintext.WriteSongSet,
			inkey:     inkey,
			outkey:    outkey,
		}

		err := transposeSong(args[0], opt)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}

	app.Run(os.Args)
}

type transposeOptions struct {
	overwrite bool
	parser    songtools.SongSetParser
	writer    songtools.SongSetWriter
	inkey     string
	outkey    string
}

func transposeSong(path string, opt *transposeOptions) error {
	noteNames, interval, err := songtools.NoteNamesAndIntervalFromKeyToKey(opt.inkey, opt.outkey)
	if err != nil {
		return fmt.Errorf("unable to get note names and interval: %v", err)
	}

	inputBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("unable to read %q: %v", path, err)
	}

	input := bytes.NewBuffer(inputBytes)
	set, err := opt.parser(input)
	if err != nil {
		return fmt.Errorf("unable to parse %q: %v", path, err)
	}

	transposed, err := songtools.TransposeSongSet(set, interval, noteNames)
	if err != nil {
		return fmt.Errorf("unable to transpose %q from %q to %q: %v", path, opt.inkey, opt.outkey, err)
	}

	if opt.overwrite {
		var output bytes.Buffer
		opt.writer(&output, transposed)
		err = ioutil.WriteFile(path, output.Bytes(), 0644)
		if err != nil {
			return fmt.Errorf("unable to write %q: %v", path, err)
		}
	} else {
		opt.writer(os.Stdout, transposed)
	}

	return nil
}
