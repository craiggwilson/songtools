package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/songtools/songtools"
	"github.com/songtools/songtools/cmd"
)

type options struct {
	inkey     string `long:"inkey" descripton:"The current key of the song. The default will try to be inferred from the input."`
	outkey    string `long:"outkey" required:"true" descripton:"The desired key of the song."`
	overwrite string `short:"w" long:"overwrite" descripton:"Write the result to the source file. When false, output is sent to stdout. The default is false."`
	format    string `short:"f" long:"format" description:"Specifies the format of the song."`
	file      string `required:"true" description:"The song file to transpose."`
}

func main() {
	parser := flags.NewNamedParser("songtool", flags.Default)

	opts := &options{}
	g := parser.AddGroup("Global", "Global options", opts)

	_, err := parser.ParseArgs(os.Args)
	if err != nil {
		os.Exit(1)
	}
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

	cmd.SetSongTitleIfNecessary(path, song)

	currentKey := songtools.Key(opt.inkey)

	if currentKey == "" && song.Key == "" {
		return fmt.Errorf("unable to get key for song %q", path)
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
