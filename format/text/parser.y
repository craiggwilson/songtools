%{
package text

import (
	"bytes"
)

%}

%union {
  empty       struct{}
  runes       []rune
  song        ast.Song
  songComponent ast.SongComponent
  songComponents   []ast.SongComponent
}

%token LEX_ERROR
%token <empty> LBRACKET RBRACKET HASH EOL
%token <bytes> TEXT

%start root

%type <song> root
%type <songComponent> song_component
%type <songComponents> song_components

%%

root:
    song_components
    {
        $$ = ast.Song{
            Components: $1,
        }
    };

song_components:
    song_component
    {
        $$ = []ast.SongComponent{$1}
    }
    | song_components song_component
    {
        $$ = append($1, $2)
    };

song_component:
    LBRACKET RBRACKET
    {
        $$ = ast.Comment{}
    };

