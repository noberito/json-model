package jsonmodel

import (
	"testing"
)

func TestCheckConstraint(t *testing.T) {
	tests := []struct {
		name     string
		val      any
		op       Op
		limit    float64
		expected bool
	}{
		// Numeric Value Tests
		{"Number Eq", 5.0, Eq, 5.0, true},
		{"Number Ne", 5.0, Ne, 6.0, true},
		{"Number Gt", 10.0, Gt, 5.0, true},
		{"Number Lt", 3.0, Lt, 5.0, true},
		{"Number Ge", 5.0, Ge, 5.0, true},
		{"Number Le", 5.0, Le, 5.0, true},
		{"Number Fail", 5.0, Gt, 10.0, false},

		// String Length Tests
		{"String Len Eq", "abc", Eq, 3.0, true},
		{"String Len Gt", "hello", Gt, 2.0, true},
		{"String Len Lt", "a", Lt, 5.0, true},
		{"String UTF8", "ñ", Eq, 1.0, true}, // Count runes, not bytes (ñ is 2 bytes, 1 char)

		// Array Size Tests
		{"Array Size Eq", []any{1, 2, 3}, Eq, 3.0, true},
		{"Array Size Ge", []any{1}, Ge, 0.0, true},

		// Object Size Tests
		{"Object Size Eq", map[string]any{"a": 1, "b": 2}, Eq, 2.0, true},

		// Invalid Types
		{"Bool Ignore", true, Eq, 1.0, false},
		{"Nil Ignore", nil, Eq, 0.0, false},
	}

	for _, tt := range tests {
		if got := CheckConstraint(tt.val, tt.op, tt.limit); got != tt.expected {
			t.Errorf("%s: CheckConstraint(%v, %v, %v) = %v; want %v",
				tt.name, tt.val, tt.op, tt.limit, got, tt.expected)
		}
	}
}