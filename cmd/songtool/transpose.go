package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/songtools/songtools"
)

func init() {
	cli.AddCommand(
		"transpose",
		"Transpose a song",
		"Transposes a song from its current key into the desired key.",
		&transposeCommand)
}

// TransposeCommand defines the options required for transposing a song.
type TransposeCommand struct {
	currentKey string `long:"currentKey" description:"The current key of the song. By default, will be discovered from the song itself."`
	format     string `short:"f" long:"format" description:"Specifies the format of the song. By default, will be discovered from the song itself."`
	args       struct {
		file flags.Filename `positional-arg-name:"file"`
		key  string         `positional-arg-name:"key"`
	} `positional-args:"true" required:"true"`
}

var transposeCommand TransposeCommand

// Execute executes the TransposeCommand.
func (cmd *TransposeCommand) Execute(args []string) error {
	var err error
	in := os.Stdin
	defer in.Close()
	if cmd.args.file != "" {
		in, err = os.Open(string(cmd.args.file))
		if err != nil {
			return fmt.Errorf("unable to open %q: %v", cmd.args.file, err)
		}
	}

	inputBytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return fmt.Errorf("unable to read: %v", err)
	}

	input := bytes.NewBuffer(inputBytes)

	format, err := findReadWriteFormat(cmd.format, string(cmd.args.file), input)
	if err != nil {
		return fmt.Errorf("unable to find format: %v", err)
	}

	song, err := format.Reader.Read(input)
	if err != nil {
		return fmt.Errorf("unable to parse: %v", err)
	}

	setSongTitleIfNecessary(string(cmd.args.file), song)

	currentKey := songtools.Key(cmd.currentKey)

	if currentKey == "" && song.Key == "" {
		return fmt.Errorf("unable to get current key")
	} else if currentKey == "" {
		currentKey = song.Key
	}

	noteNames, interval, err := songtools.NoteNamesAndIntervalFromKeyToKey(currentKey, songtools.Key(cmd.args.key))
	if err != nil {
		return fmt.Errorf("unable to get note names and interval: %v", err)
	}

	transposed, err := songtools.TransposeSong(song, interval, noteNames)
	if err != nil {
		return fmt.Errorf("unable to transpose from %q to %q: %v", currentKey, songtools.Key(cmd.args.key), err)
	}

	out := os.Stdout
	defer out.Close()

	format.Writer.Write(out, transposed)
	return nil
}
