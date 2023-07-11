package gateway_test

import (
	"bytes"
	"io"
	"net/http"

	"encoding/json"

	"github.com/NethermindEth/starknet.go/gateway"
)

// httpMock is a mock of the client.
type httpMock struct {
}

type mockResponse struct {
	reader io.Reader
}

func (m *mockResponse) Read(d []byte) (int, error) {
	return m.reader.Read(d)
}

func (m *mockResponse) Close() error {
	return nil
}

func (r *httpMock) Do(req *http.Request) (*http.Response, error) {
	if req.Method == http.MethodGet {
		switch req.URL.Path {
		case "/get_block":
			return get_block(req)
		}
	}
	return nil, nil
}

func get_block(req *http.Request) (*http.Response, error) {
	blockHash := req.URL.Query().Get("blockHash")
	block := gateway.Block{
		BlockHash: blockHash,
	}
	body, _ := json.Marshal(block)
	return &http.Response{
		Body: &mockResponse{
			reader: bytes.NewReader(body),
		},
		StatusCode: 200,
	}, nil
}
