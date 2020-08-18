package logger

import (
	"io"

	"github.com/werf/logboek/internal/stream"
	"github.com/werf/logboek/pkg/level"
	"github.com/werf/logboek/pkg/types"
)

type Logger struct {
	*Manager

	acceptedLevel     level.Level
	levelManager      map[level.Level]*Manager
	outStream         *stream.Stream
	errStream         *stream.Stream
	commonStreamState *stream.State
}

func NewLogger(outStream, errStream io.Writer) *Logger {
	l := &Logger{}

	l.commonStreamState = stream.NewStreamState()
	l.outStream = stream.NewStream(outStream, l.commonStreamState)
	l.errStream = stream.NewStream(errStream, l.commonStreamState)
	l.initLevelManager()

	return l
}

func (l *Logger) initLevelManager() {
	l.levelManager = make(map[level.Level]*Manager, len(level.List))
	for _, lvl := range level.List {
		l.levelManager[lvl] = NewManager(l, lvl)
	}

	l.Manager = l.levelManager[level.Default]
}

func (l *Logger) setCommonStreamState(state *stream.State) {
	l.commonStreamState = state
	l.outStream.State = state
	l.errStream.State = state
}

func (l *Logger) getLevelManager(lvl level.Level) *Manager {
	return l.levelManager[lvl]
}

func (l *Logger) GetLevelStream(lvl level.Level) *stream.Stream {
	if lvl == level.Error || lvl == level.Warn {
		return l.errStream
	} else {
		return l.outStream
	}
}

func (l *Logger) Error() types.ManagerInterface {
	return l.getLevelManager(level.Error)
}

func (l *Logger) Warn() types.ManagerInterface {
	return l.getLevelManager(level.Warn)
}

func (l *Logger) Default() types.ManagerInterface {
	return l.getLevelManager(level.Default)
}

func (l *Logger) Info() types.ManagerInterface {
	return l.getLevelManager(level.Info)
}

func (l *Logger) Debug() types.ManagerInterface {
	return l.getLevelManager(level.Debug)
}

func (l *Logger) AcceptedLevel() level.Level {
	return l.acceptedLevel
}

func (l *Logger) SetAcceptedLevel(lvl level.Level) {
	l.acceptedLevel = lvl
}

func (l *Logger) IsAcceptedLevel(lvl level.Level) bool {
	return l.levelManager[lvl].IsAccepted()
}

func (l *Logger) Streams() types.StreamsInterface {
	return l.commonStreamState
}

func (l *Logger) FitText(text string, options types.FitTextOptions) string {
	return l.outStream.FitText(text, options)
}

func (l *Logger) ProxyOutStream() io.Writer {
	return l.outStream.ProxyStream()
}

func (l *Logger) ProxyErrStream() io.Writer {
	return l.errStream.ProxyStream()
}

func (l *Logger) NewSubLogger(outStream, errStream io.Writer) types.LoggerInterface {
	subLogger := NewLogger(outStream, errStream)
	subLogger.setCommonStreamState(l.commonStreamState.SubState())
	subLogger.SetAcceptedLevel(l.acceptedLevel)

	for lvl, manager := range l.levelManager {
		subLogger.levelManager[lvl].style = manager.style
	}

	return subLogger
}

func (l *Logger) GetStreamsSettingsFrom(l2 types.LoggerInterface) {
	l.setCommonStreamState(l2.(*Logger).commonStreamState.SharedState())
}

func (l *Logger) Reset() {
	l.outStream.Reset()
}
