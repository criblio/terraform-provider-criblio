package hooks

import (
	"testing"
)

func TestHasTokenExpired(t *testing.T) {
	testVar := int64(0)
	if hasTokenExpired(&testVar) != true {
		t.Errorf("hasTokenExpired returned false, wanted true")
	}
}

func TestSDKInit(t *testing.T) {
	myUrl := "foobar"
	var myClient HTTPClient

	credsHook := NewClientCredentialsHook()
	url, _ := credsHook.SDKInit(myUrl, myClient)

	if url != myUrl {
		t.Errorf("creds hook init returned %s, expected %s", url, myUrl)
	}

}
