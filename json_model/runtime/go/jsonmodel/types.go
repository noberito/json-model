package jsonmodel

import (
	"math"
)

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