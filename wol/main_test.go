package main

import "testing"

func TestParsePassword(t *testing.T) {
	if _, err := parsePassword(""); err == nil {
		t.Fatal("should return an error")
	}
	if _, err := parsePassword("abc"); err == nil {
		t.Fatal("should return an error")
	}
	if pass, err := parsePassword("11223344AaBb"); err != nil {
		t.Fatal(err)
	} else if *pass != [6]byte{0x11, 0x22, 0x33, 0x44, 0xAA, 0xBB} {
		t.Fatalf("%X", *pass)
	}
}
