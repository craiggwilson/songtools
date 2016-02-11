package plaintext_test

import (
	"strings"
	"testing"

	"github.com/craiggwilson/songtools/plaintext"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	testSongText = `#Title=Test

[Verse 1]
Am      G  C
Lyrics, Oh Lyrics`
)

func TestNext(t *testing.T) {

	Convey("Subject: Next", t, func() {

		s, err := plaintext.NewScanner(strings.NewReader(testSongText))
		So(err, ShouldBeNil)

		expected := []tokenString{
			{plaintext.DirectiveToken, "Title=Test"},
			{plaintext.NewLineToken, ""},
			{plaintext.NewLineToken, ""},
			{plaintext.SectionToken, "Verse 1"},
			{plaintext.NewLineToken, ""},
			{plaintext.TextToken, "Am      G  C"},
			{plaintext.NewLineToken, ""},
			{plaintext.TextToken, "Lyrics, Oh Lyrics"},
			{plaintext.EofToken, ""},
		}

		for _, e := range expected {
			token, text, err := s.Next()
			So(err, ShouldBeNil)
			So(token, ShouldEqual, e.token)
			So(text, ShouldEqual, e.text)
		}
	})
}

type tokenString struct {
	token plaintext.Token
	text  string
}
