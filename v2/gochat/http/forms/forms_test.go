package forms

import (
	"fmt"
	"testing"
)

func TestRequired(t *testing.T) {
	m := map[string][]string{"hello": {"hi"}, "world": {"juice"}}
	f := New(m)
	got := f.AsDiv()
	fmt.Println(got)
}
