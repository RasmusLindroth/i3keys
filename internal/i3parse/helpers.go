
package i3parse

//KeysToSymbol converts code bindings to symbol bindings
func KeysToSymbol(keys []Binding) []Binding {
	for key, item := range keys {
		if item.Type == CodeBinding {
			res, err := CodeToSymbol(item)
			if err == nil {
				keys[key] = res
			}
		}
	}
	return keys
}