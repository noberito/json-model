package jsonmodel

import (
	"testing"
)

func TestIsValidDate(t *testing.T) {
	tests := []struct {
		val      any
		expected bool
	}{
		{"2023-10-27", true},
		{"2023-02-28", true},
		{"2023-02-30", false}, // Invalid date (Feb 30)
		{"2023/10/27", false}, // Wrong format
		{"not-a-date", false},
		{12345, false},        // Wrong type
	}

	for _, tt := range tests {
		if got := IsValidDate(tt.val); got != tt.expected {
			t.Errorf("IsValidDate(%v) = %v; want %v", tt.val, got, tt.expected)
		}
	}
}

func TestIsValidDateTime(t *testing.T) {
	tests := []struct {
		val      any
		expected bool
	}{
		{"2023-10-27T10:00:00Z", true},
		{"2023-10-27T10:00:00+01:00", true},
		{"2023-10-27", false}, // Missing time
		{"invalid", false},
	}

	for _, tt := range tests {
		if got := IsValidDateTime(tt.val); got != tt.expected {
			t.Errorf("IsValidDateTime(%v) = %v; want %v", tt.val, got, tt.expected)
		}
	}
}

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		val      any
		expected bool
	}{
		{"test@example.com", true},
		{"user.name@domain.co.uk", true},
		{"plainaddress", false},
		{"@example.com", false},
		{"test@", false},
		{nil, false},
	}

	for _, tt := range tests {
		if got := IsValidEmail(tt.val); got != tt.expected {
			t.Errorf("IsValidEmail(%v) = %v; want %v", tt.val, got, tt.expected)
		}
	}
}

func TestIsValidURL(t *testing.T) {
	tests := []struct {
		val      any
		expected bool
	}{
		{"https://example.com", true},
		{"http://localhost:8080", true},
		{"ftp://files.com", true},
		{"example.com", false}, // Missing scheme
		{"not a url", false},
	}

	for _, tt := range tests {
		if got := IsValidURL(tt.val); got != tt.expected {
			t.Errorf("IsValidURL(%v) = %v; want %v", tt.val, got, tt.expected)
		}
	}
}

func TestIsValidUUID(t *testing.T) {
	tests := []struct {
		val      any
		expected bool
	}{
		{"123e4567-e89b-12d3-a456-426614174000", true},
		{"123e4567-e89b-12d3-a456-4266141740001", false}, // Too long
		{"123e4567e89b12d3a456426614174000", false},      // Missing dashes
		{"zzze4567-e89b-12d3-a456-426614174000", false},  // Invalid hex chars
	}

	for _, tt := range tests {
		if got := IsValidUUID(tt.val); got != tt.expected {
			t.Errorf("IsValidUUID(%v) = %v; want %v", tt.val, got, tt.expected)
		}
	}
}

func TestIsValidRegex(t *testing.T) {
	tests := []struct {
		val      any
		expected bool
	}{
		{"^[a-z]+$", true},
		{"(abc|def)", true},
		{"[", false}, // Invalid regex syntax
		{123, false},
	}

	for _, tt := range tests {
		if got := IsValidRegex(tt.val); got != tt.expected {
			t.Errorf("IsValidRegex(%v) = %v; want %v", tt.val, got, tt.expected)
		}
	}
}