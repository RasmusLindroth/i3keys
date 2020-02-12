package i3parse

type lineType uint
type context uint

const (
	skipLine lineType = iota
	variableLine
	bindCodeLine
	bindCodeBracket
	unBindCodeLine
	bindSymLine
	bindSymBracket
	unBindSymLine
	modeLine
	unmodeLine

	mainContext context = iota
	modeContext
	bindCodeMainContext
	bindSymMainContext
	bindCodeModeContext
	bindSymModeContext
)

//Binding holds one key binding. Can only be a keysymbol
type Binding struct {
	Key       string   `json:"key"`
	Modifiers []string `json:"modifiers"`
	Command   string   `json:"command"`
	bindType  string
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

//ModifierGroup holds bindings with the same modifiers
type ModifierGroup struct {
	Modifiers []string  `json:"modifiers"`
	Bindings  []Binding `json:"bindings"`
}
