package keyboard

import (
	"errors"

	"github.com/RasmusLindroth/i3keys/pkg/xlib"
)

//ANSI holds evdev keys for an ANSI layout
var ANSI = [][]string{
	[]string{"ESC", "FK01", "FK02", "FK03", "FK04", "FK05", "FK06", "FK07", "FK08", "FK09", "FK10", "FK11", "FK12", "PRSC", "SCLK", "PAUS"},
	[]string{"TLDE", "AE01", "AE02", "AE03", "AE04", "AE05", "AE06", "AE07", "AE08", "AE09", "AE10", "AE11", "AE12", "BKSP", "INS", "HOME", "PGUP", "NMLK", "KPDV", "KPMU", "KPSU"},
	[]string{"TAB", "AD01", "AD02", "AD03", "AD04", "AD05", "AD06", "AD07", "AD08", "AD09", "AD10", "AD11", "AD12", "BKSL", "DELE", "END", "PGDN", "KP7", "KP8", "KP9", "KPAD"},
	[]string{"CAPS", "AC01", "AC02", "AC03", "AC04", "AC05", "AC06", "AC07", "AC08", "AC09", "AC10", "AC11", "RTRN", "KP4", "KP5", "KP6"},
	[]string{"LFSH", "AB01", "AB02", "AB03", "AB04", "AB05", "AB06", "AB07", "AB08", "AB09", "AB10", "RTSH", "UP", "KP1", "KP2", "KP3", "KPEN"},
	[]string{"LCTL", "LWIN", "LALT", "SPCE", "RALT", "RWIN", "MENU", "RCTL", "LEFT", "DOWN", "RGHT", "KP0", "KPDL"},
}

//ISO holds evdev keys for an ANSI layout
var ISO = [][]string{
	[]string{"ESC", "FK01", "FK02", "FK03", "FK04", "FK05", "FK06", "FK07", "FK08", "FK09", "FK10", "FK11", "FK12", "PRSC", "SCLK", "PAUS"},
	[]string{"TLDE", "AE01", "AE02", "AE03", "AE04", "AE05", "AE06", "AE07", "AE08", "AE09", "AE10", "AE11", "AE12", "BKSP", "INS", "HOME", "PGUP", "NMLK", "KPDV", "KPMU", "KPSU"},
	[]string{"TAB", "AD01", "AD02", "AD03", "AD04", "AD05", "AD06", "AD07", "AD08", "AD09", "AD10", "AD11", "AD12", "RTRN", "DELE", "END", "PGDN", "KP7", "KP8", "KP9", "KPAD"},
	[]string{"CAPS", "AC01", "AC02", "AC03", "AC04", "AC05", "AC06", "AC07", "AC08", "AC09", "AC10", "AC11", "AC12", "KP4", "KP5", "KP6"},
	[]string{"LFSH", "LSGT", "AB01", "AB02", "AB03", "AB04", "AB05", "AB06", "AB07", "AB08", "AB09", "AB10", "RTSH", "UP", "KP1", "KP2", "KP3", "KPEN"},
	[]string{"LCTL", "LWIN", "LALT", "SPCE", "RALT", "RWIN", "MENU", "RCTL", "LEFT", "DOWN", "RGHT", "KP0", "KPDL"},
}

//Keyboard holds all keys, codes and so on used for rendering
type Keyboard struct {
	Type    string     `json:"type"`
	Keys    [][]string `json:"keys"`
	Content [][]string `json:"content"` //symbol
	Codes   [][]int    `json:"codes"`   //symol int
}

//MapKeyboard returns a Keyboard matching desired layout
func MapKeyboard(layout string) (Keyboard, error) {
	var kb Keyboard

	switch layout {
	case "ANSI":
		kb.Type = "ANSI"
		kb.Keys = ANSI
	case "ISO":
		kb.Type = "ISO"
		kb.Keys = ISO
	default:
		return kb, errors.New("No keyboard with that layout")
	}

	defaultVal := "NULL"

	for i := 0; i < len(kb.Keys); i++ {
		kb.Content = append(kb.Content, []string{})
		kb.Codes = append(kb.Codes, []int{})
		for j := 0; j < len(kb.Keys[i]); j++ {
			value := defaultVal
			code := 0
			if val, ok := xlib.Evdev[kb.Keys[i][j]]; ok {
				code = xlib.KeyCodeToInt(val)

				if x, ok := xlib.KeySyms[xlib.ToHex(code)]; ok {
					value = x
				}
			}
			kb.Content[len(kb.Content)-1] = append(kb.Content[len(kb.Content)-1], value)
			kb.Codes[len(kb.Codes)-1] = append(kb.Codes[len(kb.Codes)-1], code)
		}
	}

	return kb, nil
}
