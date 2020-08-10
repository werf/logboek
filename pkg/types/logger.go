package types

import (
	"io"

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

	AcceptedLevel() level.Level
	SetAcceptedLevel(lvl level.Level)
	IsAcceptedLevel(lvl level.Level) bool

	Streams() StreamsInterface
	ProxyOutStream() io.Writer
	ProxyErrStream() io.Writer

	NewSubLogger(outStream, errStream io.Writer) LoggerInterface
}
