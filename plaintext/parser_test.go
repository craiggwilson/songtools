package plaintext_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/craiggwilson/songtools/plaintext"
	"github.com/kr/pretty"
	. "github.com/smartystreets/goconvey/convey"
)

func TestParse(t *testing.T) {

	Convey("Subject: Parse", t, func() {
		songs, err := plaintext.Parse(strings.NewReader(testSongText))
		So(err, ShouldBeNil)

		So(songs, ShouldHaveLength, 1)

		fmt.Printf("\n%# v", pretty.Formatter(songs))
	})

}
