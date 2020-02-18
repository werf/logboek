package logboek

import (
	"github.com/fatih/color"
)

var (
	// logboek logs entries with that severity or above it
	level = Default

	// severity for OutF, ErrF, w.Write can be managed only with this variable
	streamsLogLevel = level
)

func Init() error {
	return initWidth()
}

func SetLevel(l Level) {
	level = l
	streamsLogLevel = l
}

func SetQuietLevel() {
	SetLevel(quiet)
}

func EnableLogColor() {
	color.NoColor = false
}

func DisableLogColor() {
	color.NoColor = true
}

func WithFitMode(value bool, decoratedFunc func() error) error {
	oldFitModeState := isFitModeOn
	isFitModeOn = value
	err := decoratedFunc()
	isFitModeOn = oldFitModeState

	return err
}

func EnableFitMode() {
	isFitModeOn = true
}

func DisableFitMode() {
	isFitModeOn = false
}

func DisablePrettyLog() {
	RawStreamsOutputModeOn()
	disableLogProcessBorder()
	DisableFitMode()
}
