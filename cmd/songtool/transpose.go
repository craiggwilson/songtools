package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/songtools/songtools"
)

// TransposeCommand defines the options required for transposing a song.
type TransposeCommand struct {
	Format string `long:"format" description:"Specifies the format of the song. By default, will be discovered from the song itself."`
	From   string `short:"f" long:"from" description:"The current key of the song. By default, will be discovered from the song itself."`
	To     string `required:"true" short:"t" long:"to" description:"The desired key of the song."`
	Out    string `short:"o" long:"out" optional:"true" optional-value:"<input>" description:"The file to write the transposed song. If left unspecified, stdout will be used. When specified without an argument, the input file will be overwritten."`
}

var transposeCommand TransposeCommand

func init() {

	cli.AddCommand(
		"transpose",
		"Transpose a song",
		"Transposes a song from its current key into the desired key.",
		&transposeCommand)
}

// Execute executes the TransposeCommand.
func (cmd *TransposeCommand) Execute(args []string) error {

	if len(args) > 1 {
		return fmt.Errorf("too many positional arguments")
	}

	var err error
	in := os.Stdin
	file := ""
	if len(args) == 1 {
		file = args[0]
		in, err = os.Open(file)
		if err != nil {
			return fmt.Errorf("unable to open %q: %v", file, err)
		}

		if cmd.Out == "<input>" {
			cmd.Out = file
		}
	} else if cmd.Out == "<input>" {
		cmd.Out = ""
	}

	inBytes, err := ioutil.ReadAll(in)
	if err != nil {
		return fmt.Errorf("unable to read: %v", err)
	}

	input := bytes.NewBuffer(inBytes)

	format, err := findReadWriteFormat(cmd.Format, file, input)
	if err != nil {
		return fmt.Errorf("unable to find format: %v", err)
	}

	song, err := format.Reader.Read(input)
	if err != nil {
		return fmt.Errorf("unable to parse: %v", err)
	}

	setSongTitleIfNecessary(args[0], song)

	fromKey := songtools.Key(cmd.From)

	if fromKey == "" && song.Key == "" {
		return fmt.Errorf("unable to get current key")
	} else if fromKey == "" {
		fromKey = song.Key
	}

	toKey := songtools.Key(cmd.To)

	noteNames, interval, err := songtools.NoteNamesAndIntervalFromKeyToKey(fromKey, toKey)
	if err != nil {
		return fmt.Errorf("unable to get note names and interval: %v", err)
	}

	transposed, err := songtools.TransposeSong(song, interval, noteNames)
	if err != nil {
		return fmt.Errorf("unable to transpose from %q to %q: %v", fromKey, toKey, err)
	}

	transposed.Key = toKey

	out := os.Stdout
	defer out.Close()

	if cmd.Out != "" {
		out, err = os.OpenFile(cmd.Out, os.O_CREATE, 0666)
		if err != nil {
			return fmt.Errorf("unable to open %q: %v", file, err)
		}
	}

	return format.Writer.Write(out, transposed)
}
