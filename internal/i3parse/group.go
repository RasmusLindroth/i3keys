package i3parse

import "sort"

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

type sortByNumBindings []ModifierGroup

func (a sortByNumBindings) Len() int      { return len(a) }
func (a sortByNumBindings) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a sortByNumBindings) Less(i, j int) bool {
	return len(a[i].Bindings) > len(a[j].Bindings)
}

//GetModifierGroups groups bindings that have the same modifiers
func GetModifierGroups(bindings []Binding) []ModifierGroup {
	var groups []ModifierGroup
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
	sort.Sort(sortByNumBindings(groups))

	return groups
}
