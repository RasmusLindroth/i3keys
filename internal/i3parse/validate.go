package i3parse

func validateVariable(parts []string) bool {
	if len(parts) < 2 || parts[1][0] != '$' {
		return false
	}
	return true
}

func validateBinding(parts []string) bool {
	if len(parts) < 3 {
		return false
	}
	return true
}

func validateMode(parts []string) bool {
	if len(parts) > 1 && parts[0] == "mode" &&
		(parts[1] == "--pango_markup" || parts[1][0] == '"') {
		return true
	}
	return false
}

func validateInclude(parts []string) bool {
	if len(parts) < 2 {
		return false
	}
	return true
}
