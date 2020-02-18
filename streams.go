package logboek

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

var (
	outStream io.Writer = os.Stdout
	errStream io.Writer = os.Stderr

	isRawStreamsOutputModeOn = false
	isFitModeOn              = false

	streamsFitterState fitterState
)

type WriterProxy struct {
	io.Writer
}

func (p WriterProxy) Write(data []byte) (int, error) {
	if !streamsLogLevel.IsAccepted() {
		return 0, nil
	}

	msg := string(data)

	if isRawStreamsOutputModeOn {
		return logFBase(p.Writer, "%s", msg)
	}

	if isFitModeOn {
		msg, streamsFitterState = fitText(msg, streamsFitterState, ContentWidth(), true, true)
	}

	_, err := processAndLogFBase(p.Writer, "%s", msg)
	return len(data), err
}

func SetStreamsLogLevel(logLevel Level) {
	streamsLogLevel = logLevel
}

func WithStreamsLogLevel(logLevel Level, f func() error) error {
	savedStreamsLogLevel := streamsLogLevel
	streamsLogLevel = logLevel
	err := f()
	streamsLogLevel = savedStreamsLogLevel

	return err
}

func WithRawStreamsOutputModeOn(f func() error) error {
	savedIsRawOutputModeOn := isRawStreamsOutputModeOn
	isRawStreamsOutputModeOn = true
	err := f()
	isRawStreamsOutputModeOn = savedIsRawOutputModeOn

	return err
}

func GetRawStreamsOutputMode() bool {
	return isRawStreamsOutputModeOn
}

func RawStreamsOutputModeOn() {
	isRawStreamsOutputModeOn = true
}

func RawStreamsOutputModeOff() {
	isRawStreamsOutputModeOn = false
}

func GetOutStream() io.Writer {
	return WriterProxy{outStream}
}

func GetErrStream() io.Writer {
	return WriterProxy{errStream}
}

func MuteOut() {
	outStream = ioutil.Discard
}

func UnmuteOut() {
	outStream = os.Stdout
}

func MuteErr() {
	errStream = ioutil.Discard
}

func UnmuteErr() {
	errStream = os.Stderr
}

func OutF(format string, a ...interface{}) (int, error) {
	if !streamsLogLevel.IsAccepted() {
		return 0, nil
	}
	return fmt.Fprintf(GetOutStream(), format, a...)
}

func ErrF(format string, a ...interface{}) (int, error) {
	if !streamsLogLevel.IsAccepted() {
		return 0, nil
	}
	return fmt.Fprintf(GetErrStream(), format, a...)
}
