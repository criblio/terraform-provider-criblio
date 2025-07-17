package hooks

import (
	"net/http"
)

type MockHTTPClient struct {
}

func (c MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	client := &http.Client{}
	return client.Do(req)
}

type MockBody struct {
}

func (c MockBody) Read(input []byte) (int, error) {
	return 1, nil
}

func (c MockBody) Close() error {
	return nil
}
