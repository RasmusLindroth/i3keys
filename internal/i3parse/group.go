package i3parse

//ModifierGroup holds bindings with the same modifiers
type ModifierGroup struct {
	Modifiers []string  `json:"modifiers"`
	Bindings  []Binding `json:"bindings"`
}

func compareSlices(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i, val := range a {
		if val != b[i] {
			return false
		}
	}

	return true
}

//GetModifierGroups groups bindings that have the same modifiers
func GetModifierGroups(bindings []Binding, groups []ModifierGroup) []ModifierGroup {
	for _, binding := range bindings {
		match := false
		for gKey, group := range groups {
			if compareSlices(binding.Modifiers, group.Modifiers) {
				groups[gKey].Bindings = append(groups[gKey].Bindings, binding)
				match = true
			}
		}

		if !match {
			groups = append(groups, ModifierGroup{
				Modifiers: binding.Modifiers,
				Bindings:  []Binding{binding},
			})
		}
	}

	return groups
}
