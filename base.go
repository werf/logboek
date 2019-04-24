package logboek

import "fmt"

func LogLn(a ...interface{}) {
	LogF("%s", fmt.Sprintln(a...))
}

func LogF(format string, a ...interface{}) {
	colorizeFormatAndLogF(outStream, ColorizeBase, format, a...)
}

func LogHighlightLn(a ...interface{}) {
	LogHighlightF("%s", fmt.Sprintln(a...))
}

func LogHighlightF(format string, a ...interface{}) {
	colorizeFormatAndLogF(outStream, ColorizeHighlight, format, a...)
}

func LogInfoLn(a ...interface{}) {
	LogInfoF("%s", fmt.Sprintln(a...))
}

func LogInfoF(format string, a ...interface{}) {
	colorizeFormatAndLogF(outStream, ColorizeInfo, format, a...)
}

func LogErrorLn(a ...interface{}) {
	LogErrorF("%s", fmt.Sprintln(a...))
}

func LogErrorF(format string, a ...interface{}) {
	colorizeFormatAndLogF(errStream, ColorizeWarning, format, a...)
}
