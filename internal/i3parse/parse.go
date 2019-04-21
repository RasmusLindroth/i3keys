package i3parse

import (
	"bufio"
	"errors"
	"io"
	"sort"
	"strings"

	"github.com/RasmusLindroth/i3keys/internal/helpers"
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

func getLineType(parts []string) lineType {
	if len(parts) == 0 {
		return skipLine
	}

	switch parts[0] {
	case "}":
		return unmodeLine
	case "mode":
		if validateMode(parts) {
			return modeLine
		}
	case "set":
		if validateVariable(parts) {
			return variableLine
		}
	case "bindsym":
		if validateBinding(parts) {
			return bindSymLine
		}
	case "bindcode":
		if validateBinding(parts) {
			return bindCodeLine
		}
	}

	return skipLine
}

//BindingType holds a binding type
type BindingType int

const (
	//SymbolBinding = bindsym
	SymbolBinding BindingType = iota
	//CodeBinding = bindcode
	CodeBinding
)

//Binding holds one key binding. Can be keycode or keysymbol
type Binding struct {
	Key       string      `json:"key"`
	Modifiers []string    `json:"modifiers"`
	Command   string      `json:"command"`
	Type      BindingType `json:"type"`
}

//Mode holds i3 bind modes
type Mode struct {
	Name     string    `json:"name"`
	Bindings []Binding `json:"bindings"`
}

//Variable holds one variable in the config file
type Variable struct {
	Name  string `json:"name"`
	Value string `json:"value"`
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

func parseBindingParts(parts []string) ([]string, string, string) {
	var modifiers []string

	index := 1
	if parts[1] == "--release" {
		index = 2
	}

	keys := strings.Split(parts[index], "+")

	key := keys[len(keys)-1]
	modifiers = append(modifiers, keys[:len(keys)-1]...)

	var cmdParts []string
	for i := index + 1; i < len(parts) && parts[i][0] != '#'; i++ {
		cmdParts = append(cmdParts, parts[i])
	}
	cmd := strings.Join(cmdParts, " ")

	return modifiers, key, cmd
}

func parseBinding(parts []string) Binding {
	var bindType BindingType
	switch parts[0] {
	case "bindsym":
		bindType = SymbolBinding
	case "bindcode":
		bindType = CodeBinding
	}

	modifiers, key, cmd := parseBindingParts(parts)

	return Binding{Modifiers: modifiers, Key: key, Command: cmd, Type: bindType}
}

/* TODO:
See i3 documentation for group1,group2,etc.
type Group string
*/

//ParseFromRunning loads config from the running i3 instance
func ParseFromRunning() ([]Mode, []Binding, error) {
	return parse(getConfigFromRunning())
}

//ParseFromFile loads config from path
func ParseFromFile(path string) ([]Mode, []Binding, error) {
	return parse(getConfigFromFile(path))
}

func readLine(reader *bufio.Reader) (string, []string, lineType, error) {
	line, err := reader.ReadString('\n')

	if err != nil && err != io.EOF {
		return "", []string{}, skipLine, err
	}

	var lineParts = helpers.SplitBySpace(line)
	lineType := getLineType(lineParts)

	return line, lineParts, lineType, err
}

func parse(confReader io.Reader, err error) ([]Mode, []Binding, error) {
	if err != nil {
		return []Mode{}, []Binding{}, errors.New("Couldn't get the config file")
	}

	reader := bufio.NewReader(confReader)

	var modes []Mode
	var bindings []Binding
	var variables []Variable

	context := mainContext
	var readErr error
	var line string
	var lineParts []string
	var lineType lineType

	for readErr != io.EOF {
		line, lineParts, lineType, readErr = readLine(reader)

		if readErr != nil && readErr != io.EOF {
			return []Mode{}, []Binding{}, readErr
		}

		switch lineType {
		case skipLine:
			continue
		case variableLine:
			variables = append(variables, parseVariable(lineParts))
		case modeLine:
			context = modeContext
			name := parseMode(line)
			modes = append(modes, Mode{Name: name})
		case unmodeLine:
			context = mainContext
			continue
		}

		bindingLine := lineType == bindSymLine || lineType == bindCodeLine

		if context == mainContext && bindingLine {
			bindings = append(bindings, parseBinding(lineParts))
		}

		if context == modeContext && bindingLine {
			modes[len(modes)-1].Bindings = append(modes[len(modes)-1].Bindings,
				parseBinding(lineParts),
			)
		}
	}

	bindings, modes = replaceVariables(variables, bindings, modes)

	for key := range modes {
		modes[key].Bindings = sortModifiers(modes[key].Bindings)
	}
	bindings = sortModifiers(bindings)

	return modes, bindings, nil
}

func parseVariable(parts []string) Variable {
	return Variable{Name: parts[1], Value: strings.Join(parts[2:], " ")}
}

func variableNameToValue(variables []Variable, value string) string {
	for _, variable := range variables {
		if variable.Name == value {
			return variable.Value
		}
	}

	return value
}

func replaceVariables(variables []Variable, bindings []Binding, modes []Mode) ([]Binding, []Mode) {
	bindings = replaceVariablesInBindings(variables, bindings)
	modes = replaceVariablesInModes(variables, modes)

	return bindings, modes
}

func replaceVariablesInBindings(variables []Variable, bindings []Binding) []Binding {
	for key := range bindings {
		bindings[key].Key = variableNameToValue(variables, bindings[key].Key)

		for mkey := range bindings[key].Modifiers {
			bindings[key].Modifiers[mkey] = variableNameToValue(variables, bindings[key].Modifiers[mkey])
		}
	}

	return bindings
}

func replaceVariablesInModes(variables []Variable, modes []Mode) []Mode {
	for mkey, mode := range modes {
		modes[mkey].Name = variableNameToValue(variables, modes[mkey].Name)
		modes[mkey].Bindings = replaceVariablesInBindings(variables, mode.Bindings)
	}

	return modes
}

func sortModifiers(bindings []Binding) []Binding {
	for key := range bindings {
		sort.Strings(bindings[key].Modifiers)
	}

	return bindings
}
