package jsonmodel

import (
	"testing"
)

// --- Type Checker Tests ---

func TestIsInteger(t *testing.T) {
	tests := []struct {
		val      any
		expected bool
	}{
		{1.0, true},       // Standard JSON integer (unmarshaled as float64)
		{100.0, true},     // Larger integer
		{-5.0, true},      // Negative
		{0.0, true},       // Zero
		{1.5, false},      // Float
		{1.000001, false}, // Precision
		{"1", false},      // String
		{true, false},     // Bool
		{nil, false},      // Nil
	}

	for _, tt := range tests {
		if got := IsInteger(tt.val); got != tt.expected {
			t.Errorf("IsInteger(%v) = %v; want %v", tt.val, got, tt.expected)
		}
	}
}

func TestIsNumber(t *testing.T) {
	tests := []struct {
		val      any
		expected bool
	}{
		{1.0, true},
		{1.5, true},
		{-100.23, true},
		{"1.0", false},
		{nil, false},
	}

	for _, tt := range tests {
		if got := IsNumber(tt.val); got != tt.expected {
			t.Errorf("IsNumber(%v) = %v; want %v", tt.val, got, tt.expected)
		}
	}
}

func TestIsString(t *testing.T) {
	if !IsString("hello") {
		t.Error("IsString('hello') should be true")
	}
	if IsString(123) {
		t.Error("IsString(123) should be false")
	}
}

func TestIsBool(t *testing.T) {
	if !IsBool(true) {
		t.Error("IsBool(true) should be true")
	}
	if !IsBool(false) {
		t.Error("IsBool(false) should be true")
	}
	if IsBool("true") {
		t.Error("IsBool('true') should be false")
	}
}

func TestIsArray(t *testing.T) {
	arr := []any{1.0, 2.0}
	if !IsArray(arr) {
		t.Error("IsArray([]any) should be true")
	}
	if IsArray(map[string]any{}) {
		t.Error("IsArray(map) should be false")
	}
}

func TestIsObject(t *testing.T) {
	obj := map[string]any{"key": "value"}
	if !IsObject(obj) {
		t.Error("IsObject(map) should be true")
	}
	if IsObject([]any{}) {
		t.Error("IsObject([]any) should be false")
	}
}

// --- Type Converter Tests ---

func TestAsInt(t *testing.T) {
	// 1. JSON-style integers (actually float64)
	if got := AsInt(42.0); got != 42 {
		t.Errorf("AsInt(42.0) = %d; want 42", got)
	}
	// 2. Pure Go integers (e.g. manually constructed)
	if got := AsInt(10); got != 10 {
		t.Errorf("AsInt(10) = %d; want 10", got)
	}
	// 3. Zero check
	if got := AsInt(0.0); got != 0 {
		t.Errorf("AsInt(0.0) = %d; want 0", got)
	}
}

func TestAsBool(t *testing.T) {
	if got := AsBool(true); got != true {
		t.Error("AsBool(true) failed")
	}
	if got := AsBool(false); got != false {
		t.Error("AsBool(false) failed")
	}
}

func TestAsFloat(t *testing.T) {
	val := 123.456
	if got := AsFloat(val); got != val {
		t.Errorf("AsFloat(%v) = %v; want %v", val, got, val)
	}
}

func TestAsString(t *testing.T) {
	val := "test string"
	if got := AsString(val); got != val {
		t.Errorf("AsString(%v) = %v; want %v", val, got, val)
	}
}

func TestAsArray(t *testing.T) {
	val := []any{1.0, 2.0}
	got := AsArray(val)
	if len(got) != 2 {
		t.Errorf("AsArray length mismatch: got %d, want 2", len(got))
	}
	if got[0] != 1.0 {
		t.Error("AsArray content mismatch")
	}
}

func TestAsObject(t *testing.T) {
	val := map[string]any{"foo": "bar"}
	got := AsObject(val)
	if len(got) != 1 {
		t.Errorf("AsObject size mismatch: got %d, want 1", len(got))
	}
	if got["foo"] != "bar" {
		t.Error("AsObject content mismatch")
	}
}
