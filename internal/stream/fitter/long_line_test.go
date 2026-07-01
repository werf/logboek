package fitter

import (
	"strings"
	"testing"
)

// These tests lock in the current FitText behavior for long single-line input,
// so a future performance optimization (see .pp performance notes) can be
// verified against unchanged output. contentWidth == 10 (see fitter_test.go).

func TestFitText_longPlainLine(t *testing.T) {
	runFitTextTests(t, "withoutMarkedLine_%s", false, []fitTextTest{
		{
			// 25 chars, width 10 => 10/10/5, unmarked wraps with plain "\n".
			"unbrokenWord",
			strings.Repeat("a", 25),
			strings.Repeat("a", contentWidth) + "\n" +
				strings.Repeat("a", contentWidth) + "\n" +
				strings.Repeat("a", 5),
		},
		{
			// "ab " repeated: word-wrapping pads trailing space to contentWidth.
			"spaceSeparated",
			strings.Repeat("ab ", 12),
			strings.Repeat("ab ab ab  \n", 3) + "ab ab ab ",
		},
	})

	runFitTextTests(t, "withMarkedLine_%s", true, []fitTextTest{
		{
			// Marked wrap uses contentWidth-2 slice + " ↵" padding.
			"unbrokenWord",
			strings.Repeat("a", 25),
			strings.Repeat(strings.Repeat("a", contentWidth-2)+" ↵\n", 3) + "a",
		},
		{
			"spaceSeparated",
			strings.Repeat("ab ", 12),
			strings.Repeat("ab ab    ↵\n", 5) + "ab ab ",
		},
	})
}

func TestFitText_longColoredLine(t *testing.T) {
	color := func(s string) string { return "\x1b[30m" + s + "\x1b[0m" }

	runFitTextTests(t, "withoutMarkedLine_%s", false, []fitTextTest{
		{
			// Color reset before each break, restored after: each wrapped line
			// is independently wrapped in the SGR sequence.
			"unbrokenWord",
			color(strings.Repeat("a", 25)),
			color(strings.Repeat("a", contentWidth)) + "\n" +
				color(strings.Repeat("a", contentWidth)) + "\n" +
				color(strings.Repeat("a", 5)),
		},
	})

	runFitTextTests(t, "withMarkedLine_%s", true, []fitTextTest{
		{
			// Marked: " ↵" lands inside the color span before reset.
			"unbrokenWord",
			color(strings.Repeat("a", 25)),
			strings.Repeat("\x1b[30m"+strings.Repeat("a", contentWidth-2)+" ↵\x1b[0m\n", 3) +
				color("a"),
		},
	})
}

func TestFitText_longInputPreservesControlSemantics(t *testing.T) {
	runFitTextTests(t, "%s", false, []fitTextTest{
		{
			// Backspace has terminal width -1: "abc\bd" occupies 3 cells, so the
			// following long run wraps at the same boundary as 3+... chars.
			"backspace",
			"abc\bd" + strings.Repeat("e", 12),
			"abc\bdeeeeeeeee\neee",
		},
		{
			// \r\n must not get a duplicated color reset: prevCursorRune "\r"
			// suppresses the extra reset before "\n".
			"crlfNoDoubleReset",
			"\x1b[30mabc\r\ndef\x1b[0m",
			"\x1b[30mabc\x1b[0m\r\n\x1b[30mdef\x1b[0m",
		},
	})
}

// sequence.Slice cuts on rune boundaries: 11 two-byte runes at width 10 wrap
// after 10 runes, and multibyte glyphs (emoji, CJK, Cyrillic) never split
// mid-byte into replacement chars.
func TestFitText_unicodeLongWord(t *testing.T) {
	runFitTextTests(t, "%s", false, []fitTextTest{
		{
			"cyrillicRuneBoundarySplit",
			strings.Repeat("я", 11),
			"яяяяяяяяяя\nя",
		},
		{
			"emojiNotSplitMidByte",
			strings.Repeat("🛳", 11),
			"🛳🛳🛳🛳🛳🛳🛳🛳🛳🛳\n🛳",
		},
	})
}
