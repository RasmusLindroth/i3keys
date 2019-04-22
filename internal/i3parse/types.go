package i3parse

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

//ModifierGroup holds bindings with the same modifiers
type ModifierGroup struct {
	Modifiers []string  `json:"modifiers"`
	Bindings  []Binding `json:"bindings"`
}
