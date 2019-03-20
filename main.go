package logboek

func Init() error {
	return initTerminalWidth()
}

func DisablePrettyLog() {
	RawStreamsOutputModeOn()
	disableLogProcessBorder()
}
