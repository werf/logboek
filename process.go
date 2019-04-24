package logboek

import (
	"fmt"
	"strings"
	"time"
)

const (
	logProcessTimeFormat = "%.2f seconds"

	logProcessInlineProcessMsgFormat = "%s ..."
	logStateRightPartsSeparator      = " "
)

var (
	logProcessDownAndRightBorderSign     = "┌"
	logProcessVerticalBorderSign         = "│"
	logProcessVerticalAndRightBorderSign = "├"
	logProcessUpAndRightBorderSign       = "└"

	processesBorderValues             []string
	processesBorderFormattedValues    []string
	processesBorderBetweenIndentWidth = 1
	processesBorderIndentWidth        = 1
)

var (
	activeLogProcesses []*logProcessDescriptor
)

type logProcessDescriptor struct {
	StartedAt time.Time
	Msg       string
}

func disableLogProcessBorder() {
	logProcessDownAndRightBorderSign = ""
	logProcessVerticalBorderSign = "  "
	logProcessVerticalAndRightBorderSign = ""
	logProcessUpAndRightBorderSign = ""

	processesBorderIndentWidth = 0
	processesBorderBetweenIndentWidth = 0
}

type LogProcessInlineOptions struct {
	ColorizeMsgFunc func(...interface{}) string
}

func LogProcessInline(processMessage string, options LogProcessInlineOptions, processFunc func() error) error {
	return logProcessInline(processMessage, options, processFunc)
}

func logProcessInline(processMessage string, options LogProcessInlineOptions, processFunc func() error) error {
	if options.ColorizeMsgFunc == nil {
		options.ColorizeMsgFunc = ColorizeBase
	}

	processMessage = fmt.Sprintf(logProcessInlineProcessMsgFormat, processMessage)
	colorizeFormatAndLogF(outStream, options.ColorizeMsgFunc, "%s", processMessage)

	resultColorize := options.ColorizeMsgFunc
	start := time.Now()

	resultFormat := " (%s)\n"

	var err error
	if err = WithIndent(processFunc); err != nil {
		resultColorize = ColorizeFail
		resultFormat = " (%s) FAILED\n"
	}

	elapsedSeconds := fmt.Sprintf(logProcessTimeFormat, time.Since(start).Seconds())
	colorizeFormatAndLogF(outStream, resultColorize, resultFormat, elapsedSeconds)

	return err
}

func prepareLogProcessMsgLeftPart(leftPart string, colorizeFunc func(...interface{}) string, rightParts ...string) string {
	var result string

	spaceWidth := ContentWidth() - len(strings.Join(rightParts, logStateRightPartsSeparator))
	if spaceWidth > 0 {
		if spaceWidth > len(leftPart) {
			result = leftPart
		} else {
			service := " ..."
			if spaceWidth > len(" ...") {
				result = leftPart[0:spaceWidth-len(service)] + service
			} else {
				result = leftPart[0:spaceWidth]
			}
		}
	} else {
		return ""
	}

	return colorizeFunc(result)
}

type LogProcessStartOptions struct {
	ColorizeMsgFunc func(...interface{}) string
}

type LogProcessEndOptions struct {
	WithoutLogOptionalLn bool
	ColorizeMsgFunc      func(...interface{}) string
}

type LogProcessStepEndOptions struct {
	WithIndent      bool
	InfoSectionFunc func(err error)
	ColorizeMsgFunc func(...interface{}) string
}

type LogProcessOptions struct {
	WithIndent           bool
	WithoutLogOptionalLn bool
	InfoSectionFunc      func(err error)
	ColorizeMsgFunc      func(...interface{}) string
}

func LogProcessStart(processMessage string, options LogProcessStartOptions) {
	logProcessStart(processMessage, options)
}

func LogProcessEnd(options LogProcessEndOptions) {
	logProcessEnd(options)
}

func LogProcessStepEnd(processMessage string, options LogProcessStepEndOptions) {
	logProcessStepEnd(processMessage, options)
}

func LogProcessFail(options LogProcessEndOptions) {
	logProcessFail(options)
}

func LogProcess(processMessage string, options LogProcessOptions, processFunc func() error) error {
	return logProcess(processMessage, options, processFunc)
}

func logProcess(processMessage string, options LogProcessOptions, processFunc func() error) error {
	if options.ColorizeMsgFunc == nil {
		options.ColorizeMsgFunc = ColorizeBase
	}

	logProcessStart(processMessage, LogProcessStartOptions{ColorizeMsgFunc: options.ColorizeMsgFunc})

	bodyFunc := func() error {
		return processFunc()
	}

	if options.WithIndent {
		bodyFunc = decorateByWithIndent(bodyFunc)
	}

	err := bodyFunc()

	resetOptionalLnMode()

	if options.InfoSectionFunc != nil {
		applyInfoLogProcessStep(err, options.InfoSectionFunc, options.WithIndent, options.ColorizeMsgFunc)
	}

	if err != nil {
		logProcessFail(LogProcessEndOptions{WithoutLogOptionalLn: options.WithoutLogOptionalLn, ColorizeMsgFunc: options.ColorizeMsgFunc})
		return err
	}

	logProcessEnd(LogProcessEndOptions{WithoutLogOptionalLn: options.WithoutLogOptionalLn, ColorizeMsgFunc: options.ColorizeMsgFunc})
	return nil
}

func logProcessStart(processMessage string, options LogProcessStartOptions) {
	if options.ColorizeMsgFunc == nil {
		options.ColorizeMsgFunc = ColorizeBase
	}

	applyOptionalLnMode()

	headerFunc := func() error {
		return WithoutIndent(func() error {
			processAndLogLn(outStream, prepareLogProcessMsgLeftPart(processMessage, options.ColorizeMsgFunc))
			return nil
		})
	}

	headerFunc = decorateByWithExtraProcessBorder(logProcessDownAndRightBorderSign, options.ColorizeMsgFunc, headerFunc)

	_ = headerFunc()

	appendProcessBorder(logProcessVerticalBorderSign, options.ColorizeMsgFunc)

	logProcess := &logProcessDescriptor{StartedAt: time.Now(), Msg: processMessage}
	activeLogProcesses = append(activeLogProcesses, logProcess)
}

func logProcessStepEnd(processMessage string, options LogProcessStepEndOptions) {
	if options.ColorizeMsgFunc == nil {
		options.ColorizeMsgFunc = ColorizeBase
	}

	processMessageFunc := func() error {
		return WithoutIndent(func() error {
			processAndLogLn(outStream, prepareLogProcessMsgLeftPart(processMessage, options.ColorizeMsgFunc))
			return nil
		})
	}

	processMessageFunc = decorateByWithExtraProcessBorder(logProcessVerticalAndRightBorderSign, options.ColorizeMsgFunc, processMessageFunc)
	processMessageFunc = decorateByWithoutLastProcessBorder(processMessageFunc)

	_ = processMessageFunc()
}

func applyInfoLogProcessStep(userError error, infoSectionFunc func(err error), withIndent bool, colorizeMsgFunc func(...interface{}) string) {
	infoHeaderFunc := func() error {
		return WithoutIndent(func() error {
			processAndLogLn(outStream, prepareLogProcessMsgLeftPart("Info", colorizeMsgFunc))
			return nil
		})
	}

	infoHeaderFunc = decorateByWithExtraProcessBorder(logProcessVerticalAndRightBorderSign, colorizeMsgFunc, infoHeaderFunc)
	infoHeaderFunc = decorateByWithoutLastProcessBorder(infoHeaderFunc)

	_ = infoHeaderFunc()

	infoFunc := func() error {
		infoSectionFunc(userError)
		return nil
	}

	if withIndent {
		infoFunc = decorateByWithIndent(infoFunc)
	}

	infoFunc = decorateByWithExtraProcessBorder(logProcessVerticalBorderSign, colorizeMsgFunc, infoFunc)
	infoFunc = decorateByWithoutLastProcessBorder(infoFunc)

	_ = infoFunc()
}

func logProcessEnd(options LogProcessEndOptions) {
	if options.ColorizeMsgFunc == nil {
		options.ColorizeMsgFunc = ColorizeBase
	}

	popProcessBorder()

	logProcess := activeLogProcesses[len(activeLogProcesses)-1]
	activeLogProcesses = activeLogProcesses[:len(activeLogProcesses)-1]

	resetOptionalLnMode()

	elapsedSeconds := fmt.Sprintf(logProcessTimeFormat, time.Since(logProcess.StartedAt).Seconds())

	footerFunc := func() error {
		return WithoutIndent(func() error {
			timePart := fmt.Sprintf(" (%s)", elapsedSeconds)
			processAndLogF(outStream, prepareLogProcessMsgLeftPart(logProcess.Msg, options.ColorizeMsgFunc, timePart))
			colorizeFormatAndLogF(outStream, options.ColorizeMsgFunc, "%s\n", timePart)

			return nil
		})
	}

	footerFunc = decorateByWithExtraProcessBorder(logProcessUpAndRightBorderSign, options.ColorizeMsgFunc, footerFunc)

	_ = footerFunc()

	if !options.WithoutLogOptionalLn {
		LogOptionalLn()
	}
}

func logProcessFail(options LogProcessEndOptions) {
	if options.ColorizeMsgFunc == nil {
		options.ColorizeMsgFunc = ColorizeBase
	}

	popProcessBorder()

	logProcess := activeLogProcesses[len(activeLogProcesses)-1]
	activeLogProcesses = activeLogProcesses[:len(activeLogProcesses)-1]

	resetOptionalLnMode()

	elapsedSeconds := fmt.Sprintf(logProcessTimeFormat, time.Since(logProcess.StartedAt).Seconds())

	footerFunc := func() error {
		return WithoutIndent(func() error {
			timePart := fmt.Sprintf(" (%s) FAILED", elapsedSeconds)
			processAndLogF(outStream, prepareLogProcessMsgLeftPart(logProcess.Msg, ColorizeFail, timePart))
			colorizeFormatAndLogF(outStream, ColorizeFail, "%s\n", timePart)

			return nil
		})
	}

	footerFunc = decorateByWithExtraProcessBorder(logProcessUpAndRightBorderSign, options.ColorizeMsgFunc, footerFunc)

	_ = footerFunc()

	if !options.WithoutLogOptionalLn {
		LogOptionalLn()
	}
}

func decorateByWithExtraProcessBorder(colorlessBorder string, colorizeFunc func(...interface{}) string, decoratedFunc func() error) func() error {
	return func() error {
		return withExtraProcessBorder(colorlessBorder, colorizeFunc, decoratedFunc)
	}
}

func withExtraProcessBorder(colorlessValue string, colorizeFunc func(...interface{}) string, decoratedFunc func() error) error {
	appendProcessBorder(colorlessValue, colorizeFunc)
	err := decoratedFunc()
	popProcessBorder()

	return err
}

func decorateByWithoutLastProcessBorder(decoratedFunc func() error) func() error {
	return func() error {
		return withoutLastProcessBorder(decoratedFunc)
	}
}

func withoutLastProcessBorder(f func() error) error {
	oldBorderValue := processesBorderValues[len(processesBorderValues)-1]
	processesBorderValues = processesBorderValues[:len(processesBorderValues)-1]

	oldBorderFormattedValue := processesBorderFormattedValues[len(processesBorderFormattedValues)-1]
	processesBorderFormattedValues = processesBorderFormattedValues[:len(processesBorderFormattedValues)-1]

	err := f()

	processesBorderValues = append(processesBorderValues, oldBorderValue)
	processesBorderFormattedValues = append(processesBorderFormattedValues, oldBorderFormattedValue)

	return err
}

func appendProcessBorder(colorlessValue string, colorizeFunc func(...interface{}) string) {
	processesBorderValues = append(processesBorderValues, colorlessValue)
	processesBorderFormattedValues = append(processesBorderFormattedValues, colorizeFunc(colorlessValue))
}

func popProcessBorder() {
	if len(processesBorderValues) == 0 {
		return
	}

	processesBorderValues = processesBorderValues[:len(processesBorderValues)-1]
	processesBorderFormattedValues = processesBorderFormattedValues[:len(processesBorderFormattedValues)-1]
}

func formattedProcessBorders() string {
	if len(processesBorderValues) == 0 {
		return ""
	}

	return strings.Join(processesBorderFormattedValues, strings.Repeat(" ", processesBorderBetweenIndentWidth)) + strings.Repeat(" ", processesBorderIndentWidth)
}

func processBordersBlockWidth() int {
	if len(processesBorderValues) == 0 {
		return 0
	}

	return len([]rune(strings.Join(processesBorderValues, strings.Repeat(" ", processesBorderBetweenIndentWidth)))) + processesBorderIndentWidth
}
