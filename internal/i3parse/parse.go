package i3parse

import (
	"bufio"
	"errors"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/RasmusLindroth/i3keys/internal/helpers"
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

	if length > 1 && parts[0] == "mode" &&
		(parts[1] == "--pango_markup" || parts[1][0] == '"') {
		return modeLine
	}

	if length > 0 && parts[0] == "}" {
		return unmodeLine
	}

	if length > 1 {
		switch parts[0] {
		case "set":
			return variableLine
		case "bindsym":
			return bindSymLine
		case "bindcode":
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

func parseKeys(parts []string) ([]string, string, string) {
	var modifiers []string

	length := len(parts)
	if length < 3 {
		return modifiers, "", ""
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

	modifiers, key, cmd := parseKeys(parts)

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

func parse(confReader io.Reader, err error) ([]Mode, []Binding, error) {
	if err != nil {
		return []Mode{}, []Binding{}, errors.New("Couln't get the config file")
	}

	reader := bufio.NewReader(confReader)

	var modes []Mode
	var bindings []Binding
	var variables []Variable

	context := mainContext
	var readErr error
	var line string

	for readErr != io.EOF {
		line, readErr = reader.ReadString('\n')
		line = helpers.TrimSpace(line)

		if readErr != nil && readErr != io.EOF {
			return []Mode{}, []Binding{}, readErr
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
			continue
		}

		if lineType == variableLine && len(lineParts) > 2 && lineParts[1][0] == '$' {
			v := Variable{Name: lineParts[1], Value: strings.Join(lineParts[2:], " ")}
			variables = append(variables, v)
		}

		if lineType == bindSymLine || lineType == bindCodeLine {
			binding := parseBinding(lineParts)

			if context == modeContext {
				modes[len(modes)-1].Bindings = append(
					modes[len(modes)-1].Bindings,
					binding,
				)
			} else if context == mainContext {
				bindings = append(bindings, binding)
			}
		}
	}

	modes = replaceVariablesInModes(variables, modes)
	bindings = replaceVariablesInBindings(variables, bindings)

	for key := range modes {
		modes[key].Bindings = sortModifiers(modes[key].Bindings)
	}
	bindings = sortModifiers(bindings)

	return modes, bindings, nil

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
