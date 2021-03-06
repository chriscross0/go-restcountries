package restcountries

import (
	"unicode"
)

// lCFirst lowercases the first character
func lCFirst(s string) string {
	if len(s) == 0 {
		return s
	}

	r := []rune(s)
	r[0] = unicode.ToLower(r[0])
	return string(r)
}

// processFields takes a slice of strings and returns a semicolon delimited string e.g. name,capital -> name;capital;
// processFields performs a lowercase first on the input
// Used when filtering fields with the Fields option
func processFields(fields []string) string {
	out := ""
	for _, field := range fields {
		out = out + lCFirst(field) + ";"
	}

	return out
}

// processCodes takes a slice of strings and returns a semicolon delimited string e.g. CO,GB -> CO;GB;
// Used when filtering fields with the Fields option
func processCodes(fields []string) string {
	out := ""
	for _, field := range fields {
		out = out + field + ";"
	}

	return out
}
