package caigo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type StarknetGateway struct {
	Base         string `json:"base"`
	Feeder       string `json:"feeder"`
	Gateway      string `json:"gateway"`
	ChainId      string `json:"chainId"`
	client       *http.Client
	errorHandler func(e error) error
}

/*
	Instantiate a new StarkNet Gateway client
	- defaults to the GOERLI endpoints
*/
func NewGateway(opts ...GatewayOption) *StarknetGateway {
	gopts := gatewayOptions{
		chainID: GOERLI_ID,
		client:  http.DefaultClient,
	}

	for _, opt := range opts {
		opt.apply(&gopts)
	}

	var sg *StarknetGateway
	switch id := strings.ToLower(gopts.chainID); {
	case strings.Contains("main", id):
		sg = &StarknetGateway{
			Base:    MAINNET_BASE,
			Feeder:  MAINNET_BASE + "/feeder_gateway",
			Gateway: MAINNET_BASE + "/gateway",
			ChainId: MAINNET_ID,
		}
	case strings.Contains("local", id):
		fallthrough
	case strings.Contains("dev", id):
		sg = &StarknetGateway{
			Base:    LOCAL_BASE,
			Feeder:  LOCAL_BASE + "/feeder_gateway",
			Gateway: LOCAL_BASE + "/gateway",
			ChainId: GOERLI_ID,
		}
	default:
		sg = &StarknetGateway{
			Base:    GOERLI_BASE,
			Feeder:  GOERLI_BASE + "/feeder_gateway",
			Gateway: GOERLI_BASE + "/gateway",
			ChainId: GOERLI_ID,
		}
	}

	sg.client = gopts.client

	return sg
}

func (sg *StarknetGateway) newRequest(
	ctx context.Context, method, endpoint string, body interface{},
) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, sg.Feeder+endpoint, nil)
	if err != nil {
		return nil, err
	}

	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal body: %w", err)
		}
		req.Body = ioutil.NopCloser(bytes.NewBuffer(data))
		req.Header.Add("Content-Type", "application/json; charset=utf")
	}
	return req, nil
}

type Error struct {
	StatusCode int    `json:"-"`
	Body       []byte `json:"-"`

	Code    string `json:"code"`
	Message string `json:"message"`
}

// Error implements the error interface.
func (e Error) Error() string {
	return fmt.Sprintf("%d: %s %s", e.StatusCode, e.Code, e.Message)
}

// NewError creates a new Error from an API response.
func NewError(resp *http.Response) error {
	apiErr := Error{StatusCode: resp.StatusCode}
	data, err := ioutil.ReadAll(resp.Body)
	if err == nil && data != nil {
		apiErr.Body = data
		if err := json.Unmarshal(data, &apiErr); err != nil {
			apiErr.Code = "unknown_error_format"
			apiErr.Message = string(data)
		}
	}
	return &apiErr
}

func (sg *StarknetGateway) do(req *http.Request, v interface{}) error {
	resp, err := sg.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close() // nolint: errcheck

	if resp.StatusCode >= 299 {
		e := NewError(resp)
		if sg.errorHandler != nil {
			return sg.errorHandler(e)
		}
		return e
	}

	if v != nil {
		return json.NewDecoder(resp.Body).Decode(v)
	}

	return nil
}

func appendQueryValues(req *http.Request, values url.Values) {
	q := req.URL.Query()
	for k, vs := range values {
		for _, v := range vs {
			q.Add(k, v)
		}
	}
	req.URL.RawQuery = q.Encode()
}
