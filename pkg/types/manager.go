package types

import (
	"io"

	"github.com/gookit/color"
)

type ManagerInterface interface {
	ManagerLogInterface
	Stream() io.Writer

	SetStyle(style color.Style)
	Style() color.Style

	IsAccepted() bool
}

type ManagerLogInterface interface {
	LogLn(a ...interface{})
	LogF(format string, a ...interface{})
	LogLnDetails(a ...interface{})
	LogFDetails(format string, a ...interface{})
	LogLnHighlight(a ...interface{})
	LogFHighlight(format string, a ...interface{})
	LogLnWithCustomStyle(style color.Style, a ...interface{})
	LogFWithCustomStyle(style color.Style, format string, a ...interface{})
	LogOptionalLn()

	LogBlock(format string, a ...interface{}) LogBlockInterface
	LogProcessInline(format string, a ...interface{}) LogProcessInlineInterface
	LogProcess(format string, a ...interface{}) LogProcessInterface
}
