package logboek

import (
	"io"
	"os"

	"golang.org/x/net/context"

	"github.com/fatih/color"

	"github.com/werf/logboek/internal/logger"
	"github.com/werf/logboek/pkg/level"
	"github.com/werf/logboek/pkg/style"
	"github.com/werf/logboek/pkg/types"
)

var defaultLogger types.LoggerInterface

const ctxLoggerKey = "logboek_logger"

func init() {
	defaultLogger = NewLogger(os.Stdout, os.Stderr)

	defaultErrorAndWarnStyle := &style.Style{Attributes: []color.Attribute{color.FgRed, color.Bold}}
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

func LogBlock(format string, a ...interface{}) types.LogBlockInterface {
	return Default().LogBlock(format, a...)
}

func LogProcessInline(format string, a ...interface{}) types.LogProcessInlineInterface {
	return Default().LogProcessInline(format, a...)
}

func LogProcess(format string, a ...interface{}) types.LogProcessInterface {
	return Default().LogProcess(format, a...)
}

func LogLn(a ...interface{}) {
	Default().LogLn(a...)
}

func LogF(format string, a ...interface{}) {
	Default().LogF(format, a...)
}

func LogLnDetails(a ...interface{}) {
	Default().LogLnDetails(a...)
}

func LogFDetails(format string, a ...interface{}) {
	Default().LogFDetails(format, a...)
}

func LogLnHighlight(a ...interface{}) {
	Default().LogLnHighlight(a...)
}

func LogFHighlight(format string, a ...interface{}) {
	Default().LogFHighlight(format, a...)
}

func LogLnWithCustomStyle(style *style.Style, a ...interface{}) {
	Default().LogLnWithCustomStyle(style, a...)
}

func LogFWithCustomStyle(style *style.Style, format string, a ...interface{}) {
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

func ProxyOutStream() io.Writer {
	return defaultLogger.ProxyOutStream()
}

func ProxyErrStream() io.Writer {
	return defaultLogger.ProxyErrStream()
}

func NewSubLogger(outStream, errStream io.Writer) types.LoggerInterface {
	return defaultLogger.NewSubLogger(outStream, errStream)
}

func NewContext(ctx context.Context, logger types.LoggerInterface) context.Context {
	return context.WithValue(ctx, ctxLoggerKey, logger)
}

func Context(ctx context.Context) types.LoggerInterface {
	ctxValue := ctx.Value(ctxLoggerKey)
	if ctxValue == nil {
		panic("context is not bound with logboek logger")
	}

	return ctxValue.(types.LoggerInterface)
}
