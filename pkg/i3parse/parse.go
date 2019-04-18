package i3parse

import (
	"bufio"
	"errors"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/RasmusLindroth/i3keys/pkg/helpers"
	"go.i3wm.org/i3"
)

type lineType uint
type context uint

const (
	skipLine lineType = iota
	variableLine
	bindCodeLine
	bindSymLine
	modeLine
	unmodeLine

	mainContext context = iota
	modeContext
)

//getConfigFromRunning returns the i3 config file
func getConfigFromRunning() (*strings.Reader, error) {
	conf, err := i3.GetConfig()
	if err != nil {
		return nil, err
	}

	return strings.NewReader(conf.Config), nil
}

func getConfigFromFile(path string) (*os.File, error) {
	return os.Open(path)
}

func getLineType(parts []string) lineType {
	length := len(parts)

	if length > 1 && parts[0] == "set" {
		return variableLine
	}

	if length > 1 && parts[0] == "bindsym" {
		return bindSymLine
	}

	if length > 1 && parts[0] == "bindcode" {
		return bindCodeLine
	}

	if length > 1 && parts[0] == "mode" &&
		(parts[1] == "--pango_markup" || parts[1][0] == '"') {
		return modeLine
	}

	if length > 0 && parts[0] == "}" {
		return unmodeLine
	}

	return skipLine
}

//Binding holds one key binding. Can be keycode or keysymbol
type Binding struct {
	Key       string   `json:"key"`
	Modifiers []string `json:"modifiers"`
}

//Mode holds i3 bind modes
type Mode struct {
	Name    string    `json:"name"`
	Symbols []Binding `json:"symbols"`
	Codes   []Binding `json:"codes"`
}

//Variable holds one variable in the config file
type Variable struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

//KeyBindings holds multiple structs, used for return
type KeyBindings struct {
	Modes   []Mode    `json:"modes"`
	Symbols []Binding `json:"symbols"`
	Codes   []Binding `json:"codes"`
}

func parseMode(line string) string {
	start, stop := 0, 0
	for i := 0; i < len(line); i++ {
		if line[i] == '"' && start == 0 {
			start = i + 1
		}

		if line[i] == '"' && start != 0 {
			stop = i
		}
	}
	return line[start:stop]
}

func parseKeys(parts []string) ([]string, string) {

	var modifiers []string

	length := len(parts)
	if length < 3 {
		return modifiers, ""
	}
	index := 1
	if parts[1] == "--release" {
		index = 2
	}

	keys := strings.Split(parts[index], "+")
	klen := len(keys)

	key := keys[klen-1]
	for i, k := range keys {
		if i == klen-1 {
			break
		}
		modifiers = append(modifiers, k)
	}

	return modifiers, key
}

/* TODO:
See i3 documentation for group1,group2,etc.
type Group string
*/

//ParseFromRunning loads config from the running i3 instance
func ParseFromRunning() (KeyBindings, error) {
	return parse(getConfigFromRunning())
}

//ParseFromFile loads config from path
func ParseFromFile(path string) (KeyBindings, error) {
	return parse(getConfigFromFile(path))
}

func parse(confReader io.Reader, err error) (KeyBindings, error) {

	if err != nil {
		return KeyBindings{}, errors.New("Couln't get the config file")
	}

	reader := bufio.NewReader(confReader)

	var modes []Mode
	var syms []Binding
	var codes []Binding
	var variables []Variable

	context := mainContext
	var readErr error
	var line string

	for readErr != io.EOF {
		line, readErr = reader.ReadString('\n')
		line = helpers.TrimSpace(line)

		if readErr != nil && readErr != io.EOF {
			return KeyBindings{}, readErr
		}

		if len(line) == 0 {
			continue
		}

		if line[0] == '#' {
			continue
		}

		var lineParts = helpers.SplitBySpace(line)

		lineType := getLineType(lineParts)

		switch lineType {
		case skipLine:
			continue
		case modeLine:
			context = modeContext
			name := parseMode(line)
			modes = append(modes, Mode{Name: name})
		case unmodeLine:
			context = mainContext
		}

		if lineType == variableLine && len(lineParts) > 2 && lineParts[1][0] == '$' {
			v := Variable{Name: lineParts[1], Value: strings.Join(lineParts[2:], " ")}
			variables = append(variables, v)
		}

		if lineType == bindSymLine || lineType == bindCodeLine {
			modifiers, key := parseKeys(lineParts)

			binding := Binding{Key: key, Modifiers: modifiers}

			if context == modeContext && lineType == bindSymLine {
				modes[len(modes)-1].Symbols = append(
					modes[len(modes)-1].Symbols,
					binding,
				)
			} else if context == modeContext && lineType == bindSymLine {
				modes[len(modes)-1].Codes = append(
					modes[len(modes)-1].Codes,
					binding,
				)
			} else if context == mainContext && lineType == bindSymLine {
				syms = append(syms, binding)
			} else if context == mainContext && lineType == bindCodeLine {
				codes = append(codes, binding)
			}
		}

	}

	allBindings := KeyBindings{
		Modes:   modes,
		Symbols: syms,
		Codes:   codes,
	}

	allBindings = replaceVariables(variables, allBindings)
	allBindings = sortModifiers(allBindings)

	return allBindings, nil

}

func replaceVariablesInBindings(variables []Variable, bindings []Binding) []Binding {
	for _, variable := range variables {

		for key, binding := range bindings {
			if binding.Key == variable.Name {
				bindings[key].Key = variable.Value
			}

			for mkey, mod := range binding.Modifiers {
				if mod == variable.Name {
					bindings[key].Modifiers[mkey] = variable.Value
				}
			}
		}
	}

	return bindings
}

func replaceVariablesInModes(variables []Variable, modes []Mode) []Mode {

	for mkey, mode := range modes {
		for _, variable := range variables {
			if variable.Name == mode.Name {
				modes[mkey].Name = variable.Value
			}
		}

		modes[mkey].Symbols = replaceVariablesInBindings(variables, mode.Symbols)
		modes[mkey].Codes = replaceVariablesInBindings(variables, mode.Codes)
	}

	return modes
}

func replaceVariables(variables []Variable, bindings KeyBindings) KeyBindings {
	bindings.Symbols = replaceVariablesInBindings(variables, bindings.Symbols)
	bindings.Codes = replaceVariablesInBindings(variables, bindings.Codes)
	bindings.Modes = replaceVariablesInModes(variables, bindings.Modes)

	return bindings
}

func sortSingleModifiers(bindings []Binding) []Binding {
	for key := range bindings {
		sort.Strings(bindings[key].Modifiers)
	}

	return bindings
}

func sortModifiers(bindings KeyBindings) KeyBindings {
	bindings.Symbols = sortSingleModifiers(bindings.Symbols)
	bindings.Codes = sortSingleModifiers(bindings.Codes)
	for key := range bindings.Modes {
		bindings.Modes[key].Symbols = sortSingleModifiers(bindings.Modes[key].Symbols)
		bindings.Modes[key].Codes = sortSingleModifiers(bindings.Modes[key].Codes)
	}

	return bindings
}
