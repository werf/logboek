package stream

import (
	stylePkg "github.com/werf/logboek/pkg/style"
	"github.com/werf/logboek/pkg/types"
)

type LogBlock struct {
	manager    types.ManagerInterface
	stream     *Stream
	title      string
	isDisabled bool
	options    *LogBlockOptions
}

func (l *LogBlock) Disable() types.LogBlockInterface {
	l.isDisabled = true
	return l
}

func (l *LogBlock) Enable() types.LogBlockInterface {
	l.isDisabled = false
	return l
}

func (l *LogBlock) Options(f func(options types.LogBlockOptionsInterface)) types.LogBlockInterface {
	f(l.options)
	return l
}

func (l *LogBlock) Do(f func()) {
	_ = l.DoError(func() error {
		f()
		return nil
	})
}

func (l *LogBlock) DoError(f func() error) error {
	if l.isDisabled {
		return nil
	} else if !l.manager.IsAccepted() {
		if l.options.disableIfLevelNotAccepted {
			return nil
		} else {
			return f()
		}
	} else if l.stream.IsMuted() {
		return f()
	}

	return l.stream.logBlock(l.title, l.options, f)
}

type LogBlockOptions struct {
	disableIfLevelNotAccepted bool
	withIndent                bool
	withoutLogOptionalLn      bool
	style                     *stylePkg.Style
}

func (opts *LogBlockOptions) DisableIfLevelNotAccepted() {
	opts.disableIfLevelNotAccepted = true
}

func (opts *LogBlockOptions) WithIndent() {
	opts.withIndent = true
}

func (opts *LogBlockOptions) WithoutLogOptionalLn() {
	opts.withoutLogOptionalLn = true
}

func (opts *LogBlockOptions) Style(s *stylePkg.Style) {
	opts.style = s
}

type LogProcessInline struct {
	manager    types.ManagerInterface
	title      string
	options    *LogProcessInlineOptions
	isDisabled bool
	stream     *Stream
}

func (l *LogProcessInline) Disable() types.LogProcessInlineInterface {
	l.isDisabled = true
	return l
}

func (l *LogProcessInline) Enable() types.LogProcessInlineInterface {
	l.isDisabled = false
	return l
}

func (l *LogProcessInline) Options(f func(options types.LogProcessInlineOptionsInterface)) types.LogProcessInlineInterface {
	f(l.options)
	return l
}

func (l *LogProcessInline) Do(f func()) {
	_ = l.DoError(func() error {
		f()
		return nil
	})
}

func (l *LogProcessInline) DoError(f func() error) error {
	if l.isDisabled {
		return nil
	} else if !l.manager.IsAccepted() {
		if l.options.disableIfLevelNotAccepted {
			return nil
		} else {
			return f()
		}
	} else if l.stream.IsMuted() {
		return f()
	}

	return l.stream.logProcessInline(l.title, l.options, f)
}

type LogProcessInlineOptions struct {
	disableIfLevelNotAccepted bool
	style                     *stylePkg.Style
}

func (opts *LogProcessInlineOptions) DisableIfLevelNotAccepted() {
	opts.disableIfLevelNotAccepted = true
}

func (opts *LogProcessInlineOptions) Style(s *stylePkg.Style) {
	opts.style = s
}

type LogProcess struct {
	manager    types.ManagerInterface
	title      string
	options    *LogProcessOptions
	isDisabled bool
	stream     *Stream
	isStarted  bool
	isLaunched bool
}

func (l *LogProcess) Disable() types.LogProcessInterface {
	l.isDisabled = true
	return l
}

func (l *LogProcess) Enable() types.LogProcessInterface {
	l.isDisabled = false
	return l
}

func (l *LogProcess) Options(f func(options types.LogProcessOptionsInterface)) types.LogProcessInterface {
	f(l.options)
	return l
}

func (l *LogProcess) Do(f func()) {
	_ = l.DoError(func() error {
		f()
		return nil
	})
}

func (l *LogProcess) DoError(f func() error) error {
	if l.isStarted {
		panic("process has been already started")
	} else if l.isLaunched {
		panic("process has been already launched")
	}

	l.isLaunched = true

	if l.isDisabled {
		return nil
	} else if !l.manager.IsAccepted() {
		if l.options.disableIfLevelNotAccepted {
			return nil
		} else {
			return f()
		}
	} else if l.stream.IsMuted() {
		return f()
	}

	return l.stream.logProcess(l.title, l.options, f)
}

func (l *LogProcess) Start() {
	if l.isStarted {
		panic("process has been already started")
	} else if l.isLaunched {
		panic("process has been already launched")
	}

	l.isStarted = true

	if l.isDisabled || !l.manager.IsAccepted() || l.stream.IsMuted() {
		return
	}

	l.stream.logProcessStart(l.title, *l.options)
}

func (l *LogProcess) StepEnd(format string, a ...interface{}) {
	if !l.isStarted {
		panic("process has not been started yet")
	} else if l.isLaunched {
		panic("process has been already launched")
	}

	if l.isDisabled || !l.manager.IsAccepted() || l.stream.IsMuted() {
		return
	}

	l.stream.logProcessStepEnd(stylePkg.SimpleFormat(format, a...), *l.options)
}

func (l *LogProcess) End() {
	if !l.isStarted {
		panic("process has not been started yet")
	} else if l.isLaunched {
		panic("process has been already launched")
	}

	l.isLaunched = true

	if l.isDisabled || !l.manager.IsAccepted() || l.stream.IsMuted() {
		return
	}

	l.stream.logProcessEnd(*l.options)
}

func (l *LogProcess) Fail() {
	if !l.isStarted {
		panic("process has not been started yet")
	} else if l.isLaunched {
		panic("process has been already launched")
	}

	l.isLaunched = true

	if l.isDisabled || !l.manager.IsAccepted() || l.stream.IsMuted() {
		return
	}

	l.stream.logProcessFail(*l.options)
}

type LogProcessOptions struct {
	disableIfLevelNotAccepted bool
	withIndent                bool
	withoutLogOptionalLn      bool
	withoutElapsedTime        bool
	infoSectionFunc           func(error)
	successInfoSectionFunc    func()
	style                     *stylePkg.Style
}

func (opts *LogProcessOptions) DisableIfLevelNotAccepted() {
	opts.disableIfLevelNotAccepted = true
}

func (opts *LogProcessOptions) WithIndent() {
	opts.withIndent = true
}

func (opts *LogProcessOptions) WithoutLogOptionalLn() {
	opts.withoutLogOptionalLn = true
}

func (opts *LogProcessOptions) WithoutElapsedTime() {
	opts.withoutElapsedTime = true
}

func (opts *LogProcessOptions) InfoSectionFunc(f func(err error)) {
	opts.infoSectionFunc = f
}

func (opts *LogProcessOptions) SuccessInfoSectionFunc(f func()) {
	opts.successInfoSectionFunc = f
}

func (opts *LogProcessOptions) Style(style *stylePkg.Style) {
	opts.style = style
}
