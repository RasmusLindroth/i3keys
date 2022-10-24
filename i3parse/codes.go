package i3parse

import (
	"fmt"
	"strconv"

	"github.com/RasmusLindroth/i3keys/xlib"
)

// CodeToSymbol returns a code binding with the symbol equivalent
func CodeToSymbol(key string) (string, error) {
	i, err := strconv.Atoi(key)

	if err != nil {
		return "", fmt.Errorf("couldn't parse string %s to int", key)
	}

	hex := xlib.KeyCodeToHex(i)
	name, ok := xlib.KeySyms[hex]
	if !ok {
		return "", fmt.Errorf("keycode %s not in keysymdef.h", key)
	}

	return name, nil
}
