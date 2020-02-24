package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/RasmusLindroth/i3keys/internal/svg"
	"github.com/RasmusLindroth/i3keys/internal/text"
	"github.com/RasmusLindroth/i3keys/internal/web"
)

const version string = "0.0.11"

func helpText(exitCode int) {
	fmt.Print("Usage:\n\n\ti3keys [-s] <command> [arguments]\n")
	fmt.Print("\tAdd the flag -s for sway\n\n")
	fmt.Print("The commands are:\n\n")
	fmt.Print("\tweb [port]\n\t\tstart the web ui and listen on random port or [port]\n\n")
	fmt.Print("\ttext <layout> [mods]\n\t\toutput available keybindings in the terminal\n\n")
	fmt.Print("\tsvg <layout> [dest] [mods]\n\t\toutputs one SVG file for each modifier group\n\n")
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
	sway := false
	if os.Args[1] == "-s" {
		sway = true
		sIndex = 2
	}
	cmd := os.Args[sIndex]

	port := "-1"
	if cmd == "web" && len(os.Args) >= 2+sIndex {
		port = os.Args[1+sIndex]
	}

	layoutCheck := len(os.Args) > 1+sIndex && (strings.ToUpper(os.Args[1+sIndex]) != "ISO" && strings.ToUpper(os.Args[1+sIndex]) != "ANSI")

	if cmd == "text" && len(os.Args) < 2+sIndex || (cmd == "text" && layoutCheck) {
		fmt.Println("You need to set the <layout> to ISO or ANSI")
		os.Exit(2)
	}

	if (cmd == "svg" && len(os.Args) < 2+sIndex) ||
		(cmd == "svg" && layoutCheck) {
		fmt.Println("You need to set the <layout> to ISO or ANSI")
		os.Exit(2)
	}

	switch cmd {
	case "web":
		web.Output(sway, port)
	case "text":
		if len(os.Args) < 3+sIndex {
			text.Output(sway, os.Args[1+sIndex], "")
		} else {
			text.Output(sway, os.Args[1+sIndex], os.Args[2+sIndex])
		}
	case "svg":
		if len(os.Args) < 3+sIndex {
			svg.Output(sway, os.Args[1+sIndex], "", "")
		} else if len(os.Args) < 4+sIndex {
			svg.Output(sway, os.Args[1+sIndex], os.Args[2+sIndex], "")
		} else {
			svg.Output(sway, os.Args[1+sIndex], os.Args[2+sIndex], os.Args[3+sIndex])
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
