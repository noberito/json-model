package jsonmodel

import (
	"fmt"
	"strings"
)

// Path represents the current location in the JSON tree (Linked List).
// Mirrors jm_path_t in json-model.h
type Path struct {
	Parent   *Path
	Name     string // Key name (for objects)
	Index    int    // Array index (if Name is empty and Index >= 0)
}

func (p *Path) String() string {
	if p == nil {
		return "$"
	}
	var parts []string
	for cur := p; cur != nil; cur = cur.Parent {
		if cur.Name != "" {
			parts = append([]string{"." + cur.Name}, parts...)
		} else if cur.Index >= 0 {
			parts = append([]string{fmt.Sprintf("[%d]", cur.Index)}, parts...)
		}
	}
	return "$" + strings.Join(parts, "")
}

// Report collects validation errors.
// Mirrors jm_report_t in json-model.h
type Report struct {
	Errors []string
}

func (r *Report) Add (msg string, path *Path) {
	// In C, paths are reconstructed on error. Here we format immediately.
	r.Errors = append(r.Errors, fmt.Sprintf("%s: %s", path.String(), msg))
}

func (r *Report) HasErrors() bool {
	return len(r.Errors) > 0
}