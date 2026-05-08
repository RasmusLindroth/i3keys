package xlib

// #cgo LDFLAGS: -lX11 -lXtst
// #include <stdlib.h>
// #include <X11/Xlib.h>
// #include <X11/keysym.h>
// #include <X11/keysymdef.h>
// #include <X11/extensions/XTest.h>
import "C"
import (
	"fmt"
	"os"
	"unsafe"
)

var (
	x11display = C.CString(os.Getenv("DISPLAY"))
	dpy        = C.XOpenDisplay(x11display)
)

//ToHex converts int to hexadecimal with max two 0 padding
func ToHex(val int) string {
	str := fmt.Sprintf("%x", val)

	for len(str) < 4 {
		str = "0" + str
	}

	return "0x" + str
}

/*
KeyCodeToHex returns the matching hexadecimal number. Used to lookup keysyms as
they are written in /usr/include/X11/keysymdef.h
*/
func KeyCodeToHex(code int) string {
	val := KeyCodeToInt(code)

	return ToHex(val)
}

//KeyCodeToInt returns the integer representation of the keycode
func KeyCodeToInt(code int) int {

	var keysymsPerKeycodeReturn C.int
	var keysym *C.ulong
	defer C.XFree(unsafe.Pointer(keysym))

	keysym = C.XGetKeyboardMapping(dpy,
		C.uchar(code),
		1,
		&keysymsPerKeycodeReturn)
	return int(*keysym)
}

//CanUse returns a bool if the X11 display is found
func CanUse() bool {
	return dpy != nil
}

/* Almost all under is taken from https://github.com/Unrud/remote-touchpad/blob/147a4b2874fc87b9a8dddace0005fb8785ae311c/backend_x11.go */

//GetModifiers grabs the name of the keys that corresponds the each modifier
func GetModifiers() map[string][]string {
	var modifierIndices = [...]uint{C.ShiftMapIndex, C.ControlMapIndex, C.Mod1MapIndex, C.Mod2MapIndex, C.Mod3MapIndex, C.Mod4MapIndex, C.Mod5MapIndex}
	modKeymap := C.XGetModifierMapping(dpy)
	defer C.XFreeModifiermap(modKeymap)
	modKeycodes := make(map[uint][]int)
	for _, modIndex := range modifierIndices {
		modKeycodes[1<<uint(modIndex)] = []int{}
		for i := 0; i < int(modKeymap.max_keypermod); i++ {
			keycode := *(*C.KeyCode)(unsafe.Pointer(uintptr(unsafe.Pointer(modKeymap.modifiermap)) +
				uintptr(uint(modIndex)*uint(modKeymap.max_keypermod)+uint(i))))
			if keycode != 0 {
				modKeycodes[1<<uint(modIndex)] = append(modKeycodes[1<<uint(modIndex)], int(C.uint(keycode)))
			}
		}
	}

	modKeyName := make(map[string][]string)
	for key, val := range modKeycodes {
		var name string
		switch key {
		case 1:
			name = "Shift"
		case 4:
			name = "Ctrl"
		case 8:
			name = "Mod1"
		case 16:
			name = "Mod2"
		case 32:
			name = "Mod3"
		case 64:
			name = "Mod4"
		case 128:
			name = "Mod5"
		}
		modKeyName[name] = []string{}

		for i := 0; i < len(val); i++ {
			symbolCode := KeyCodeToHex(val[i])

			if content, ok := KeySyms[symbolCode]; ok {
				modKeyName[name] = append(modKeyName[name], content)
			}
		}
	}

	return modKeyName
}
