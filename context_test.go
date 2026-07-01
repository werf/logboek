package logboek

import (
	"os"
	"testing"

	"golang.org/x/net/context"
)

func TestContextFallsBackToDefaultLogger(t *testing.T) {
	cases := map[string]context.Context{
		"nil context":            nil,
		"background context":     context.Background(),
		"context without logger": context.WithValue(context.Background(), "unrelated", "x"),
	}
	for name, ctx := range cases {
		if got := Context(ctx); got != DefaultLogger() {
			t.Errorf("%s: Context() = %v, want DefaultLogger()", name, got)
		}
	}
}

func TestContextReturnsBoundLogger(t *testing.T) {
	l := NewLogger(os.Stdout, os.Stderr)
	ctx := NewContext(context.Background(), l)
	if got := Context(ctx); got != l {
		t.Errorf("Context() = %v, want bound logger %v", got, l)
	}
}

func TestMustContextPanicsWithoutLogger(t *testing.T) {
	for name, ctx := range map[string]context.Context{
		"nil context":        nil,
		"background context": context.Background(),
	} {
		t.Run(name, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Error("MustContext() did not panic")
				}
			}()
			MustContext(ctx)
		})
	}
}

func TestMustContextReturnsBoundLogger(t *testing.T) {
	l := NewLogger(os.Stdout, os.Stderr)
	ctx := NewContext(context.Background(), l)
	if got := MustContext(ctx); got != l {
		t.Errorf("MustContext() = %v, want bound logger %v", got, l)
	}
}
