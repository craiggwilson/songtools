package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/songtools/songtools"
	"github.com/songtools/songtools/format"
	_ "github.com/songtools/songtools/format/chordsOverLyrics" // formats are registered in the init functions.
	_ "github.com/songtools/songtools/format/html"             // formats are registered in the init functions.
)

var cli = flags.NewNamedParser("songtool", flags.Default)

type option struct {
	CurrentFormat string `long:"currentFormat" description:"Specifies the format of the song. By default, will be discovered from the song itself."`
	CurrentKey    string `long:"currentKey" description:"The current key of the song. By default, will be discovered from the song itself."`
	ToFormat      string `short:"f" long:"format" description:"The desired format of the song."`
	ToKey         string `short:"k" long:"key" description:"The desired key of the song."`
	Out           string `short:"o" long:"out" optional:"true" optional-value:"<unset>" description:"The file to write the transposed song. If left unspecified, stdout will be used. When specified without an argument, the input file will be overwritten."`
}

func main() {

	var opt option
	cli.AddGroup("Global", "Global options", &opt)

	args, err := cli.Parse()
	if err != nil {
		os.Exit(1)
	}

	if err = opt.execute(args); err != nil {
		println(err.Error())
		os.Exit(2)
	}
}

func (cmd *option) execute(args []string) error {

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
	}

	inBytes, err := ioutil.ReadAll(in)
	if err != nil {
		return fmt.Errorf("unable to read: %v", err)
	}

	input := bytes.NewBuffer(inBytes)

	readFormat, err := findReadFormat(cmd.CurrentFormat, file, input)
	if err != nil {
		return fmt.Errorf("unable to find input format for %q: %v", file, err)
	}

	writeFormat := readFormat
	if cmd.ToFormat != "" {
		var ok bool
		if writeFormat, ok = format.ByName(cmd.ToFormat); !ok {
			return fmt.Errorf("unable to find output format %q", cmd.ToFormat)
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

	if cmd.ToKey != "" {
		fromKey := songtools.Key(cmd.CurrentKey)

		if fromKey == "" && song.Key == "" {
			return fmt.Errorf("unable to get current key")
		} else if fromKey == "" {
			fromKey = song.Key
		}

		toKey := songtools.Key(cmd.ToKey)

		noteNames, interval, err := songtools.NoteNamesAndIntervalFromKeyToKey(fromKey, toKey)
		if err != nil {
			return fmt.Errorf("unable to get note names and interval: %v", err)
		}

		song, err = songtools.TransposeSong(song, interval, noteNames)
		if err != nil {
			return fmt.Errorf("unable to transpose from %q to %q: %v", fromKey, toKey, err)
		}

		song.Key = toKey
	}

	out := os.Stdout
	defer out.Close()

	if cmd.Out == "<unset>" {
		name := song.Title
		if name == "" && file == "" {
			return fmt.Errorf("'out' was specified, but the song does not have a title and the input was not a file")
		} else if name == "" {
			_, file := filepath.Split(file)
			ext := filepath.Ext(file)
			name = strings.TrimSuffix(file, ext)
		}

		if len(writeFormat.Extensions) > 0 {
			name += "." + writeFormat.Extensions[0]
		}

		cmd.Out = name
	}

	if cmd.Out != "" {
		out, err = os.OpenFile(cmd.Out, os.O_CREATE, 0666)
		if err != nil {
			return fmt.Errorf("unable to open %q: %v", cmd.Out, err)
		}
	}

	return writeFormat.Writer.Write(out, song)
}

func findReadFormat(name, path string, buffer *bytes.Buffer) (*format.Format, error) {
	return findFormat(name, path, buffer, func(f *format.Format) bool {
		return f.CanRead()
	})
}

func findFormat(name, path string, buffer *bytes.Buffer, filter func(*format.Format) bool) (*format.Format, error) {
	formats := format.Formats{}

	if name != "" {
		f, ok := format.ByName(name)
		if !ok {
			return nil, fmt.Errorf("unable to find format %q", name)
		}

		formats = append(formats, f)
	} else if path != "" {
		ext := filepath.Ext(path)
		if ext != "" {
			formats = format.RegisteredFormats()
			formats = formats.Filter(func(f *format.Format) bool {
				for _, e := range f.Extensions {
					if strings.ToLower(e) == ext {
						return true
					}
				}

				return false
			})
		}
	}

	formats = formats.Filter(filter)

	if len(formats) > 0 {
		return formats[0], nil
	}

	return nil, fmt.Errorf("unable to find format")
}

// SetSongTitleIfNecessary sets the song title to the path if the song doesn't already have a title.
func setSongTitleIfNecessary(path string, s *songtools.Song) {
	if path != "" && s.Title == "" {
		_, filename := filepath.Split(path)
		ext := filepath.Ext(path)
		if ext != "" && strings.HasSuffix(filename, ext) {
			filename = strings.TrimSuffix(filename, ext)
		}
		s.Title = filename
	}
}
