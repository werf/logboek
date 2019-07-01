package logboek

import (
	"fmt"
	"strings"
	"testing"
)

type fitTextTest struct {
	name     string
	data     string
	expected string
}

func runFitTextTests(t *testing.T, testNameFormat string, withMarkedLine bool, tests []fitTextTest) {
	for _, test := range tests {
		t.Run(fmt.Sprintf(testNameFormat, test.name), func(t *testing.T) {
			result := FitText(test.data, FitTextOptions{MarkWrappedLine: withMarkedLine})
			if test.expected != result {
				t.Errorf("\n[EXPECTED]: %q\n[GOT]: %q", test.expected, result)
			}
		})
	}
}

func TestFitText_sentence(t *testing.T) {
	contentWidth := 10
	SetWidth(contentWidth)

	runFitTextTests(t, "withoutMarkedLine_%s", false, []fitTextTest{
		{
			"short",
			"foo bar",
			"foo bar",
		},
		{
			"equal",
			"foo bar da",
			"foo bar da",
		},
		{
			"bigger",
			"foo bar data",
			"foo bar   \ndata",
		},
		{
			"biggerWithLongWord",
			"foo bar " + strings.Repeat("l", contentWidth+1),
			"foo bar ll\nlllllllll",
		},
		{
			"doubleLongWords",
			strings.Repeat("l", contentWidth) + " " + strings.Repeat("l", contentWidth+1),
			"llllllllll\n lllllllll\nll",
		},
	})

	runFitTextTests(t, "withMarkedLine_%s", true, []fitTextTest{
		{
			"short",
			"foo bar",
			"foo bar",
		},
		{
			"equal",
			"foo bar da",
			"foo bar da",
		},
		{
			"bigger",
			"foo bar data",
			"foo bar  ↵\ndata",
		},
		{
			"biggerWithLongWord",
			"foo bar " + strings.Repeat("l", contentWidth+1),
			"foo bar  ↵\nllllllll ↵\nlll",
		},
		{
			"doubleLongWords",
			strings.Repeat("l", contentWidth) + " " + strings.Repeat("l", contentWidth+1),
			"llllllll ↵\nll lllll ↵\nllllll",
		},
	})
}

func TestFitText_word(t *testing.T) {
	contentWidth := 10
	SetWidth(contentWidth)

	runFitTextTests(t, "withoutMarkedLine_%s", false, []fitTextTest{
		{
			"short",
			strings.Repeat("1", contentWidth-1),
			strings.Repeat("1", contentWidth-1),
		},
		{
			"color_short",
			"\x1b[30m" + strings.Repeat("1", contentWidth-1) + "\x1b[0m",
			"\x1b[30m" + strings.Repeat("1", contentWidth-1) + "\x1b[0m",
		},
		{
			"equal",
			strings.Repeat("1", contentWidth),
			strings.Repeat("1", contentWidth),
		},
		{
			"color_equal",
			"\x1b[30m" + strings.Repeat("1", contentWidth) + "\x1b[0m",
			"\x1b[30m" + strings.Repeat("1", contentWidth) + "\x1b[0m",
		},
		{
			"bigger",
			strings.Repeat("1", contentWidth+1),
			strings.Repeat("1", contentWidth) + "\n1",
		},
		{
			"color_bigger",
			"\x1b[30m" + strings.Repeat("1", contentWidth+1) + "\x1b[0m",
			"\x1b[30m" + strings.Repeat("1", contentWidth) + "\x1b[0m" + "\n" + "\x1b[30m" + strings.Repeat("1", 1) + "\x1b[0m",
		},
	})

	runFitTextTests(t, "withMarkedLine_%s", true, []fitTextTest{
		{
			"short",
			strings.Repeat("1", contentWidth-1),
			strings.Repeat("1", contentWidth-1),
		},
		{
			"color_short",
			"\x1b[30m" + strings.Repeat("1", contentWidth-1) + "\x1b[0m",
			"\x1b[30m" + strings.Repeat("1", contentWidth-1) + "\x1b[0m",
		},
		{
			"equal",
			strings.Repeat("1", contentWidth),
			"1111111111",
		},
		{
			"color_equal",
			"\x1b[30m" + strings.Repeat("1", contentWidth) + "\x1b[0m",
			"\x1b[30m1111111111\x1b[0m",
		},
		{
			"bigger",
			strings.Repeat("1", contentWidth+1),
			strings.Repeat("1", contentWidth-2) + " ↵\n" + strings.Repeat("1", 3),
		},
		{
			"color_bigger",
			"\x1b[30m" + strings.Repeat("1", contentWidth+1) + "\x1b[0m",
			"\x1b[30m11111111 ↵\x1b[0m\n\x1b[30m111\x1b[0m",
		},
	})
}

func TestWrapperState_splitSequencesStack(t *testing.T) {
	contentWidth := 10

	tests := []struct {
		name            string
		data            string
		markWrappedLine bool
		expectedLines   string
	}{
		{
			"nonMarkedLines",
			"1234567890123456789012345678901234567890",
			false,
			`1234567890
1234567890
1234567890
1234567890`,
		},
		{
			"markedLines",
			"1234567890123456789012345678901234567890",
			true,
			`12345678 ↵
90123456 ↵
78901234 ↵
56789012 ↵
34567890`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ss := wrapperState{}
			ss.WriteData(test.data)

			result := ss.splitSequenceStack(contentWidth, test.markWrappedLine)
			if test.expectedLines != result {
				t.Errorf("\n[EXPECTED]: %q\n[GOT]: %q", test.expectedLines, result)
			}
		})
	}
}

func Test_markLine(t *testing.T) {
	contentWidth := 10

	tests := []struct {
		name     string
		data     string
		expected string
	}{
		{
			"empty",
			"",
			strings.Repeat(" ", contentWidth-1) + "↵",
		},
		{
			"short",
			"012345",
			"012345" + strings.Repeat(" ", contentWidth-len("012345")-1) + "↵",
		},
		{
			"equal",
			strings.Repeat("0", contentWidth),
			strings.Repeat("0", contentWidth) + "↵",
		},
		{
			"contentWidthMinus1",
			strings.Repeat("0", contentWidth-1),
			strings.Repeat("0", contentWidth-1) + "↵",
		},
		{
			"contentWidthMinus2",
			strings.Repeat("0", contentWidth-2),
			strings.Repeat("0", contentWidth-2) + " " + "↵",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := markLine(test.data, len(test.data), contentWidth)
			if test.expected != result {
				t.Errorf("\n[EXPECTED]: %q\n[GOT]: %q", test.expected, result)
			}
		})
	}
}
