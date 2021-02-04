package util

import (
	"fmt"
	"strings"

	"github.com/gookit/color"
)

func ColorizeF(style color.Style, format string, a ...interface{}) string {
	var resultLines []string
	for _, line := range strings.Split(fmt.Sprintf(format, a...), "\n") {
		if line == "" {
			resultLines = append(resultLines, line)
		} else {
			resultLines = append(resultLines, style.Sprint(line))
		}
	}

	return strings.Join(resultLines, "\n")
}
