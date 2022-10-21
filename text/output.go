package text

import (
	"fmt"
	"log"
	"strings"

	"github.com/RasmusLindroth/i3keys/helpers"
	"github.com/RasmusLindroth/i3keys/i3parse"
	"github.com/RasmusLindroth/i3keys/keyboard"
	"github.com/RasmusLindroth/i3keys/xlib"
)

func printKeyboards(keyboards []keyboard.Keyboard, groups []i3parse.ModifierGroup, prefix string) {
	for kbIndex, kb := range keyboards {
		kbName := kb.Name
		if kbName == "" {
			kbName = "no modifiers"
		}
		dots := "-"
		for i := 0; i < len(kbName); i++ {
			dots = dots + "-"
		}

		fmt.Printf("%s%s:\n%s%s\n", prefix, kbName, prefix, dots)

		for _, keyRow := range kb.Keys {
			var unused []string
			for _, key := range keyRow {
				if !key.InUse {
					unused = append(unused, key.Symbol)
				}
			}
			unusedStr := strings.Join(unused, ", ")
			fmt.Printf("%s%s\n", prefix, unusedStr)
		}
		if kbIndex+1 < len(groups) {
			fmt.Printf("\n\n")
		}
	}
}

// Output prints the keyboards to os.Stdout
func Output(wm string, layout string, filter string) {
	modes, _, err := i3parse.ParseFromRunning(wm, true)

	if err != nil {
		log.Fatalln(err)
	}

	toFilter := filter != ""
	filterMods := helpers.HandleFilterArgs(filter)

	layout = strings.ToUpper(layout)

	modifiers := xlib.GetModifiers()

	for _, mode := range modes {
		groups := i3parse.GetModifierGroups(mode.Bindings)
		var mKeyboards []keyboard.Keyboard

		for _, group := range groups {
			if kb, err := keyboard.MapKeyboard(layout, group, modifiers); err == nil {
				if !toFilter || helpers.CompareSlices(group.Modifiers, filterMods) {
					mKeyboards = append(mKeyboards, kb)
				}
			}
		}

		modeName := mode.Name
		if modeName == "" {
			modeName = "default"
		}
		fmt.Printf("\n\nMode: %s\n", modeName)
		printKeyboards(mKeyboards, groups, "\t")
	}
}
