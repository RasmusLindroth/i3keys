package web

import (
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/RasmusLindroth/i3keys/i3parse"
	"github.com/RasmusLindroth/i3keys/keyboard"
	"github.com/RasmusLindroth/i3keys/xlib"
)

// Output starts the server at desired port
func Output(wm string, port string) {
	modes, _, err := i3parse.ParseFromRunning(wm, true)

	if err != nil {
		log.Fatalln(err)
	}

	modifiers := xlib.GetModifiers()

	layouts := map[string][]modeKeyboards{"ISO": {}, "ANSI": {}}

	for _, mode := range modes {
		groups := i3parse.GetModifierGroups(mode.Bindings)

		tmpModes := map[string]modeKeyboards{}
		for lt := range layouts {
			tmpModes[lt] = modeKeyboards{Name: mode.Name}
		}

		for _, group := range groups {
			for lt := range layouts {
				kb, err := keyboard.MapKeyboard(lt, group, modifiers)
				if err != nil {
					log.Fatalln(err)
				}
				tmpMode := tmpModes[lt]
				tmpMode.Keyboards = append(tmpMode.Keyboards, kb)
				tmpModes[lt] = tmpMode
			}
		}
		for lt := range layouts {
			layouts[lt] = append(layouts[lt], tmpModes[lt])
		}
	}

	handler := New(layouts)

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
