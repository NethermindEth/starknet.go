package rpc

import (
	"context"
	"encoding/json"

	"github.com/nsf/jsondiff"
)

type spy struct {
	callCloser
	s []byte
}

func NewSpy(client callCloser) *spy {
	return &spy{
		callCloser: client,
		s:          []byte{},
	}
}

func (s *spy) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	raw := json.RawMessage{}
	err := s.callCloser.CallContext(ctx, &raw, method, args...)
	if err != nil {
		return err
	}
	err = json.Unmarshal(raw, &result)
	s.s = []byte(raw)
	return err
}

func (s *spy) Close() {
	s.callCloser.Close()
}

func (s *spy) Compare(o interface{}) (string, error) {
	b, err := json.Marshal(o)
	if err != nil {
		return "", err
	}
	diff, _ := jsondiff.Compare(s.s, b, &jsondiff.Options{})
	return diff.String(), nil
}
