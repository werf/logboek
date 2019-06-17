package logboek

import (
	"reflect"
	"strings"
	"testing"
)

func TestSequence_NewSequence(t *testing.T) {
	tests := []string{"", "data", "\x1b[0m"}

	for _, testData := range tests {
		t.Run(testData, func(t *testing.T) {
			s := newSequence(testData)

			if s.kind != unknownSequenceKind {
				t.Errorf("expectedLines unknown sequence kind")
			}

			if testData == "" {
				if s.IsEmpty() != true {
					t.Errorf("expectedLines empty sequence")
				}
			} else {
				if s.IsEmpty() != false {
					t.Errorf("expectedLines non empty sequence")
				}
			}

			if s.String() != testData {
				t.Errorf("\n[EXPECTED]: %s\n[GOT]: %s", testData, s.String())
			}

			if s.TWidth() != len(testData) {
				t.Errorf("\n[EXPECTED]: %d\n[GOT]: %d", s.TWidth(), len(testData))
			}
		})
	}
}

func TestSequence_SetKind_base(t *testing.T) {
	for _, kind := range []sequenceKind{plainSequenceKind, controlSequenceKind} {
		s := newSequence("")
		s.SetKind(kind)

		if s.kind != kind {
			t.Errorf("\n[EXPECTED]: %v\n[GOT]: %v", s.kind, kind)
		}
	}
}

func TestSequence_SetKind_panic(t *testing.T) {
	s := newSequence("")
	s.SetKind(plainSequenceKind)

	assertPanic(t, "expectedLines panic during reset sequence kind", func() {
		s.SetKind(plainSequenceKind)
	})
}

func TestSequence_TWidth(t *testing.T) {
	tests := []struct {
		data           string
		kind           sequenceKind
		expectedTWidth int
	}{
		{
			"data",
			plainSequenceKind,
			len("data"),
		},
		{
			"data",
			controlSequenceKind,
			0,
		},
	}

	for _, test := range tests {
		t.Run(string(test.kind), func(t *testing.T) {
			s := newSequence(test.data)
			s.SetKind(test.kind)

			if s.TWidth() != test.expectedTWidth {
				t.Errorf("\n[EXPECTED]: %d\n[GOT]: %d", s.TWidth(), test.expectedTWidth)
			}
		})
	}
}

func TestSequence_Append(t *testing.T) {
	s := newSequence("12")
	s.Append("34")

	if s.String() != "1234" {
		t.Errorf("\n[EXPECTED]: %v\n[GOT]: %v", "1234", s.String())
	}
}

func assertPanic(t *testing.T, msg string, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf(msg)
		}
	}()

	f()
}

func TestSequence_Slice(t *testing.T) {
	for _, test := range sliceTests {
		t.Run(test.name, func(t *testing.T) {
			s := newSequence(test.data)
			s.SetKind(test.kind)

			result, rest := s.Slice(sliceTestMaxTWidth)

			if test.result != result {
				t.Errorf("\n[EXPECTED]: %s\n[GOT]: %s", test.result, result)
			}

			if test.rest != rest {
				t.Errorf("\n[EXPECTED]: %d\n[GOT]: %d", test.rest, rest)
			}

			if test.newData != s.data {
				t.Errorf("\n[EXPECTED]: %s\n[GOT]: %s", test.newData, s.data)
			}
		})
	}
}

type sliceTest struct {
	name    string
	data    string
	kind    sequenceKind
	result  string
	rest    int
	newData string
}

var (
	sliceTestMaxTWidth = 10
	sliceTests         = []sliceTest{
		{
			name:    "emptyPlain",
			data:    "",
			kind:    plainSequenceKind,
			result:  "",
			rest:    sliceTestMaxTWidth,
			newData: "",
		},
		{
			name:    "emptyControl",
			data:    "",
			kind:    controlSequenceKind,
			result:  "",
			rest:    sliceTestMaxTWidth,
			newData: "",
		},
		{
			name:    "shortPlain",
			data:    "data",
			kind:    plainSequenceKind,
			result:  "data",
			rest:    sliceTestMaxTWidth - len("data"),
			newData: "",
		},
		{
			name:    "control",
			data:    "\x1b[0m",
			kind:    controlSequenceKind,
			result:  "\x1b[0m",
			rest:    sliceTestMaxTWidth,
			newData: "",
		},
		{
			name:    "equalPlain",
			data:    strings.Repeat("l", sliceTestMaxTWidth),
			kind:    plainSequenceKind,
			result:  strings.Repeat("l", sliceTestMaxTWidth),
			rest:    0,
			newData: "",
		},
		{
			name:    "longPlain",
			data:    strings.Repeat("l", sliceTestMaxTWidth+5),
			kind:    plainSequenceKind,
			result:  strings.Repeat("l", sliceTestMaxTWidth),
			rest:    0,
			newData: strings.Repeat("l", 5),
		},
	}
)

func TestSequenceStack_newSequenceStack(t *testing.T) {
	ss := newSequenceStack()

	if len(ss.sequences) != 0 {
		t.Errorf("expectedLines empty sequence stack")
	}
}

func TestSequenceStack_TopSequence(t *testing.T) {
	ss := newSequenceStack()

	if ss.TopSequence() == nil {
		t.Errorf("expectedLines autocreated sequence in stack")
	}
}

func TestSequenceStack_WriteData(t *testing.T) {
	ss := newSequenceStack()
	data := "data"
	ss.WriteData(data)

	if len(ss.sequences) != 1 {
		t.Errorf("expectedLines top sequence in stack (got %d)", len(ss.sequences))
	}

	if ss.TopSequence().kind != unknownSequenceKind {
		t.Errorf("\n[EXPECTED]: %v\n[GOT]: %v", unknownSequenceKind, ss.TopSequence().kind)
	}

	if ss.TopSequence().String() != data {
		t.Errorf("\n[EXPECTED]: %s\n[GOT]: %s", data, ss.TopSequence().String())
	}

	if ss.String() != data {
		t.Errorf("\n[EXPECTED]: %s\n[GOT]: %s", data, ss.String())
	}
}

func TestSequenceStack_WriteControlData(t *testing.T) {
	ss := newSequenceStack()
	data := "data"
	ss.WriteControlData(data)

	if len(ss.sequences) != 2 {
		t.Errorf("expectedLines commited and autocreated sequences in stack")
	}

	if ss.sequences[0].kind != controlSequenceKind {
		t.Errorf("\n[EXPECTED]: %v\n[GOT]: %v", controlSequenceKind, ss.sequences[0].kind)
	}

	if ss.String() != data {
		t.Errorf("\n[EXPECTED]: %s\n[GOT]: %s", data, ss.String())
	}

	if ss.TWidth() != 0 {
		t.Errorf("\n[EXPECTED]: %d\n[GOT]: %d", 0, ss.TWidth())
	}
}

func TestSequenceStack_WritePlainData(t *testing.T) {
	ss := newSequenceStack()
	data := "data"
	ss.WritePlainData(data)

	if len(ss.sequences) != 2 {
		t.Errorf("expectedLines commited and autocreated sequences in stack")
	}

	if ss.sequences[0].kind != plainSequenceKind {
		t.Errorf("\n[EXPECTED]: %v\n[GOT]: %v", plainSequenceKind, ss.sequences[0].kind)
	}

	if ss.String() != data {
		t.Errorf("\n[EXPECTED]: %s\n[GOT]: %s", data, ss.String())
	}

	if ss.TWidth() != len(data) {
		t.Errorf("\n[EXPECTED]: %d\n[GOT]: %d", len(data), ss.TWidth())
	}
}

func TestSequenceStack_Merge_empty(t *testing.T) {
	ss1 := newSequenceStack()
	ss2 := newSequenceStack()

	ss1.Merge(ss2)

	if len(ss1.sequences) != 1 {
		t.Errorf("expectedLines auto top sequence in stack")
	}
}

func TestSequenceStack_Merge_autokind(t *testing.T) {
	ss1 := newSequenceStack()
	ss2 := newSequenceStack()

	ss1.WriteData("data1")
	ss2.WriteData("data2")

	ss1.Merge(ss2)

	if ss1.sequences[0].kind != plainSequenceKind {
		t.Errorf("expectedLines auto plain kind")
	}

	if ss1.sequences[1].kind != plainSequenceKind {
		t.Errorf("expectedLines auto plain kind")
	}

	if len(ss1.sequences) != 3 {
		t.Errorf("expectedLines auto top sequence in stack (got %d)", len(ss1.sequences))
	}
}

func TestSequenceStack_Slice(t *testing.T) {
	for _, test := range sliceTests {
		t.Run(test.name, func(t *testing.T) {
			ss := newSequenceStack()
			ss.NewSequence(test.data)
			ss.CommitTopSequence(test.kind)

			result, rest := ss.Slice(sliceTestMaxTWidth)

			if test.result != result {
				t.Errorf("\n[EXPECTED]: %s\n[GOT]: %s", test.result, result)
			}

			if test.rest != rest {
				t.Errorf("\n[EXPECTED]: %d\n[GOT]: %d", test.rest, rest)
			}

			if test.newData != ss.String() {
				t.Errorf("\n[EXPECTED]: %s\n[GOT]: %s", test.newData, ss.String())
			}

			if test.newData != test.data {
				switch test.newData {
				case "":
					if len(ss.sequences) != 1 {
						t.Errorf("expectedLines only top empty sequence in stack (got %d)", len(ss.sequences))
					}
				default:
					if test.newData != ss.String() {
						t.Errorf("\n[EXPECTED]: %s\n[GOT]: %s", test.newData, ss.String())
					}
				}

			}
		})
	}
}

func TestSequenceStack_Slices(t *testing.T) {
	tests := []struct {
		name           string
		data           string
		sliceSize      int
		expectedSlices []string
		expectedRest   int
	}{
		{
			"short",
			"12",
			5,
			[]string{"12"},
			3,
		},
		{
			"equal",
			"12345",
			5,
			[]string{"12345"},
			0,
		},
		{
			"long",
			"12345678901",
			5,
			[]string{"12345", "67890", "1"},
			4,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ss := newSequenceStack()
			ss.WritePlainData(test.data)

			slices, rest := ss.Slices(test.sliceSize)
			if !reflect.DeepEqual(test.expectedSlices, slices) {
				t.Errorf("\n[EXPECTED]: %+v\n[GOT]: %+v", test.expectedSlices, slices)
			}

			if test.expectedRest != rest {
				t.Errorf("\n[EXPECTED]: %d\n[GOT]: %d", test.expectedRest, rest)
			}
		})
	}
}
