package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/RasmusLindroth/i3keys/internal/svg"
	"github.com/RasmusLindroth/i3keys/internal/text"
	"github.com/RasmusLindroth/i3keys/internal/web"
)

const version string = "0.0.7"

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
		web.Output(os.Args[2])
	case "text":
		text.Output(os.Args[2])
	case "svg":
		if len(os.Args) < 4 {
			svg.Output(os.Args[2], "")
		} else {
			svg.Output(os.Args[2], os.Args[3])
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
