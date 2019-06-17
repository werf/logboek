package logboek

import "strings"

type FitTextOptions struct {
	ExtraIndentWidth int
	Width            int
	MaxWidth         int
	MarkWrappedLine  bool
}

func FitText(text string, options FitTextOptions) string {
	var lineWidth int
	if options.Width != 0 {
		lineWidth = options.Width
	} else {
		tw := width
		if options.MaxWidth != 0 && tw > options.MaxWidth {
			lineWidth = options.MaxWidth
		} else {
			lineWidth = tw
		}
	}

	return fitTextWithIndent(text, lineWidth, options.ExtraIndentWidth, options.MarkWrappedLine)
}

func fitTextWithIndent(text string, lineWidth, extraIndentWidth int, markWrappedLine bool) string {
	var result string
	var resultLines []string

	contentWidth := lineWidth - serviceWidth() - extraIndentWidth

	fittedText, _ := fitText(text, fitterState{}, contentWidth, markWrappedLine, false)
	for _, line := range strings.Split(fittedText, "\n") {
		indent := strings.Repeat(" ", extraIndentWidth)
		resultLines = append(resultLines, strings.Join([]string{indent, line}, ""))
	}

	result = strings.Join(resultLines, "\n")

	return result
}
