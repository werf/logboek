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
	Log(a ...interface{})
	LogLn(a ...interface{})
	LogF(format string, a ...interface{})
	LogDetails(a ...interface{})
	LogLnDetails(a ...interface{})
	LogFDetails(format string, a ...interface{})
	LogHighlight(a ...interface{})
	LogLnHighlight(a ...interface{})
	LogFHighlight(format string, a ...interface{})
	LogWithCustomStyle(style color.Style, a ...interface{})
	LogLnWithCustomStyle(style color.Style, a ...interface{})
	LogFWithCustomStyle(style color.Style, format string, a ...interface{})
	LogOptionalLn()

	LogBlock(headerOrFormat string, a ...interface{}) LogBlockInterface
	LogProcessInline(headerOrFormat string, a ...interface{}) LogProcessInlineInterface
	LogProcess(headerOrFormat string, a ...interface{}) LogProcessInterface
}
