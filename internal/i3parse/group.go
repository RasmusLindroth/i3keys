package i3parse

import "sort"

import "github.com/RasmusLindroth/i3keys/internal/helpers"

type sortByNumModifiers []ModifierGroup

func (a sortByNumModifiers) Len() int      { return len(a) }
func (a sortByNumModifiers) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a sortByNumModifiers) Less(i, j int) bool {
	return len(a[i].Modifiers) < len(a[j].Modifiers)
}

//GetModifierGroups groups bindings that have the same modifiers
func GetModifierGroups(bindings []Binding) []ModifierGroup {
	var groups []ModifierGroup
	for _, binding := range bindings {
		match := false
		for gKey, group := range groups {
			if helpers.CompareSlices(binding.Modifiers, group.Modifiers) {
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
	sort.Sort(sortByNumModifiers(groups))

	return groups
}
