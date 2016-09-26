package main

import "testing"

func TestNopFormatterDoesNothing(t *testing.T) {
	src := []byte{68, 68, 68}
	fmted, err := NopFormatter().Format(src)
	if err != nil {
		t.Fatal(err)
	}

	if !same(fmted, src) {
		t.Errorf("want %v, got: %v\n", src, fmted)
	}
}

func TestJSONFormatter(t *testing.T) {
	src := []byte(`{"json":true}`)
	want := []byte(`{
  "json": true
}`)

	f := &JSONFormatter{
		Prefix: "",
		Indent: "  ",
	}
	fmted, err := f.Format(src)
	if err != nil {
		t.Fatal(err)
	}

	if !same(fmted, want) {
		t.Errorf("want %v, got: %v\n", want, fmted)
	}
}

func same(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
