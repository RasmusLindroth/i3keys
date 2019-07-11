package web

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/RasmusLindroth/i3keys/internal/i3parse"
	"github.com/RasmusLindroth/i3keys/internal/keyboard"
	"github.com/RasmusLindroth/i3keys/internal/xlib"
)

type modeKeyboards struct {
	Name      string
	Keyboards []keyboard.Keyboard
}

//Output starts the server at desired port
func Output(port string) {
	modes, keys, err := i3parse.ParseFromRunning()

	if err != nil {
		log.Fatalln(err)
	}

	keys = i3parse.KeysToSymbol(keys)

	modifiers := xlib.GetModifiers()
	groups := i3parse.GetModifierGroups(keys)

	var isoKeyboards []keyboard.Keyboard
	var ansiKeyboards []keyboard.Keyboard
	for _, group := range groups {
		kbISO, err := keyboard.MapKeyboard("ISO", group, modifiers)
		if err != nil {
			log.Fatalln(err)
		}
		isoKeyboards = append(isoKeyboards, kbISO)

		kbANSI, err := keyboard.MapKeyboard("ANSI", group, modifiers)
		if err != nil {
			log.Fatalln(err)
		}
		ansiKeyboards = append(ansiKeyboards, kbANSI)
	}

	var isoModes []modeKeyboards
	var ansiModes []modeKeyboards

	for _, mode := range modes {
		keys := i3parse.KeysToSymbol(mode.Bindings)
		groups := i3parse.GetModifierGroups(keys)

		isoMode := modeKeyboards{Name: mode.Name}
		ansiMode := modeKeyboards{Name: mode.Name}

		for _, group := range groups {
			kbISO, err := keyboard.MapKeyboard("ISO", group, modifiers)
			if err != nil {
				log.Fatalln(err)
			}
			isoMode.Keyboards = append(isoMode.Keyboards, kbISO)

			kbANSI, err := keyboard.MapKeyboard("ANSI", group, modifiers)
			if err != nil {
				log.Fatalln(err)
			}
			ansiMode.Keyboards = append(ansiMode.Keyboards, kbANSI)
		}
		isoModes = append(isoModes, isoMode)
		ansiModes = append(ansiModes, ansiMode)
	}

	ISOkeyboardJSON, _ := json.Marshal(isoKeyboards)
	ANSIkeyboardJSON, _ := json.Marshal(ansiKeyboards)
	ISOmodesJSON, _ := json.Marshal(isoModes)
	ANSImodesJSON, _ := json.Marshal(ansiModes)

	js := fmt.Sprintf(
		"let generatedISO = %s;\n"+
			"let generatedANSI = %s;\n"+
			"let generatedISOmodes = %s;\n"+
			"let generatedANSImodes = %s;\n",
		ISOkeyboardJSON, ANSIkeyboardJSON, ISOmodesJSON, ANSImodesJSON,
	)

	handler := New(js)

	fmt.Printf("Starting server at http://localhost:%s\nGo there "+
		"to see all of your available keys.\n\n", port)
	err = handler.Start(port)
	if err != nil {
		log.Fatalln(err)
	}
}
