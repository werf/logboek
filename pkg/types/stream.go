package types

import (
	"github.com/gookit/color"
)

type StreamsInterface interface {
	Mute()
	Unmute()
	IsMuted() bool

	EnableGitlabCollapsibleSections()
	DisableGitlabCollapsibleSections()
	IsGitlabCollapsibleSections() bool

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

	DoWithTag(value string, style color.Style, f func())
	DoErrorWithTag(value string, style color.Style, f func() error) error
	SetTag(value string)
	SetTagStyle(style color.Style)
	SetTagWithStyle(value string, style color.Style)
	ResetTag()

	EnablePrefixWithTime()
	DisablePrefixWithTime()
	IsPrefixWithTimeEnabled() bool
	ResetPrefixTime()
	SetPrefix(value string)
	SetPrefixStyle(style color.Style)
	ResetPrefix()

	EnableLogProcessBorder()
	DisableLogProcessBorder()
	IsLogProcessBorderEnabled() bool
}
