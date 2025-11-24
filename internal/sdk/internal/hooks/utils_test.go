package hooks

import (
	"os"
	"testing"
)

func TestIsFile(t *testing.T) {
	//fail on check(directory)
	ok, err := isFile("/var/tmp")
	if ok {
		t.Fatal("isFile() returned true, should have returned false on '/tmp'")
	}
	if err != nil {
		t.Fatal("isFile() returned unexpected error: ", err)
	}

	creds := []byte("hello\ngo\n")
	path := "/var/tmp/.cribl"

	err = os.WriteFile(path, creds, 0644)
	if err != nil {
		t.Errorf("Could not write temporary config file: %s", err)
	}

	//pass on check(file)
	ok, err = isFile(path)
	if !ok {
		t.Fatal("isFile() returned false, should have returned true on ", path)
	}
	if err != nil {
		t.Fatal("isFile() returned unexpected error: ", err)
	}

	os.Remove(path)
}
