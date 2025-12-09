package hooks

import (
	"testing"
)

func TestMockBody(t *testing.T) {
	myMock := MockBody{}

	output, err := myMock.Read([]byte("Hello"))
	if err != nil {
		t.Errorf("got unexpected error from myMock.Read: %s", err)
	}

	expected := 1
	if output != expected {
		t.Errorf("got wrong output from myMock.Read, expected '%d' and got '%d'", expected, output)
	}

	if err = myMock.Close(); err != nil {
		t.Errorf("got unexpected error from myMock.Close: %s", err)
	}
}
