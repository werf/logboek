package logboek

import (
	"fmt"
	"io"
	"strings"
	"time"
)

const (
	logProcessTimeFormat = "%.2f seconds"

	logStateRightPartsSeparator = " "
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

type LevelLogBlockOptions struct {
	WithIndent           bool
	WithoutLogOptionalLn bool
	Style                *Style
}

type LogBlockOptions struct {
	LevelLogBlockOptions
	Level Level
}

func LogBlock(blockMessage string, options LogBlockOptions, blockFunc func() error) error {
	stream := options.Level.Stream()

	if !options.Level.IsAccepted() {
		return blockFunc()
	}

	style := options.Style
	if options.Style == nil {
		style = options.Level.Style()
	}

	titleFunc := func() error {
		processAndLogLn(stream, formatWithStyle(style, blockMessage))
		return nil
	}

	applyOptionalLnMode()

	bodyFunc := func() error {
		return blockFunc()
	}

	if options.WithIndent {
		bodyFunc = decorateByWithIndent(bodyFunc)
	}

	_ = decorateByWithExtraProcessBorder(logProcessDownAndRightBorderSign, style, titleFunc)()
	err := decorateByWithExtraProcessBorder(logProcessVerticalBorderSign, style, bodyFunc)()

	resetOptionalLnMode()

	_ = decorateByWithExtraProcessBorder(logProcessUpAndRightBorderSign, style, titleFunc)()

	if !options.WithoutLogOptionalLn {
		LogOptionalLn()
	}

	return err
}

type LevelLogProcessInlineOptions struct {
	Style *Style
}

type LogProcessInlineOptions struct {
	LevelLogProcessInlineOptions
	Level Level
}

func LogProcessInline(processMessage string, options LogProcessInlineOptions, processFunc func() error) error {
	if !options.Level.IsAccepted() {
		return processFunc()
	}

	return logProcessInline(processMessage, options, processFunc)
}

func logProcessInline(processMessage string, options LogProcessInlineOptions, processFunc func() error) error {
	stream := options.Level.Stream()

	style := options.Style
	if options.Style == nil {
		style = options.Level.Style()
	}

	progressDots := "..."
	maxLength := ContentWidth() - len(" ") - len(progressDots) - len(fmt.Sprintf(logProcessTimeFormat, 1234.0))
	if len(processMessage) > maxLength {
		processMessage = processMessage[:maxLength-1]
	}

	processMessage = processMessage + " " + progressDots
	formatAndLogF(stream, style, "%s", processMessage)

	resultStyle := style
	start := time.Now()

	resultFormat := " (%s)\n"

	var err error
	if err = WithIndent(processFunc); err != nil {
		resultStyle = StyleByName(FailStyleName)
		resultFormat = " (%s) FAILED\n"
	}

	elapsedSeconds := fmt.Sprintf(logProcessTimeFormat, time.Since(start).Seconds())
	formatAndLogF(stream, resultStyle, resultFormat, elapsedSeconds)

	return err
}

func prepareLogProcessMsgLeftPart(leftPart string, style *Style, rightParts ...string) string {
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

	return formatWithStyle(style, result)
}

type LevelLogProcessStartOptions struct {
	Style *Style
}

type LogProcessStartOptions struct {
	LevelLogProcessStartOptions
	Level Level
}

type LevelLogProcessEndOptions struct {
	WithoutLogOptionalLn bool
	WithoutElapsedTime   bool
	Style                *Style
}

type LogProcessEndOptions struct {
	LevelLogProcessEndOptions
	Level Level
}

type LevelLogProcessFailOptions struct {
	LevelLogProcessEndOptions
}

type LogProcessFailOptions struct {
	LevelLogProcessFailOptions
	Level Level
}

type LevelLogProcessStepEndOptions struct {
	WithIndent      bool
	InfoSectionFunc func(err error)
	Style           *Style
}

type LogProcessStepEndOptions struct {
	LevelLogProcessStepEndOptions
	Level Level
}

type LevelLogProcessOptions struct {
	WithIndent             bool
	WithoutLogOptionalLn   bool
	WithoutElapsedTime     bool
	InfoSectionFunc        func(err error)
	SuccessInfoSectionFunc func()
	Style                  *Style
}

type LogProcessOptions struct {
	LevelLogProcessOptions
	Level Level
}

func LogProcessStart(processMessage string, options LogProcessStartOptions) {
	if !options.Level.IsAccepted() {
		return
	}

	logProcessStart(processMessage, options)
}

func LogProcessEnd(options LogProcessEndOptions) {
	if !options.Level.IsAccepted() {
		return
	}

	logProcessEnd(options)
}

func LogProcessStepEnd(processMessage string, options LogProcessStepEndOptions) {
	if !options.Level.IsAccepted() {
		return
	}

	logProcessStepEnd(processMessage, options)
}

func LogProcessFail(options LogProcessFailOptions) {
	if !options.Level.IsAccepted() {
		return
	}

	logProcessFail(options)
}

func LogProcess(processMessage string, options LogProcessOptions, processFunc func() error) error {
	if !options.Level.IsAccepted() {
		return processFunc()
	}

	return logProcess(processMessage, options, processFunc)
}

func logProcess(processMessage string, options LogProcessOptions, processFunc func() error) error {
	stream := options.Level.Stream()

	style := options.Style
	if options.Style == nil {
		style = options.Level.Style()
	}

	logProcessStart(
		processMessage,
		LogProcessStartOptions{
			LevelLogProcessStartOptions: LevelLogProcessStartOptions{
				Style: style,
			},
			Level: options.Level,
		},
	)

	bodyFunc := func() error {
		return processFunc()
	}

	if options.WithIndent {
		bodyFunc = decorateByWithIndent(bodyFunc)
	}

	err := bodyFunc()

	resetOptionalLnMode()

	if options.InfoSectionFunc != nil {
		applyInfoLogProcessStep(stream, err, options.InfoSectionFunc, options.WithIndent, style)
	}

	if options.SuccessInfoSectionFunc != nil && err == nil {
		infoSectionFunc := func(_ error) {
			options.SuccessInfoSectionFunc()
		}

		applyInfoLogProcessStep(stream, err, infoSectionFunc, options.WithIndent, style)
	}

	if err != nil {
		logProcessFail(LogProcessFailOptions{
			LevelLogProcessFailOptions: LevelLogProcessFailOptions{
				LevelLogProcessEndOptions: LevelLogProcessEndOptions{
					WithoutLogOptionalLn: options.WithoutLogOptionalLn,
					WithoutElapsedTime:   options.WithoutElapsedTime,
					Style:                style,
				},
			},
			Level: options.Level,
		})
		return err
	}

	logProcessEnd(
		LogProcessEndOptions{
			LevelLogProcessEndOptions: LevelLogProcessEndOptions{
				WithoutLogOptionalLn: options.WithoutLogOptionalLn,
				WithoutElapsedTime:   options.WithoutElapsedTime,
				Style:                style,
			},
			Level: options.Level,
		},
	)
	return nil
}

func logProcessStart(processMessage string, options LogProcessStartOptions) {
	stream := options.Level.Stream()

	style := options.Style
	if options.Style == nil {
		style = options.Level.Style()
	}

	applyOptionalLnMode()

	headerFunc := func() error {
		return WithoutIndent(func() error {
			processAndLogLn(stream, prepareLogProcessMsgLeftPart(processMessage, style))
			return nil
		})
	}

	headerFunc = decorateByWithExtraProcessBorder(logProcessDownAndRightBorderSign, style, headerFunc)

	_ = headerFunc()

	appendProcessBorder(logProcessVerticalBorderSign, style)

	logProcess := &logProcessDescriptor{StartedAt: time.Now(), Msg: processMessage}
	activeLogProcesses = append(activeLogProcesses, logProcess)
}

func logProcessStepEnd(processMessage string, options LogProcessStepEndOptions) {
	stream := options.Level.Stream()

	style := options.Style
	if options.Style == nil {
		style = options.Level.Style()
	}

	processMessageFunc := func() error {
		return WithoutIndent(func() error {
			processAndLogLn(stream, prepareLogProcessMsgLeftPart(processMessage, style))
			return nil
		})
	}

	processMessageFunc = decorateByWithExtraProcessBorder(logProcessVerticalAndRightBorderSign, style, processMessageFunc)
	processMessageFunc = decorateByWithoutLastProcessBorder(processMessageFunc)

	_ = processMessageFunc()
}

func applyInfoLogProcessStep(stream io.Writer, userError error, infoSectionFunc func(err error), withIndent bool, style *Style) {
	infoHeaderFunc := func() error {
		return WithoutIndent(func() error {
			processAndLogLn(stream, prepareLogProcessMsgLeftPart("Info", style))
			return nil
		})
	}

	infoHeaderFunc = decorateByWithExtraProcessBorder(logProcessVerticalAndRightBorderSign, style, infoHeaderFunc)
	infoHeaderFunc = decorateByWithoutLastProcessBorder(infoHeaderFunc)

	_ = infoHeaderFunc()

	infoFunc := func() error {
		infoSectionFunc(userError)
		return nil
	}

	if withIndent {
		infoFunc = decorateByWithIndent(infoFunc)
	}

	infoFunc = decorateByWithExtraProcessBorder(logProcessVerticalBorderSign, style, infoFunc)
	infoFunc = decorateByWithoutLastProcessBorder(infoFunc)

	_ = infoFunc()
}

func logProcessEnd(options LogProcessEndOptions) {
	stream := options.Level.Stream()

	style := options.Style
	if options.Style == nil {
		style = options.Level.Style()
	}

	popProcessBorder()

	logProcess := activeLogProcesses[len(activeLogProcesses)-1]
	activeLogProcesses = activeLogProcesses[:len(activeLogProcesses)-1]

	resetOptionalLnMode()

	elapsedSeconds := fmt.Sprintf(logProcessTimeFormat, time.Since(logProcess.StartedAt).Seconds())

	footerFunc := func() error {
		return WithoutIndent(func() error {
			timePart := ""
			if !options.WithoutElapsedTime {
				timePart = fmt.Sprintf(" (%s)", elapsedSeconds)
			}

			processAndLogF(stream, prepareLogProcessMsgLeftPart(logProcess.Msg, style, timePart))
			formatAndLogF(stream, style, "%s\n", timePart)

			return nil
		})
	}

	footerFunc = decorateByWithExtraProcessBorder(logProcessUpAndRightBorderSign, style, footerFunc)

	_ = footerFunc()

	if !options.WithoutLogOptionalLn {
		LogOptionalLn()
	}
}

func logProcessFail(options LogProcessFailOptions) {
	stream := options.Level.Stream()

	style := options.Style
	if options.Style == nil {
		style = options.Level.Style()
	}

	popProcessBorder()

	logProcess := activeLogProcesses[len(activeLogProcesses)-1]
	activeLogProcesses = activeLogProcesses[:len(activeLogProcesses)-1]

	resetOptionalLnMode()

	elapsedSeconds := fmt.Sprintf(logProcessTimeFormat, time.Since(logProcess.StartedAt).Seconds())

	footerFunc := func() error {
		return WithoutIndent(func() error {
			timePart := " FAILED"
			if !options.WithoutElapsedTime {
				timePart = fmt.Sprintf(" (%s) FAILED", elapsedSeconds)
			}

			processAndLogF(stream, prepareLogProcessMsgLeftPart(logProcess.Msg, StyleByName(FailStyleName), timePart))
			formatAndLogF(stream, StyleByName(FailStyleName), "%s\n", timePart)

			return nil
		})
	}

	footerFunc = decorateByWithExtraProcessBorder(logProcessUpAndRightBorderSign, style, footerFunc)

	_ = footerFunc()

	if !options.WithoutLogOptionalLn {
		LogOptionalLn()
	}
}

func decorateByWithExtraProcessBorder(colorlessBorder string, style *Style, decoratedFunc func() error) func() error {
	return func() error {
		return withExtraProcessBorder(colorlessBorder, style, decoratedFunc)
	}
}

func withExtraProcessBorder(colorlessValue string, style *Style, decoratedFunc func() error) error {
	appendProcessBorder(colorlessValue, style)
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

func appendProcessBorder(colorlessValue string, style *Style) {
	processesBorderValues = append(processesBorderValues, colorlessValue)
	processesBorderFormattedValues = append(processesBorderFormattedValues, formatWithStyle(style, colorlessValue))
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
