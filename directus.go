package directusapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

type API[R, W any] struct {
	Scheme         string
	Host           string
	Namespace      string
	CollectionName string
	BearerToken    string
	HTTPClient     *http.Client
}

func (d API[R, W]) CreateToken(ctx context.Context, email, password string) (string, error) {
	u := fmt.Sprintf("%s://%s/%s/auth/authenticate", d.Scheme, d.Host, d.Namespace)

	body := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		email,
		password,
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("marshal auth request: %w", err)
	}

	req, _ := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		u,
		bytes.NewBuffer(bodyBytes),
	)
	req.Header.Set("Content-Type", "application/json")
	resp, err := d.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("directus api execute request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBytes, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("unexpected status %s: %s", resp.Status, string(respBytes))
	}

	var respBody struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return "", fmt.Errorf("decoding json response: %w", err)
	}
	return respBody.Data.Token, nil
}

func (d API[R, W]) TotalItems(ctx context.Context, q query) (int, error) {
	// https://docs.directus.io/reference/query/#aggregation-grouping
	panic("not implemented")
}

func (d API[R, W]) Insert(ctx context.Context, item W) (R, error) {
	var empty R
	u := fmt.Sprintf("%s://%s/%s/items/%s", d.Scheme, d.Host, d.Namespace, d.CollectionName)

	bodyBytes, err := json.Marshal(item)
	if err != nil {
		return empty, fmt.Errorf("marshal item: %w", err)
	}

	req, _ := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		u,
		bytes.NewBuffer(bodyBytes),
	)

	queryValues := url.Values{}

	fields := d.jsonFieldsR()
	queryValues.Set("fields", strings.Join(fields, ","))

	req.URL.RawQuery = queryValues.Encode()

	req.Header.Set("Authorization", "Bearer "+d.BearerToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := d.HTTPClient.Do(req)
	if err != nil {
		return empty, fmt.Errorf("directus api execute request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBytes, _ := ioutil.ReadAll(resp.Body)
		return empty, fmt.Errorf("unexpected status %s: %s", resp.Status, string(respBytes))
	}

	var respBody struct {
		Data R `json:"data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return empty, fmt.Errorf("decoding json response: %w", err)
	}
	return respBody.Data, nil
}

func (d API[R, W]) Create(ctx context.Context, partials map[string]any) (R, error) {
	var empty R
	u := fmt.Sprintf("%s://%s/%s/items/%s", d.Scheme, d.Host, d.Namespace, d.CollectionName)

	bodyBytes, err := json.Marshal(partials)
	if err != nil {
		return empty, fmt.Errorf("marshal partials: %w", err)
	}

	req, _ := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		u,
		bytes.NewBuffer(bodyBytes),
	)

	queryValues := url.Values{}

	fields := d.jsonFieldsR()
	queryValues.Set("fields", strings.Join(fields, ","))

	req.URL.RawQuery = queryValues.Encode()

	req.Header.Set("Authorization", "Bearer "+d.BearerToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := d.HTTPClient.Do(req)
	if err != nil {
		return empty, fmt.Errorf("directus api execute request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBytes, _ := ioutil.ReadAll(resp.Body)
		return empty, fmt.Errorf("unexpected status %s: %s", resp.Status, string(respBytes))
	}

	var respBody struct {
		Data R `json:"data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return empty, fmt.Errorf("decoding json response: %w", err)
	}
	return respBody.Data, nil

}

func (d API[R, W]) GetByID(ctx context.Context, id string) (R, error) {
	u := fmt.Sprintf("%s://%s/%s/items/%s/%s", d.Scheme, d.Host, d.Namespace, d.CollectionName, id)
	req, _ := http.NewRequest(http.MethodGet, u, nil)

	queryValues := url.Values{}

	fields := d.jsonFieldsR()
	queryValues.Set("fields", strings.Join(fields, ","))

	req.URL.RawQuery = queryValues.Encode()

	req.Header.Set("Authorization", "Bearer "+d.BearerToken)

	var empty R
	resp, err := d.HTTPClient.Do(req)
	if err != nil {
		return empty, fmt.Errorf("directus api execute request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBytes, _ := ioutil.ReadAll(resp.Body)
		return empty, fmt.Errorf("unexpected status %s: %s", resp.Status, string(respBytes))
	}

	var respBody struct {
		Data R `json:"data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return empty, fmt.Errorf("decoding json response: %w", err)
	}
	return respBody.Data, nil
}

func (d API[R, W]) Update(ctx context.Context, id string, partials map[string]any) (R, error) {
	var empty R
	u := fmt.Sprintf("%s://%s/%s/items/%s/%s", d.Scheme, d.Host, d.Namespace, d.CollectionName, id)

	bodyBytes, err := json.Marshal(partials)
	if err != nil {
		return empty, fmt.Errorf("marshal partials: %w", err)
	}

	req, _ := http.NewRequestWithContext(
		ctx,
		http.MethodPatch,
		u,
		bytes.NewBuffer(bodyBytes),
	)

	queryValues := url.Values{}

	fields := d.jsonFieldsR()
	queryValues.Set("fields", strings.Join(fields, ","))

	req.URL.RawQuery = queryValues.Encode()

	req.Header.Set("Authorization", "Bearer "+d.BearerToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := d.HTTPClient.Do(req)
	if err != nil {
		return empty, fmt.Errorf("directus api execute request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBytes, _ := ioutil.ReadAll(resp.Body)
		return empty, fmt.Errorf("unexpected status %s: %s", resp.Status, string(respBytes))
	}

	var respBody struct {
		Data R `json:"data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return empty, fmt.Errorf("decoding json response: %w", err)
	}
	return respBody.Data, nil
}

func (d API[R, W]) Set(ctx context.Context, id string, item W) (R, error) {
	var empty R
	u := fmt.Sprintf("%s://%s/%s/items/%s/%s", d.Scheme, d.Host, d.Namespace, d.CollectionName, id)

	bodyBytes, err := json.Marshal(item)
	if err != nil {
		return empty, fmt.Errorf("marshal item: %w", err)
	}

	req, _ := http.NewRequestWithContext(
		ctx,
		http.MethodPatch,
		u,
		bytes.NewBuffer(bodyBytes),
	)

	queryValues := url.Values{}

	fields := d.jsonFieldsR()
	queryValues.Set("fields", strings.Join(fields, ","))

	req.URL.RawQuery = queryValues.Encode()

	req.Header.Set("Authorization", "Bearer "+d.BearerToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := d.HTTPClient.Do(req)
	if err != nil {
		return empty, fmt.Errorf("directus api execute request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBytes, _ := ioutil.ReadAll(resp.Body)
		return empty, fmt.Errorf("unexpected status %s: %s", resp.Status, string(respBytes))
	}

	var respBody struct {
		Data R `json:"data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return empty, fmt.Errorf("decoding json response: %w", err)
	}
	return respBody.Data, nil
}

func (d API[R, W]) Delete(ctx context.Context, id string) error {
	u := fmt.Sprintf("%s://%s/%s/items/%s/%s", d.Scheme, d.Host, d.Namespace, d.CollectionName, id)
	req, _ := http.NewRequest(http.MethodDelete, u, nil)

	req.Header.Set("Authorization", "Bearer "+d.BearerToken)

	resp, err := d.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("directus api execute request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		respBytes, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status %s: %s", resp.Status, string(respBytes))
	}
	return nil
}

func (d API[R, W]) Items(ctx context.Context, q query) ([]R, error) {
	u := fmt.Sprintf("%s://%s/%s/items/%s", d.Scheme, d.Host, d.Namespace, d.CollectionName)
	req, _ := http.NewRequest(http.MethodGet, u, nil)

	queryValues := url.Values{}

	fields := d.jsonFieldsR()
	queryValues.Set("fields", strings.Join(fields, ","))
	buildQuery(queryValues, q)

	req.URL.RawQuery = queryValues.Encode()

	req.Header.Set("Authorization", "Bearer "+d.BearerToken)

	resp, err := d.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("directus api query: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBytes, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status %s: %s", resp.Status, string(respBytes))
	}

	var respBody struct {
		Data []R `json:"data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return nil, fmt.Errorf("decoding json response: %w", err)
	}
	return respBody.Data, nil
}

func buildQuery(queryValues url.Values, q query) {
	for k, v := range q.eqFilter {
		queryValues.Set(
			fmt.Sprintf("filter[%s][eq]", k),
			v,
		)
	}
	if len(q.sort) > 0 {
		queryValues.Set("sort", strings.Join(q.sort, ","))
	}
	if q.limit != nil {
		queryValues.Set("limit", fmt.Sprint(q.limit))
	}
	if q.offset != nil {
		queryValues.Set("offset", fmt.Sprint(q.offset))
	}
}

var fieldsR []string

func (d API[R, W]) jsonFieldsR() []string {
	if fieldsR == nil {
		var x R
		t := reflect.TypeOf(x)
		fieldsR = iterateFields(t, "")
	}
	return fieldsR
}

func iterateFields(t reflect.Type, prefix string) []string {
	fields := []string{}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tagVal := ""
		if v, ok := f.Tag.Lookup("json"); ok {
			tagVal = v
		} else {
			tagVal = f.Name
		}
		switch f.Type.Kind() {
		case reflect.Struct:
			p := prefix
			if p != "" {
				p = p + "." + tagVal
			} else {
				p = tagVal
			}
			fields = append(fields, iterateFields(f.Type, p)...)
		case reflect.Map:
			panic("map is not implemented")
		default:
			// field is not nested
			v := tagVal
			if prefix != "" {
				v = prefix + "." + tagVal
			}
			fields = append(fields, v)
		}
	}
	return fields
}
