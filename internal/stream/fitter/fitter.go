package fitter

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

type State struct {
	wrapperState
	controlSequenceState
	colorState
}

func NewState() State {
	return State{}
}

type wrapperState struct {
	sequenceStack
}

func (ws *wrapperState) Apply(contentWidth int, markLines bool) string {
	var result string

	if ws.sequenceStack.IsEmpty() {
		return ""
	}

	if ws.sequenceStack.TWidth() <= contentWidth {
		result = ws.sequenceStack.String()
	} else {
		result = ws.splitSequenceStack(contentWidth, markLines)
	}

	ws.resetSequenceStack()

	return result
}

func (ws *wrapperState) splitSequenceStack(contentWidth int, markLines bool) string {
	var sliceTWidth int
	if markLines {
		sliceTWidth = contentWidth - 2
	} else {
		sliceTWidth = contentWidth
	}

	if sliceTWidth < 1 {
		sliceTWidth = 1
	}

	slices, _ := ws.sequenceStack.Slices(sliceTWidth)

	var formattedSlices []string
	if len(slices) == 0 {
		return ""
	} else if len(slices) > 1 {
		for _, slice := range slices[:len(slices)-1] {
			if markLines {
				formattedSlices = append(formattedSlices, markLine(slice, sliceTWidth, contentWidth))
			} else {
				formattedSlices = append(formattedSlices, slice)
			}
		}
	}

	lastSlice := slices[len(slices)-1]
	if strings.HasPrefix(lastSlice, " ") {
		lastSlice = lastSlice[1:]
	}

	if lastSlice != "" {
		formattedSlices = append(formattedSlices, lastSlice)
	}

	return strings.Join(formattedSlices, "\n")
}

func (ws *wrapperState) resetSequenceStack() {
	ws.sequenceStack = newSequenceStack()
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

func markLine(line string, twidth, contentWidth int) string {
	var padding int
	if twidth <= contentWidthWithoutMarkSign(contentWidth, true) {
		padding = contentWidth - twidth - 1
	}

	return line + strings.Repeat(" ", padding) + "↵"
}

func (s *State) parseColorCodes() []int {
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

func (s *State) generateColorControlSequence() string {
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

func (s *State) addColorControlSequenceCode(newColorCode int) {
	for i, colorCode := range s.colorControlSequenceCodes {
		if colorCode == newColorCode {
			s.colorControlSequenceCodes = append(s.colorControlSequenceCodes[:i], s.colorControlSequenceCodes[i+1:]...)
			break
		}
	}

	s.colorControlSequenceCodes = append(s.colorControlSequenceCodes, newColorCode)
}

func (s *State) resetColorCodes() {
	s.colorControlSequenceCodes = []int{}
}

func contentWidthWithoutMarkSign(contentWidth int, markWrappedLine bool) int {
	if markWrappedLine {
		return contentWidth - 1
	}

	return contentWidth
}

func FitText(text string, s *State, contentWidth int, markWrappedLine bool, cacheIncompleteLine bool) string {
	var b strings.Builder

	for _, r := range []rune(text) {
		b.WriteString(runFitterWrapper(r, s, contentWidth, markWrappedLine))
		ignoreControlSequenceTWidth(r, s)
	}

	if !cacheIncompleteLine {
		b.WriteString(s.wrapperState.Apply(contentWidth, markWrappedLine))
	}

	return addRequiredColorControlSequences(b.String(), s)
}

func runFitterWrapper(r rune, s *State, contentWidth int, markWrappedLine bool) string {
	var result string

	switch string(r) {
	case "\b":
		s.wrapperState.sequenceStack.WriteControlData(string(r))
	case "\n", "\r":
		carriage := string(r)
		result += s.wrapperState.Apply(contentWidth, markWrappedLine)
		result += carriage
	case " ":
		s.wrapperState.sequenceStack.WritePlainData(" ")
	default:
		s.wrapperState.sequenceStack.WriteData(string(r))
	}

	return result
}

func ignoreControlSequenceTWidth(r rune, s *State) {
	processControlSequenceFunc := func(s *State, _ string) {
		s.wrapperState.sequenceStack.CommitTopSequenceAsControl()
	}

	processEscapeSequenceCodeFunc := func(s *State) {
		s.wrapperState.sequenceStack.DivideLastSign()
	}

	processFitterControlSequence(r, s, processEscapeSequenceCodeFunc, processControlSequenceFunc)
}

func processFitterControlSequence(r rune, s *State, processEscapeSequenceCodeFunc func(f *State), processControlSequenceFunc func(f *State, code string)) {
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

func addRequiredColorControlSequences(fittedText string, s *State) string {
	var b strings.Builder

	for _, r := range []rune(fittedText) {
		switch string(r) {
		case "\n", "\r":
			if string(s.prevCursorRune) == "\r" {
				b.WriteRune(r)
			} else {
				if s.isColorLine {
					b.WriteString(resetColorControlSequence)
				}

				b.WriteRune(r)
			}
		default:
			if string(s.prevCursorRune) == "\r" || string(s.prevCursorRune) == "\n" {
				b.WriteString(s.generateColorControlSequence())
			}

			b.WriteRune(r)
		}

		s.prevCursorRune = r

		processControlSequenceFunc := func(s *State, code string) {
			if string(r) == "m" {
				processColorControlSequence(s)
			}
		}

		processFitterControlSequence(r, s, nil, processControlSequenceFunc)
	}

	return b.String()
}

func processColorControlSequence(s *State) {
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
