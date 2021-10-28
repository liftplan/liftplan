package fto

import "testing"

func TestFormFields(t *testing.T) {
	t.Parallel()
	f := FormFields()
	if _, err := f.Render(); err != nil {
		t.Error(err)
	}
	if f.Name() != "Beyond 5/3/1" {
		t.Error("unexpected Name")
	}
	if f.ShortCode() != "fto" {
		t.Error("unexpected ShortCode")
	}
	if _, err := f.Elaborate(); err != nil {
		t.Error(err)
	}
}
