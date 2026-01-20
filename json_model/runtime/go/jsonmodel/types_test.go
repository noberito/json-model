package jsonmodel

import (
	"testing"
)

func TestIsInteger(t *testing.T) {
	tests := []struct {
		val      any
		expected bool
	}{
		{1.0, true},       // Standard JSON integer (unmarshaled as float64)
		{100.0, true},     
		{-5.0, true},      
		{0.0, true},       
		{1.5, false},      // Float
		{1.000001, false}, 
		{"1", false},      // String
		{true, false},     // Bool
		{nil, false},      
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