package types

import (
	"github.com/gookit/color"
)

type LogBlockInterface interface {
	Options(func(options LogBlockOptionsInterface)) LogBlockInterface
	Disable() LogBlockInterface
	Enable() LogBlockInterface
	Do(func())
	DoError(func() error) error
}

type LogBlockOptionsInterface interface {
	DisableIfLevelNotAccepted()
	Mute()
	WithIndent()
	WithoutLogOptionalLn()
	Style(color.Style)
}

type LogProcessInlineInterface interface {
	Options(func(options LogProcessInlineOptionsInterface)) LogProcessInlineInterface
	Disable() LogProcessInlineInterface
	Enable() LogProcessInlineInterface
	Do(func())
	DoError(func() error) error
}

type LogProcessInlineOptionsInterface interface {
	DisableIfLevelNotAccepted()
	Mute()
	Style(color.Style)
}

type LogProcessInterface interface {
	Options(func(options LogProcessOptionsInterface)) LogProcessInterface
	Disable() LogProcessInterface
	Enable() LogProcessInterface
	Do(func())
	DoError(func() error) error
	Start()
	StepEnd(format string, a ...interface{})
	End()
	Fail()
}

type LogProcessOptionsInterface interface {
	DisableIfLevelNotAccepted()
	Mute()
	WithIndent()
	WithoutLogOptionalLn()
	WithoutElapsedTime()
	InfoSectionFunc(func(err error))
	SuccessInfoSectionFunc(func())
	Style(color.Style)
}
