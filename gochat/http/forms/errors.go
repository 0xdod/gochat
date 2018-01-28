package forms

type Errors map[string][]string

// Add is used to add an error message with a specified key.
func (e Errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

// Get retrieves the first error message for a field.
func (e Errors) Get(field string) string {
	es := e[field]
	if len(es) == 0 {
		return ""
	}
	return es[0]
}
