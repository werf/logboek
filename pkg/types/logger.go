package types

import (
	"io"

	"github.com/werf/logboek/pkg/level"
	"github.com/werf/logboek/pkg/style"
)

type LoggerInterface interface {
	ManagerLogInterface

	Error() ManagerInterface
	Warn() ManagerInterface
	Default() ManagerInterface
	Info() ManagerInterface
	Debug() ManagerInterface

	FitText(text string, options FitTextOptions) string
	Colorize(style *style.Style, format string, a ...interface{}) string

	AcceptedLevel() level.Level
	SetAcceptedLevel(lvl level.Level)
	IsAcceptedLevel(lvl level.Level) bool

	Streams() StreamsInterface
	ProxyOutStream() io.Writer
	ProxyErrStream() io.Writer

	NewSubLogger(outStream, errStream io.Writer) LoggerInterface
	GetStreamsSettingsFrom(l LoggerInterface)

	Reset()
}
