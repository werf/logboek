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
