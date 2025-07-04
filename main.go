package logboek

import (
	"io"
	"os"

	"golang.org/x/net/context"

	"github.com/gookit/color"

	"github.com/werf/logboek/internal/logger"
	"github.com/werf/logboek/pkg/level"
	"github.com/werf/logboek/pkg/types"
)

var defaultLogger types.LoggerInterface

const ctxLoggerKey = "logboek_logger"

func init() {
	defaultLogger = NewLogger(os.Stdout, os.Stderr)

	defaultErrorAndWarnStyle := color.Style{color.FgRed, color.Bold}
	defaultLogger.Error().SetStyle(defaultErrorAndWarnStyle)
	defaultLogger.Warn().SetStyle(defaultErrorAndWarnStyle)
}

func NewLogger(outStream, errStream io.Writer) types.LoggerInterface {
	return logger.NewLogger(outStream, errStream)
}

func DefaultLogger() types.LoggerInterface {
	return defaultLogger
}

func Error() types.ManagerInterface {
	return defaultLogger.Error()
}

func Warn() types.ManagerInterface {
	return defaultLogger.Warn()
}

func Default() types.ManagerInterface {
	return defaultLogger.Default()
}

func Info() types.ManagerInterface {
	return defaultLogger.Info()
}

func Debug() types.ManagerInterface {
	return defaultLogger.Debug()
}

func LogBlock(headerOrFormat string, a ...interface{}) types.LogBlockInterface {
	return Default().LogBlock(headerOrFormat, a...)
}

func LogProcessInline(headerOrFormat string, a ...interface{}) types.LogProcessInlineInterface {
	return Default().LogProcessInline(headerOrFormat, a...)
}

func LogProcess(headerOrFormat string, a ...interface{}) types.LogProcessInterface {
	return Default().LogProcess(headerOrFormat, a...)
}

func Log(a ...interface{}) {
	Default().Log(a...)
}

func LogLn(a ...interface{}) {
	Default().LogLn(a...)
}

func LogF(format string, a ...interface{}) {
	Default().LogF(format, a...)
}

func LogDetails(a ...interface{}) {
	Default().LogDetails(a...)
}

func LogLnDetails(a ...interface{}) {
	Default().LogLnDetails(a...)
}

func LogFDetails(format string, a ...interface{}) {
	Default().LogFDetails(format, a...)
}

func LogHighlight(a ...interface{}) {
	Default().LogHighlight(a...)
}

func LogLnHighlight(a ...interface{}) {
	Default().LogLnHighlight(a...)
}

func LogFHighlight(format string, a ...interface{}) {
	Default().LogFHighlight(format, a...)
}

func LogWithCustomStyle(style color.Style, a ...interface{}) {
	Default().LogWithCustomStyle(style, a...)
}

func LogLnWithCustomStyle(style color.Style, a ...interface{}) {
	Default().LogLnWithCustomStyle(style, a...)
}

func LogFWithCustomStyle(style color.Style, format string, a ...interface{}) {
	Default().LogFWithCustomStyle(style, format, a...)
}

func LogOptionalLn() {
	Default().LogOptionalLn()
}

func AcceptedLevel() level.Level {
	return defaultLogger.AcceptedLevel()
}

func SetAcceptedLevel(lvl level.Level) {
	defaultLogger.SetAcceptedLevel(lvl)
}

func IsAcceptedLevel(lvl level.Level) bool {
	return defaultLogger.IsAcceptedLevel(lvl)
}

func Streams() types.StreamsInterface {
	return defaultLogger.Streams()
}

func FitText(text string, options types.FitTextOptions) string {
	return defaultLogger.FitText(text, options)
}

func Colorize(style color.Style, a ...interface{}) string {
	return defaultLogger.Colorize(style, a...)
}

func ColorizeLn(style color.Style, a ...interface{}) string {
	return defaultLogger.ColorizeLn(style, a...)
}

func ColorizeF(style color.Style, format string, a ...interface{}) string {
	return defaultLogger.ColorizeF(style, format, a...)
}

func OutStream() io.Writer {
	return defaultLogger.OutStream()
}

func ErrStream() io.Writer {
	return defaultLogger.ErrStream()
}

func NewSubLogger(outStream, errStream io.Writer) types.LoggerInterface {
	return defaultLogger.NewSubLogger(outStream, errStream)
}

func NewContext(ctx context.Context, logger types.LoggerInterface) context.Context {
	return context.WithValue(ctx, ctxLoggerKey, logger)
}

func Context(ctx context.Context) types.LoggerInterface {
	if ctx == nil || ctx == context.Background() {
		return DefaultLogger()
	}

	if ctxValue := ctx.Value(ctxLoggerKey); ctxValue != nil {
		if lgr, ok := ctxValue.(types.LoggerInterface); ok {
			return lgr
		}
	}

	return DefaultLogger()
}

func MustContext(ctx context.Context) types.LoggerInterface {
	if ctx == nil || ctx == context.Background() {
		panic("context is not bound with logboek logger")
	}

	ctxValue := ctx.Value(ctxLoggerKey)
	if ctxValue == nil {
		panic("context is not bound with logboek logger")
	}

	return ctxValue.(types.LoggerInterface)
}

func Reset() {
	defaultLogger.Reset()
}

func ResetState() {
	defaultLogger.ResetState()
}

func ResetModes() {
	defaultLogger.ResetModes()
}
