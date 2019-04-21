package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/RasmusLindroth/i3keys/internal/i3parse"
	"github.com/RasmusLindroth/i3keys/internal/keyboard"
	"github.com/RasmusLindroth/i3keys/internal/web"
	"github.com/RasmusLindroth/i3keys/internal/xlib"
	flag "github.com/spf13/pflag"
)

func main() {
	var port string
	flag.StringVarP(&port, "port", "p", "", "port for the web ui to listen on")
	flag.Parse()

	if port == "" {
		fmt.Fprintf(os.Stderr, "You need to set --port e.g. i3keys --port 8080\n")
		os.Exit(1)
	}

	_, keys, err := i3parse.ParseFromRunning()

	for key, item := range keys {
		if item.Type == i3parse.CodeBinding {
			res, err := i3parse.CodeToSymbol(item)
			if err == nil {
				keys[key] = res
			}
		}
	}

	if err != nil {
		log.Fatalln(err)
	}

	var groups []i3parse.ModifierGroup
	groups = i3parse.GetModifierGroups(keys, groups)

	groupsJSON, _ := json.Marshal(groups)

	kbISO, err := keyboard.MapKeyboard("ISO")

	if err != nil {
		log.Fatalln(err)
	}

	kbANSI, err := keyboard.MapKeyboard("ANSI")

	if err != nil {
		log.Fatalln(err)
	}

	ISOkeyboardJSON, _ := json.Marshal(kbISO)
	ANSIkeyboardJSON, _ := json.Marshal(kbANSI)
	blacklistJSON, _ := json.Marshal(web.Blacklist)
	modifiers, _ := json.Marshal(xlib.GetModifiers())

	js := fmt.Sprintf("let blacklist = %s;\nlet groups = %s;\nlet generatedISO = %s;\nlet generatedANSI = %s;\n let modifierList = %s;", blacklistJSON, groupsJSON, ISOkeyboardJSON, ANSIkeyboardJSON, modifiers)

	handler := web.New(js)

	fmt.Printf("Starting server at http://localhost:%s\nGo there "+
		"to see all of your available keys.", port)
	err = handler.Start(port)
	if err != nil {
		log.Fatalln(err)
	}

}
