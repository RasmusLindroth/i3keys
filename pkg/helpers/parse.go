package helpers

import (
	"strings"
	"unicode"
)

//TabsToSpaces replaces all tabs with spaces
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

//TrimSpace removes spaces, before and after. And removes newlines
func TrimSpace(s string) string {
	s = strings.TrimSuffix(s, "\n")
	return strings.TrimSpace(s)
}

//SplitBySpace splits by space and remove empty parts
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
