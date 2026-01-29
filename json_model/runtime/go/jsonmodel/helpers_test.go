package jsonmodel

import (
	"testing"
)

// TestExtendPath verifies the linked-list path nesting
func TestExtendPath(t *testing.T) {
	parent := &Path{Name: "user", Index: -1}
	child := ExtendPath(parent, "email")

	if child.Parent != parent {
		t.Errorf("Expected parent to be %v, got %v", parent, child.Parent)
	}
	if child.Name != "email" {
		t.Errorf("Expected name 'email', got %s", child.Name)
	}
}

// TestSelectPath ensures conditional reporting logic works
func TestSelectPath(t *testing.T) {
	path := &Path{Name: "test"}

	if SelectPath(path, true) != path {
		t.Error("SelectPath(true) should return the path instance")
	}
	if SelectPath(path, false) != nil {
		t.Error("SelectPath(false) should return nil")
	}
}

// TestLen validates the reflection-based length checker
func TestLen(t *testing.T) {
	testCases := []struct {
		input    interface{}
		expected int
	}{
		{map[string]int{"a": 1, "b": 2}, 2},
		{[]string{"one", "two", "three"}, 3},
		{"hello world", 11},
		{42, 0},     // Integers have no length
		{nil, 0},    // Nil has no length
	}

	for _, tc := range testCases {
		res := Len(tc.input)
		if res != tc.expected {
			t.Errorf("Len(%v) failed: expected %d, got %d", tc.input, tc.expected, res)
		}
	}
}

// TestObjectHasPropVal checks safe map extraction
func TestObjectHasPropVal(t *testing.T) {
	data := map[string]interface{}{
		"name": "Hobbes",
		"age":  6,
	}
	var target interface{}

	// Case 1: Key exists
	if !ObjectHasPropVal(data, "name", &target) {
		t.Fatal("Should have found key 'name'")
	}
	if target != "Hobbes" {
		t.Errorf("Expected 'Hobbes', got %v", target)
	}

	// Case 2: Key missing
	if ObjectHasPropVal(data, "gender", &target) {
		t.Error("Should NOT have found key 'gender'")
	}

	// Case 3: Input is not a map
	if ObjectHasPropVal("not-a-map", "name", &target) {
		t.Error("Should return false when input is a string")
	}
}