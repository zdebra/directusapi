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

const tagName = "directus"

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
		// todo: custom json marshal based on custom struct tag with reflection
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
		// todo: custom json unmarshal based on custom tag
		err = json.NewDecoder(resp.Body).Decode(dest)
		if err != nil {
			return fmt.Errorf("decoding json response: %w", err)
		}
	}

	return nil
}

func jsonMarshal(inp any) ([]byte, error) {
	jsonFieldsMap, err := mapStruct(inp)
	if err != nil {
		return nil, fmt.Errorf("struct to map: %w", err)
	}
	return json.Marshal(jsonFieldsMap)
}

func mapStruct(inp any) (OrderedMap, error) {
	output := OrderedMap{}
	structVal := reflect.ValueOf(inp)
	structType := structVal.Type()
	// iterate through struct fields
	for i := 0; i < structVal.NumField(); i++ {
		fieldVal := structVal.Field(i)
		fieldType := structType.Field(i)

		directusFieldName, ok := fieldType.Tag.Lookup(tagName)
		if !ok {
			// field has no tag, skipping the field
			continue
		}
		key := directusFieldName
		val, err := valFromReflectVal(fieldVal)
		if err != nil {
			return OrderedMap{}, fmt.Errorf("value for item: %w", err)
		}
		output = append(output, KeyVal{key, val})
	}
	return output, nil
}

func valFromReflectVal(refVal reflect.Value) (any, error) {
	switch refVal.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return refVal.Int(), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return refVal.Uint(), nil
	case reflect.Float32, reflect.Float64:
		return refVal.Float(), nil
	case reflect.Bool:
		return refVal.Bool(), nil
	case reflect.String:
		return refVal.String(), nil
	case reflect.Slice, reflect.Array:
		size := refVal.Len()
		items := []any{}
		for j := 0; j < size; j++ {
			collectionItemVal := refVal.Index(j)
			collectionVal, err := valFromReflectVal(collectionItemVal)
			if err != nil {
				return nil, fmt.Errorf("value from nested collection: %w", err)
			}
			items = append(items, collectionVal)
		}
		return items, nil
	case reflect.Struct:
		return mapStruct(refVal.Interface())
	case reflect.Map:
		return mapStringMap(refVal)
	case reflect.Pointer:
		panic("not implemented " + refVal.String())
	default:
		panic("unsupported field type")
	}
}

func mapStringMap(refVal reflect.Value) (OrderedMap, error) {
	out := OrderedMap{}

	iter := refVal.MapRange()
	for iter.Next() {
		k := iter.Key()
		if k.Kind() != reflect.String {
			return nil, fmt.Errorf("unsupported key type")
		}
		keyStr := k.String()

		v, err := valFromReflectVal(iter.Value())
		if err != nil {
			return out, fmt.Errorf("value from map value of key %q: %w", keyStr, err)
		}
		out = append(out, KeyVal{keyStr, v})
	}
	return out, nil
}
