package i3parse

import (
	"bufio"
	"errors"
	"io"
	"log"
	"sort"
	"strings"

	"github.com/RasmusLindroth/i3keys/internal/helpers"
)

func getLineType(parts []string, c context) lineType {
	if len(parts) == 0 {
		return skipLine
	}

	switch parts[0] {
	case "}":
		switch c {
		case modeContext:
			return unmodeLine
		case bindCodeMainContext:
			fallthrough
		case bindCodeModeContext:
			return unBindCodeLine
		case bindSymMainContext:
			fallthrough
		case bindSymModeContext:
			return unBindSymLine
		}
	case "mode":
		if validateMode(parts) {
			return modeLine
		}
	case "set":
		if validateVariable(parts) {
			return variableLine
		}
	case "bindsym":
		if validateBinding(parts) && parts[len(parts)-1] != "{" {
			return bindSymLine
		}
		if parts[len(parts)-1] == "{" {
			return bindSymBracket
		}
	case "bindcode":
		if validateBinding(parts) && parts[len(parts)-1] != "{" {
			return bindCodeLine
		}
		if parts[len(parts)-1] == "{" {
			return bindCodeBracket
		}
	}

	return skipLine
}

/* TODO:
See i3 documentation for group1,group2,etc.
type Group string
*/

//ParseFromRunning loads config from the running i3 instance
func ParseFromRunning(sway bool) ([]Mode, []Binding, error) {
	if sway {
		return parse(getConfigFromRunningSway())
	}
	return parse(getConfigFromRunningi3())
}

//ParseFromFile loads config from path
func ParseFromFile(path string) ([]Mode, []Binding, error) {
	return parse(getConfigFromFile(path))
}

func readLine(reader *bufio.Reader, c context) (string, []string, lineType, error) {
	line, err := reader.ReadString('\n')

	if err != nil && err != io.EOF {
		return "", []string{}, skipLine, err
	}

	var lineParts = helpers.SplitBySpace(line)
	if c == bindCodeMainContext || c == bindSymMainContext ||
		c == bindCodeModeContext || c == bindSymModeContext {
		if len(lineParts) > 0 && lineParts[0] != "}" {
			switch c {
			case bindCodeMainContext:
				fallthrough
			case bindCodeModeContext:
				lineParts = append([]string{"bindcode"}, lineParts...)
			case bindSymMainContext:
				fallthrough
			case bindSymModeContext:
				lineParts = append([]string{"bindsym"}, lineParts...)
			}
		}
	}
	lineType := getLineType(lineParts, c)

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
		line, lineParts, lineType, readErr = readLine(reader, context)

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
		case bindCodeBracket:
			if context == mainContext {
				context = bindCodeMainContext
			} else {
				context = bindCodeModeContext
			}
			continue
		case bindSymBracket:
			if context == mainContext {
				context = bindSymMainContext
			} else {
				context = bindSymModeContext
			}
			continue
		case unmodeLine:
			fallthrough
		case unBindCodeLine:
			fallthrough
		case unBindSymLine:
			if context == bindSymMainContext || context == bindCodeMainContext ||
				context == modeContext {
				context = mainContext
			} else {
				context = modeContext
			}
			continue
		}

		bindingLine := lineType == bindSymLine || lineType == bindCodeLine

		binding, err := parseBinding(lineParts)
		if err != nil {
			log.Println(err)
			continue
		}

		isMainContext := context == mainContext || context == bindCodeMainContext || context == bindSymMainContext
		if isMainContext && bindingLine {
			bindings = append(bindings, binding)
		}

		isModeContext := context == modeContext || context == bindCodeModeContext || context == bindSymModeContext
		if isModeContext && bindingLine {
			modes[len(modes)-1].Bindings = append(modes[len(modes)-1].Bindings,
				binding,
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
	for i := 1; i < len(parts); i++ {
		if !strings.HasPrefix(parts[i], "--") {
			break
		}
		index++
	}

	keys := strings.Split(parts[index], "+")

	key := keys[len(keys)-1]
	modifiers = append(modifiers, keys[:len(keys)-1]...)
	for i, mod := range modifiers {
		if mod[0] == '$' {
			continue
		}
		modifiers[i] = strings.Title(mod)
	}

	var cmdParts []string
	for i := index + 1; i < len(parts) && parts[i][0] != '#'; i++ {
		cmdParts = append(cmdParts, parts[i])
	}
	cmd := strings.Join(cmdParts, " ")

	return modifiers, key, cmd
}

func parseBinding(parts []string) (Binding, error) {
	var bindType string
	switch parts[0] {
	case "bindsym":
		bindType = "symbol"
	case "bindcode":
		bindType = "code"
	}

	modifiers, key, cmd := parseBindingParts(parts)

	variable := key[0] == '$'

	var err error
	if bindType == "code" && variable == false {
		key, err = CodeToSymbol(key)
	}

	binding := Binding{Modifiers: modifiers, Key: key, Command: cmd, bindType: bindType}
	return binding, err
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
	var nb []Binding
	for key := range bindings {
		prev := bindings[key].Key
		bindings[key].Key = variableNameToValue(variables, bindings[key].Key)

		if bindings[key].Key != prev && bindings[key].bindType == "code" {
			k, err := CodeToSymbol(bindings[key].Key)
			if err != nil {
				log.Println(err)
				continue
			}
			bindings[key].Key = k
		}

		for mkey := range bindings[key].Modifiers {
			bindings[key].Modifiers[mkey] = variableNameToValue(variables, bindings[key].Modifiers[mkey])
		}

		nb = append(nb, bindings[key])
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
		var a []string
		var b []string
		for _, m := range bindings[key].Modifiers {
			if len(m) > 2 && m[:3] == "Mod" {
				a = append(a, m)
				continue
			}
			b = append(b, m)
		}
		sort.Strings(a)
		sort.Strings(b)
		bindings[key].Modifiers = append(a, b...)
	}

	return bindings
}
