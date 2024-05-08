package forms

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
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
		return false
	}

	return true
}

// Required checks if form fields are not empty
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field can't be empty")
		}
	}
}

// MinLenght checks if form field in the post has minimum number oh characters
func (f *Form) MinLength(field string, length int, r *http.Request) bool {
	value := f.Get(field)

	if len(strings.TrimSpace(value)) < length {
		f.Errors.Add(field, fmt.Sprintf("This field should have at least %d characters", length))
		return false
	}

	return true
}

// IsEmail checks for valid email address
func (f *Form) IsEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "Invalid email address")
	}
}

// IsPhone checks for valid phone number
func (f *Form) IsPhone(field string) {
	value := f.Get(field)

	if !govalidator.IsNumeric(value) || !govalidator.StringLength(value, "10", "10") {
		f.Errors.Add(field, "Invalid phone number")
	}
}

// Valid returns true if there are no error and the form is valid, otherwise returns false
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
