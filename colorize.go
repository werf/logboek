package logboek

import (
	"fmt"
	"io"
	"strings"

	"github.com/fatih/color"
)

var (
	noneFormat []color.Attribute

	baseFormat = noneFormat

	highlightFormat = []color.Attribute{color.Bold}
	infoFormat      = []color.Attribute{color.FgHiBlue}
	warningFormat   = []color.Attribute{color.FgRed, color.Bold}

	successFormat = []color.Attribute{color.FgGreen, color.Bold}
	failFormat    = warningFormat
)

func SetBaseFormat(attributes []color.Attribute) {
	baseFormat = attributes
}

func SetHighlightFormat(attributes []color.Attribute) {
	highlightFormat = attributes
}

func SetInfoFormat(attributes []color.Attribute) {
	infoFormat = attributes
}

func SetWarningFormat(attributes []color.Attribute) {
	warningFormat = attributes
}

func SetFailFormat(attributes []color.Attribute) {
	failFormat = attributes
}

func colorizeFormatAndLogF(w io.Writer, colorizeFunc func(...interface{}) string, format string, args ...interface{}) {
	var msg string
	if len(args) > 0 {
		msg = colorizeBaseF(colorizeFunc, format, args...)
	} else {
		msg = colorizeBaseF(colorizeFunc, "%s", format)
	}

	msg = FitText(msg, FitTextOptions{MarkWrappedLine: true})

	processAndLogF(w, msg)
}

func colorizeBaseF(colorizeFunc func(...interface{}) string, format string, args ...interface{}) string {
	var colorizeLines []string
	lines := strings.Split(fmt.Sprintf(format, args...), "\n")
	for _, line := range lines {
		if line == "" {
			colorizeLines = append(colorizeLines, line)
		} else {
			colorizeLines = append(colorizeLines, colorizeFunc(line))
		}
	}

	return strings.Join(colorizeLines, "\n")
}

func ColorizeNone(a ...interface{}) string {
	return colorize(noneFormat, a...)
}

func ColorizeBase(a ...interface{}) string {
	return colorize(baseFormat, a...)
}

func ColorizeHighlight(a ...interface{}) string {
	return colorize(highlightFormat, a...)
}

func ColorizeInfo(a ...interface{}) string {
	return colorize(infoFormat, a...)
}

func ColorizeWarning(a ...interface{}) string {
	return colorize(warningFormat, a...)
}

func ColorizeSuccess(a ...interface{}) string {
	return colorize(successFormat, a...)
}

func ColorizeFail(a ...interface{}) string {
	return colorize(failFormat, a...)
}

func colorize(attributes []color.Attribute, a ...interface{}) string {
	return color.New(attributes...).Sprint(a...)
}
