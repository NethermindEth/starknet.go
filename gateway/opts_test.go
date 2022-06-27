package gateway

import (
	"fmt"
	"net/http"
	"testing"
)

func TestWithErrorHandler(t *testing.T) {
	spy := false
	spyPtr := &spy

	opt := WithErrorHandler(func(err error) error {
		*spyPtr = true
		fmt.Print("spy triggered\n")
		return err
	})

	client := NewClient(opt)
	var result string
	request, _ := http.NewRequest(http.MethodGet, "https://httpstat.us/400", nil)
	err := client.do(request, &result)

	if err == nil {
		t.Errorf("should fail due to http-400")
	}
	if !*spyPtr {
		t.Errorf("spy should be true, got %t", *spyPtr)
	}
}
