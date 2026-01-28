package main

import (
	"strings"
	"testing"

	// Use the correct import path for your library
	"jsonmodel/json_model/runtime/go/jsonmodel"
)

func TestCheckPerson(t *testing.T) {
	// Define a struct to hold test cases
	tests := []struct {
		name        string // Description of the test case
		input       any    // The JSON data (simulated as map[string]any)
		shouldPass  bool   // Expected boolean result
		errorSubstr string // A snippet of the error message we expect (if failing)
	}{
		{
			name: "Valid Person",
			input: map[string]any{
				"name": "Alice",
				"age":  30.0,
			},
			shouldPass: true,
		},
		{
			name: "Valid Person with Website",
			input: map[string]any{
				"name":    "Bob",
				"age":     25.0,
				"website": "https://example.com",
			},
			shouldPass: true,
		},
		{
			name: "Fail: Name Too Short",
			input: map[string]any{
				"name": "A", // Length < 2
				"age":  25.0,
			},
			shouldPass:  false,
			errorSubstr: "length must be >= 2",
		},
		{
			name: "Fail: Name Wrong Type",
			input: map[string]any{
				"name": 123.0, // Number instead of string
				"age":  25.0,
			},
			shouldPass:  false,
			errorSubstr: "expected string",
		},
		{
			name: "Fail: Missing Name",
			input: map[string]any{
				"age": 20.0,
			},
			shouldPass:  false,
			errorSubstr: "missing property 'name'",
		},
		{
			name: "Fail: Age Negative",
			input: map[string]any{
				"name": "Charlie",
				"age":  -5.0,
			},
			shouldPass:  false,
			errorSubstr: "must be >= 0",
		},
		{
			name: "Fail: Age Wrong Type",
			input: map[string]any{
				"name": "Charlie",
				"age":  "twenty", // String instead of number
			},
			shouldPass:  false,
			errorSubstr: "expected integer",
		},
		{
			name: "Fail: Invalid Website URL",
			input: map[string]any{
				"name":    "Dave",
				"age":     20.0,
				"website": "not-a-url",
			},
			shouldPass:  false,
			errorSubstr: "invalid URL",
		},
		{
			name:       "Fail: Not an Object",
			input:      "Just a string", // Root input is not a map
			shouldPass: false,
			errorSubstr: "expected object",
		},
	}

	// Loop through all test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// 1. Setup Report
			report := &jsonmodel.Report{}

			// 2. Run Function
			// We pass nil for path since these inputs represent the root object
			valid := CheckPerson(tc.input, nil, report)

			// 3. Verify Boolean Result
			if valid != tc.shouldPass {
				t.Errorf("CheckPerson() returned %v, want %v", valid, tc.shouldPass)
			}

			// 4. Verify Error Messages (if failure expected)
			if !tc.shouldPass && tc.errorSubstr != "" {
				// Combine all errors into one string for easy searching
				allErrors := strings.Join(report.Errors, "; ")
				
				if !strings.Contains(allErrors, tc.errorSubstr) {
					t.Errorf("Expected error containing '%s', got: %v", tc.errorSubstr, report.Errors)
				}
			}
		})
	}
}