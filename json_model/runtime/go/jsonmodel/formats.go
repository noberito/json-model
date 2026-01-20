package jsonmodel

import (
	"net/mail"
	"net/url"
	"regexp"
	"time"
)

// IsValidDate checks "YYYY-MM-DD".
func IsValidDate(v any) bool {
	s, ok := v.(string)
	if !ok {
		return false
	}
	_, err := time.Parse("2006-01-02", s)
	return err == nil
}

// IsValidDateTime checks "YYYY-MM-DDTHH:mm:ss...".
func IsValidDateTime(v any) bool {
	s, ok := v.(string)
	if !ok {
		return false
	}
	// Try standard RFC3339 (JSON standard for dates)
	_, err := time.Parse(time.RFC3339, s)
	return err == nil
}

// IsValidEmail checks for basic email structure.
func IsValidEmail(v any) bool {
	s, ok := v.(string)
	if !ok {
		return false
	}
	_, err := mail.ParseAddress(s)
	return err == nil
}

// IsValidURL checks if the string is a valid URI.
func IsValidURL(v any) bool {
	s, ok := v.(string)
	if !ok {
		return false
	}
	_, err := url.ParseRequestURI(s)
	return err == nil
}

// IsValidUUID checks for 8-4-4-4-12 hex structure.
var uuidRegex = regexp.MustCompile(`^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$`)

func IsValidUUID(v any) bool {
	s, ok := v.(string)
	if !ok {
		return false
	}
	return uuidRegex.MatchString(s)
}

// IsValidRegex checks if the string is a valid regex pattern.
func IsValidRegex(v any) bool {
	s, ok := v.(string)
	if !ok {
		return false
	}
	_, err := regexp.Compile(s)
	return err == nil
}