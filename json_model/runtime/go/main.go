package main

import (
	"fmt"
	"jsonmodel/json_model/runtime/go/jsonmodel" // Import the local package created above
)

// CheckPerson is the function you would generate
func CheckPerson(v any, path *jsonmodel.Path, r *jsonmodel.Report) bool {
	// 1. Check Object Type
	obj, ok := v.(map[string]any)
	if !ok {
		r.Add("expected object", path)
		return false
	}

	valid := true

	// 2. Validate "name"
	if val, exists := obj["name"]; exists {
		p := &jsonmodel.Path{Parent: path, Name: "name"}
		if !jsonmodel.IsString(val) {
			r.Add("expected string", p)
			valid = false
		} else {
			// Constraint: Length >= 2
			if !jsonmodel.CheckConstraint(val, jsonmodel.Ge, 2) {
				r.Add("length must be >= 2", p)
				valid = false
			}
		}
	} else {
		r.Add("missing property 'name'", path)
		valid = false
	}

	// 3. Validate "age"
	if val, exists := obj["age"]; exists {
		p := &jsonmodel.Path{Parent: path, Name: "age"}
		if !jsonmodel.IsInteger(val) {
			r.Add("expected integer", p)
			valid = false
		} else {
			// Constraint: Value >= 0
			if !jsonmodel.CheckConstraint(val, jsonmodel.Ge, 0) {
				r.Add("must be >= 0", p)
				valid = false
			}
		}
	}

	// 4. Validate "website" (Optional)
	if val, exists := obj["website"]; exists {
		p := &jsonmodel.Path{Parent: path, Name: "website"}
		if !jsonmodel.IsValidURL(val) {
			r.Add("invalid URL", p)
			valid = false
		}
	}

	return valid
}

func main() {
	// Example Data
	// In Go, JSON numbers become float64, objects become map[string]any
	data := map[string]any{
		"name":    "A",    // Too short (Fail)
		"age":     25.0,   // Integer (Pass)
		"website": "not-a-url", // (Fail)
	}

	report := &jsonmodel.Report{}
	success := CheckPerson(data, nil, report)

	if success {
		fmt.Println("Validation Passed")
	} else {
		fmt.Println("Validation Failed:")
		for _, err := range report.Errors {
			fmt.Println(" -", err)
		}
	}
}