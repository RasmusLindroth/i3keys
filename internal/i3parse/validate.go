package i3parse

func validateVariable(parts []string) bool {
	if len(parts) < 2 || parts[1][0] != '$' {
		return false
	}
	return true
}

func validateBinding(parts []string) bool {
	return len(parts) >= 3
}

func validateMode(parts []string) bool {
	if len(parts) > 1 && parts[0] == "mode" &&
		(parts[1] == "--pango_markup" || parts[1][0] == '"') {
		return true
	}
	return false
}

func validateInclude(parts []string) bool {
	return len(parts) >= 2
}
