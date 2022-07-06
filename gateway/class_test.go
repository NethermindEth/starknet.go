package gateway

import (
	"context"
	"testing"
)

func TestClassByHash(t *testing.T) {
	gw := NewClient()

	class, err := gw.ClassByHash(context.Background(), "0x25ec026985a3bf9d0cc1fe17326b245dfdc3ff89b8fde106542a3ea56c5a918")
	if err != nil {
		t.Errorf("Could not pull class definition by hash: %v", err)
	}

	if len(class.Program) == 0 {
		t.Errorf("Could not unmarshall class program: ")
	}
}

func TestClassHashAt(t *testing.T) {
	gw := NewClient()

	classHash, err := gw.ClassHashAt(context.Background(), "0x0126dd900b82c7fc95e8851f9c64d0600992e82657388a48d3c466553d4d9246")
	if err != nil {
		t.Errorf("Could not pull class hash: %v", err)
	}

	if classHash.String() != "0x25ec026985a3bf9d0cc1fe17326b245dfdc3ff89b8fde106542a3ea56c5a918" {
		t.Errorf("Pulled incorrect class hash: %v", err)
	}
}
