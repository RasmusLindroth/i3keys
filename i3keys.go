package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/RasmusLindroth/i3keys/internal/svg"
	"github.com/RasmusLindroth/i3keys/internal/text"
	"github.com/RasmusLindroth/i3keys/internal/web"
)

const version string = "0.0.8"

func helpText(exitCode int) {
	fmt.Printf("Usage:\n\n\ti3keys <command> [arguments]\n\n")
	fmt.Printf("The commands are:\n\n")
	fmt.Print("\tweb <port>\n\t\tstart the web ui and listen on <port>\n\n")
	fmt.Print("\ttext <layout> [mods]\n\t\toutput available keybindings in the terminal\n\n")
	fmt.Print("\tsvg <layout> [dest] [mods]\n\t\toutputs one SVG file for each modifier group\n\n")
	fmt.Print("\tversion\n\t\tprint i3keys version\n\n")
	fmt.Printf("Arguments:\n\n")
	fmt.Print("\t<layout>\n\t\tis required. Can be ISO or ANSI\n\n")
	fmt.Print("\t[mods]\n\t\tis optional. Can be a single modifier or a group of modifiers. Group them with a plus sign, e.g. Mod4+Ctrl\n\n")
	fmt.Print("\t[dest]\n\t\tis optional. Where to output files, defaults to the current directory\n\n")
	os.Exit(exitCode)
}

func main() {
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
		web.Output(os.Args[2])
	case "text":
		if len(os.Args) < 4 {
			text.Output(os.Args[2], "")
		} else {
			text.Output(os.Args[2], os.Args[3])
		}
	case "svg":
		if len(os.Args) < 4 {
			svg.Output(os.Args[2], "", "")
		} else if len(os.Args) < 5 {
			svg.Output(os.Args[2], os.Args[3], "")
		} else {
			svg.Output(os.Args[2], os.Args[3], os.Args[4])
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
