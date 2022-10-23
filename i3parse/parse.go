package i3parse

import (
	"bufio"
	"errors"
	"io"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/RasmusLindroth/i3keys/helpers"
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
	case "include":
		if validateInclude(parts) {
			return includeLine
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

// ParseFromRunning loads config from the running i3 instance
func ParseFromRunning(wm string, logError bool) ([]Mode, []Variable, error) {
	switch wm {
	case "i3":
		r, err := getConfigFromRunningi3()
		return parse(r, logError, err)
	case "sway":
		r, err := getConfigFromRunningSway()
		return parse(r, logError, err)
	default:
		r, err := getAutoWM()
		return parse(r, logError, err)
	}
}

// ParseFromFile loads config from path
func ParseFromFile(path string, logError bool) ([]Mode, []Variable, error) {
	r, err := getConfigFromFile(path)
	return parse(r, logError, err)
}

func readLine(reader *bufio.Reader, c context, variables []Variable) (string, []string, lineType, error) {
	line, err := reader.ReadString('\n')

	if err != nil && err != io.EOF {
		return "", []string{}, skipLine, err
	}
	line = strings.TrimSpace(line)
	var lineParts = helpers.SplitBySpace(line)

	for len(lineParts) > 0 && lineParts[len(lineParts)-1] == "\\" {
		nl, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return "", []string{}, skipLine, err
		}
		nl = strings.TrimSpace(nl)
		line = line[:len(line)-1] + nl
		lineParts = helpers.SplitBySpace(line)
	}

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

	sort.Sort(SortByVarbyLen(variables))
	if len(lineParts) > 0 && lineParts[0] != "set" {
		for i := range lineParts {
			for _, v := range variables {
				if strings.Contains(lineParts[i], v.Name) {
					lineParts[i] = strings.ReplaceAll(lineParts[i], v.Name, v.Value)
				}
			}
		}
		lineParts = helpers.SplitBySpace(strings.Join(lineParts, " "))
	}

	lineType := getLineType(lineParts, c)

	return line, lineParts, lineType, err
}

func parseConfig(confReader io.Reader, confPath string, variables []Variable, logError bool, err error) ([]Mode, []Variable, []string, error) {
	if err != nil {
		return []Mode{}, []Variable{}, []string{}, errors.New("couldn't get the config file")
	}

	reader := bufio.NewReader(confReader)

	modes := []Mode{{}}

	var includes []helpers.Include

	context := mainContext
	var readErr error
	var lineParts []string
	var lineType lineType

	for readErr != io.EOF {
		_, lineParts, lineType, readErr = readLine(reader, context, variables)

		if readErr != nil && readErr != io.EOF {
			return []Mode{}, []Variable{}, []string{}, readErr
		}

		switch lineType {
		case skipLine:
			continue
		case variableLine:
			variables = append(variables, parseVariable(lineParts))
			continue
		case modeLine:
			context = modeContext
			name := parseMode(strings.Join(lineParts, " "))
			modes = append(modes, Mode{Name: name})
			continue
		case includeLine:
			inc := helpers.Include{
				ParentPath: confPath,
				Path:       strings.Join(lineParts[1:], " "),
			}
			includes = append(includes, inc)
			continue
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
		if err != nil && logError {
			log.Println(err)
			continue
		}

		isMainContext := context == mainContext || context == bindCodeMainContext || context == bindSymMainContext
		if isMainContext && bindingLine {
			modes[0].Bindings = append(modes[0].Bindings, binding)
		}

		isModeContext := context == modeContext || context == bindCodeModeContext || context == bindSymModeContext
		if isModeContext && bindingLine {
			l := len(modes) - 1
			modes[l].Bindings = append(modes[l].Bindings, binding)
		}
	}

	var includePaths []string
	for _, incl := range includes {
		matches, err := helpers.GetPaths(incl)
		if err != nil && logError {
			log.Printf("couldn't parse the following include \"%s\" got error %v", incl, err)
			continue
		}
		includePaths = append(includePaths, matches...)
	}

	return modes, variables, includePaths, nil
}

func parse(confReader io.Reader, logError bool, err error) ([]Mode, []Variable, error) {
	configPath, _ := helpers.GetSwayDefaultConfig()
	modes, variables, includes, err := parseConfig(confReader, configPath, []Variable{}, logError, err)
	if err != nil {
		return []Mode{}, []Variable{}, errors.New("couldn't get the config file")
	}
	var parsedIncludes []string
	for j := 0; j < len(includes); j++ {
		incl := includes[j]
		done := false
		for _, ap := range parsedIncludes {
			if ap == incl {
				done = true
			}
		}
		if done {
			continue
		}
		f, ferr := os.Open(incl)
		if err != nil && logError {
			log.Printf("couldn't open the included file %s, got err: %v\n", incl, ferr)
		}
		m, v, i, perr := parseConfig(f, incl, variables, logError, err)
		if err != nil && logError {
			log.Printf("couldn't parse the included file %s, got err: %v\n", incl, perr)
		}
		// add modes merging existing bindings
		for iNew := range m {
			found := false
			for iOld := range modes {
				if m[iNew].Name == modes[iOld].Name {
					found = true
					modes[iOld].Bindings = append(modes[iOld].Bindings, m[iNew].Bindings...) // duplicates?
					break
				}
			}
			if !found {
				modes = append(modes, m[iNew])
			}
		}
		variables = v // NOTE: variables are updated in parseConfig
		includes = append(includes, i...)
		parsedIncludes = append(parsedIncludes, incl)
	}

	for key := range modes {
		modes[key].Bindings = sortModifiers(modes[key].Bindings)
	}

	return modes, variables, nil
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

	variable := false
	if len(key) > 0 {
		variable = key[0] == '$'
	}

	var err error
	if bindType == "code" && !variable {
		key, err = CodeToSymbol(key)
	}

	binding := Binding{Modifiers: modifiers, Key: key, Command: cmd, bindType: bindType}
	return binding, err
}

func parseVariable(parts []string) Variable {
	return Variable{Name: parts[1], Value: strings.Join(parts[2:], " ")}
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
