package style

import (
	"github.com/gookit/color"
)

const (
	// predefined styles
	HighlightName = "logboek_highlight"
	DetailsName   = "logboek_details"
	NoneName      = "logboek_none"

	// internal styles
	ProcessFailName = "logboek_process_fail"
)

func init() {
	predefinedStyles := map[string]color.Style{
		HighlightName:   {color.Bold},
		DetailsName:     {color.FgBlue, color.Bold},
		NoneName:        {},
		ProcessFailName: {color.FgRed, color.Bold},
	}

	for name, s := range predefinedStyles {
		color.AddStyle(name, s)
	}
}

func Details() color.Style {
	return color.GetStyle(DetailsName)
}

func Highlight() color.Style {
	return color.GetStyle(HighlightName)
}

func None() color.Style {
	return color.GetStyle(NoneName)
}
