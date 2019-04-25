package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/RasmusLindroth/i3keys/internal/i3parse"
	"github.com/RasmusLindroth/i3keys/internal/keyboard"
	"github.com/RasmusLindroth/i3keys/internal/svg"
	"github.com/RasmusLindroth/i3keys/internal/web"
	"github.com/RasmusLindroth/i3keys/internal/xlib"
)

const version string = "0.0.4"

func keysToSymbol(keys []i3parse.Binding) []i3parse.Binding {
	for key, item := range keys {
		if item.Type == i3parse.CodeBinding {
			res, err := i3parse.CodeToSymbol(item)
			if err == nil {
				keys[key] = res
			}
		}
	}
	return keys
}

func webOutput(port string) {
	_, keys, err := i3parse.ParseFromRunning()

	if err != nil {
		log.Fatalln(err)
	}

	keys = keysToSymbol(keys)

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

	ISOkeyboardJSON, _ := json.Marshal(isoKeyboards)
	ANSIkeyboardJSON, _ := json.Marshal(ansiKeyboards)

	js := fmt.Sprintf("let generatedISO = %s;\nlet generatedANSI = %s;", ISOkeyboardJSON, ANSIkeyboardJSON)

	handler := web.New(js)

	fmt.Printf("Starting server at http://localhost:%s\nGo there "+
		"to see all of your available keys.\n\n", port)
	err = handler.Start(port)
	if err != nil {
		log.Fatalln(err)
	}
}

func textOutput(layout string) {
	_, keys, err := i3parse.ParseFromRunning()

	if err != nil {
		log.Fatalln(err)
	}

	layout = strings.ToUpper(layout)

	keys = keysToSymbol(keys)

	groups := i3parse.GetModifierGroups(keys)
	modifiers := xlib.GetModifiers()

	var keyboards []keyboard.Keyboard
	for _, group := range groups {
		kb, err := keyboard.MapKeyboard(layout, group, modifiers)
		if err == nil {
			keyboards = append(keyboards, kb)
		}
	}

	fmt.Printf("Available keybindings per modifier group\n\n")
	for kbIndex, kb := range keyboards {
		dots := "-"
		for i := 0; i < len(kb.Name); i++ {
			dots = dots + "-"
		}

		fmt.Printf("%s:\n%s\n", kb.Name, dots)

		for _, keyRow := range kb.Keys {
			var unused []string
			for _, key := range keyRow {
				if key.InUse == false {
					unused = append(unused, key.Symbol)
				}
			}
			unusedStr := strings.Join(unused, ", ")
			fmt.Printf("%s\n", unusedStr)
		}
		if kbIndex+1 < len(groups) {
			fmt.Printf("\n\n")
		}
	}
}

func svgOutput(layout string, dest string) {
	_, keys, err := i3parse.ParseFromRunning()

	if err != nil {
		log.Fatalln(err)
	}

	keys = keysToSymbol(keys)

	layout = strings.ToUpper(layout)
	if dest == "" {
		dest = filepath.Join("./")
	}

	groups := i3parse.GetModifierGroups(keys)
	modifiers := xlib.GetModifiers()

	for _, group := range groups {
		kb, err := keyboard.MapKeyboard(layout, group, modifiers)

		if err != nil {
			log.Fatalln(err)
		}

		data := svg.Generate(layout, kb)

		fname := "no-modifiers"
		if len(group.Modifiers) > 0 {
			fname = strings.Join(group.Modifiers, "-")
		}
		fname = fname + "-" + layout + ".svg"

		file, err := os.Create(filepath.Join(dest, fname))

		if err != nil {
			log.Fatalln(err)
		}

		file.Write(data)
	}
}

func helpText(exitCode int) {
	fmt.Printf("Usage:\n\n\ti3keys <command> [arguments]\n\n")
	fmt.Printf("The commands are:\n\n")
	fmt.Println("\tweb <port>            start the web ui and listen on <port>")
	fmt.Println("\ttext <layout>         output available keybindings in the terminal. <layout> can be ISO or ANSI")
	fmt.Println("\tsvg <layout> [dest]   outputs one SVG file for each modifier group. <layout> can be ISO or ANSI, [dest] defaults to current directory")
	fmt.Println("\tversion               print i3keys version")
	os.Exit(exitCode)
}

func main() {
	webCmd := flag.NewFlagSet("web", flag.ExitOnError)
	webCmd.String("port", "", "port to listen on")

	if len(os.Args) == 1 {
		helpText(2)
	}

	cmd := os.Args[1]

	if cmd == "web" && len(os.Args) < 3 {
		fmt.Println("You need to set the <port>")
		os.Exit(2)
	}

	layoutCheck := len(os.Args) > 2 && (strings.ToUpper(os.Args[2]) != "ISO" && strings.ToUpper(os.Args[2]) != "ANSI")

	if cmd == "text" && len(os.Args) < 3 || (cmd == "text" && layoutCheck) {
		fmt.Println("You need to set the <layout> to ISO or ANSI")
		os.Exit(2)
	}

	if (cmd == "svg" && len(os.Args) < 3) ||
		(cmd == "svg" && layoutCheck) {
		fmt.Println("You need to set the <layout> to ISO or ANSI")
		os.Exit(2)
	}

	switch cmd {
	case "web":
		webOutput(os.Args[2])
	case "text":
		textOutput(os.Args[2])
	case "svg":
		if len(os.Args) < 4 {
			svgOutput(os.Args[2], "")
		} else {
			svgOutput(os.Args[2], os.Args[3])
		}
	case "version":
		fmt.Printf("i3keys version %s by Rasmus Lindroth\n", version)
		os.Exit(0)
	case "help":
		helpText(0)
	default:
		helpText(2)
	}
}
