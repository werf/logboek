package stream

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"
)

// newWrappingStream builds a Stream over buf with line wrapping enabled and a
// fixed width, so ContentWidth()/ServiceWidth() are deterministic across
// environments (no prefix/process/tag/indent => ServiceWidth == 0 =>
// ContentWidth == width). Style is nil (no SGR decoration).
func newWrappingStream(buf io.Writer, width int) *Stream {
	s := NewStream(buf, NewStreamState())
	s.EnableLineWrapping()
	s.SetWidth(width)
	return s
}

// TestFormatAndLogF_longLineCrossChunk locks in FormatAndLogF behavior for a
// single long line that exceeds chunkSize (256 runes), exercising the
// multi-chunk loop with shared State.
func TestFormatAndLogF_longLineCrossChunk(t *testing.T) {
	const width = 10
	in := strings.Repeat("a", 300) // > 256 to cross the chunk boundary

	// width 10, marked wrap uses contentWidth-2 slice + " ↵": 37 full "aaaaaaaa ↵"
	// lines then "aaaa" (300 = 37*8 + 4).
	expected := strings.Repeat("aaaaaaaa ↵\n", 37) + "aaaa"

	t.Run("cacheIncompleteLine=false", func(t *testing.T) {
		var buf bytes.Buffer
		newWrappingStream(&buf, width).FormatAndLogF(nil, false, "%s", in)
		if got := buf.String(); got != expected {
			t.Errorf("\n[EXPECTED]: %q\n[GOT]: %q", expected, got)
		}
	})

	t.Run("cacheIncompleteLine=true_flushesOnNextWrite", func(t *testing.T) {
		var buf bytes.Buffer
		s := newWrappingStream(&buf, width)

		// With cacheIncompleteLine=true, the trailing incomplete line is held
		// in State and not emitted yet.
		s.FormatAndLogF(nil, true, "%s", in)
		if got := buf.String(); got != "" {
			t.Errorf("expected cached incomplete line to withhold output, got %q", got)
		}

		// A subsequent non-cached write flushes the cached tail joined to it.
		s.FormatAndLogF(nil, false, "%s", "END")
		if got, want := buf.String(), expected+"END"; got != want {
			t.Errorf("\n[EXPECTED]: %q\n[GOT]: %q", want, got)
		}
	})
}

// BenchmarkFormatAndLogF_longLine quantifies long-line stream formatting
// time/alloc growth. Baseline only, not a CI gate — n=100000 allocates
// multiple GB (manual runs: go test -bench . -benchmem ./internal/stream/).
func BenchmarkFormatAndLogF_longLine(b *testing.B) {
	for _, n := range []int{1000, 10000, 100000} {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			s := newWrappingStream(io.Discard, 10)
			text := strings.Repeat("a", n)
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				s.FormatAndLogF(nil, false, "%s", text)
			}
		})
	}
}
