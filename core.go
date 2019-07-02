package logboek

import (
	"fmt"
	"io"
	"strings"
)

func processAndLogLn(w io.Writer, a ...interface{}) {
	processAndLogF(w, fmt.Sprintln(a...))
}

func processAndLogF(w io.Writer, format string, a ...interface{}) {
	_, err := processAndLogFBase(w, format, a...)
	if err != nil {
		panic(err)
	}
}

func processAndLogFBase(w io.Writer, format string, a ...interface{}) (int, error) {
	var msg string
	if len(a) != 0 {
		msg = fmt.Sprintf(format, a...)
	} else {
		msg = format
	}

	var formattedMsg string
	for _, r := range []rune(msg) {
		switch string(r) {
		case "\r", "\n":
			formattedMsg += processNewLineAndRemoveCarriage(string(r))
		default:
			formattedMsg += processDefault()
		}

		formattedMsg += string(r)
	}

	return logFBase(w, "%s", formattedMsg)
}

var (
	isCursorOnNewLine              = true
	isPrevCursorStateOnRemoveCaret = false
)

func processNewLineAndRemoveCarriage(carriage string) string {
	var result string

	if isCursorOnNewLine && !isPrevCursorStateOnRemoveCaret {
		result += processService()
	}

	isPrevCursorStateOnRemoveCaret = carriage == "\r"
	isCursorOnNewLine = true

	return result
}

func processDefault() string {
	var result string

	result += processOptionalLnMode()

	if isCursorOnNewLine {
		result += processService()
		result += strings.Repeat(" ", indentWidth)

		isCursorOnNewLine = false
	}

	isPrevCursorStateOnRemoveCaret = false

	return result
}

func processService() string {
	var result string

	result += formattedPrefix()
	result += formattedProcessBorders()
	result += formattedTag()

	return result
}

func logFBase(w io.Writer, format string, a ...interface{}) (int, error) {
	return fmt.Fprintf(w, format, a...)
}

var prefix string
var prefixColorizeFunc func(...interface{}) string

func SetPrefix(value string, colorizeFunc func(...interface{}) string) {
	prefix = value
	prefixColorizeFunc = colorizeFunc
}

func ResetPrefix() {
	prefix = ""
	prefixColorizeFunc = nil
}

func formattedPrefix() string {
	if prefix == "" {
		return ""
	}

	if prefixColorizeFunc == nil {
		return prefix
	}

	return prefixColorizeFunc(prefix)
}

func prefixWidth() int {
	return len([]rune(prefix))
}

var indentWidth int

func WithIndent(f func() error) error {
	IndentUp()
	err := f()
	IndentDown()

	return err
}

func WithoutIndent(decoratedFunc func() error) error {
	oldIndentWidth := indentWidth
	indentWidth = 0
	err := decoratedFunc()
	indentWidth = oldIndentWidth

	return err
}

func IndentUp() {
	indentWidth += 2
	resetOptionalLnMode()
}

func IndentDown() {
	if indentWidth == 0 {
		return
	}

	indentWidth -= 2
	resetOptionalLnMode()
}

func decorateByWithIndent(decoratedFunc func() error) func() error {
	return func() error {
		return WithIndent(decoratedFunc)
	}
}

var isOptionalLnModeOn bool

func LogOptionalLn() {
	isOptionalLnModeOn = true
}

func resetOptionalLnMode() {
	isOptionalLnModeOn = false
}

func applyOptionalLnMode() {
	_, _ = logFBase(outStream, processOptionalLnMode())
}

func processOptionalLnMode() string {
	var result string

	if isOptionalLnModeOn {
		result += processService()
		result += "\n"

		resetOptionalLnMode()
		isCursorOnNewLine = true
	}

	return result
}

func ContentWidth() int {
	return width - serviceWidth()
}

func serviceWidth() int {
	return prefixWidth() + processBordersBlockWidth() + tagPartWidth + indentWidth
}
