package logboek

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

const (
	HighlightStyleName = "highlight"
	DetailsStyleName   = "details"
	FailStyleName      = "fail"
)

var (
	styles = map[string]*Style{}
)

func init() {
	predefinedStyles := map[string]*Style{
		HighlightStyleName: {Attributes: []color.Attribute{color.Bold}},
		DetailsStyleName:   {Attributes: []color.Attribute{color.FgBlue, color.Bold}},
		FailStyleName:      {Attributes: []color.Attribute{color.FgRed, color.Bold}},
	}

	for name, style := range predefinedStyles {
		SetStyle(name, style)
	}
}

func SetStyle(name string, style *Style) {
	styles[name] = style
}

func StyleByName(name string) *Style {
	return styles[name]
}

func DetailsStyle() *Style {
	return StyleByName(DetailsStyleName)
}

func HighlightStyle() *Style {
	return StyleByName(HighlightStyleName)
}

type Style struct {
	Attributes []color.Attribute
}

func (s *Style) Colorize(format string, a ...interface{}) string {
	var colorizedLines []string
	lines := strings.Split(simpleFormat(format, a...), "\n")
	for _, line := range lines {
		if line == "" {
			colorizedLines = append(colorizedLines, line)
		} else {
			colorizedLines = append(colorizedLines, color.New(s.Attributes...).Sprint(line))
		}
	}

	return strings.Join(colorizedLines, "\n")
}

func formatWithStyle(style *Style, format string, a ...interface{}) string {
	if style == nil {
		return simpleFormat(format, a...)
	} else {
		return style.Colorize(format, a...)
	}
}

func simpleFormat(format string, a ...interface{}) string {
	var msg string
	if len(a) > 0 {
		msg = fmt.Sprintf(format, a...)
	} else {
		msg = fmt.Sprintf("%s", format)
	}

	return msg
}
