package hooks

import (
	"testing"
)

func TestHasTokenExpired(t *testing.T){
    testVar := int64(0)
    if hasTokenExpired(&testVar) != true {
        t.Errorf("hasTokenExpired returned false, wanted true")
    }
}
