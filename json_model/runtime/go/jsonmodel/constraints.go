package jsonmodel

import (
	"unicode/utf8"
)

type Op int

const (
	Eq Op = iota // ==
	Ne           // !=
	Le           // <=
	Lt           // <
	Ge           // >=
	Gt           // >
)

// CheckConstraint applies a constraint (op, limit) to a value v.
// Mirrors jm_check_constraint in json-model.c
func CheckConstraint(v any, op Op, limit float64) bool {
	var val float64

	// Determine what we are comparing (Value vs Size/Length)
	switch t := v.(type) {
	case float64:
		val = t
	case string:
		// Go strings are UTF-8 bytes; utf8.RuneCountInString gives character count
		val = float64(utf8.RuneCountInString(t))
	case []any:
		val = float64(len(t))
	case map[string]any:
		val = float64(len(t))
	default:
		// Constraints don't apply to bool/null in JSON Model
		return false
	}

	// Perform the comparison
	switch op {
	case Eq:
		return val == limit
	case Ne:
		return val != limit
	case Le:
		return val <= limit
	case Lt:
		return val < limit
	case Ge:
		return val >= limit
	case Gt:
		return val > limit
	}
	return false
}