package logboek

import "strings"

var (
	tagPartWidth int

	tagValue       string
	tagStyle       *Style
	tagIndentWidth = 2
)

func WithTag(value string, style *Style, f func() error) error {
	savedTag := tagValue
	savedStyle := tagStyle

	SetTag(value, style)
	err := f()
	SetTag(savedTag, savedStyle)

	return err
}

func SetTag(value string, style *Style) {
	if value != "" {
		tagPartWidth = len(value) + tagIndentWidth
	} else {
		tagPartWidth = 0
	}

	tagValue = value
	tagStyle = style
}

func formattedTag() string {
	if len(tagValue) == 0 {
		return ""
	}

	return strings.Join([]string{
		tagStyle.Colorize(tagValue),
		strings.Repeat(" ", tagIndentWidth),
	}, "")
}
