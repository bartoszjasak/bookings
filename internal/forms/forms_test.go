package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)

	form := New(r.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Error("The post form is not falid")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)

	form := New(r.PostForm)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("The form was valid when Required field was missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "b")
	postedData.Add("c", "c")

	r, _ = http.NewRequest("POST", "/whatever", nil)
	r.PostForm = postedData
	form = New(postedData)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("Shows, doeasnt heve required field when it does")
	}
}

func TestForm_Has(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)

	form := New(r.PostForm)

	has := form.Has("a")
	if has {
		t.Error("returned true when given field is not in form")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	form = New(postedData)
	has = form.Has("a")
	if !has {
		t.Error("Given field is not in form")
	}
}

func TestForm_MinLenght(t *testing.T) {
	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "abcd")
	form := New(postedData)

	isMinLength := form.MinLenght("a", 3)
	if isMinLength {
		t.Error("MinLenght should return false ")
	}

	isMinLength = form.MinLenght("b", 3)
	if !isMinLength {
		t.Error("MinLenght should return true ")
	}
}

func TestForm_IsEmail(t *testing.T) {
	postedData := url.Values{}
	postedData.Add("v", "b@b.com")
	postedData.Add("nv", "abcd")
	form := New(postedData)

	form.IsEmail("v")
	if !form.Valid() {
		t.Error("IsEmail should return true ")
	}

	isError := form.Errors.Get("nv")
	if isError != "" {
		t.Error("Shouldn't have error but got one")
	}

	form.IsEmail("nv")
	if form.Valid() {
		t.Error("IsEmail should return false ")
	}

	isError = form.Errors.Get("nv")
	if isError == "" {
		t.Error("Should have error but didnt got one")
	}
}
