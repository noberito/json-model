package jsonmodel

import (
	"testing"
)

func TestPathString(t *testing.T) {
	// Test nil path (root)
	var p *Path
	if s := p.String(); s != "$" {
		t.Errorf("expected '$', got '%s'", s)
	}

	// Test object property: $.users
	p1 := &Path{Parent: nil, Name: "users", Index: -1}
	if s := p1.String(); s != "$.users" {
		t.Errorf("expected '$.users', got '%s'", s)
	}

	// Test array index: $.users[0]
	p2 := &Path{Parent: p1, Name: "", Index: 0}
	if s := p2.String(); s != "$.users[0]" {
		t.Errorf("expected '$.users[0]', got '%s'", s)
	}

	// Test nested property: $.users[0].name
	p3 := &Path{Parent: p2, Name: "name", Index: -1}
	if s := p3.String(); s != "$.users[0].name" {
		t.Errorf("expected '$.users[0].name', got '%s'", s)
	}
}

func TestReport(t *testing.T) {
	r := &Report{}
	if r.HasErrors() {
		t.Error("new report should be empty")
	}

	// Simulate an error at $.age
	p := &Path{Parent: nil, Name: "age", Index: -1}
	r.Add("must be positive", p)

	if !r.HasErrors() {
		t.Error("report should have errors after adding one")
	}

	if len(r.Errors) != 1 {
		t.Errorf("expected 1 error, got %d", len(r.Errors))
	}

	expected := "$.age: must be positive"
	if r.Errors[0] != expected {
		t.Errorf("expected '%s', got '%s'", expected, r.Errors[0])
	}
}