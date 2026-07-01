package fitter

import (
	"strings"
	"unicode/utf8"
)

type sequenceKind int

const (
	plainSequenceKind = iota
	controlSequenceKind
)

// sequence holds one run of text. It has two lifecycle phases that never
// interleave (verified: wrapperState.Apply always resets the stack after
// Slices, so a sequence is either being built OR being sliced, never both):
//   - BUILD: per-rune Append into the strings.Builder (amortized O(1), keeps
//     phase 1 / PR #77 from regressing to O(n²) concat).
//   - SLICE: materialize the Builder to str ONCE (zero-copy), then advance a
//     byte offset off for each Slice (O(1), no tail copy) with a cached twidth.
type sequence struct {
	data    strings.Builder // build-phase accumulator
	str     string          // slice-phase view, materialized once from data
	strSet  bool            // str materialized this cycle
	off     int             // byte offset into str
	twidth  int             // cached rune width of str[off:]
	twValid bool
	kind    sequenceKind
}

func newSequence(data string) *sequence {
	s := &sequence{}
	s.data.WriteString(data)
	s.kind = plainSequenceKind

	return s
}

// content returns the current logical content, materializing the build buffer
// on first read of the slice phase and reading str[off:] thereafter.
func (s *sequence) content() string {
	if !s.strSet {
		s.str = s.data.String() // zero-copy in Go >=1.10
		s.strSet = true
	}
	return s.str[s.off:]
}

func (s *sequence) Append(data string) {
	if s.off != 0 { // ponytail: Append never follows Slice (verified); assert to catch future misuse
		panic("sequence: Append after Slice")
	}
	s.data.WriteString(data)
	s.strSet = false
	s.twValid = false
}

// setData replaces the buffer contents (Builder has no in-place truncate).
func (s *sequence) setData(data string) {
	s.data.Reset()
	s.data.WriteString(data)
	s.str = ""
	s.strSet = false
	s.off = 0
	s.twValid = false
}

func (s *sequence) SetKind(kind sequenceKind) {
	s.kind = kind
}

func (s *sequence) String() string {
	return s.content()
}

func (s *sequence) TWidth() int {
	if s.kind == controlSequenceKind {
		if s.content() == "\b" {
			return -1
		}
		return 0
	}

	if !s.twValid {
		s.twidth = utf8.RuneCountInString(s.content())
		s.twValid = true
	}
	return s.twidth
}

func (s *sequence) Slice(maxTWidth int) (string, int) {
	content := s.content()
	difference := maxTWidth - s.TWidth()
	if difference <= 0 {
		isASCII := s.twidth == len(content) // byte==rune => pure ASCII view
		if isASCII {
			// ASCII fast-path: byte==rune so cutting at maxTWidth is exact and O(1);
			// the tail of an ASCII string is ASCII, so the cache stays valid.
			result := content[:maxTWidth]
			s.off += maxTWidth
			s.twidth -= maxTWidth
			return result, 0
		}
		// multibyte: cut on a rune boundary so emoji/wide chars never split mid-byte.
		byteLen := 0
		for i := 0; i < maxTWidth; i++ {
			_, size := utf8.DecodeRuneInString(content[byteLen:])
			byteLen += size
		}
		s.off += byteLen
		s.twValid = false
		return content[:byteLen], 0
	}

	s.off = len(s.str)
	s.twidth = 0
	s.twValid = true
	return content, difference
}

func (s *sequence) IsEmpty() bool {
	if s.strSet {
		return s.off >= len(s.str)
	}
	return s.data.Len() == 0
}

type sequenceStack struct {
	sequences []*sequence
}

func newSequenceStack() sequenceStack {
	return sequenceStack{}
}

func (ss *sequenceStack) String() string {
	var b strings.Builder
	for _, s := range ss.sequences {
		b.WriteString(s.String())
	}

	return b.String()
}

func (ss *sequenceStack) TWidth() int {
	var result int
	for _, s := range ss.sequences {
		result += s.TWidth()
	}

	return result
}

func (ss *sequenceStack) CommitTopSequenceAsPlain() {
	ss.CommitTopSequence(plainSequenceKind)
}

func (ss *sequenceStack) CommitTopSequenceAsControl() {
	ss.CommitTopSequence(controlSequenceKind)
}

func (ss *sequenceStack) CommitTopSequence(kind sequenceKind) {
	ss.commitTopSequence(kind)
	_ = ss.NewSequence("")
}

func (ss *sequenceStack) commitTopSequence(kind sequenceKind) {
	ss.TopSequence().SetKind(kind)
}

func (ss *sequenceStack) NewSequence(data string) *sequence {
	topSequence := newSequence(data)
	ss.sequences = append(ss.sequences, topSequence)

	return topSequence
}

func (ss *sequenceStack) WritePlainData(data string) {
	ss.WriteData(data)
	ss.CommitTopSequenceAsPlain()
}

func (ss *sequenceStack) WriteControlData(data string) {
	ss.WriteData(data)
	ss.CommitTopSequenceAsControl()
}

func (ss *sequenceStack) WriteData(data string) {
	if len(ss.sequences) == 0 || ss.TopSequence().IsEmpty() {
		ss.NewSequence(data)
	} else {
		ss.TopSequence().Append(data)
	}
}

func (ss *sequenceStack) DivideLastSign() {
	if len(ss.sequences) == 0 {
		panic("empty sequence stack")
	}

	data := ss.TopSequence().String()
	if len(data) == 0 {
		panic("empty top sequence")
	}

	sign := data[len(data)-1]
	ss.TopSequence().setData(data[:len(data)-1])

	ss.CommitTopSequenceAsPlain()
	ss.WriteData(string(sign))
}

func (ss *sequenceStack) Merge(ss2 sequenceStack) {
	if !ss.IsEmpty() {
		ss.commitTopSequence(plainSequenceKind)
	}

	if !ss2.IsEmpty() {
		ss2.commitTopSequence(plainSequenceKind)
	}

	ss.sequences = append(ss.sequences, ss2.sequences...)

	ss.NewSequence("")
}

func (ss *sequenceStack) TopSequence() *sequence {
	if ss.IsEmpty() {
		return ss.NewSequence("")
	}

	return ss.sequences[len(ss.sequences)-1]
}

func (ss *sequenceStack) IsEmpty() bool {
	return len(ss.sequences) == 0
}

func (ss *sequenceStack) Slice(sliceTWidth int) (string, int) {
	var newSequences []*sequence

	rest := sliceTWidth

	var b strings.Builder
	for ind, s := range ss.sequences {
		if s.TWidth() == 0 {
			b.WriteString(s.String())
			continue
		}

		if rest == 0 {
			newSequences = append(newSequences, ss.sequences[ind:]...)
			break
		} else {
			if s.TWidth() > rest && s.TWidth() <= sliceTWidth {
				newSequences = append(newSequences, ss.sequences[ind:]...)
				break
			}

			var part string
			part, rest = s.Slice(rest)
			b.WriteString(part)

			if !s.IsEmpty() {
				newSequences = append(newSequences, ss.sequences[ind:]...)
				break
			}
		}
	}

	result := b.String()
	ss.sequences = newSequences
	if len(ss.sequences) == 0 {
		ss.NewSequence("")
	}

	return result, rest
}

func (ss *sequenceStack) Slices(sliceTWidth int) ([]string, int) {
	var result []string

	if ss.IsEmpty() {
		return result, sliceTWidth
	}

	for {
		slice, rest := ss.Slice(sliceTWidth)

		if ss.TWidth() == 0 {
			result = append(result, slice)
			return result, rest
		} else {
			result = append(result, slice+strings.Repeat(" ", rest))
		}
	}
}
