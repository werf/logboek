package types

import "github.com/werf/logboek/pkg/style"

type StreamsInterface interface {
	Mute()
	Unmute()
	IsMuted() bool

	DisablePrettyLog()

	DoWithProxyStreamDataFormatting(func())
	DoWithoutProxyStreamDataFormatting(func())
	DoErrorWithProxyStreamDataFormatting(func() error) error
	DoErrorWithoutProxyStreamDataFormatting(func() error) error
	EnableProxyStreamDataFormatting()
	DisableProxyStreamDataFormatting()
	IsProxyStreamDataFormattingEnabled() bool

	EnableLineWrapping()
	DisableLineWrapping()
	IsLineWrappingEnabled() bool

	EnableStyle()
	DisableStyle()
	IsStyleEnabled() bool

	Width() int
	SetWidth(value int)
	ContentWidth() int
	ServiceWidth() int

	DoWithIndent(func())
	DoErrorWithIndent(func() error) error
	DoWithoutIndent(func())
	DoErrorWithoutIndent(func() error) error
	IncreaseIndent()
	DecreaseIndent()
	ResetIndent()

	DoWithTag(value string, style *style.Style, f func())
	DoErrorWithTag(value string, style *style.Style, f func() error) error
	SetTag(value string)
	SetTagStyle(style *style.Style)
	SetTagWithStyle(value string, style *style.Style)
	ResetTag()

	EnablePrefixWithTime()
	DisablePrefixWithTime()
	IsPrefixWithTimeEnabled() bool
	ResetPrefixTime()
	SetPrefix(value string)
	SetPrefixStyle(style *style.Style)
	ResetPrefix()

	EnableLogProcessBorder()
	DisableLogProcessBorder()
}
