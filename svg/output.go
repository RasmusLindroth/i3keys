package svg

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/RasmusLindroth/i3keys/helpers"
	"github.com/RasmusLindroth/i3keys/i3parse"
	"github.com/RasmusLindroth/i3keys/keyboard"
	"github.com/RasmusLindroth/i3keys/xlib"
)

func sanitizeDirName(s string) string {
	r := ""
	prevChar := '-'
	for _, c := range s {
		switch {
		case unicode.IsSpace(c) || unicode.IsPunct(c):
			if prevChar != '-' {
				r += "-"
				prevChar = '-'
			}
		case unicode.IsDigit(c) || unicode.IsLetter(c):
			r += string(c)
			prevChar = c
		}
	}
	r = filepath.Base(filepath.Clean(r))
	return r
}
func createGroup(layout string, dest string, group i3parse.ModifierGroup, modifiers map[string][]string) {
	kb, err := keyboard.MapKeyboard(layout, group, modifiers)

	if err != nil {
		log.Fatalln(err)
	}

	data := Generate(layout, kb)

	fname := "no-modifiers"
	if len(group.Modifiers) > 0 {
		fname = strings.Join(group.Modifiers, "-")
	}
	fname = fname + "-" + layout + ".svg"

	file, err := os.Create(filepath.Join(dest, fname))
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	file.Write(data)
}

// Output generates svg-files of the keyboards at the desired location
func Output(wm string, layout string, dest string, filter string) {
	modes, _, err := i3parse.ParseFromRunning(wm, true)

	if err != nil {
		log.Fatalln(err)
	}

	toFilter := filter != ""
	filterMods := helpers.HandleFilterArgs(filter)

	layout = strings.ToUpper(layout)

	modifiers := xlib.GetModifiers()

	if dest == "" {
		dest = filepath.Join("./")
	}

	err = os.MkdirAll(dest, os.ModePerm)
	if err != nil {
		log.Fatalln(err)
	}

	for i, mode := range modes {
		groups := i3parse.GetModifierGroups(mode.Bindings)

		modeDir := fmt.Sprintf("mode%d", i)
		if len(mode.Name) > 0 {
			modeDir += "-" + sanitizeDirName(mode.Name)
		}

		if len(modeDir) > 50 {
			modeDir = modeDir[0:50]
		}

		mDest := filepath.Join(dest, modeDir)
		err = os.MkdirAll(mDest, os.ModePerm)
		if err != nil {
			log.Fatalln(err)
		}

		if file, err := os.Create(filepath.Join(mDest, "mode-full-name.txt")); err == nil {
			file.Write([]byte(mode.Name))
			file.Close()
		} else {
			log.Fatalln(err)
		}

		for _, group := range groups {
			if !toFilter || helpers.CompareSlices(group.Modifiers, filterMods) {
				createGroup(layout, mDest, group, modifiers)
			}
		}
	}
}
