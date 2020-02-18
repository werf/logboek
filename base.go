package logboek

import (
	"fmt"
	"io"
)

func LogLn(a ...interface{}) {
	Default.LogLn(a...)
}

func LogF(format string, a ...interface{}) {
	Default.LogF(format, a...)
}

func LogInfoLn(a ...interface{}) {
	Info.LogLn(a...)
}

func LogInfoF(format string, a ...interface{}) {
	Info.LogF(format, a...)
}

func LogErrorLn(a ...interface{}) {
	Error.LogLn(a...)
}

func LogErrorF(format string, a ...interface{}) {
	Error.LogF(format, a...)
}

func LogWarnLn(a ...interface{}) {
	Warn.LogLn(a...)
}

func LogWarnF(format string, a ...interface{}) {
	Warn.LogF(format, a...)
}

type LogOptions struct {
	Level Level
	Style *Style
}

type LogLnOptions struct {
	LogOptions
}

func LogLnWithOptions(options LogLnOptions, a ...interface{}) {
	style := options.Style
	if style == nil {
		style = options.Level.Style()
	}

	logLnCustom(options.Level.Stream(), options.Level, style, a...)
}

type LogFOptions struct {
	LogOptions
}

func LogFWithOptions(options LogFOptions, format string, a ...interface{}) {
	style := options.Style
	if style == nil {
		style = options.Level.Style()
	}

	logFCustom(options.Level.Stream(), options.Level, style, format, a...)
}

func logLnCustom(stream io.Writer, logLevel Level, style *Style, a ...interface{}) {
	logFCustom(stream, logLevel, style, "%s", fmt.Sprintln(a...))
}

func logFCustom(stream io.Writer, logLevel Level, style *Style, format string, a ...interface{}) {
	if !logLevel.IsAccepted() {
		return
	}

	formatAndLogF(stream, style, format, a...)
}
