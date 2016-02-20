package plain_test

import (
	"fmt"
	"strings"
	"testing"

	"bytes"

	"github.com/kr/pretty"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/songtools/songtools"
)

const (
	testSongText = `#Title=You've Got a Friend
#Author=Carol King

            Em       B7               Em        B7     Em
When you're down and troubled and you need some lovin' care
    Am       D                G
And nothing, nothing is going right
F#m                 B7              Em     B7      Em
Close your eyes and think of me and soon I will be there
   Am         D                 Am
To brighen up even your darkest nights

[Chorus]
         G                         C 
You just call out my name, and you know, wherever I am
       G                  Am    D7
I come running to see you again
G                       G7    C                   Am
Winter spring summer or fall, all you go to do is call
         C               Am
And I'll be there, yes I will
A            G
You've got a friend


If the sky above you grows dark and full of clouds,
and that old north wind begins to blow
Keep your head together, and call my name out  loud
Soon you'll hear me knocking at your door

{Chorus}

[Bridge]
F                              C
Now ain't it good to know that you've got a friend
     G                G7
When people can be so cold
        C                 Fm             Em                    Am7
They'll hurt you, yes and desert you and take your soul if you let them
    A7            D
But don't you let them

{Chorus}

A            G       D            Em
You've got a friend. You've got a friend
`
)

func TestParse(t *testing.T) {

	Convey("Subject: Parse", t, func() {
		set, err := plaintext.ParseSongSet(strings.NewReader(testSongText))
		So(err, ShouldBeNil)

		So(set.Songs, ShouldHaveLength, 1)

		fmt.Printf("\n%# v", pretty.Formatter(set))

		w := bytes.NewBufferString("")
		plaintext.WriteSongSet(w, set)

		println(w.String())

		noteNames, interval, _ := songtools.NoteNamesAndIntervalFromKeyToKey("Em", "F#m")
		transposed, _ := songtools.TransposeSongSet(set, interval, noteNames)

		w = bytes.NewBufferString("")
		plaintext.WriteSongSet(w, transposed)

		println(w.String())
	})

}
