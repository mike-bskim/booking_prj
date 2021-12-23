package forms

import (
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Error("got invalid when should have been valid")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form shows valid when required fields missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "ab")
	postedData.Add("b", "a")
	postedData.Add("c", "a")

	// r, _ = http.NewRequest("POST", "/whatever", nil)
	// r.PostForm = postedData
	form = New(postedData)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("shows does not have required fields when it does")
	}
}

func TestForm_Has(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	has := form.Has("whatever")
	if has {
		t.Error("form shows has field when it does not")
	}

	postedData := url.Values{}
	postedData.Add("a", "ab") // key:value

	form = New(postedData)

	has = form.Has("a")
	if !has {
		t.Error("form shows does not have field when it should")
	}
}

func TestMinLength(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	form.MinLength("whatever", 3)
	if form.Valid() {
		t.Error("form shows min length for non-existent field")
	}

	isError := form.Errors.Get("whatever")
	if isError == "" {
		t.Error("should have an error, but did not get one>>>", isError)
	}
	// MinLength 에서 error 처리하여 minlength 확인필요 없음
	// if minlength {
	// 	t.Error("form shows has field when it does not")
	// }
	postedData := url.Values{}
	postedData.Add("some_field", "12345") // key:value
	form = New(postedData)

	form.MinLength("some_field", 100)
	if form.Valid() {
		t.Error("form shows min length is shorter than 100")
	}

	postedData = url.Values{}
	postedData.Add("some_field", "12345") // key:value
	form = New(postedData)

	form.MinLength("some_field", 1)
	if !form.Valid() {
		t.Error("form shows min length is longer than 1")
	}
	isError = form.Errors.Get("some_field")
	if isError != "" {
		t.Error("should not have an error, but get one")
	}

}

func TestIsEmail(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	form.IsEmail("whatever")
	if form.Valid() {
		t.Error("form shows valid email for non-existent field")
	}

	postedData = url.Values{}
	postedData.Add("email", "a@a.com") // key:value
	form = New(postedData)

	form.IsEmail("email")
	if !form.Valid() {
		t.Error("invaild email")
	}

	postedData = url.Values{}
	postedData.Add("email", "a@com") // key:value
	form = New(postedData)

	form.IsEmail("email")
	if form.Valid() {
		t.Error("Vaild email")
	}

}
