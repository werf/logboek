package stream

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/werf/logboek/internal/stream/fitter"
	stylePkg "github.com/werf/logboek/pkg/style"
	"github.com/werf/logboek/pkg/types"
)

const defaultWidth = 140

type Stream struct {
	io.Writer
	*StateAndModes
}

func NewStream(w io.Writer, state *StateAndModes) *Stream {
	s := &Stream{
		Writer:        w,
		StateAndModes: state,
	}
	s.initWidth()
	return s
}

func (s *Stream) initWidth() {
	f, ok := s.Writer.(*os.File)
	if ok && terminal.IsTerminal(int(f.Fd())) {
		width, _, err := terminal.GetSize(int(f.Fd()))
		if err != nil {
			panic(fmt.Sprintf("get terminal size failed: %s", err))
		}

		if width != 0 {
			s.width = width
			return
		}
	}

	s.width = defaultWidth
}

func (s *Stream) FitText(text string, options types.FitTextOptions) string {
	var lineWidth int
	if options.Width != 0 {
		lineWidth = options.Width
	} else {
		tw := s.width
		if options.MaxWidth != 0 && tw > options.MaxWidth {
			lineWidth = options.MaxWidth
		} else {
			lineWidth = tw
		}

		lineWidth -= s.ServiceWidth()
	}

	return fitTextWithIndent(text, lineWidth, options.ExtraIndentWidth, options.MarkWrappedLine)
}

func fitTextWithIndent(text string, lineWidth, extraIndentWidth int, markWrappedLine bool) string {
	var result string
	var resultLines []string

	contentWidth := lineWidth - extraIndentWidth

	fittedText := fitter.FitText(text, &fitter.State{}, contentWidth, markWrappedLine, false)
	for _, line := range strings.Split(fittedText, "\n") {
		indent := strings.Repeat(" ", extraIndentWidth)
		resultLines = append(resultLines, strings.Join([]string{indent, line}, ""))
	}

	result = strings.Join(resultLines, "\n")

	return result
}

func (s *Stream) ProxyStream() io.Writer {
	return proxyStream{s}
}

func (s *Stream) FormatAndLogF(style *stylePkg.Style, format string, a ...interface{}) {
	s.StateAndModes.mutex.Lock()
	defer s.StateAndModes.mutex.Unlock()

	msg := s.formatWithStyle(style, format, a...)

	if s.IsLineWrappingEnabled() {
		msg = s.FitText(msg, types.FitTextOptions{MarkWrappedLine: true})
	}

	s.processAndLogF(msg)
}

func (s *Stream) processAndLogLn(a ...interface{}) {
	s.processAndLogF(fmt.Sprintln(a...))
}

func (s *Stream) processAndLogF(format string, a ...interface{}) {
	_, err := s.processAndLogFBase(format, a...)
	if err != nil {
		panic(err)
	}
}

func (s *Stream) processAndLogFBase(format string, a ...interface{}) (int, error) {
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
			formattedMsg += s.processNewLineAndRemoveCarriage(string(r))
		default:
			formattedMsg += s.processDefault()
		}

		formattedMsg += string(r)
	}

	return s.logFBase("%s", formattedMsg)
}

func (s *Stream) processNewLineAndRemoveCarriage(carriage string) string {
	var result string

	if s.isCursorOnNewLine && !s.isPrevCursorStateOnRemoveCaret {
		result += s.processService()
	}

	s.isPrevCursorStateOnRemoveCaret = carriage == "\r"
	s.isCursorOnNewLine = true

	return result
}

func (s *Stream) processDefault() string {
	var result string

	result += s.processOptionalLn()

	if s.isCursorOnNewLine {
		result += s.processService()
		result += strings.Repeat(" ", s.indentWidth)

		s.isCursorOnNewLine = false
	}

	s.isPrevCursorStateOnRemoveCaret = false

	return result
}

func (s *Stream) logFBase(format string, a ...interface{}) (int, error) {
	return fmt.Fprintf(s.Writer, format, a...)
}

func (s *Stream) applyOptionalLn() {
	_, _ = s.logFBase(s.processOptionalLn())
}

func (s *Stream) Reset() {
	s.ResetState()
	s.ResetModes()
}

func (s *Stream) ResetState() {
	s.endAllActiveProcesses()
	s.StateAndModes.resetState()
}

func (s *Stream) ResetModes() {
	s.StateAndModes.resetModes()
}

func splitData(data []byte, chunkSize int) [][]rune {
	buf := bytes.Runes(data)
	chunks := make([][]rune, 0, len(buf)/chunkSize+1)

	for len(buf) >= chunkSize {
		var chunk []rune
		chunk, buf = buf[:chunkSize], buf[chunkSize:]
		chunks = append(chunks, chunk)
	}

	if len(buf) > 0 {
		chunks = append(chunks, buf)
	}

	return chunks
}
