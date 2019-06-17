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

	streamsFitterState fitterState
)

type WriterProxy struct {
	io.Writer
}

func (p WriterProxy) Write(data []byte) (int, error) {
	msg := string(data)

	if isRawStreamsOutputModeOn {
		return logFBase(p.Writer, "%s", msg)
	}

	msg, streamsFitterState = fitText(msg, streamsFitterState, ContentWidth(), true, true)

	_, err := processAndLogFBase(p.Writer, "%s", msg)
	return len(data), err
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
	return fmt.Fprintf(GetOutStream(), format, a...)
}

func ErrF(format string, a ...interface{}) (int, error) {
	return fmt.Fprintf(GetErrStream(), format, a...)
}
