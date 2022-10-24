package json

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/RasmusLindroth/i3keys/i3parse"
	"github.com/RasmusLindroth/i3keys/keyboard"
	"github.com/RasmusLindroth/i3keys/xlib"
)

type KeyboardModesJSON struct {
	Name           string
	ModifierGroups []KeyboardModifierJSON
}

type KeyboardModifierJSON struct {
	Modifiers []string
	Keys      []KeyJSON
}

type KeyJSON struct {
	Row      int
	Modifier bool
	InUse    bool
	Symbol   string
	Command  string
}

// Output full json for the requested layout
func Output(wm string, layout string) {
	modes, _, err := i3parse.ParseFromRunning(wm, false)

	if err != nil {
		log.Fatalln(err)
	}

	modifiers := xlib.GetModifiers()

	var jsonBoards []KeyboardModesJSON
	for _, mode := range modes {
		groups := i3parse.GetModifierGroups(mode.Bindings)

		modeBoard := KeyboardModesJSON{
			Name: mode.Name,
		}

		for _, group := range groups {
			kb, err := keyboard.MapKeyboard(layout, group, modifiers)
			if err != nil {
				log.Fatalln(err)
			}
			b := KeyboardModifierJSON{
				Modifiers: kb.Modifiers,
			}
			for i, keyRow := range kb.Keys {
				for _, key := range keyRow {
					binding := KeyJSON{
						Row:      i,
						Modifier: key.Modifier,
						InUse:    key.InUse,
						Symbol:   key.Symbol,
						Command:  key.Binding.Command,
					}
					b.Keys = append(b.Keys, binding)
				}
			}
			modeBoard.ModifierGroups = append(modeBoard.ModifierGroups, b)
		}
		jsonBoards = append(jsonBoards, modeBoard)
	}
	jb, err := json.Marshal(jsonBoards)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%s\n", jb)
}
