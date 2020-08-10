package style

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

const (
	HighlightName = "highlight"
	DetailsName   = "details"
	FailName      = "fail"
	NoneName      = "none"
)

var (
	styles = map[string]*Style{}
)

func init() {
	predefinedStyles := map[string]*Style{
		HighlightName: {Attributes: []color.Attribute{color.Bold}},
		DetailsName:   {Attributes: []color.Attribute{color.FgBlue, color.Bold}},
		FailName:      {Attributes: []color.Attribute{color.FgRed, color.Bold}},
		NoneName:      {Attributes: []color.Attribute{}},
	}

	for name, s := range predefinedStyles {
		Set(name, s)
	}
}

type Style struct {
	Attributes []color.Attribute
}

func (s *Style) Colorize(format string, a ...interface{}) string {
	var colorizedLines []string
	lines := strings.Split(SimpleFormat(format, a...), "\n")
	for _, line := range lines {
		if line == "" {
			colorizedLines = append(colorizedLines, line)
		} else {
			c := color.New(s.Attributes...)
			c.EnableColor()
			colorizedLines = append(colorizedLines, c.Sprint(line))
		}
	}

	return strings.Join(colorizedLines, "\n")
}

func SimpleFormat(format string, a ...interface{}) string {
	var msg string
	if len(a) > 0 {
		msg = fmt.Sprintf(format, a...)
	} else {
		msg = fmt.Sprintf("%s", format)
	}

	return msg
}

func Set(name string, style *Style) {
	styles[name] = style
}

func Get(name string) *Style {
	return styles[name]
}

func Details() *Style {
	return Get(DetailsName)
}

func Highlight() *Style {
	return Get(HighlightName)
}

func None() *Style {
	return Get(NoneName)
}
