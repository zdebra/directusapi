package directusapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"reflect"
)

type request struct {
	ctx    context.Context
	method string
	url    string
	qv     map[string]string
	body   any
}

func (a *API[R, W, PK]) executeRequest(r request, expectedStatus int, dest any) error {
	if dest != nil && reflect.ValueOf(dest).Kind() != reflect.Ptr {
		return fmt.Errorf("dest has to be a pointer")
	}

	var b io.Reader
	if r.body != nil {
		bodyBytes, err := json.Marshal(r.body)
		if err != nil {
			return fmt.Errorf("marshal request body: %w", err)
		}
		b = bytes.NewBuffer(bodyBytes)
	}

	req, err := http.NewRequestWithContext(
		r.ctx,
		r.method,
		r.url,
		b,
	)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	queryValues := url.Values{}

	for k, v := range r.qv {
		queryValues.Set(k, v)
	}

	req.URL.RawQuery = queryValues.Encode()

	req.Header.Set("Authorization", "Bearer "+a.BearerToken)
	req.Header.Set("Content-Type", "application/json")

	if a.debug {
		reqDump, _ := httputil.DumpRequestOut(req, true)
		fmt.Println("--- Request start ---")
		fmt.Println(string(reqDump))
		fmt.Println("--- Request end ---")
	}

	resp, err := a.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("execute request: %v", err)
	}
	defer resp.Body.Close()

	if a.debug {
		respDump, _ := httputil.DumpResponse(resp, true)
		fmt.Println("--- Response start ---")
		fmt.Println(string(respDump))
		fmt.Println("--- Response end ---")
	}

	if resp.StatusCode != expectedStatus {
		respBytes, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status %s: %s", resp.Status, string(respBytes))
	}

	if dest != nil {
		err = json.NewDecoder(resp.Body).Decode(dest)
		if err != nil {
			return fmt.Errorf("decoding json response: %w", err)
		}
	}

	return nil
}
