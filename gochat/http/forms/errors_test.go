package forms

import "testing"

func TestAdd(t *testing.T) {
	errs := Errors{}
	errs.Add("title", "money")

	if errs["title"][0] != "money" {
		t.Errorf("expected %s, but got %s", "money", errs["title"][0])
	}
}

func TestGet(t *testing.T) {
	errs := Errors{}
	want := "money"
	errs["title"] = []string{want}
	got := errs.Get("title")
	if got != want {
		t.Errorf("expected %q, but got %q", want, got)
	}
}
