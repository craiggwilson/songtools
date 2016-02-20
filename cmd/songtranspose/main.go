package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/codegangsta/cli"
	"github.com/songtools/songtools"
	"github.com/songtools/songtools/cmd"
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
			Name:  "fmt, f",
			Usage: "Specifies the format to be used. Valid options are (plain).",
		},
	}
	app.Action = func(c *cli.Context) {
		write := c.Bool("write")
		inkey := c.String("inkey")
		outkey := c.String("outkey")

		args := c.Args()
		if len(args) == 0 {
			fmt.Fprintln(os.Stderr, "must specify a file to transpose")
			os.Exit(1)
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
	inputBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("unable to read %q: %v", path, err)
	}
	input := bytes.NewBuffer(inputBytes)

	format, err := cmd.FindReadWriteFormat(opt.fmt, path, input)
	if err != nil {
		return fmt.Errorf("unable to find format for %q: %v", path, err)
	}

	set, err := format.Reader.Read(input)
	if err != nil {
		return fmt.Errorf("unable to parse %q: %v", path, err)
	}

	inkey := opt.inkey
	if inkey == "" {
		chords := set.Songs[0].Chords()
		fmt.Printf("% v\n", chords)
		os.Exit(22)
	}

	noteNames, interval, err := songtools.NoteNamesAndIntervalFromKeyToKey(inkey, opt.outkey)
	if err != nil {
		return fmt.Errorf("unable to get note names and interval: %v", err)
	}

	transposed, err := songtools.TransposeSongSet(set, interval, noteNames)
	if err != nil {
		return fmt.Errorf("unable to transpose %q from %q to %q: %v", path, opt.inkey, opt.outkey, err)
	}

	if opt.overwrite {
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
