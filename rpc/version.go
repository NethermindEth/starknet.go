package rpc

import (
	"encoding/json"
	"fmt"

	"github.com/Masterminds/semver/v3"
)

func (p *providerWrapper) Version() RPCVersion {
	return p.version
}

// @todo add tests for this type

type RPCVersion int

const (
	RPCVersion9 RPCVersion = iota
	RPCVersion10
)

func (v RPCVersion) String() string {
	switch v {
	case RPCVersion9:
		return "0.9.0"
	case RPCVersion10:
		return "0.10.1"
	}
	return ""
}

func (v RPCVersion) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.String())
}

func (v *RPCVersion) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	semversion, err := semver.NewVersion(s)
	if err != nil {
		return fmt.Errorf("failed to parse version to semver: %w", err)
	}

	switch {
	case semversion.Compare(semver.MustParse(RPCVersion9.String())) == 0:
		*v = RPCVersion9
	case semversion.Compare(semver.MustParse(RPCVersion10.String())) == 0:
		*v = RPCVersion10
	default:
		return fmt.Errorf("invalid RPC version: %s", s)
	}
	return nil
}
