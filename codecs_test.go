package goffx

import "testing"

func TestStringEncrypt(t *testing.T) {
	e := String("secret-key", "abc", 6)
	result, err := e.Encrypt("aaabbb")
	if err != nil {
		t.Fatalf("Encrypt is error: %v", err)
	}
	if result != "acbacc" {
		t.Fatalf("Encrypt result is %s, expect is acbacc", result)
	}
}

func TestStringDecrypt(t *testing.T) {
	e := String("secret-key", "abc", 6)
	result, err := e.Decrypt("acbacc")
	if err != nil {
		t.Fatalf("Decrypt is error: %v", err)
	}
	if result != "aaabbb" {
		t.Fatalf("Encrypt result is %s, expect is aaabbb", result)
	}
}

func TestIntegerEncrypt(t *testing.T) {
	e := Integer("secret-key", 4)
	result, err := e.Encrypt(1234)
	if err != nil {
		t.Fatalf("Encrypt is error: %v", err)
	}
	if result != 6103 {
		t.Fatalf("Encrypt result is %d, expect is 6103", result)
	}
}

func TestIntegerDecrypt(t *testing.T) {
	e := Integer("secret-key", 4)
	result, err := e.Decrypt(6103)
	if err != nil {
		t.Fatalf("Decrypt is error: %v", err)
	}
	if result != 1234 {
		t.Fatalf("Encrypt result is %d, expect is 1234", result)
	}
}
