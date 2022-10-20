package keyboard

import (
	"errors"
	"strings"

	"github.com/RasmusLindroth/i3keys/i3parse"
	"github.com/RasmusLindroth/i3keys/xlib"
)

/*
// ANSI holds evdev keys for an ANSI layout

	var ANSI = [][]string{
		{"ESC", "FK01", "FK02", "FK03", "FK04", "FK05", "FK06", "FK07", "FK08", "FK09", "FK10", "FK11", "FK12", "PRSC", "SCLK", "PAUS"},
		{"TLDE", "AE01", "AE02", "AE03", "AE04", "AE05", "AE06", "AE07", "AE08", "AE09", "AE10", "AE11", "AE12", "BKSP", "INS", "HOME", "PGUP", "NMLK", "KPDV", "KPMU", "KPSU"},
		{"TAB", "AD01", "AD02", "AD03", "AD04", "AD05", "AD06", "AD07", "AD08", "AD09", "AD10", "AD11", "AD12", "BKSL", "DELE", "END", "PGDN", "KP7", "KP8", "KP9", "KPAD"},
		{"CAPS", "AC01", "AC02", "AC03", "AC04", "AC05", "AC06", "AC07", "AC08", "AC09", "AC10", "AC11", "RTRN", "KP4", "KP5", "KP6"},
		{"LFSH", "AB01", "AB02", "AB03", "AB04", "AB05", "AB06", "AB07", "AB08", "AB09", "AB10", "RTSH", "UP", "KP1", "KP2", "KP3", "KPEN"},
		{"LCTL", "LWIN", "LALT", "SPCE", "RALT", "RWIN", "MENU", "RCTL", "LEFT", "DOWN", "RGHT", "KP0", "KPDL"},
	}

// ISO holds evdev keys for an ANSI layout

	var ISO = [][]string{
		{"ESC", "FK01", "FK02", "FK03", "FK04", "FK05", "FK06", "FK07", "FK08", "FK09", "FK10", "FK11", "FK12", "PRSC", "SCLK", "PAUS"},
		{"TLDE", "AE01", "AE02", "AE03", "AE04", "AE05", "AE06", "AE07", "AE08", "AE09", "AE10", "AE11", "AE12", "BKSP", "INS", "HOME", "PGUP", "NMLK", "KPDV", "KPMU", "KPSU"},
		{"TAB", "AD01", "AD02", "AD03", "AD04", "AD05", "AD06", "AD07", "AD08", "AD09", "AD10", "AD11", "AD12", "RTRN", "DELE", "END", "PGDN", "KP7", "KP8", "KP9", "KPAD"},
		{"CAPS", "AC01", "AC02", "AC03", "AC04", "AC05", "AC06", "AC07", "AC08", "AC09", "AC10", "AC11", "AC12", "KP4", "KP5", "KP6"},
		{"LFSH", "LSGT", "AB01", "AB02", "AB03", "AB04", "AB05", "AB06", "AB07", "AB08", "AB09", "AB10", "RTSH", "UP", "KP1", "KP2", "KP3", "KPEN"},
		{"LCTL", "LWIN", "LALT", "SPCE", "RALT", "RWIN", "MENU", "RCTL", "LEFT", "DOWN", "RGHT", "KP0", "KPDL"},
	}
*/
type KbMap struct {
	Symbols [][]string
	Sizes   [][]string
}

var kbMaps = map[string]KbMap{
	"ANSI": {
		[][]string{
			{"ESC", "FK01", "FK02", "FK03", "FK04", "FK05", "FK06", "FK07", "FK08", "FK09", "FK10", "FK11", "FK12", "PRSC", "SCLK", "PAUS"},
			{"TLDE", "AE01", "AE02", "AE03", "AE04", "AE05", "AE06", "AE07", "AE08", "AE09", "AE10", "AE11", "AE12", "BKSP", "INS", "HOME", "PGUP", "NMLK", "KPDV", "KPMU", "KPSU"},
			{"TAB", "AD01", "AD02", "AD03", "AD04", "AD05", "AD06", "AD07", "AD08", "AD09", "AD10", "AD11", "AD12", "BKSL", "DELE", "END", "PGDN", "KP7", "KP8", "KP9", "KPAD"},
			{"CAPS", "AC01", "AC02", "AC03", "AC04", "AC05", "AC06", "AC07", "AC08", "AC09", "AC10", "AC11", "RTRN", "KP4", "KP5", "KP6"},
			{"LFSH", "AB01", "AB02", "AB03", "AB04", "AB05", "AB06", "AB07", "AB08", "AB09", "AB10", "RTSH", "UP", "KP1", "KP2", "KP3", "KPEN"},
			{"LCTL", "LWIN", "LALT", "SPCE", "RALT", "RWIN", "MENU", "RCTL", "LEFT", "DOWN", "RGHT", "KP0", "KPDL"},
		},
		[][]string{
			{"single", "emptySingle", "single", "single", "single", "single", "emptySmall", "single", "single", "single", "single", "emptySmall", "single", "single", "single", "single", "emptySmall", "single", "single", "single"},
			{"single", "single", "single", "single", "single", "single", "single", "single", "single", "single", "single", "single", "single", "double", "emptySmall", "single", "single", "single", "emptySmall", "single", "single", "single", "single"},
			{"onehalf", "single", "single", "single", "single", "single", "single", "single", "single", "single", "single", "single", "single", "enterUp", "emptySmall", "single", "single", "single", "emptySmall", "single", "single", "single", "doubleY"},
			{"semidouble", "single", "single", "single", "single", "single", "single", "single", "single", "single", "single", "single", "single", "enterDown", "emptySmall", "emptySingle", "emptySingle", "emptySingle", "emptySmall", "single", "single", "single"},
			{"modifier", "single", "single", "single", "single", "single", "single", "single", "single", "single", "single", "single", "large", "emptySmall", "emptySingle", "single", "emptySingle", "emptySmall", "single", "single", "single", "doubleY"},
			{"modifier", "modifier", "modifier", "space", "modifier", "modifier", "modifier", "modifier", "emptySmall", "single", "single", "single", "emptySmall", "double", "single"},
		},
	},
	"ISO": {
		[][]string{
			{"ESC", "FK01", "FK02", "FK03", "FK04", "FK05", "FK06", "FK07", "FK08", "FK09", "FK10", "FK11", "FK12", "PRSC", "SCLK", "PAUS"},
			{"TLDE", "AE01", "AE02", "AE03", "AE04", "AE05", "AE06", "AE07", "AE08", "AE09", "AE10", "AE11", "AE12", "BKSP", "INS", "HOME", "PGUP", "NMLK", "KPDV", "KPMU", "KPSU"},
			{"TAB", "AD01", "AD02", "AD03", "AD04", "AD05", "AD06", "AD07", "AD08", "AD09", "AD10", "AD11", "AD12", "RTRN", "DELE", "END", "PGDN", "KP7", "KP8", "KP9", "KPAD"},
			{"CAPS", "AC01", "AC02", "AC03", "AC04", "AC05", "AC06", "AC07", "AC08", "AC09", "AC10", "AC11", "AC12", "KP4", "KP5", "KP6"},
			{"LFSH", "LSGT", "AB01", "AB02", "AB03", "AB04", "AB05", "AB06", "AB07", "AB08", "AB09", "AB10", "RTSH", "UP", "KP1", "KP2", "KP3", "KPEN"},
			{"LCTL", "LWIN", "LALT", "SPCE", "RALT", "RWIN", "MENU", "RCTL", "LEFT", "DOWN", "RGHT", "KP0", "KPDL"},
		},
		[][]string{
			{"single", "emptySingle", "single", "single", "single", "single", "emptySmall", "single", "single", "single", "single", "emptySmall", "single", "single", "single", "single", "emptySmall", "single", "single", "single"},
			{"single", "single", "single", "single", "single", "single", "single", "single", "single", "single", "single", "single", "single", "double", "emptySmall", "single", "single", "single", "emptySmall", "single", "single", "single", "single"},
			{"onehalf", "single", "single", "single", "single", "single", "single", "single", "single", "single", "single", "single", "single", "onehalf", "emptySmall", "single", "single", "single", "emptySmall", "single", "single", "single", "doubleY"},
			{"semidouble", "single", "single", "single", "single", "single", "single", "single", "single", "single", "single", "single", "semilarge", "emptySmall", "emptySingle", "emptySingle", "emptySingle", "emptySmall", "single", "single", "single"},
			{"semilarge", "single", "single", "single", "single", "single", "single", "single", "single", "single", "single", "large", "emptySmall", "emptySingle", "single", "emptySingle", "emptySmall", "single", "single", "single", "doubleY"},
			{"modifier", "modifier", "modifier", "space", "modifier", "modifier", "modifier", "modifier", "emptySmall", "single", "single", "single", "emptySmall", "double", "single"},
		},
	},
}

// Keyboard holds one keyboard for one modifier group
type Keyboard struct {
	Name      string
	Modifiers []string
	Keys      [][]Key
}

// Key holds one key. Used for rendering
type Key struct {
	Binding    i3parse.Binding
	Modifier   bool
	InUse      bool
	Symbol     string
	SymbolCode int
	Identifier string
}

func bindingMatch(symbol string, symbolCode int, identifier string, group i3parse.ModifierGroup, modifiers map[string][]string) Key {
	println("bindingMatch: ", symbol, symbolCode, identifier, strings.Join(group.Modifiers, "+"), modifiers)
	rKey := Key{
		Binding:    i3parse.Binding{},
		Modifier:   false,
		InUse:      false,
		Symbol:     symbol,
		SymbolCode: symbolCode,
		Identifier: identifier,
	}

	for _, key := range group.Bindings {
		if symbol == key.Key {
			rKey.InUse = true
			rKey.Binding = key
			return rKey
		}
	}

	for _, modifier := range group.Modifiers {
		if mlist, ok := modifiers[modifier]; ok {
			for _, mval := range mlist {
				if mval == symbol {
					rKey.Modifier = true
					return rKey
				}
			}
		}
	}

	return rKey
}

// MapKeyboard returns a Keyboard matching desired layout
func MapKeyboard(layout string, group i3parse.ModifierGroup, modifiers map[string][]string) (Keyboard, error) {
	var kb [][]Key
	var kbMap, ok = kbMaps[layout]
	if !ok {
		return Keyboard{}, errors.New("no keyboard with that layout")
	}
	kSym := kbMap.Symbols

	defaultVal := "NULL"

	for i := 0; i < len(kSym); i++ {
		var row []Key

		for j := 0; j < len(kSym[i]); j++ {
			symbol := defaultVal
			symbolCode := 0
			if val, ok := xlib.Evdev[kSym[i][j]]; ok {
				symbolCode = xlib.KeyCodeToInt(val)

				if x, ok := xlib.KeySyms[xlib.ToHex(symbolCode)]; ok {
					symbol = x
				}
			}

			k := bindingMatch(symbol, symbolCode, kSym[i][j], group, modifiers)
			row = append(row, k)
		}
		kb = append(kb, row)
	}
	name := strings.Join(group.Modifiers, "+")
	/*
		if name == "" {
			name = "No modifiers"
		}
	*/
	return Keyboard{Name: name, Keys: kb, Modifiers: group.Modifiers}, nil
}
