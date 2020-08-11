package stream

import (
	"os"
	"strings"
	"time"

	"github.com/fatih/color"

	"github.com/werf/logboek/internal/stream/fitter"
	stylePkg "github.com/werf/logboek/pkg/style"
)

type State struct {
	width               int
	indentWidth         int
	isOptionalLnEnabled bool

	fitter.State
	streamMode
	cursorState
	processState
	tagState
	prefixState
}

func NewStreamState() *State {
	return &State{
		streamMode:   newStreamMode(),
		cursorState:  newCursorState(),
		processState: newProcessState(),
		prefixState:  newPrefixState(),
	}
}

func (s *State) SubState() *State {
	return &State{
		width:      s.ContentWidth(),
		streamMode: s.streamMode,
	}
}

func (s *State) DisablePrettyLog() {
	s.DisableProxyStreamDataFormatting()
	s.DisableLogProcessBorder()
	s.DisableLineWrapping()
}

type streamMode struct {
	isMuted                            bool
	isStyleEnabled                     bool
	isLineWrappingEnabled              bool
	isProxyStreamDataFormattingEnabled bool
	isGitlabCollapsibleSectionsEnabled bool
}

func newStreamMode() streamMode {
	return streamMode{
		isStyleEnabled:                     !color.NoColor,
		isLineWrappingEnabled:              true,
		isProxyStreamDataFormattingEnabled: true,
		isGitlabCollapsibleSectionsEnabled: os.Getenv("GITLAB_CI") == "true",
	}
}

func (s *State) Mute() {
	s.isMuted = true
}

func (s *State) Unmute() {
	s.isMuted = false
}

func (s *State) IsMuted() bool {
	return s.isMuted
}

func (s *State) EnableGitlabCollapsibleSections() {
	s.isGitlabCollapsibleSectionsEnabled = true
}

func (s *State) DisableGitlabCollapsibleSections() {
	s.isGitlabCollapsibleSectionsEnabled = false
}

func (s *State) IsGitlabCollapsibleSections() bool {
	return s.isGitlabCollapsibleSectionsEnabled
}

func (s *State) DoWithProxyStreamDataFormatting(f func()) {
	_ = s.doErrorWithProxyStreamDataFormatting(true, func() error {
		f()
		return nil
	})
}

func (s *State) DoErrorWithProxyStreamDataFormatting(f func() error) error {
	return s.doErrorWithProxyStreamDataFormatting(true, f)
}

func (s *State) DoWithoutProxyStreamDataFormatting(f func()) {
	_ = s.doErrorWithProxyStreamDataFormatting(false, func() error {
		f()
		return nil
	})
}

func (s *State) DoErrorWithoutProxyStreamDataFormatting(f func() error) error {
	return s.doErrorWithProxyStreamDataFormatting(false, f)
}

func (s *State) doErrorWithProxyStreamDataFormatting(value bool, f func() error) error {
	savedValue := s.isProxyStreamDataFormattingEnabled
	s.isProxyStreamDataFormattingEnabled = value
	err := f()
	s.isProxyStreamDataFormattingEnabled = savedValue

	return err
}

func (s *State) EnableProxyStreamDataFormatting() {
	s.isProxyStreamDataFormattingEnabled = true
}

func (s *State) DisableProxyStreamDataFormatting() {
	s.isProxyStreamDataFormattingEnabled = false
}

func (s *State) IsProxyStreamDataFormattingEnabled() bool {
	return s.isProxyStreamDataFormattingEnabled
}

func (s *State) EnableLineWrapping() {
	s.isLineWrappingEnabled = true
}

func (s *State) DisableLineWrapping() {
	s.isLineWrappingEnabled = false
}

func (s *State) IsLineWrappingEnabled() bool {
	return s.isLineWrappingEnabled
}

func (s *State) EnableStyle() {
	s.isStyleEnabled = true
}

func (s *State) DisableStyle() {
	s.isStyleEnabled = false
}

func (s *State) IsStyleEnabled() bool {
	return s.isStyleEnabled
}

func (s *State) processService() string {
	var result string

	result += s.formattedPrefix()
	result += s.formattedProcessBorders()
	result += s.formattedTag()

	return result
}

func (s *State) Width() int {
	return s.width
}

func (s *State) SetWidth(value int) {
	s.width = value
}

func (s *State) ContentWidth() int {
	return s.width - s.ServiceWidth()
}

func (s *State) ServiceWidth() int {
	return s.prefixWidth() + s.processBordersBlockWidth() + s.tagPartWidth() + s.indentWidth
}

func (s *State) DoWithIndent(f func()) {
	_ = s.DoErrorWithIndent(func() error {
		f()
		return nil
	})
}

func (s *State) DoErrorWithIndent(f func() error) error {
	s.IncreaseIndent()
	err := f()
	s.DecreaseIndent()

	return err
}

func (s *State) DoWithoutIndent(f func()) {
	_ = s.DoErrorWithoutIndent(func() error {
		f()
		return nil
	})
}

func (s *State) DoErrorWithoutIndent(f func() error) error {
	savedIndentWidth := s.indentWidth
	s.indentWidth = 0
	err := f()
	s.indentWidth = savedIndentWidth

	return err
}

func (s *State) IncreaseIndent() {
	s.indentWidth += 2
	s.DisableOptionalLn()
}

func (s *State) DecreaseIndent() {
	if s.indentWidth == 0 {
		return
	}

	s.indentWidth -= 2
	s.DisableOptionalLn()
}

func (s *State) ResetIndent() {
	s.indentWidth = 0
}

func (s *State) decorateByDoErrorWithIndent(f func() error) func() error {
	return func() error {
		return s.DoErrorWithIndent(f)
	}
}

func (s *State) EnableOptionalLn() {
	s.isOptionalLnEnabled = true
}

func (s *State) DisableOptionalLn() {
	s.isOptionalLnEnabled = false
}

type cursorState struct {
	isCursorOnNewLine              bool
	isPrevCursorStateOnRemoveCaret bool
}

func newCursorState() cursorState {
	return cursorState{
		isCursorOnNewLine:              true,
		isPrevCursorStateOnRemoveCaret: false,
	}
}

type processState struct {
	logProcessDownAndRightBorderSign     string
	logProcessVerticalBorderSign         string
	logProcessVerticalAndRightBorderSign string
	logProcessUpAndRightBorderSign       string
	processesBorderBetweenIndentWidth    int
	processesBorderIndentWidth           int
	processesBorderValues                []string
	processesBorderFormattedValues       []string
	activeLogProcesses                   []*logProcessDescriptor
	isGitlabCollapsibleSectionActive     bool
}

func newProcessState() processState {
	s := processState{}
	s.EnableLogProcessBorder()
	return s
}

func (s *processState) EnableLogProcessBorder() {
	s.logProcessDownAndRightBorderSign = "┌"
	s.logProcessVerticalBorderSign = "│"
	s.logProcessVerticalAndRightBorderSign = "├"
	s.logProcessUpAndRightBorderSign = "└"
	s.processesBorderBetweenIndentWidth = 1
	s.processesBorderIndentWidth = 1
}

func (s *processState) DisableLogProcessBorder() {
	s.logProcessDownAndRightBorderSign = ""
	s.logProcessVerticalBorderSign = "  "
	s.logProcessVerticalAndRightBorderSign = ""
	s.logProcessUpAndRightBorderSign = ""
	s.processesBorderIndentWidth = 0
	s.processesBorderBetweenIndentWidth = 0
}

func (s *State) decorateByWithExtraProcessBorder(colorlessBorder string, style *stylePkg.Style, decoratedFunc func() error) func() error {
	return func() error {
		return s.withExtraProcessBorder(colorlessBorder, style, decoratedFunc)
	}
}

func (s *State) withExtraProcessBorder(colorlessValue string, style *stylePkg.Style, decoratedFunc func() error) error {
	s.appendProcessBorder(colorlessValue, style)
	err := decoratedFunc()
	s.popProcessBorder()

	return err
}

func (s *State) decorateByWithoutLastProcessBorder(decoratedFunc func() error) func() error {
	return func() error {
		return s.withoutLastProcessBorder(decoratedFunc)
	}
}

func (s *State) withoutLastProcessBorder(f func() error) error {
	oldBorderValue := s.processesBorderValues[len(s.processesBorderValues)-1]
	s.processesBorderValues = s.processesBorderValues[:len(s.processesBorderValues)-1]

	oldBorderFormattedValue := s.processesBorderFormattedValues[len(s.processesBorderFormattedValues)-1]
	s.processesBorderFormattedValues = s.processesBorderFormattedValues[:len(s.processesBorderFormattedValues)-1]

	err := f()

	s.processesBorderValues = append(s.processesBorderValues, oldBorderValue)
	s.processesBorderFormattedValues = append(s.processesBorderFormattedValues, oldBorderFormattedValue)

	return err
}

func (s *State) appendProcessBorder(colorlessValue string, style *stylePkg.Style) {
	s.processesBorderValues = append(s.processesBorderValues, colorlessValue)
	s.processesBorderFormattedValues = append(s.processesBorderFormattedValues, s.formatWithStyle(style, colorlessValue))
}

func (s *State) popProcessBorder() {
	if len(s.processesBorderValues) == 0 {
		return
	}

	s.processesBorderValues = s.processesBorderValues[:len(s.processesBorderValues)-1]
	s.processesBorderFormattedValues = s.processesBorderFormattedValues[:len(s.processesBorderFormattedValues)-1]
}

func (s *State) formattedProcessBorders() string {
	if len(s.processesBorderValues) == 0 {
		return ""
	}

	return strings.Join(s.processesBorderFormattedValues, strings.Repeat(" ", s.processesBorderBetweenIndentWidth)) + strings.Repeat(" ", s.processesBorderIndentWidth)
}

func (s *State) processBordersBlockWidth() int {
	if len(s.processesBorderValues) == 0 {
		return 0
	}

	return len([]rune(strings.Join(s.processesBorderValues, strings.Repeat(" ", s.processesBorderBetweenIndentWidth)))) + s.processesBorderIndentWidth
}

type tagState struct {
	tagValue     string
	tagStyle     *stylePkg.Style
	tagPartWidth int
}

const tagIndentWidth = 2

func (s *State) DoWithTag(value string, style *stylePkg.Style, f func()) {
	_ = s.DoErrorWithTag(value, style, func() error {
		f()
		return nil
	})
}

func (s *State) DoErrorWithTag(value string, style *stylePkg.Style, f func() error) error {
	savedTag := s.tagValue
	savedStyle := s.tagStyle
	s.SetTagWithStyle(value, style)
	err := f()
	s.SetTagWithStyle(savedTag, savedStyle)

	return err
}

func (s *State) SetTag(value string) {
	s.tagValue = value
}

func (s *State) SetTagStyle(style *stylePkg.Style) {
	s.tagStyle = style
}

func (s *State) SetTagWithStyle(value string, style *stylePkg.Style) {
	s.SetTagStyle(style)
	s.SetTag(value)
}

func (s *State) ResetTag() {
	s.tagState = tagState{}
}

func (s *State) tagPartWidth() int {
	if s.tagValue != "" {
		return len(s.tagValue) + tagIndentWidth
	}

	return 0
}

func (s *State) formattedTag() string {
	if len(s.tagValue) == 0 {
		return ""
	}

	return strings.Join([]string{
		s.formatWithStyle(s.tagStyle, s.tagValue),
		strings.Repeat(" ", tagIndentWidth),
	}, "")
}

type prefixState struct {
	prefix                  string
	prefixStyle             *stylePkg.Style
	isPrefixWithTimeEnabled bool
	prefixTime              time.Time
}

func newPrefixState() prefixState {
	return prefixState{
		prefixTime: time.Now(),
	}
}

func (s *State) EnablePrefixWithTime() {
	s.isPrefixWithTimeEnabled = true
}

func (s *State) DisablePrefixWithTime() {
	s.isPrefixWithTimeEnabled = false
}

func (s *State) IsPrefixWithTimeEnabled() bool {
	return s.isPrefixWithTimeEnabled
}

func (s *State) ResetPrefixTime() {
	s.prefixTime = time.Now()
}

func (s *State) SetPrefix(value string) {
	s.ResetPrefix()
	s.prefix = value
}

func (s *State) SetPrefixStyle(style *stylePkg.Style) {
	s.prefixStyle = style
}

func (s *State) ResetPrefix() {
	s.prefix = ""
	s.prefixStyle = nil
	s.isPrefixWithTimeEnabled = false
}

func (s *State) formattedPrefix() string {
	if s.preparePrefixValue() == "" {
		return ""
	}

	return s.formatWithStyle(s.prefixStyle, s.preparePrefixValue())
}

func (s *State) preparePrefixValue() string {
	if s.isPrefixWithTimeEnabled {
		timeString := time.Since(s.prefixTime).String()
		timeStringRunes := []rune(timeString)
		if len(timeStringRunes) > 12 {
			timeString = string(timeStringRunes[:12])
		} else {
			timeString += strings.Repeat(" ", 12-len(timeStringRunes))
		}

		timeString += " "
		return timeString
	}

	return s.prefix
}

func (s *State) prefixWidth() int {
	return len([]rune(s.preparePrefixValue()))
}

func (s *State) processOptionalLn() string {
	var result string

	if s.isOptionalLnEnabled {
		result += s.processService()
		result += "\n"

		s.DisableOptionalLn()
		s.isCursorOnNewLine = true
	}

	return result
}

func (s *State) formatWithStyle(style *stylePkg.Style, format string, a ...interface{}) string {
	if !s.isStyleEnabled || style == nil {
		return stylePkg.SimpleFormat(format, a...)
	} else {
		return style.Colorize(format, a...)
	}
}
