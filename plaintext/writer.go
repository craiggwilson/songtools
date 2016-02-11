package plaintext

import (
	"fmt"
	"io"

	"github.com/craiggwilson/songtools"
)

func WriteSong(w io.Writer, song *songtools.Song) error {

	for k, v := range song.Attributes {
		_, err := io.WriteString(w, fmt.Sprintf("#%v=%v\n", k, v))
		if err != nil {
			return fmt.Errorf("unable to write directive(%v): %v", k, err)
		}
	}

	for _, section := range song.Sections {
		if section.Kind != "" {
			_, err := io.WriteString(w, fmt.Sprintf("[%v]\n", section.Kind))
			if err != nil {
				return fmt.Errorf("unable to write section kind(%v): %v", section.Kind, err)
			}
		}

		for _, line := range section.Lines {
            
		}
	}

}
