package stmt_test

import (
	"strings"
	"testing"

	"github.com/laacin/inyorm/internal/stmt"
)

func TestPlaceholderGen_Write(t *testing.T) {
	t.Run("Simple stmt.", func(t *testing.T) {
		var sb strings.Builder
		ph := &stmt.PlaceholderGen{Kind: stmt.Simple}

		ph.Write(&sb, 1, "a", true)

		if got := sb.String(); got != "(?, ?, ?)" {
			t.Errorf("unexpected result: %s", got)
		}
		if len(ph.Values()) != 3 {
			t.Errorf("expected 3 values, got %d", len(ph.Values()))
		}
	})

	t.Run("Numbered stmt.", func(t *testing.T) {
		var sb strings.Builder
		ph := &stmt.PlaceholderGen{Kind: stmt.Numbered}

		ph.Write(&sb, 10, 20)

		if got := sb.String(); got != "($1, $2)" {
			t.Errorf("unexpected result: %s", got)
		}
	})

	t.Run("Stringify mode", func(t *testing.T) {
		var sb strings.Builder
		ph := &stmt.PlaceholderGen{Stringify: true}

		ph.Write(&sb, "x", 5, true, 3.14)

		if got := sb.String(); got != "('x', 5, true, 3.14)" {
			t.Errorf("unexpected result: %s", got)
		}
	})

	t.Run("Stringify stress", func(t *testing.T) {
		var sb strings.Builder
		ph := &stmt.PlaceholderGen{Stringify: true}

		vals := []any{
			42,
			"hello",
			3.14,
			true,
			nil,
			int64(100),
			false,
			"world",
		}

		ph.Write(&sb, vals...)

		expect := "(42, 'hello', 3.14, true, NULL, 100, false, 'world')"
		if got := sb.String(); got != expect {
			t.Errorf("unexpected result: %s", got)
		}
	})
}
