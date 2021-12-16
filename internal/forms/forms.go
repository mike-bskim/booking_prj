package forms

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

// Form creates a custom form struct, embeds a url.Values object
type Form struct {
	url.Values
	Errors errors
}

// Valid returns true if there are no errors, otherwise false
func (f *Form) Valid() bool {
	log.Println("forms.go >> Valid", (f.Errors))
	return len(f.Errors) == 0
}

// New initializes a form struct
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

/*
"first_name"
"last_name"
"email"
"phone"
*/

// Required checks if form field is not empty in post
func (f *Form) Required(fields ...string) {

	for _, field := range fields {
		vaule := f.Get(field)
		if strings.TrimSpace(vaule) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}

}

// Has checks if form field is in post and not empty, , r *http.Request
func (f *Form) Has(field string) bool {
	// x := r.Form.Get(field)
	x := f.Get(field)
	log.Printf("Has >>> [%s]", x)
	if x == "" {
		// f.Errors.Add(field, "This field cannot be blank") // 다른 형식 폼을 검증하기 위해서 오류메시지 부분 삭제.
		return false
	}
	return true
}

// MinLength checks for string minimum length
func (f *Form) MinLength(field string, length int, r *http.Request) bool {
	// x := f.Get(field)
	x := r.Form.Get(field)
	if len(x) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must be more than %d characters", length))
		return false
	}
	return true
}

// IsEmail checks validation of email
func (f *Form) IsEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "Invaid email address")
	}
}
