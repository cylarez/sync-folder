package file

import "testing"

func TestNewSharedFile(t *testing.T) {
    path := "test.txt"
    f := NewSharedFile(path)
    if f.Path != path {
        t.Errorf("File not create correctly")
    }
}

func TestGetFileCount(t *testing.T) {
    expected := GetFileCount() + 1
    _ = NewSharedFile("new-test.text")
    val := GetFileCount()
    if val != expected {
        t.Errorf("Wrong file count get %d but expected %d", val, expected)
    }
}
