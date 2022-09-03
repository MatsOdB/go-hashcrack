package main

import (
	"testing"
)

func TestGetChars(t *testing.T) {
	if getChars("L") != "abcdefghijklmnopqrstuvwxyz" {
		t.Error("Did not return correct characters")
	}
	if getChars("U") != "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		t.Error("Did not return correct characters")
	}
	if getChars("D") != "0123456789" {
		t.Error("Did not return correct characters")
	}
	if getChars("C") != "\"'|@#&-!;,?.:/%[]{}()<>\\~=+*_$ " {
		t.Error("Did not return correct characters")
	}
	if getChars("LUDC") != "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789\"'|@#&-!;,?.:/%[]{}()<>\\~=+*_$ " {
		t.Error("Did not return correct characters")
	}
}

func TestNext(t *testing.T) {
	allowedChars = getChars("LUDC")
	if next("a", 1) != "b" {
		t.Error("Did not return correct next sequence")
	}
	if next("a", 2) != "c" {
		t.Error("Did not return correct next sequence")
	}
	if next("a", 3) != "d" {
		t.Error("Did not return correct next sequence")
	}
	if next("a", 26) != "A" {
		t.Error("Did not return correct next sequence")
	}
	if next(" ", 1) != "aa" {
		t.Error("Did not return correct next sequence")
	}
}

func TestCrack(t *testing.T) {
	if crack("d7cb62855cc3a04933d835db565be339b4727bab711fb4d7bc277538709b1d32", "LUDC", 3, false) != "psw" {
		t.Error("Did not return correct password")
	}
	if crack("b7fb0394c7183fd5cac17fb41961c826212a185070e4c1d2f4920e51c1dee35f", "LUDC", 3, false) != "pswd" {
		t.Error("Did not return correct password")
	}
	if crack("d74ff0ee8da3b9806b18c877dbf29bbde50b5bd8e4dad7a3a725000feb82e8f1", "LUDC", 3, false) != "pass" {
		t.Error("Did not return correct password")
	}
}
