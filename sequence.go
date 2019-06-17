package logboek

type sequenceKind int

const (
	unknownSequenceKind = iota
	plainSequenceKind
	controlSequenceKind
)

type sequence struct {
	data string
	kind sequenceKind
}

func newSequence(data string) *sequence {
	s := &sequence{}
	s.data = data
	s.kind = unknownSequenceKind

	return s
}

func (s *sequence) Append(data string) {
	s.data += data
}

func (s *sequence) SetKind(kind sequenceKind) {
	if s.kind != unknownSequenceKind {
		panic("sequence kind already exists")
	}

	s.kind = kind
}

func (s *sequence) String() string {
	return s.data
}

func (s *sequence) TWidth() int {
	if s.kind == controlSequenceKind {
		if s.data == "\b" {
			return -1
		}
		return 0
	}

	return len(s.data)
}

func (s *sequence) Slice(maxTWidth int) (string, int) {
	var result string
	var rest int

	difference := maxTWidth - s.TWidth()
	if difference <= 0 {
		result = s.data[:maxTWidth]
		s.data = s.data[maxTWidth:]
		rest = 0
	} else {
		result = s.data
		s.data = ""
		rest = difference
	}

	return result, rest
}

func (s *sequence) IsEmpty() bool {
	return len(s.data) == 0
}

type sequenceStack struct {
	sequences []*sequence
}

func newSequenceStack() sequenceStack {
	return sequenceStack{}
}

func (ss *sequenceStack) String() string {
	var result string
	for _, s := range ss.sequences {
		result += s.String()
	}

	return result
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
	if !ss.IsEmpty() && ss.TopSequence().kind == unknownSequenceKind {
		ss.TopSequence().SetKind(plainSequenceKind)
	}

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

func (ss *sequenceStack) DivideListSign() {
	if len(ss.sequences) == 0 {
		panic("empty sequence stack")
	}

	data := ss.TopSequence().data
	if len(data) == 0 {
		panic("empty top sequence")
	}

	sign := data[len(data)-1]
	ss.TopSequence().data = data[:len(data)-1]

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
	var result string
	var newSequences []*sequence

	for _, s := range ss.sequences {
		if sliceTWidth == 0 {
			newSequences = append(newSequences, s)
		} else if s.TWidth() == 0 {
			result += s.String()
		} else {
			var part string
			part, sliceTWidth = s.Slice(sliceTWidth)
			result += part

			if !s.IsEmpty() {
				newSequences = append(newSequences, s)
			}
		}

	}

	ss.sequences = newSequences
	if len(ss.sequences) == 0 {
		ss.NewSequence("")
	}

	return result, sliceTWidth
}

func (ss *sequenceStack) Slices(sliceTWidth int) ([]string, int) {
	var result []string

	if ss.IsEmpty() {
		return result, sliceTWidth
	}

	for {
		slice, rest := ss.Slice(sliceTWidth)
		result = append(result, slice)

		if ss.TWidth() == 0 {
			return result, rest
		}
	}
}
