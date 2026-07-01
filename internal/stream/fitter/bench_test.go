package fitter

import (
	"fmt"
	"strings"
	"testing"
)

// BenchmarkFitText_longLine quantifies the near-quadratic time/alloc growth of
// FitText on long single-line input. It is a baseline for a future
// optimization, not a CI gate — the n=100000 case allocates multiple GB and is
// slow, intended for manual baseline runs (go test -bench . -benchmem).
func BenchmarkFitText_longLine(b *testing.B) {
	for _, n := range []int{1000, 10000, 100000} {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			text := strings.Repeat("a", n)
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				FitText(text, &State{}, contentWidth, true, false)
			}
		})
	}
}
