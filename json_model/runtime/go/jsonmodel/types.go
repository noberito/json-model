package jsonmodel

import (
	"math"
)

// --- Type Checkers ---

// IsInteger checks if v is a number with no fractional part.
func IsInteger(v any) bool {
	f, ok := v.(float64)
	return ok && f == math.Trunc(f)
}

// IsNumber checks if v is any JSON number.
func IsNumber(v any) bool {
	_, ok := v.(float64)
	return ok
}

func IsString(v any) bool {
	_, ok := v.(string)
	return ok
}

func IsBool(v any) bool {
	_, ok := v.(bool)
	return ok
}

func IsArray(v any) bool {
	_, ok := v.([]any)
	return ok
}

func IsObject(v any) bool {
	_, ok := v.(map[string]any)
	return ok
}

// --- Type Converters (Casters) ---

func AsBool(v any) bool {
	return v.(bool)
}

func AsInt(v any) int {
	// JSON numbers are often float64 by default in Go unmarshaling
	if f, ok := v.(float64); ok {
		return int(f)
	}
	// Fallback if data was manually constructed as int
	if i, ok := v.(int); ok {
		return i
	}
	return 0
}

func AsFloat(v any) float64 {
	return v.(float64)
}

func AsString(v any) string {
	return v.(string)
}

func AsArray(v any) []any {
	return v.([]any)
}

func AsObject(v any) map[string]any {
	return v.(map[string]any)
}
