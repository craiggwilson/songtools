package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/codegangsta/cli"
	"github.com/songtools/songtools"
	"github.com/songtools/songtools/cmd"
	"github.com/songtools/songtools/format"
)

func main() {
	app := cli.NewApp()
	app.Name = "songtranspose"
	app.Version = cmd.Version
	app.Author = "Craig Wilson"
	app.Usage = "Transposes a song."
	app.UsageText = fmt.Sprintf("%v [flags] [path]", app.Name)
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
			Name:  "fmt, f",
			Usage: fmt.Sprintf("Specifies the format the song is currently in. Valid options are %v.", format.FilteredNames(true, true)),
		},
	}
	app.Action = func(c *cli.Context) {
		write := c.Bool("write")
		inkey := c.String("inkey")
		outkey := c.String("outkey")

		args := c.Args()
		if len(args) == 0 {
			args = append(args, "")
		}

		if outkey == "" {
			fmt.Fprintln(os.Stderr, "\"outkey\" is a required argument")
			os.Exit(2)
		}

		opt := &transposeOptions{
			overwrite: write,
			fmt:       c.String("fmt"),
			inkey:     inkey,
			outkey:    outkey,
		}

		err := transposeSong(args[0], opt)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(3)
		}
	}

	app.Run(os.Args)
}

type transposeOptions struct {
	overwrite bool
	fmt       string
	inkey     string
	outkey    string
}

func transposeSong(path string, opt *transposeOptions) error {
	var inputBytes []byte
	var err error
	if path != "" {
		inputBytes, err = ioutil.ReadFile(path)
		if err != nil {
			return fmt.Errorf("unable to read %q: %v", path, err)
		}
	} else {
		inputBytes, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("unable to read stdin: %v", err)
		}
	}

	input := bytes.NewBuffer(inputBytes)

	format, err := cmd.FindReadWriteFormat(opt.fmt, path, input)
	if err != nil {
		return fmt.Errorf("unable to find format for %q: %v", path, err)
	}

	song, err := format.Reader.Read(input)
	if err != nil {
		return fmt.Errorf("unable to parse %q: %v", path, err)
	}

	currentKey := songtools.Key(opt.inkey)

	if currentKey == "" {
		var ok bool
		if currentKey, ok = song.Key(); !ok {
			return fmt.Errorf("unable to get key for song %q", path)
		}
	}

	noteNames, interval, err := songtools.NoteNamesAndIntervalFromKeyToKey(currentKey, songtools.Key(opt.outkey))
	if err != nil {
		return fmt.Errorf("unable to get note names and interval: %v", err)
	}

	transposed, err := songtools.TransposeSong(song, interval, noteNames)
	if err != nil {
		return fmt.Errorf("unable to transpose %q from %q to %q: %v", path, currentKey, songtools.Key(opt.outkey), err)
	}

	if path != "" && opt.overwrite {
		var output bytes.Buffer
		format.Writer.Write(&output, transposed)
		err = ioutil.WriteFile(path, output.Bytes(), 0644)
		if err != nil {
			return fmt.Errorf("unable to write %q: %v", path, err)
		}
	} else {
		format.Writer.Write(os.Stdout, transposed)
	}

	return nil
}
