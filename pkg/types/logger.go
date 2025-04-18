package types

import (
	"io"

	"github.com/gookit/color"

	"github.com/werf/logboek/pkg/level"
)

type LoggerInterface interface {
	ManagerLogInterface

	Error() ManagerInterface
	Warn() ManagerInterface
	Default() ManagerInterface
	Info() ManagerInterface
	Debug() ManagerInterface

	FitText(text string, options FitTextOptions) string
	Colorize(style color.Style, a ...interface{}) string
	ColorizeF(style color.Style, format string, a ...interface{}) string
	ColorizeLn(style color.Style, a ...interface{}) string

	AcceptedLevel() level.Level
	SetAcceptedLevel(lvl level.Level)
	IsAcceptedLevel(lvl level.Level) bool

	Streams() StreamsInterface
	OutStream() io.Writer
	ErrStream() io.Writer
	SetErrorStreamRedirection(lvl ...level.Level)

	NewSubLogger(outStream, errStream io.Writer) LoggerInterface
	GetStreamsSettingsFrom(l LoggerInterface)

	Reset()
	ResetState()
	ResetModes()
}
