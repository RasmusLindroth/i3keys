package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/RasmusLindroth/i3keys/i3parse"
	"github.com/RasmusLindroth/i3keys/keyboard"
	"github.com/RasmusLindroth/i3keys/xlib"
)

type modeKeyboards struct {
	Name      string
	Keyboards []keyboard.Keyboard
}

// Output starts the server at desired port
func Output(wm string, port string) {
	modes, keys, err := i3parse.ParseFromRunning(wm, true)

	if err != nil {
		log.Fatalln(err)
	}

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
		groups := i3parse.GetModifierGroups(mode.Bindings)

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

	if port == "-1" {
		//get the kernel to give us a free TCP port
		addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
		if err != nil {
			log.Fatalln(err)
		}

		l, err := net.ListenTCP("tcp", addr)
		if err != nil {
			log.Fatalln(err)
		}
		port = strconv.Itoa(l.Addr().(*net.TCPAddr).Port)
		l.Close()
	}

	fmt.Printf("Starting server at http://localhost:%s\nGo there "+
		"to see all of your available keys.\n\n", port)
	err = handler.Start(port)
	if err != nil {
		log.Fatalln(err)
	}
}
