package logboek

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

const (
	isControlSequenceNoneProcessed = iota
	isControlSequenceEscapeSequenceProcessed
	isControlSequenceOpenBorderProcessed
	isControlSequenceParametersProcessed

	resetColorControlSequence = "\x1b[0m"
	escapeSequenceCode        = 27

	resetColorCode = 0
)

type fitterState struct {
	wrapperState
	controlSequenceState
	colorState
}

type wrapperState struct {
	lineSequenceStack sequenceStack
	wordSequenceStack sequenceStack
}

func (ws *wrapperState) ResetLineSequences() {
	ws.lineSequenceStack = newSequenceStack()
}

func (ws *wrapperState) ResetWordSequences() {
	ws.wordSequenceStack = newSequenceStack()
}

func (ws *wrapperState) String() string {
	return ws.lineSequenceStack.String() + ws.wordSequenceStack.String()
}

func (ws *wrapperState) TWidth() int {
	return ws.lineSequenceStack.TWidth() + ws.wordSequenceStack.TWidth()
}

type controlSequenceState struct {
	controlSequenceBytes       []rune
	controlSequenceCursorState int
}

type colorState struct {
	isColorLine               bool
	prevCursorRune            rune
	colorControlSequenceCodes []int
}

func stripLongSequenceString(sequenceStack sequenceStack, contentWidth int, markLines, markLastLine bool) string {
	var result string

	var sliceTWidth int
	if markLines {
		sliceTWidth = contentWidth - 2
	} else {
		sliceTWidth = contentWidth
	}

	slices, rest := sequenceStack.Slices(sliceTWidth)

	if len(slices) == 0 {
		return result
	} else if len(slices) > 1 {
		for _, slice := range slices[:len(slices)-1] {
			if markLines {
				result += markLine(slice, sliceTWidth, contentWidth)
			} else {
				result += slice
			}
			result += "\n"
		}
	}

	lastSlice := slices[len(slices)-1]
	if markLines && markLastLine {
		result += markLine(lastSlice, sliceTWidth-rest, contentWidth)
	} else {
		result += lastSlice
	}

	return result
}

func markLine(line string, twidth, contentWidth int) string {
	var padding int
	if twidth <= contentWidthWithoutMarkSign(contentWidth, true) {
		padding = contentWidth - twidth - 1
	}

	return line + strings.Repeat(" ", padding) + "↵"
}

func (s *fitterState) parseColorCodes() []int {
	preparedColorFormatsPart := string(s.controlSequenceBytes[:len(s.controlSequenceBytes)-1])
	preparedColorFormatsPart = string(bytes.TrimPrefix([]byte(preparedColorFormatsPart), []byte{escapeSequenceCode, []byte("[")[0]}))

	colorCodesStrings := strings.Split(preparedColorFormatsPart, ";")
	var colorCodes []int
	for _, colorCodeString := range colorCodesStrings {
		if colorCodeString == "" {
			continue
		}

		cd, err := strconv.Atoi(colorCodeString)
		if err != nil {
			panic(err)
		}

		colorCodes = append(colorCodes, cd)
	}

	return colorCodes
}

func (s *fitterState) generateColorControlSequence() string {
	var result string

	if len(s.colorControlSequenceCodes) != 0 {
		result = "\x1b["

		var colorCodesStrings []string
		for _, colorCode := range s.colorControlSequenceCodes {
			colorCodesStrings = append(colorCodesStrings, fmt.Sprintf("%d", colorCode))
		}
		result += strings.Join(colorCodesStrings, ";")

		result += "m"
	}

	return result
}

func (s *fitterState) addColorControlSequenceCode(newColorCode int) {
	for i, colorCode := range s.colorControlSequenceCodes {
		if colorCode == newColorCode {
			s.colorControlSequenceCodes = append(s.colorControlSequenceCodes[:i], s.colorControlSequenceCodes[i+1:]...)
			break
		}
	}

	s.colorControlSequenceCodes = append(s.colorControlSequenceCodes, newColorCode)
}

func (s *fitterState) resetColorCodes() {
	s.colorControlSequenceCodes = []int{}
}

func contentWidthWithoutMarkSign(contentWidth int, markWrappedLine bool) int {
	if markWrappedLine {
		return contentWidth - 1
	}

	return contentWidth
}

func fitText(text string, s fitterState, contentWidth int, markWrappedLine bool, cacheIncompleteLine bool) (string, fitterState) {
	var result string

	for _, r := range []rune(text) {
		result += runFitterWrapper(r, &s, contentWidth, markWrappedLine)
		ignoreControlSequenceTWidth(r, &s)
	}

	if !cacheIncompleteLine {
		result += processFitterCachedLineAndWord(&s, contentWidth, markWrappedLine)
	}

	result = addRequiredColorControlSequences(result, &s)

	return result, s
}

func runFitterWrapper(r rune, s *fitterState, contentWidth int, markWrappedLine bool) string {
	var result string

	switch string(r) {
	case "\b":
		if !s.wordSequenceStack.IsEmpty() {
			s.wordSequenceStack.WriteControlData(string(r))
		} else {
			s.lineSequenceStack.WriteControlData(string(r))
		}
	case "\n", "\r":
		carriage := string(r)

		if s.wrapperState.TWidth() <= contentWidth {
			result += s.wrapperState.String()
		} else if s.wrapperState.TWidth() > contentWidth {
			result += stripLongSequenceString(s.lineSequenceStack, contentWidth, markWrappedLine, true)
			result += "\n"
			result += stripLongSequenceString(s.wordSequenceStack, contentWidth, markWrappedLine, false)
		}

		result += carriage

		s.ResetLineSequences()
		s.ResetWordSequences()
	case " ":
		space := string(r)
		spaceTWidth := len(space)

		if s.wrapperState.TWidth()+spaceTWidth > contentWidthWithoutMarkSign(contentWidth, markWrappedLine) {
			result += stripLongSequenceString(s.lineSequenceStack, contentWidth, markWrappedLine, true)
			result += "\n"

			s.wrapperState.ResetLineSequences()
		}

		s.lineSequenceStack.Merge(s.wordSequenceStack)
		s.lineSequenceStack.WritePlainData(space)

		s.ResetWordSequences()
	default:
		s.wordSequenceStack.WriteData(string(r))
	}

	return result
}

func ignoreControlSequenceTWidth(r rune, s *fitterState) {
	processControlSequenceFunc := func(s *fitterState, _ string) {
		s.wordSequenceStack.CommitTopSequenceAsControl()
	}

	processEscapeSequenceCodeFunc := func(s *fitterState) {
		s.wordSequenceStack.DivideListSign()
	}

	processFitterControlSequence(r, s, processEscapeSequenceCodeFunc, processControlSequenceFunc)
}

func processFitterControlSequence(r rune, s *fitterState, processEscapeSequenceCodeFunc func(f *fitterState), processControlSequenceFunc func(f *fitterState, code string)) {
	switch s.controlSequenceCursorState {
	case isControlSequenceNoneProcessed:
		switch r {
		case escapeSequenceCode:
			s.controlSequenceBytes = []rune{r}
			s.controlSequenceCursorState = isControlSequenceEscapeSequenceProcessed

			if processEscapeSequenceCodeFunc != nil {
				processEscapeSequenceCodeFunc(s)
			}
		}
	case isControlSequenceEscapeSequenceProcessed:
		switch string(r) {
		case "[":
			s.controlSequenceBytes = append(s.controlSequenceBytes, r)
			s.controlSequenceCursorState = isControlSequenceOpenBorderProcessed
		}
	case isControlSequenceOpenBorderProcessed, isControlSequenceParametersProcessed:
		if unicode.IsNumber(r) || string(r) == ";" {
			s.controlSequenceBytes = append(s.controlSequenceBytes, r)
			s.controlSequenceCursorState = isControlSequenceParametersProcessed
		} else {
			if unicode.IsLetter(r) {
				s.controlSequenceBytes = append(s.controlSequenceBytes, r)

				if processControlSequenceFunc != nil {
					processControlSequenceFunc(s, string(r))
				}
			}

			s.controlSequenceCursorState = isControlSequenceNoneProcessed
		}
	default:
		s.controlSequenceCursorState = isControlSequenceNoneProcessed
	}
}

func processFitterCachedLineAndWord(s *fitterState, contentWidth int, markWrappedLine bool) string {
	var result string

	if s.wrapperState.TWidth() > contentWidthWithoutMarkSign(contentWidth, markWrappedLine) {
		if s.lineSequenceStack.String() != "" {
			result += stripLongSequenceString(s.lineSequenceStack, contentWidth, markWrappedLine, true)

			if s.wordSequenceStack.String() != "" {
				result += "\n"
				result += stripLongSequenceString(s.wordSequenceStack, contentWidth, markWrappedLine, false)

				s.ResetLineSequences()
				s.ResetWordSequences()
			}
		} else if s.wordSequenceStack.String() != "" {
			result += stripLongSequenceString(s.wordSequenceStack, contentWidth, markWrappedLine, false)

			s.ResetWordSequences()
		}
	} else {
		result += s.wrapperState.String()
	}

	return result
}

func addRequiredColorControlSequences(fittedText string, s *fitterState) string {
	var result string

	for _, r := range []rune(fittedText) {
		switch string(r) {
		case "\n", "\r":
			if string(s.prevCursorRune) == "\r" {
				result += string(r)
			} else {
				if s.isColorLine {
					result += resetColorControlSequence
				}

				result += string(r)
			}
		default:
			if string(s.prevCursorRune) == "\r" || string(s.prevCursorRune) == "\n" {
				result += s.generateColorControlSequence()
			}

			result += string(r)
		}

		s.prevCursorRune = r

		processControlSequenceFunc := func(s *fitterState, code string) {
			if string(r) == "m" {
				processColorControlSequence(s)
			}
		}

		processFitterControlSequence(r, s, nil, processControlSequenceFunc)
	}

	return result
}

func processColorControlSequence(s *fitterState) {
	colorCodes := s.parseColorCodes()
	for _, colorCode := range colorCodes {
		if isResetColorCode(colorCode) {
			s.resetColorCodes()
			s.isColorLine = false
		} else {
			s.addColorControlSequenceCode(colorCode)
			s.isColorLine = true
		}
	}
}

func isResetColorCode(code int) bool {
	return code == resetColorCode
}
