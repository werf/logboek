package main

import (
	"C"

	"github.com/werf/logboek"
	"github.com/werf/logboek/pkg/types"
)

//go:generate go build -o ../_logboek.so -buildmode=c-shared github.com/werf/logboek/c_lib

var processes []types.LogProcessInterface

//export Init
func Init() *C.char {
	return nil
}

//export DisablePrettyLog
func DisablePrettyLog() {
	logboek.DefaultLogger().Streams().DisablePrettyLog()
}

//export EnableFitMode
func EnableFitMode() {
	logboek.DefaultLogger().Streams().EnableLineWrapping()
}

//export DisableFitMode
func DisableFitMode() {
	logboek.DefaultLogger().Streams().DisableLineWrapping()
}

//export EnableLogColor
func EnableLogColor() {
	logboek.DefaultLogger().Streams().EnableStyle()
}

//export DisableLogColor
func DisableLogColor() {
	logboek.DefaultLogger().Streams().DisableStyle()
}

//export SetTerminalWidth
func SetTerminalWidth(width C.int) {
	logboek.DefaultLogger().Streams().SetWidth(int(width))
}

//export IndentUp
func IndentUp() {
	logboek.DefaultLogger().Streams().IncreaseIndent()
}

//export IndentDown
func IndentDown() {
	logboek.DefaultLogger().Streams().DecreaseIndent()
}

//export OptionalLnModeOn
func OptionalLnModeOn() {
	logboek.DefaultLogger().LogOptionalLn()
}

//export Log
func Log(data *C.char) {
	logboek.DefaultLogger().Default().LogF("%s", C.GoString(data))
}

//export LogHighlight
func LogHighlight(data *C.char) {
	logboek.DefaultLogger().Default().LogFHighlight("%s", C.GoString(data))
}

//export LogService
func LogService(data *C.char) {
	logboek.Default().LogF("%s", C.GoString(data))
}

//export LogInfo
func LogInfo(data *C.char) {
	logboek.DefaultLogger().Default().LogFDetails("%s", C.GoString(data))
}

//export LogError
func LogError(data *C.char) {
	logboek.Warn().LogF("%s", C.GoString(data))
}

//export LogProcessStart
func LogProcessStart(msg *C.char) {
	processes = append(processes, logboek.DefaultLogger().Default().LogProcess(C.GoString(msg)))
	processes[len(processes)-1].Start()
}

//export LogProcessEnd
func LogProcessEnd(withoutLogOptionalLn bool) {
	processes[len(processes)-1].Options(func(options types.LogProcessOptionsInterface) {
		if withoutLogOptionalLn {
			options.WithoutLogOptionalLn()
		}
	}).End()
	processes = processes[:len(processes)-1]
}

//export LogProcessStepEnd
func LogProcessStepEnd(msg *C.char) {
	processes[len(processes)-1].StepEnd(C.GoString(msg))
	processes = processes[:len(processes)-1]
}

//export LogProcessFail
func LogProcessFail(withoutLogOptionalLn bool) {
	processes[len(processes)-1].Options(func(options types.LogProcessOptionsInterface) {
		if withoutLogOptionalLn {
			options.WithoutLogOptionalLn()
		}
	}).Fail()
	processes = processes[:len(processes)-1]
}

//export FitText
func FitText(text *C.char, extraIndentWidth, maxWidth int, markWrappedFile bool) *C.char {
	return C.CString(logboek.FitText(C.GoString(text), types.FitTextOptions{
		ExtraIndentWidth: extraIndentWidth,
		MaxWidth:         maxWidth,
		MarkWrappedLine:  markWrappedFile,
	}))
}

//export GetRawStreamsOutputMode
func GetRawStreamsOutputMode() bool {
	return logboek.DefaultLogger().Streams().IsProxyStreamDataFormattingEnabled()
}

//export RawStreamsOutputModeOn
func RawStreamsOutputModeOn() {
	logboek.DefaultLogger().Streams().DisableProxyStreamDataFormatting()
}

//export RawStreamsOutputModeOff
func RawStreamsOutputModeOff() {
	logboek.DefaultLogger().Streams().EnableProxyStreamDataFormatting()
}

//export MuteOut
func MuteOut() {
	logboek.DefaultLogger().Streams().Mute()
}

//export UnmuteOut
func UnmuteOut() {
	logboek.DefaultLogger().Streams().Unmute()
}

//export MuteErr
func MuteErr() {
	logboek.DefaultLogger().Streams().Mute()
}

//export UnmuteErr
func UnmuteErr() {
	logboek.DefaultLogger().Streams().Unmute()
}

//export Out
func Out(msg *C.char) {
	logboek.DefaultLogger().Default().LogF("%s", C.GoString(msg))
}

//export Err
func Err(msg *C.char) {
	logboek.DefaultLogger().Error().LogF("%s", C.GoString(msg))
}

func main() {}
