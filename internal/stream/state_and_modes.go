package stream

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gookit/color"

	"github.com/werf/logboek/internal/stream/fitter"
	stylePkg "github.com/werf/logboek/pkg/style"
)

type StateAndModes struct {
	copyable
	mutex sync.Mutex
}

type copyable struct {
	width int

	modes
	baseState
	fitter.State
	cursorState
	processState
	tagState
	prefixState
}

func NewStreamState() *StateAndModes {
	s := &StateAndModes{}
	s.initModes()
	s.initState()
	return s
}

func (s *StateAndModes) initModes() {
	s.modes = newModes()
}

func (s *StateAndModes) initState() {
	s.baseState = newBaseState()
	s.State = fitter.NewState()
	s.cursorState = newCursorState()
	s.processState = newProcessState()
	s.prefixState = newPrefixState()
}

func (s *StateAndModes) reset() {
	s.resetState()
	s.resetModes()
}

func (s *StateAndModes) resetState() {
	s.initState()
}

func (s *StateAndModes) resetModes() {
	s.initModes()
}

func (s *StateAndModes) SubState() *StateAndModes {
	ss := s.SharedState()
	ss.width = s.ContentWidth()
	return ss
}

func (s *StateAndModes) SharedState() *StateAndModes {
	ss := s.clone()
	ss.isOptionalLnEnabled.Store(false)
	ss.State = fitter.State{}
	ss.cursorState = newCursorState()
	ss.processState = newProcessState()
	return ss
}

func (s *StateAndModes) DisablePrettyLog() {
	s.DisableProxyStreamDataFormatting()
	s.DisableLogProcessBorder()
	s.DisableLineWrapping()
}

type baseState struct {
	indentWidth         int
	isOptionalLnEnabled atomic.Bool
}

func newBaseState() baseState {
	return baseState{}
}

type modes struct {
	isMuted                            bool
	isStyleEnabled                     bool
	isLineWrappingEnabled              bool
	isProxyStreamDataFormattingEnabled bool
	isGitlabCollapsibleSectionsEnabled bool
	isPrefixDurationEnabled            bool
	isPrefixTimeEnabled                bool
	isLogProcessBorderEnabled          bool
}

func newModes() modes {
	return modes{
		isStyleEnabled:                     true,
		isLineWrappingEnabled:              true,
		isProxyStreamDataFormattingEnabled: true,
		isGitlabCollapsibleSectionsEnabled: os.Getenv("GITLAB_CI") == "true",
		isLogProcessBorderEnabled:          true,
	}
}

func (s *StateAndModes) Mute() {
	s.isMuted = true
}

func (s *StateAndModes) Unmute() {
	s.isMuted = false
}

func (s *StateAndModes) IsMuted() bool {
	return s.isMuted
}

func (s *StateAndModes) EnableGitlabCollapsibleSections() {
	s.isGitlabCollapsibleSectionsEnabled = true
}

func (s *StateAndModes) DisableGitlabCollapsibleSections() {
	s.isGitlabCollapsibleSectionsEnabled = false
}

func (s *StateAndModes) IsGitlabCollapsibleSections() bool {
	return s.isGitlabCollapsibleSectionsEnabled
}

func (s *StateAndModes) EnableLogProcessBorder() {
	s.isLogProcessBorderEnabled = true
}

func (s *StateAndModes) DisableLogProcessBorder() {
	s.isLogProcessBorderEnabled = true
}

func (s *StateAndModes) IsLogProcessBorderEnabled() bool {
	return s.isLogProcessBorderEnabled
}

func (s *StateAndModes) DoWithProxyStreamDataFormatting(f func()) {
	_ = s.doErrorWithProxyStreamDataFormatting(true, func() error {
		f()
		return nil
	})
}

func (s *StateAndModes) DoErrorWithProxyStreamDataFormatting(f func() error) error {
	return s.doErrorWithProxyStreamDataFormatting(true, f)
}

func (s *StateAndModes) DoWithoutProxyStreamDataFormatting(f func()) {
	_ = s.doErrorWithProxyStreamDataFormatting(false, func() error {
		f()
		return nil
	})
}

func (s *StateAndModes) DoErrorWithoutProxyStreamDataFormatting(f func() error) error {
	return s.doErrorWithProxyStreamDataFormatting(false, f)
}

func (s *StateAndModes) doErrorWithProxyStreamDataFormatting(value bool, f func() error) error {
	savedValue := s.isProxyStreamDataFormattingEnabled
	s.isProxyStreamDataFormattingEnabled = value
	err := f()
	s.isProxyStreamDataFormattingEnabled = savedValue

	return err
}

func (s *StateAndModes) EnableProxyStreamDataFormatting() {
	s.isProxyStreamDataFormattingEnabled = true
}

func (s *StateAndModes) DisableProxyStreamDataFormatting() {
	s.isProxyStreamDataFormattingEnabled = false
}

func (s *StateAndModes) IsProxyStreamDataFormattingEnabled() bool {
	return s.isProxyStreamDataFormattingEnabled
}

func (s *StateAndModes) EnableLineWrapping() {
	s.isLineWrappingEnabled = true
}

func (s *StateAndModes) DisableLineWrapping() {
	s.isLineWrappingEnabled = false
}

func (s *StateAndModes) IsLineWrappingEnabled() bool {
	return s.isLineWrappingEnabled
}

func (s *StateAndModes) EnableStyle() {
	color.ForceColor()
	s.isStyleEnabled = true
}

func (s *StateAndModes) DisableStyle() {
	s.isStyleEnabled = false
}

func (s *StateAndModes) IsStyleEnabled() bool {
	return s.isStyleEnabled
}

func (s *StateAndModes) processService() string {
	var result string

	result += s.formattedPrefix()
	result += s.formattedProcessBorders()
	result += s.formattedTag()

	return result
}

func (s *StateAndModes) Width() int {
	return s.width
}

func (s *StateAndModes) SetWidth(value int) {
	s.width = value
}

func (s *StateAndModes) ContentWidth() int {
	return s.width - s.ServiceWidth()
}

func (s *StateAndModes) ServiceWidth() int {
	return s.prefixWidth() + s.processBordersBlockWidth() + s.tagPartWidth() + s.indentWidth
}

func (s *StateAndModes) DoWithIndent(f func()) {
	_ = s.DoErrorWithIndent(func() error {
		f()
		return nil
	})
}

func (s *StateAndModes) DoErrorWithIndent(f func() error) error {
	s.IncreaseIndent()
	err := f()
	s.DecreaseIndent()

	return err
}

func (s *StateAndModes) DoWithoutIndent(f func()) {
	_ = s.DoErrorWithoutIndent(func() error {
		f()
		return nil
	})
}

func (s *StateAndModes) DoErrorWithoutIndent(f func() error) error {
	savedIndentWidth := s.indentWidth
	s.indentWidth = 0
	err := f()
	s.indentWidth = savedIndentWidth

	return err
}

func (s *StateAndModes) IncreaseIndent() {
	s.indentWidth += 2
	s.DisableOptionalLn()
}

func (s *StateAndModes) DecreaseIndent() {
	if s.indentWidth == 0 {
		return
	}

	s.indentWidth -= 2
	s.DisableOptionalLn()
}

func (s *StateAndModes) ResetIndent() {
	s.indentWidth = 0
}

func (s *StateAndModes) decorateByDoErrorWithIndent(f func() error) func() error {
	return func() error {
		return s.DoErrorWithIndent(f)
	}
}

func (s *StateAndModes) EnableOptionalLn() {
	s.isOptionalLnEnabled.Store(true)
}

func (s *StateAndModes) DisableOptionalLn() {
	s.isOptionalLnEnabled.Store(false)
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
	processesBorderValues          []string
	processesBorderFormattedValues []string
	activeLogProcesses             []*logProcessDescriptor
}

func newProcessState() processState {
	return processState{}
}

func (s *StateAndModes) LogProcessDownAndRightBorderSign() string {
	if s.isLogProcessBorderEnabled {
		return "┌"
	}

	return ""
}

func (s *StateAndModes) LogProcessVerticalBorderSign() string {
	if s.isLogProcessBorderEnabled {
		return "│"
	}

	return ""
}

func (s *StateAndModes) LogProcessVerticalAndRightBorderSign() string {
	if s.isLogProcessBorderEnabled {
		return "├"
	}

	return ""
}

func (s *StateAndModes) LogProcessUpAndRightBorderSign() string {
	if s.isLogProcessBorderEnabled {
		return "└"
	}

	return ""
}

func (s *StateAndModes) ProcessesBorderBetweenIndentWidth() int {
	if s.isLogProcessBorderEnabled {
		return 1
	}

	return 0
}

func (s *StateAndModes) ProcessesBorderIndentWidth() int {
	if s.isLogProcessBorderEnabled {
		return 1
	}

	return 0
}

func (s *StateAndModes) decorateByWithExtraProcessBorder(colorlessBorder string, style color.Style, decoratedFunc func() error) func() error {
	return func() error {
		return s.withExtraProcessBorder(colorlessBorder, style, decoratedFunc)
	}
}

func (s *StateAndModes) withExtraProcessBorder(colorlessValue string, style color.Style, decoratedFunc func() error) error {
	s.appendProcessBorder(colorlessValue, style)
	err := decoratedFunc()
	s.popProcessBorder()

	return err
}

func (s *StateAndModes) decorateByWithoutLastProcessBorder(decoratedFunc func() error) func() error {
	return func() error {
		return s.withoutLastProcessBorder(decoratedFunc)
	}
}

func (s *StateAndModes) withoutLastProcessBorder(f func() error) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	oldBorderValue := s.processesBorderValues[len(s.processesBorderValues)-1]
	s.processesBorderValues = s.processesBorderValues[:len(s.processesBorderValues)-1]

	oldBorderFormattedValue := s.processesBorderFormattedValues[len(s.processesBorderFormattedValues)-1]
	s.processesBorderFormattedValues = s.processesBorderFormattedValues[:len(s.processesBorderFormattedValues)-1]
	err := f()

	s.processesBorderValues = append(s.processesBorderValues, oldBorderValue)
	s.processesBorderFormattedValues = append(s.processesBorderFormattedValues, oldBorderFormattedValue)

	return err
}

func (s *StateAndModes) appendProcessBorder(colorlessValue string, style color.Style) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.processesBorderValues = append(s.processesBorderValues, colorlessValue)
	s.processesBorderFormattedValues = append(s.processesBorderFormattedValues, s.FormatWithStyle(style, colorlessValue))
}

func (s *StateAndModes) popProcessBorder() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if len(s.processesBorderValues) == 0 {
		return
	}

	s.processesBorderValues = s.processesBorderValues[:len(s.processesBorderValues)-1]
	s.processesBorderFormattedValues = s.processesBorderFormattedValues[:len(s.processesBorderFormattedValues)-1]
}

func (s *StateAndModes) formattedProcessBorders() string {
	if len(s.processesBorderValues) == 0 {
		return ""
	}

	return strings.Join(s.processesBorderFormattedValues, strings.Repeat(" ", s.ProcessesBorderBetweenIndentWidth())) + strings.Repeat(" ", s.ProcessesBorderIndentWidth())
}

func (s *StateAndModes) processBordersBlockWidth() int {
	if len(s.processesBorderValues) == 0 {
		return 0
	}

	return len([]rune(strings.Join(s.processesBorderValues, strings.Repeat(" ", s.ProcessesBorderBetweenIndentWidth())))) + s.ProcessesBorderIndentWidth()
}

type tagState struct {
	tagValue     string
	tagStyle     color.Style
	tagPartWidth int
}

const tagIndentWidth = 2

func (s *StateAndModes) DoWithTag(value string, style color.Style, f func()) {
	_ = s.DoErrorWithTag(value, style, func() error {
		f()
		return nil
	})
}

func (s *StateAndModes) DoErrorWithTag(value string, style color.Style, f func() error) error {
	savedTag := s.tagValue
	savedStyle := s.tagStyle
	s.SetTagWithStyle(value, style)
	err := f()
	s.SetTagWithStyle(savedTag, savedStyle)

	return err
}

func (s *StateAndModes) SetTag(value string) {
	s.tagValue = value
}

func (s *StateAndModes) SetTagStyle(style color.Style) {
	s.tagStyle = style
}

func (s *StateAndModes) SetTagWithStyle(value string, style color.Style) {
	s.SetTagStyle(style)
	s.SetTag(value)
}

func (s *StateAndModes) ResetTag() {
	s.tagState = tagState{}
}

func (s *StateAndModes) tagPartWidth() int {
	if s.tagValue != "" {
		return len(s.tagValue) + tagIndentWidth
	}

	return 0
}

func (s *StateAndModes) formattedTag() string {
	if len(s.tagValue) == 0 {
		return ""
	}

	return strings.Join([]string{
		s.FormatWithStyle(s.tagStyle, s.tagValue),
		strings.Repeat(" ", tagIndentWidth),
	}, "")
}

type prefixState struct {
	prefix                  string
	prefixStyle             color.Style
	prefixDurationStartTime time.Time
	prefixTimeFormat        string
}

func newPrefixState() prefixState {
	return prefixState{
		prefixDurationStartTime: time.Now(),
		prefixTimeFormat:        time.RFC3339,
	}
}

func (s *StateAndModes) EnablePrefixDuration() {
	s.disablePrefix()
	s.isPrefixDurationEnabled = true
}

func (s *StateAndModes) IsPrefixDurationEnabled() bool {
	return s.isPrefixDurationEnabled
}

func (s *StateAndModes) DisablePrefixDuration() {
	s.isPrefixDurationEnabled = false
}

func (s *StateAndModes) ResetPrefixDurationStartTime() {
	s.prefixDurationStartTime = time.Now()
}

func (s *StateAndModes) SetPrefixTimeFormat(format string) {
	s.prefixTimeFormat = format
}

func (s *StateAndModes) EnablePrefixTime() {
	s.disablePrefix()
	s.isPrefixTimeEnabled = true
}

func (s *StateAndModes) IsPrefixTimeEnabled() bool {
	return s.isPrefixTimeEnabled
}

func (s *StateAndModes) DisablePrefixTime() {
	s.isPrefixTimeEnabled = false
}

func (s *StateAndModes) SetPrefix(value string) {
	s.disablePrefix()
	s.prefix = value
}

func (s *StateAndModes) SetPrefixStyle(style color.Style) {
	s.prefixStyle = style
}

func (s *StateAndModes) DisablePrefix() {
	s.disablePrefix()
}

func (s *StateAndModes) disablePrefix() {
	s.prefix = ""
	s.isPrefixDurationEnabled = false
	s.isPrefixTimeEnabled = false
}

func (s *StateAndModes) formattedPrefix() string {
	if s.preparePrefixValue() == "" {
		return ""
	}

	return s.FormatWithStyle(s.prefixStyle, s.preparePrefixValue())
}

func (s *StateAndModes) preparePrefixValue() string {
	switch {
	case s.isPrefixDurationEnabled:
		timeString := time.Since(s.prefixDurationStartTime).String()
		timeStringRunes := []rune(timeString)
		if len(timeStringRunes) > 12 {
			timeString = string(timeStringRunes[:12])
		} else {
			timeString += strings.Repeat(" ", 12-len(timeStringRunes))
		}

		timeString += " "
		return timeString
	case s.isPrefixTimeEnabled:
		return time.Now().Format(s.prefixTimeFormat) + " "
	default:
		return s.prefix
	}
}

func (s *StateAndModes) prefixWidth() int {
	return len([]rune(s.preparePrefixValue()))
}

func (s *StateAndModes) processOptionalLn() string {
	var result string

	if s.isOptionalLnEnabled.Load() {
		result += s.processService()
		result += "\n"

		s.DisableOptionalLn()
		s.isCursorOnNewLine = true
	}

	return result
}

func (s *StateAndModes) FormatWithStyle(style color.Style, format string, a ...interface{}) string {
	if !s.isStyleEnabled || style == nil {
		style = stylePkg.None()
	}

	var resultLines []string
	for _, line := range strings.Split(fmt.Sprintf(format, a...), "\n") {
		if line == "" {
			resultLines = append(resultLines, line)
		} else {
			resultLines = append(resultLines, style.Sprint(line))
		}
	}

	return strings.Join(resultLines, "\n")
}

func (s *StateAndModes) clone() *StateAndModes {
	sClone := &StateAndModes{}
	sClone.copyable = s.copyable.clone()
	return sClone
}

func (s *copyable) clone() copyable {
	sClone := *s
	return sClone
}
