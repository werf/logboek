package level

type Level int

const (
	Error Level = iota - 2
	Warn
	Default // 0
	Info
	Debug
)

var List = []Level{Error, Warn, Default, Info, Debug}
