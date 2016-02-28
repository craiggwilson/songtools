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
        * {
            font-family: monospace;
        }
        
        .song {
            margin: 24px;
        }
        
        h1 {
            font-size: 26px;
            font-weight: bold;
            border-bottom-width: 1px;
            border-bottom-style: dotted;
            margin-bottom: 2px;
            padding-bottom: 10px;
        }
        
        h2 {
            font-size: 14px;
            font-weight: bold;
            margin: 0;
            padding: 0;
            text-decoration: underline;
        }
        
        ul {
            padding: 0;
            margin: 0;
        }
        
        ul li {
            display: inline;
            font-weight: bold;
        }
        
        ul.song-authors li {
            font-weight: normal;
        }
        
        ul li:nth-of-type(1n+2)::before {
            content: "/ ";
        }
        
        ul.song-authors::before {
            content: "Author(s): "
        }
        
        section {
            margin: 16px 0;	
        }
        
        .song-key::before {
            content: "Key: "
        }
        
        .song-comment {
            font-style: italic;
            font-weight: bold;
            margin: 0 0;	
        }
        
        .song-line-group {
            white-space: pre;
        }
        
        .song-chord-line {
            font-weight: bold;
        }
    </style>
</head>
<body class='song'>
    <header>
    {{if .Title}}
        <h1 class='song-title'>{{.Title}}</h1>
    {{end}}
    {{if .Subtitles}}
        <ul class='song-subtitles'>
        {{range .Subtitles}}
            <li>{{.}}</li>
        {{end}}
        </ul>
    {{end}}
    {{if .Authors}}
        <ul class='song-authors'>
        {{range .Authors}}
            <li>{{.}}</li>
        {{end}}
        </ul>
    {{end}}
    {{if .Key}}
        <div class='song-key'>{{.Key}}</div>
    {{end}}
    </header>
    <div class='song-content'>
        {{Content .}}
    </div>
</body>
</html>`
)

// WriteSong writes a single song to the writer.
func WriteSong(w io.Writer, s *songtools.Song) error {
	funcs := make(map[string]interface{})
	funcs["Content"] = writeContent
	t := template.Must(template.New("song").Funcs(funcs).Parse(songTemplate))

	return t.ExecuteTemplate(w, "song", s)
}

func writeContent(s *songtools.Song) string {
	buf := ""
	for _, n := range s.Nodes {
		buf += writeSongNode(n)
	}
	return buf
}

func writeSongNode(n songtools.SongNode) string {
	buf := ""
	switch typedN := n.(type) {
	case *songtools.Comment:
		buf += writeComment(typedN)
	case *songtools.Section:

		anyChords := len(typedN.Chords()) > 0
		for i, sn := range typedN.Nodes {
			if i == 0 {
				if typedN.Kind != "" {
					// only use the first word for the css class
					kind := strings.ToLower(strings.Split(string(typedN.Kind), " ")[0])
					buf += "<section class='song-" + kind + "'>"

					kind = string(typedN.Kind)
					cont := false
					if c, ok := sn.(*songtools.Comment); ok {
						kind += " " + c.Text
						cont = true
					}

					buf += "<h2 class='song-section-kind'>" + kind + "</h2>"
					if cont {
						continue
					}
				} else {
					buf += "<section class='song-verse'>"
				}
			}

			buf += writeSectionNode(sn, anyChords)
		}

		buf += "</section>"
	}

	return buf
}

func writeSectionNode(n songtools.SectionNode, blankLineForNoChords bool) string {
	switch typedN := n.(type) {
	case *songtools.Comment:
		return writeComment(typedN)
	case *songtools.Line:
		return writeLine(typedN, blankLineForNoChords)
	default:
		return ""
	}
}

func writeComment(c *songtools.Comment) string {
	buf := ""
	if !c.Hidden {
		buf += "<div class='song-comment'>"
		buf += c.Text
		buf += "</div>"
	}
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
