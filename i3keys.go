package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/RasmusLindroth/i3keys/json"
	"github.com/RasmusLindroth/i3keys/svg"
	"github.com/RasmusLindroth/i3keys/text"
	"github.com/RasmusLindroth/i3keys/web"
)

const version string = "0.0.15"

func helpText(exitCode int) {
	fmt.Print("Usage:\n\n\ti3keys [-i|-s] <command> [arguments]\n")
	fmt.Print("\tAdd the flag -i for i3 and -s for Sway if you don't want autodetection\n\n")
	fmt.Print("The commands are:\n\n")
	fmt.Print("\tweb [port]\n\t\tstart the web ui and listen on random port or [port]\n\n")
	fmt.Print("\tsvg <layout> [dest] [mods]\n\t\toutputs one SVG file for each modifier group\n\n")
	fmt.Print("\tjson <layout>\n\t\toutput all keybindings as json\n\n")
	fmt.Print("\ttext <layout> [mods]\n\t\toutput available keybindings in the terminal\n\n")
	fmt.Print("\tversion\n\t\tprint i3keys version\n\n")
	fmt.Print("Arguments:\n\n")
	fmt.Print("\t<layout>\n\t\tis required. Can be ISO or ANSI\n\n")
	fmt.Print("\t[mods]\n\t\tis optional. Can be a single modifier or a group of modifiers. Group them with a plus sign, e.g. Mod4+Ctrl\n\n")
	fmt.Print("\t[dest]\n\t\tis optional. Where to output files, defaults to the current directory\n\n")
	os.Exit(exitCode)
}

func main() {
	if len(os.Args) == 1 {
		helpText(2)
	}

	sIndex := 1
	wm := ""
	switch os.Args[1] {
	case "-i":
		wm = "i3"
		sIndex = 2
	case "-s":
		wm = "sway"
		sIndex = 2
	}
	cmd := os.Args[sIndex]

	port := "-1"
	if cmd == "web" && len(os.Args) >= 2+sIndex {
		port = os.Args[1+sIndex]
	}

	layout := ""
	if cmd == "text" || cmd == "json" || cmd == "svg" {
		if len(os.Args) <= 1+sIndex {
			fmt.Println("You need to set the <layout> to ISO or ANSI")
			os.Exit(2)
		}
		layout = strings.ToUpper(os.Args[1+sIndex])
		if layout != "ISO" && layout != "ANSI" {
			fmt.Println("You need to set the <layout> to ISO or ANSI")
			os.Exit(2)
		}
	}

	switch cmd {
	case "web":
		web.Output(wm, port)
	case "text":
		if len(os.Args) < 3+sIndex {
			text.Output(wm, os.Args[1+sIndex], "")
		} else {
			text.Output(wm, os.Args[1+sIndex], os.Args[2+sIndex])
		}
	case "json":
		json.Output(wm, layout)
	case "svg":
		if len(os.Args) < 3+sIndex {
			svg.Output(wm, os.Args[1+sIndex], "", "")
		} else if len(os.Args) < 4+sIndex {
			svg.Output(wm, os.Args[1+sIndex], os.Args[2+sIndex], "")
		} else {
			svg.Output(wm, os.Args[1+sIndex], os.Args[2+sIndex], os.Args[3+sIndex])
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
