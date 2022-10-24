package helpers

import (
	"sort"
	"strings"
	"unicode"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// TabsToSpaces replaces all tabs with spaces
func TabsToSpaces(s string) string {
	var res string

	for _, c := range s {
		if unicode.IsSpace(c) {
			res += " "
			continue
		}
		res += string(c)
	}

	return res
}

// TrimSpace removes spaces, before and after. And removes newlines
func TrimSpace(s string) string {
	s = strings.TrimSuffix(s, "\n")
	return strings.TrimSpace(s)
}

// SplitBySpace splits by space and remove empty parts
func SplitBySpace(s string) []string {
	s = TabsToSpaces(s)
	var parts []string

	tmpParts := strings.Split(s, " ")
	for _, x := range tmpParts {
		if x != "" {
			parts = append(parts, x)
		}
	}

	return parts
}

// HandleFilterArgs sorts and splits strings like Mod4+Shift
func HandleFilterArgs(s string) []string {
	var a []string
	var b []string

	parts := strings.Split(s, "+")
	for _, item := range parts {
		p := cases.Title(language.Und, cases.NoLower).String(item)
		if len(p) > 2 && p[:3] == "Mod" {
			a = append(a, p)
			continue
		}
		b = append(b, p)
	}
	sort.Strings(a)
	sort.Strings(b)
	return append(a, b...)
}

// CompareSlices compares if two slices are equal
func CompareSlices(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i, val := range a {
		if val != b[i] {
			return false
		}
	}

	return true
}
