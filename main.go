package logboek

import (
	"github.com/fatih/color"
)

func Init() error {
	return initWidth()
}

func EnableLogColor() {
	color.NoColor = false
}

func DisableLogColor() {
	color.NoColor = true
}

func DisablePrettyLog() {
	RawStreamsOutputModeOn()
	disableLogProcessBorder()
}
