package forms

import (
	"net/http"
	"net/url"
)

// Form creates a custom form struct, and embeds a url.Values object
type Form struct {
	url.Values
	Errors errors
}

// New initializes a form struct
func New(data url.Values) *Form {
	return &Form{data, errors(make(map[string][]string))}
}

// Has checks if form field is in post and not empty
func (f *Form) Has(field string, r *http.Request) bool {
	x := r.Form.Get(field)

	if x == "" {
		f.Errors.Add(field, "This field can't be empty")
		return false
	}

	return true
}

// Valid returns true if there are no error and the form is valid, otherwise returns false
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
