package stream

import (
	"fmt"
	"strings"
	"time"

	stylePkg "github.com/werf/logboek/pkg/style"
	"github.com/werf/logboek/pkg/types"
)

const (
	logProcessTimeFormat        = "%.2f seconds"
	logStateRightPartsSeparator = " "
	progressDots                = "..."
)

func (s *Stream) NewLogBlock(manager types.ManagerInterface, format string, a ...interface{}) *LogBlock {
	return &LogBlock{manager: manager, stream: s, title: stylePkg.SimpleFormat(format, a...), options: &LogBlockOptions{}}
}

func (s *Stream) NewLogProcessInline(manager types.ManagerInterface, format string, a ...interface{}) *LogProcessInline {
	return &LogProcessInline{manager: manager, stream: s, title: stylePkg.SimpleFormat(format, a...), options: &LogProcessInlineOptions{}}
}

func (s *Stream) NewLogProcess(manager types.ManagerInterface, format string, a ...interface{}) *LogProcess {
	return &LogProcess{manager: manager, stream: s, title: stylePkg.SimpleFormat(format, a...), options: &LogProcessOptions{}}
}

func (s *Stream) logBlock(blockMessage string, options *LogBlockOptions, blockFunc func() error) error {
	style := options.style
	if options.style == nil {
		style = stylePkg.None()
	}

	titleFunc := func() error {
		s.processAndLogLn(s.formatWithStyle(style, blockMessage))
		return nil
	}

	s.applyOptionalLn()

	bodyFunc := func() error {
		return blockFunc()
	}

	if options.withIndent {
		bodyFunc = s.decorateByDoErrorWithIndent(bodyFunc)
	}

	_ = s.decorateByWithExtraProcessBorder(
		s.logProcessDownAndRightBorderSign,
		style,
		titleFunc,
	)()

	err := s.decorateByWithExtraProcessBorder(
		s.logProcessVerticalBorderSign,
		style,
		bodyFunc,
	)()

	s.DisableOptionalLn()

	_ = s.decorateByWithExtraProcessBorder(
		s.logProcessUpAndRightBorderSign,
		style,
		titleFunc,
	)()

	if !options.withoutLogOptionalLn {
		s.EnableOptionalLn()
	}

	return err
}

func (s *Stream) logProcessInline(processMessage string, options *LogProcessInlineOptions, processFunc func() error) error {
	style := options.style
	if options.style == nil {
		style = stylePkg.None()
	}

	maxLength := s.ContentWidth() - len(" ") - len(progressDots) - len(fmt.Sprintf(logProcessTimeFormat, 1234.0))
	if maxLength > 0 && len(processMessage) > maxLength {
		processMessage = processMessage[:maxLength-1]
	} else {
		processMessage = ""
	}

	processMessage = processMessage + " " + progressDots
	s.FormatAndLogF(style, "%s", processMessage)

	resultStyle := style
	start := time.Now()

	resultFormat := " (%s)\n"

	err := s.DoErrorWithIndent(processFunc)
	if err != nil {
		resultStyle = stylePkg.Get(stylePkg.FailName)
		resultFormat = " (%s) FAILED\n"
	}

	elapsedSeconds := fmt.Sprintf(logProcessTimeFormat, time.Since(start).Seconds())
	s.FormatAndLogF(resultStyle, resultFormat, elapsedSeconds)

	return err
}

func (s *Stream) prepareLogProcessMsgLeftPart(leftPart string, style *stylePkg.Style, rightParts ...string) string {
	var result string

	spaceWidth := s.ContentWidth() - len(strings.Join(rightParts, logStateRightPartsSeparator))
	if spaceWidth > 0 {
		if spaceWidth > len(leftPart) {
			result = leftPart
		} else {
			service := " " + progressDots
			if spaceWidth > len([]rune(service)) {
				result = leftPart[0:spaceWidth-len(service)] + service
			} else {
				result = leftPart[0:spaceWidth]
			}
		}
	} else {
		return ""
	}

	return s.formatWithStyle(style, result)
}

func (s *Stream) logProcess(processMessage string, options *LogProcessOptions, processFunc func() error) error {
	style := options.style
	if options.style == nil {
		style = stylePkg.None()
	}

	s.logProcessStart(
		processMessage,
		LogProcessOptions{
			style: style,
		},
	)

	bodyFunc := func() error {
		return processFunc()
	}

	if options.withIndent {
		bodyFunc = s.decorateByDoErrorWithIndent(bodyFunc)
	}

	err := bodyFunc()

	s.DisableOptionalLn()

	if options.infoSectionFunc != nil {
		s.applyInfoLogProcessStep(err, options.infoSectionFunc, options.withIndent, style)
	}

	if options.successInfoSectionFunc != nil && err == nil {
		infoSectionFunc := func(_ error) {
			options.successInfoSectionFunc()
		}

		s.applyInfoLogProcessStep(err, infoSectionFunc, options.withIndent, style)
	}

	if err != nil {
		s.logProcessFail(
			LogProcessOptions{
				withoutLogOptionalLn: options.withoutLogOptionalLn,
				withoutElapsedTime:   options.withoutElapsedTime,
				style:                style,
			})

		return err
	}

	s.logProcessEnd(
		LogProcessOptions{
			withoutLogOptionalLn: options.withoutLogOptionalLn,
			withoutElapsedTime:   options.withoutElapsedTime,
			style:                style,
		},
	)

	return nil
}

type logProcessDescriptor struct {
	StartedAt time.Time
	Msg       string
}

func (s *Stream) logProcessStart(processMessage string, options LogProcessOptions) {
	style := options.style
	if options.style == nil {
		style = stylePkg.None()
	}

	s.applyOptionalLn()

	headerFunc := func() error {
		return s.DoErrorWithoutIndent(func() error {
			s.processAndLogLn(s.prepareLogProcessMsgLeftPart(processMessage, style))
			return nil
		})
	}

	headerFunc = s.decorateByWithExtraProcessBorder(s.logProcessDownAndRightBorderSign, style, headerFunc)

	_ = headerFunc()

	s.appendProcessBorder(s.logProcessVerticalBorderSign, style)

	logProcess := &logProcessDescriptor{StartedAt: time.Now(), Msg: processMessage}
	s.activeLogProcesses = append(s.activeLogProcesses, logProcess)
}

func (s *Stream) logProcessStepEnd(processMessage string, options LogProcessOptions) {
	style := options.style
	if options.style == nil {
		style = stylePkg.None()
	}

	processMessageFunc := func() error {
		return s.DoErrorWithoutIndent(func() error {
			s.processAndLogLn(s.prepareLogProcessMsgLeftPart(processMessage, style))
			return nil
		})
	}

	processMessageFunc = s.decorateByWithExtraProcessBorder(s.logProcessVerticalAndRightBorderSign, style, processMessageFunc)
	processMessageFunc = s.decorateByWithoutLastProcessBorder(processMessageFunc)

	_ = processMessageFunc()
}

func (s *Stream) applyInfoLogProcessStep(userError error, infoSectionFunc func(err error), withIndent bool, style *stylePkg.Style) {
	infoHeaderFunc := func() error {
		return s.DoErrorWithoutIndent(func() error {
			s.processAndLogLn(s.prepareLogProcessMsgLeftPart("Info", style))
			return nil
		})
	}

	infoHeaderFunc = s.decorateByWithExtraProcessBorder(s.logProcessVerticalAndRightBorderSign, style, infoHeaderFunc)
	infoHeaderFunc = s.decorateByWithoutLastProcessBorder(infoHeaderFunc)

	_ = infoHeaderFunc()

	infoFunc := func() error {
		infoSectionFunc(userError)
		return nil
	}

	if withIndent {
		infoFunc = s.decorateByDoErrorWithIndent(infoFunc)
	}

	infoFunc = s.decorateByWithExtraProcessBorder(s.logProcessVerticalBorderSign, style, infoFunc)
	infoFunc = s.decorateByWithoutLastProcessBorder(infoFunc)

	_ = infoFunc()
}

func (s *Stream) logProcessEnd(options LogProcessOptions) {
	style := options.style
	if options.style == nil {
		style = stylePkg.None()
	}

	s.popProcessBorder()

	logProcess := s.activeLogProcesses[len(s.activeLogProcesses)-1]
	s.activeLogProcesses = s.activeLogProcesses[:len(s.activeLogProcesses)-1]

	s.DisableOptionalLn()

	elapsedSeconds := fmt.Sprintf(logProcessTimeFormat, time.Since(logProcess.StartedAt).Seconds())

	footerFunc := func() error {
		return s.DoErrorWithoutIndent(func() error {
			timePart := ""
			if !options.withoutElapsedTime {
				timePart = fmt.Sprintf(" (%s)", elapsedSeconds)
			}

			s.processAndLogF(s.prepareLogProcessMsgLeftPart(logProcess.Msg, style, timePart))
			s.FormatAndLogF(style, "%s\n", timePart)

			return nil
		})
	}

	footerFunc = s.decorateByWithExtraProcessBorder(s.logProcessUpAndRightBorderSign, style, footerFunc)

	_ = footerFunc()

	if !options.withoutLogOptionalLn {
		s.EnableOptionalLn()
	}
}

func (s *Stream) logProcessFail(options LogProcessOptions) {
	style := options.style
	if options.style == nil {
		style = stylePkg.None()
	}

	s.popProcessBorder()

	logProcess := s.activeLogProcesses[len(s.activeLogProcesses)-1]
	s.activeLogProcesses = s.activeLogProcesses[:len(s.activeLogProcesses)-1]

	s.DisableOptionalLn()

	elapsedSeconds := fmt.Sprintf(logProcessTimeFormat, time.Since(logProcess.StartedAt).Seconds())

	footerFunc := func() error {
		return s.DoErrorWithoutIndent(func() error {
			timePart := " FAILED"
			if !options.withoutElapsedTime {
				timePart = fmt.Sprintf(" (%s) FAILED", elapsedSeconds)
			}

			s.processAndLogF(s.prepareLogProcessMsgLeftPart(logProcess.Msg, stylePkg.Get(stylePkg.FailName), timePart))
			s.FormatAndLogF(stylePkg.Get(stylePkg.FailName), "%s\n", timePart)

			return nil
		})
	}

	footerFunc = s.decorateByWithExtraProcessBorder(s.logProcessUpAndRightBorderSign, style, footerFunc)

	_ = footerFunc()

	if !options.withoutLogOptionalLn {
		s.EnableOptionalLn()
	}
}
