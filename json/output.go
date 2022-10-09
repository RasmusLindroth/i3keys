package json

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/RasmusLindroth/i3keys/i3parse"
	"github.com/RasmusLindroth/i3keys/keyboard"
	"github.com/RasmusLindroth/i3keys/xlib"
)

type modeKeyboards struct {
	Name      string
	Keyboards []keyboard.Keyboard
}

// Output full json for the requested layout
func Output(wm string, layout string) {
	modes, keys, err := i3parse.ParseFromRunning(wm)

	if err != nil {
		log.Fatalln(err)
	}

	modifiers := xlib.GetModifiers()
	groups := i3parse.GetModifierGroups(keys)

	layoutModeDefault := modeKeyboards{Name: "(default)"}
	for _, group := range groups {
		kb, err := keyboard.MapKeyboard(layout, group, modifiers)
		if err != nil {
			log.Fatalln(err)
		}
		layoutModeDefault.Keyboards = append(layoutModeDefault.Keyboards, kb)

	}

	var layoutModes []modeKeyboards
	layoutModes = append(layoutModes, layoutModeDefault)

	for _, mode := range modes {
		groups := i3parse.GetModifierGroups(mode.Bindings)

		layoutMode := modeKeyboards{Name: mode.Name}

		for _, group := range groups {
			kb, err := keyboard.MapKeyboard(layout, group, modifiers)
			if err != nil {
				log.Fatalln(err)
			}
			layoutMode.Keyboards = append(layoutMode.Keyboards, kb)
		}
		layoutModes = append(layoutModes, layoutMode)
	}

	layoutJSON, _ := json.Marshal(layoutModes)
	fmt.Printf("%s\n", layoutJSON)
}
