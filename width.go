package logboek

import (
	"fmt"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

const DefaultWidth = 140

var width = DefaultWidth

func initWidth() error {
	if terminal.IsTerminal(int(os.Stdout.Fd())) {
		w, _, err := terminal.GetSize(int(os.Stdout.Fd()))
		if err != nil {
			return fmt.Errorf("get terminal size failed: %s", err)
		}

		if w == 0 {
			w = DefaultWidth
		}

		SetWidth(w)
	}

	return nil
}

func SetWidth(value int) {
	width = value
}
