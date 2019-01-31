package jhttp_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"

	"bitbucket.org/creachadair/jrpc2/handler"
	"bitbucket.org/creachadair/jrpc2/jhttp"
	"bitbucket.org/creachadair/jrpc2/server"
)

func Example() {
	cli, wait := server.Local(handler.Map{
		"Test": handler.New(func(ctx context.Context, ss ...string) (string, error) {
			return strings.Join(ss, " "), nil
		}),
	}, nil)
	defer wait()

	b := jhttp.NewClientBridge(cli)
	defer b.Close()

	hsrv := httptest.NewServer(b)
	defer hsrv.Close()

	rsp, err := http.Post(hsrv.URL, "application/json", strings.NewReader(`{
  "jsonrpc": "2.0",
  "id": 10235,
  "method": "Test",
  "params": ["full", "plate", "and", "packing", "steel"]
}`))
	if err != nil {
		log.Fatalf("POST request failed: %v", err)
	}
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Fatalf("Reading response body: %v", err)
	}

	fmt.Println(string(body))
	// Output:
	// {"jsonrpc":"2.0","id":10235,"result":"full plate and packing steel"}
}