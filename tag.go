package logboek

import "strings"

var (
	tagPartWidth int

	tagValue        string
	tagColorizeFunc func(...interface{}) string
	tagIndentWidth  = 2
)

func WithTag(value string, colorizeFunc func(...interface{}) string, f func() error) error {
	savedTag := tagValue
	savedColorizeFunc := tagColorizeFunc

	SetTag(value, colorizeFunc)
	err := f()
	SetTag(savedTag, savedColorizeFunc)

	return err
}

func SetTag(value string, colorizeFunc func(...interface{}) string) {
	if value != "" {
		tagPartWidth = len(value) + tagIndentWidth
	} else {
		tagPartWidth = 0
	}

	tagValue = value
	tagColorizeFunc = colorizeFunc
}

func formattedTag() string {
	if len(tagValue) == 0 {
		return ""
	}

	return strings.Join([]string{
		tagColorizeFunc(tagValue),
		strings.Repeat(" ", tagIndentWidth),
	}, "")
}
