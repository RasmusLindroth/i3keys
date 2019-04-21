package keyboard

//Keys groups keys that belongs together
var Keys = map[string][]string{
	"0-9": []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"},

	"a-z": []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j",
		"k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x",
		"y", "z"},

	"nordic": []string{"aring", "adiaeresis", "odiaeresis", "oslash", "ae"},

	"arrow": []string{"Left", "Up", "Right", "Down"},

	"common": []string{"BackSpace", "Tab", "Return", "Pause",
		"Scroll_Lock", "Escape", "Delete", "Prior", "Next", "End", "Insert",
		"Menu", "Break", "space", "comma", "period", "slash", "semicolon",
		"backslash", "bracketleft", "bracketright", "plus", "minus", "equal",
		"less", "greater", "apostrophe", "asterisk", "grave", "section"},

	"uncommon": []string{"Clear", "Sys_Req", "Select", "Print", "Begin",
		"Find", "Cancel", "Help", "Execute", "Undo", "Redo"},

	"function": []string{"F1", "F2", "F3", "F4", "F5", "F6", "F7", "F8", "F9",
		"F10"},

	"numpad": []string{"KP_0", "KP_1", "KP_2", "KP_3", "KP_4", "KP_5", "KP_6",
		"KP_7", "KP_8", "KP_9"},

	"numpad_other": []string{"KP_Space", "KP_Tab", "KP_Enter", "KP_F1", "KP_F2",
		"KP_F3", "KP_F4", "KP_Home", "KP_Left", "KP_Up", "KP_Right", "KP_Down",
		"KP_Prior", "KP_Next", "KP_End", "KP_Begin", "KP_Insert", "KP_Delete",
		"KP_Equal", "KP_Multiply", "KP_Add", "KP_Separator", "KP_Subtract",
		"KP_Decimal", "KP_Divide"},
}
