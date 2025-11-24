package hooks

import (
	"testing"
	"fmt"
)

func TestTrimPath(t *testing.T) {
	example := "/api/v1/bar"
	output := trimPath(example)
	expected := "bar"
	if output != expected {
		t.Fatal(fmt.Sprintf("got wrong output from trimPath, expected '%s' and got '%s'", expected, output))
        }
}
