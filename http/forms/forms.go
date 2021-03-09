package forms

import (
	"fmt"
	"net/url"
	"strings"
	"unicode/utf8"
)

type Form struct {
	url.Values
	Errors
}

func New(data url.Values) *Form {
	return &Form{data, Errors{}}
}

// Required is a method that checks that specific fields in the form
// data are present and not blank. If any fields fail this check, add the
// appropriate message to the form errors.
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Values.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

// MaxLength checks to see that a specific field in the form contains
// a maximum number of characters. If the check fails then add the
// appropriate message to the form errors.
func (f *Form) MaxLength(field string, d int) {
	value := f.Values.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) > d {
		f.Errors.Add(field, fmt.Sprintf("This field is too long (maximum is %d characters)", d))
	}
}

// PermittedValues method checks that a specific field in the form matches one of a set of specific
// permitted values. If the check fails
// then add the appropriate message to the form errors.
func (f *Form) PermittedValues(field string, opts ...string) {
	value := f.Values.Get(field)
	if value == "" {
		return
	}
	for _, opt := range opts {
		if value == opt {
			return
		}
	}
	f.Errors.Add(field, "This field is invalid")
}

// Valid returns true if there are no errors.
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// as div
// <div>
//   <p><label></p>
//   <p><label class=errors></p>
//   <p><input /></p>
// </div>

// func (f *Form) AsDiv() string {
// 	str := ""
// 	for k, v := range f.Values {
// 		str += "<div>"
// 		for _, v := range f.Errors {
// 			str += fmt.Sprintf("<p><label>%s</label></p>", v[0])
// 		}

// 		str += fmt.Sprintf("<p><label>%s</label></p>", k)
// 		str += fmt.Sprintf(`<p><input type="text" value=%q></p>`, v[0])
// 		str += "</div>"
// 	}
// 	return str
// }
