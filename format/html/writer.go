package html

import (
	"io"
	"strings"
	"text/template"

	"github.com/songtools/songtools"
)

const (
	songTemplate = `<!DOCTYPE html>
<html>
<head>
    <title>{{.Title}}</title>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, intial-scale=1.0">
    <style>
        .song {
            margin: 2em;
        }
        
        .song-title {
            font-size: 1.5em;
        }
        
        .song-section {
            margin: 2em 0;	
        }
        
        .song-section-kind {
            font-size: 1.2em;
            font-weight: bold;
        }
        
        .song-chorus {
            padding-left: 1em;
            font-style: italic;
        }
        
        .song-comment {
            font-style: italic;
            font-weight: bold;
            margin: 0 0;	
        }
        
        .song-chord-line, .song-lyric-line {
            white-space: pre;
            font-family: monospace;
        }
    </style>
</head>
<body>
    <div class='song'>
        {{Song .}}
    </div>
</body>
</html>`
)

type songToWrite struct {
	Title string
	Song  *songtools.Song
}

// WriteSong writes a single song to the writer.
func WriteSong(w io.Writer, s *songtools.Song) error {
	funcs := make(map[string]interface{})
	funcs["Song"] = writeSong
	t := template.Must(template.New("song").Funcs(funcs).Parse(songTemplate))

	title, _ := s.Title()
	sToW := songToWrite{title, s}

	return t.ExecuteTemplate(w, "song", sToW)
}

func writeSong(s songToWrite) string {
	buf := ""
	for _, n := range s.Song.Nodes {
		buf += writeSongNode(n)
	}

	return buf
}

func writeSongNode(n songtools.SongNode) string {
	buf := ""
	switch typedN := n.(type) {
	case *songtools.Comment:
		buf += writeComment(typedN)
	case *songtools.Directive:
		buf += writeDirective(typedN)
	case *songtools.Section:
		if typedN.Kind != "" {
			// only use the first word for the css class
			kind := strings.ToLower(strings.Split(string(typedN.Kind), " ")[0])
			buf += "<div class='song-section song-" + kind + "'>"
			buf += "<div class='song-section-kind'>" + string(typedN.Kind) + "</div>"
		} else {
			buf += "<div class='song-section song-verse'>"
		}

		anyChords := len(typedN.Chords()) > 0

		for _, sn := range typedN.Nodes {
			buf += writeSectionNode(sn, anyChords)
		}

		buf += "</div>"
	default:
		panic("Unknown node")
	}

	return buf
}

func writeSectionNode(n songtools.SectionNode, blankLineForNoChords bool) string {
	switch typedN := n.(type) {
	case *songtools.Comment:
		return writeComment(typedN)
	case *songtools.Directive:
		return writeDirective(typedN)
	case *songtools.Line:
		return writeLine(typedN, blankLineForNoChords)
	default:
		panic("Unknown node")
	}
}

func writeComment(c *songtools.Comment) string {
	buf := "<div class='song-hidden-comment'>"
	buf += c.Text
	buf += "</div>"
	return buf
}

func writeDirective(d *songtools.Directive) string {
	buf := "<div class='song-" + strings.Replace(strings.ToLower(d.Name), " ", "_", -1) + "'>"
	buf += d.Value
	buf += "</div>"
	return buf
}

func writeLine(l *songtools.Line, blankLineForNoChords bool) string {
	buf := "<div class='song-line-group'>"
	if l.Chords != nil {
		buf += "<div class='song-chord-line'>"
		pos := 0
		for i := 0; i < len(l.Chords); i++ {
			buf += strings.Repeat(" ", l.ChordPositions[i]-pos)
			buf += "<span class='song-chord'>"
			buf += l.Chords[i].Name
			buf += "</span>"
			pos = l.ChordPositions[i] + len(l.Chords[i].Name)
		}
		buf += "</div>"
	} else if blankLineForNoChords {
		buf += "<div class='song-chord-line'> </div>"
	}

	buf += "<div class='song-lyric-line'>"
	buf += l.Text
	buf += "</div></div>"
	return buf
}
