package types

import (
	stylePkg "github.com/werf/logboek/pkg/style"
)

type ManagerInterface interface {
	ManagerLogInterface

	SetStyle(style *stylePkg.Style)
	Style() *stylePkg.Style

	IsAccepted() bool
}

type ManagerLogInterface interface {
	LogLn(a ...interface{})
	LogF(format string, a ...interface{})
	LogLnDetails(a ...interface{})
	LogFDetails(format string, a ...interface{})
	LogLnHighlight(a ...interface{})
	LogFHighlight(format string, a ...interface{})
	LogLnWithCustomStyle(style *stylePkg.Style, a ...interface{})
	LogFWithCustomStyle(style *stylePkg.Style, format string, a ...interface{})
	LogOptionalLn()

	LogBlock(format string, a ...interface{}) LogBlockInterface
	LogProcessInline(format string, a ...interface{}) LogProcessInlineInterface
	LogProcess(format string, a ...interface{}) LogProcessInterface
}
