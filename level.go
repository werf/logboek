package logboek

import (
	"io"

	"github.com/fatih/color"
)

type Level int

const (
	quiet Level = iota - 3
	Error
	Warn
	Default // 0
	Info
	Debug
)

var (
	levelStreams = map[Level]io.Writer{}
	levelStyles  = map[Level]*Style{}
)

func init() {
	levelStyles[Error] = &Style{Attributes: []color.Attribute{color.FgRed, color.Bold}}
	levelStyles[Warn] = &Style{Attributes: []color.Attribute{color.FgRed, color.Bold}}
}

func (l Level) LogLn(a ...interface{}) {
	logLnCustom(l.Stream(), l, l.Style(), a...)
}

func (l Level) LogF(format string, a ...interface{}) {
	logFCustom(l.Stream(), l, l.Style(), format, a...)
}

func (l Level) LogLnDetails(a ...interface{}) {
	l.LogLnWithCustomStyle(DetailsStyle(), a...)
}

func (l Level) LogFDetails(format string, a ...interface{}) {
	l.LogFWithCustomStyle(DetailsStyle(), format, a...)
}

func (l Level) LogLnHighlight(a ...interface{}) {
	l.LogLnWithCustomStyle(HighlightStyle(), a...)
}

func (l Level) LogFHighlight(format string, a ...interface{}) {
	l.LogFWithCustomStyle(HighlightStyle(), format, a...)
}

func (l Level) LogLnWithCustomStyle(style *Style, a ...interface{}) {
	logLnCustom(l.Stream(), l, style, a...)
}

func (l Level) LogFWithCustomStyle(style *Style, format string, a ...interface{}) {
	logFCustom(l.Stream(), l, style, format, a...)
}

func (l Level) LogBlock(blockMessage string, options LevelLogBlockOptions, processFunc func() error) error {
	return LogBlock(
		blockMessage,
		LogBlockOptions{
			LevelLogBlockOptions: options,
			Level:                l,
		},
		processFunc,
	)
}

func (l Level) LogProcessInline(processMessage string, options LevelLogProcessInlineOptions, processFunc func() error) error {
	return LogProcessInline(
		processMessage,
		LogProcessInlineOptions{
			LevelLogProcessInlineOptions: options,
			Level:                        l,
		},
		processFunc,
	)
}

func (l Level) LogProcess(processMessage string, options LevelLogProcessOptions, processFunc func() error) error {
	return LogProcess(
		processMessage,
		LogProcessOptions{
			LevelLogProcessOptions: options,
			Level:                  l,
		},
		processFunc,
	)
}

func (l Level) LogProcessStart(processMessage string, options LevelLogProcessStartOptions) {
	LogProcessStart(
		processMessage,
		LogProcessStartOptions{
			LevelLogProcessStartOptions: options,
			Level:                       l,
		},
	)
}

func (l Level) LogProcessStepEnd(processMessage string, options LevelLogProcessStepEndOptions) {
	LogProcessStepEnd(
		processMessage,
		LogProcessStepEndOptions{
			LevelLogProcessStepEndOptions: options,
			Level:                         l,
		},
	)
}

func (l Level) LogProcessEnd(options LevelLogProcessEndOptions) {
	LogProcessEnd(
		LogProcessEndOptions{
			LevelLogProcessEndOptions: options,
			Level:                     l,
		},
	)
}

func (l Level) LogProcessFail(options LevelLogProcessFailOptions) {
	LogProcessFail(
		LogProcessFailOptions{
			LevelLogProcessFailOptions: options,
			Level:                      l,
		},
	)
}

func (l Level) Stream() io.Writer {
	stream, ok := levelStreams[l]
	if ok && stream != nil {
		return stream
	} else if l == Error || l == Warn {
		return errStream
	} else {
		return outStream
	}
}

func (l Level) SetStream(stream io.Writer) {
	levelStreams[l] = stream
}

func (l Level) ResetStream() {
	levelStreams[l] = nil
}

func (l Level) Style() *Style {
	return levelStyles[l]
}

func (l Level) SetStyle(style *Style) {
	levelStyles[l] = style
}

func (l Level) IsAccepted() bool {
	return l <= level
}

func (l Level) String() string {
	switch l {
	case quiet:
		return "quiet"
	case Error:
		return "error"
	case Warn:
		return "warn"
	case Default:
		return "default"
	case Info:
		return "info"
	case Debug:
		return "debug"
	default:
		panic("runtime error")
	}
}

func WithLevel(level Level, f func()) {
	savedLevel := level
	savedStreamsLogLevel := streamsLogLevel
	SetLevel(level)
	f()
	level = savedLevel
	streamsLogLevel = savedStreamsLogLevel
}
