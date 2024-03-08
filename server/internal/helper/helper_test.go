package helper

import (
	"testing"
)

func TestStrToInt(t *testing.T) {
	var expected int32 = 99
	res := StrToInt("99")
	if res != expected {
		t.Errorf("got %v want %v", res, expected)
	}
}

func TestBToMb(t *testing.T) {
	var expected uint64 = 10
	res := BToMb(10485760)
	if res != expected {
		t.Errorf("got %v want %v", res, expected)
	}
}

func TestHashFileContent(t *testing.T) {
	hash, err := HashFileContent("./../../README.md")
	if err != nil {
		t.Errorf("Error while hashing file content %s", err)
	}
	if hash == "" {
		t.Errorf("Empty hash")
	}
}
