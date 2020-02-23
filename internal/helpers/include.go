package helpers

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

//Equal to os/exec Expand but changed ${} to $()
func replaceDollar(s string, mapping func(string) string) string {
	var buf []byte
	// $() is all ASCII, so bytes are fine for this operation.
	i := 0
	for j := 0; j < len(s); j++ {
		if s[j] == '$' && j+1 < len(s) {
			if buf == nil {
				buf = make([]byte, 0, 2*len(s))
			}
			buf = append(buf, s[i:j]...)
			name, w := getShellName(s[j+1:])
			if name == "" && w > 0 {
				// Encountered invalid syntax; eat the
				// characters.
			} else if name == "" {
				// Valid syntax, but $ was not followed by a
				// name. Leave the dollar character untouched.
				buf = append(buf, s[j])
			} else {
				buf = append(buf, mapping(name)...)
			}
			j += w
			i = j + 1
		}
	}
	if buf == nil {
		return s
	}
	return string(buf) + s[i:]
}

//Equal to os/exec Expand but changed ${} to $()
func getShellName(s string) (string, int) {
	switch {
	case s[0] == '(':
		if len(s) > 2 && isShellSpecialVar(s[1]) && s[2] == ')' {
			return s[1:2], 3
		}
		// Scan to closing brace
		for i := 1; i < len(s); i++ {
			if s[i] == ')' {
				if i == 1 {
					return "", 2 // Bad syntax; eat "$()""
				}
				return s[1:i], i + 1
			}
		}
		return "", 1 // Bad syntax; eat "$("
	case isShellSpecialVar(s[0]):
		return s[0:1], 1
	}
	// Scan alphanumerics.
	var i int
	for i = 0; i < len(s) && isAlphaNum(s[i]); i++ {
	}
	return s[:i], i
}

//Equal to os/exec Expand but changed ${} to $()
func isShellSpecialVar(c uint8) bool {
	switch c {
	case '*', '#', '$', '@', '!', '?', '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return true
	}
	return false
}

//Equal to os/exec Expand but changed ${} to $()
func isAlphaNum(c uint8) bool {
	return c == '_' || '0' <= c && c <= '9' || 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z'
}

func runCMD(program string) string {
	a := strings.Split(program, " ")
	out, err := exec.Command(a[0], a[1:]...).Output()
	if err != nil {
		log.Printf("couldn't expand the following command in include: %s\n", program)
		return ""
	}
	return strings.TrimSuffix(string(out), "\n")
}

func replaceAccent(s string, mapping func(string) string) string {
	r := ""
	for i := 0; i < len(s); i++ {
		if s[i] == '`' && (i > 0 || s[i-1] != '\\') && i+1 < len(s) {
			buf := ""
			closed := false
			for j := i + 1; j < len(s) && !closed; j++ {
				if s[j] != '`' {
					buf += string(s[j])
				} else {
					closed = true
					i = j
				}
			}
			if closed {
				r += mapping(buf)
			}
		} else {
			r += string(s[i])
		}
	}
	return r
}

func ExpandCommand(s string) string {
	s = replaceDollar(s, runCMD)
	s = replaceAccent(s, runCMD)
	return s
}

func checkPath(s string) []string {
	var r []string
	info, err := os.Stat(s)
	if err != nil {
		return []string{}
	}
	if !info.IsDir() {
		cp := filepath.Clean(s)
		if cp[0] == '~' {
			home, _ := os.LookupEnv("HOME")
			cp = home + cp[1:]
		}
		r = append(r, cp)
	}
	return r
}

func GetPaths(s string) ([]string, error) {
	s = ExpandCommand(s)
	matches, err := filepath.Glob(s)
	if err != nil {
		return nil, err
	}
	var paths []string
	for _, m := range matches {
		paths = append(paths, checkPath(m)...)
	}
	return paths, nil
}
