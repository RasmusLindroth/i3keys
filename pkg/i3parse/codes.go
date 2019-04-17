package i3parse

import (
	"errors"
	"strconv"

	"github.com/RasmusLindroth/i3keys/pkg/xlib"
)

//CodesToSymbols returns multiple bindings converted from code to symbol binding
func CodesToSymbols(codes []Binding) []Binding {
	for i, code := range codes {
		sym, err := CodeToSymbol(code)

		if err == nil {
			codes[i] = sym
		}
	}

	return codes
}

//CodeToSymbol returns a code binding with the symbol equivalent
func CodeToSymbol(code Binding) (Binding, error) {
	i, err := strconv.Atoi(code.Key)

	if err != nil {
		return Binding{}, errors.New("Couldn't parse string to int")
	}

	hex := xlib.KeyCodeToHex(i)
	if name, ok := xlib.KeySyms[hex]; ok {
		code.Key = name

		return code, nil
	}

	return Binding{}, errors.New("Keycode not in keysymdef.h")
}
