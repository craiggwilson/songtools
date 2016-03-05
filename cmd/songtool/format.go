package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/songtools/songtools/format"
)

// FormatCommand defines the options required for formatting a song.
type FormatCommand struct {
	From string `short:"f" long:"from" description:"Specifies the input format of the song. By default, will be discovered from the song itself."`
	To   string `short:"t" long:"to" description:"Specified the output format of the song. When not specified, will use the 'from' format."`
	Out  string `short:"o" long:"out" optional:"true" optional-value:"<input>" description:"The file to write the transposed song. If left unspecified, stdout will be used. When specified without an argument, the input file will be overwritten."`
}

var formatCommand FormatCommand

func init() {

	cmd, _ := cli.AddCommand(
		"format",
		"Formats a song",
		"Formats a song. Can be used to format a song using the same format, or change the format of a song.",
		&formatCommand)

	fromArg := cmd.FindOptionByLongName("from")
	fromArg.Choices = format.FilteredNames(true, false)

	toArg := cmd.FindOptionByLongName("to")
	toArg.Choices = format.FilteredNames(false, true)
}

// Execute executes the FormatCommand.
func (cmd *FormatCommand) Execute(args []string) error {
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

	readFormat, err := findReadFormat(cmd.From, file, input)
	if err != nil {
		return fmt.Errorf("unable to find input format for %q: %v", file, err)
	}

	writeFormat := readFormat
	if cmd.To != "" {
		var ok bool
		if writeFormat, ok = format.ByName(cmd.To); !ok {
			return fmt.Errorf("unable to find output format %q", cmd.To)
		}
	}

	if writeFormat.Writer == nil {
		return fmt.Errorf("the input format %q is unable to be used for writing", readFormat.Name)
	}

	song, err := readFormat.Reader.Read(input)
	if err != nil {
		return fmt.Errorf("unable to parse %q: %v", file, err)
	}

	setSongTitleIfNecessary(file, song)

	out := os.Stdout
	defer out.Close()

	if cmd.Out != "" {
		out, err = os.OpenFile(cmd.Out, os.O_CREATE, 0666)
		if err != nil {
			return fmt.Errorf("unable to open %q: %v", cmd.Out, err)
		}
	}

	return writeFormat.Writer.Write(out, song)
}
