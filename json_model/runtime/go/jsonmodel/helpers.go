package jsonmodel

import "reflect"

// ExtendPath creates a new path segment for a property
func ExtendPath(p *Path, name string) *Path {
	return &Path{Parent: p, Name: name, Index: -1}
}

// SelectPath returns the current path if the condition is met (used for reporting)
func SelectPath(p *Path, condition bool) *Path {
	if condition {
		return p
	}
	return nil
}

// Len returns the length of a map, slice, or string
func Len(v interface{}) int {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Map, reflect.Slice, reflect.String:
		return rv.Len()
	default:
		return 0
	}
}

// ObjectHasPropVal checks if a key exists in a map and assigns it to dst
func ObjectHasPropVal(obj interface{}, prop string, dst *interface{}) bool {
	m, ok := obj.(map[string]interface{})
	if !ok {
		return false
	}
	val, exists := m[prop]
	if exists {
		*dst = val
	}
	return exists
}