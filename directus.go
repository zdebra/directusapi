package directusapi

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

type PrimaryKey interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~string
}

type API[R, W any, PK PrimaryKey] struct {
	Scheme         string
	Host           string
	Namespace      string
	CollectionName string
	BearerToken    string
	HTTPClient     *http.Client
	debug          bool
}

func (d API[R, W, PK]) CreateToken(ctx context.Context, email, password string) (string, error) {
	u := fmt.Sprintf("%s://%s/%s/auth/authenticate", d.Scheme, d.Host, d.Namespace)

	body := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		email,
		password,
	}

	req := request{
		ctx,
		http.MethodPost,
		u,
		nil,
		body,
	}

	var respBody struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}

	err := d.executeRequest(req, http.StatusOK, &respBody)
	if err != nil {
		return "", fmt.Errorf("execute create token request: %w", err)
	}
	return respBody.Data.Token, nil
}

func (d API[R, W, PK]) Insert(ctx context.Context, item W) (R, error) {
	var empty R
	u := fmt.Sprintf("%s://%s/%s/items/%s", d.Scheme, d.Host, d.Namespace, d.CollectionName)

	req := request{
		ctx,
		http.MethodPost,
		u,
		map[string]string{
			"fields": strings.Join(d.jsonFieldsR(), ","),
		},
		item,
	}
	var respBody struct {
		Data R `json:"data"`
	}
	err := d.executeRequest(req, http.StatusOK, &respBody)
	if err != nil {
		return empty, fmt.Errorf("execute insert request: %w", err)
	}
	return respBody.Data, nil
}

func (d API[R, W, PK]) Create(ctx context.Context, partials map[string]any) (R, error) {
	var empty R
	u := fmt.Sprintf("%s://%s/%s/items/%s", d.Scheme, d.Host, d.Namespace, d.CollectionName)

	req := request{
		ctx,
		http.MethodPost,
		u,
		map[string]string{
			"fields": strings.Join(d.jsonFieldsR(), ","),
		},
		partials,
	}

	var respBody struct {
		Data R `json:"data"`
	}
	err := d.executeRequest(req, http.StatusOK, &respBody)
	if err != nil {
		return empty, fmt.Errorf("execute create request: %w", err)
	}
	return respBody.Data, nil

}

func (d API[R, W, PK]) GetByID(ctx context.Context, id PK) (R, error) {
	u := fmt.Sprintf("%s://%s/%s/items/%s/%v", d.Scheme, d.Host, d.Namespace, d.CollectionName, id)

	req := request{
		ctx,
		http.MethodGet,
		u,
		map[string]string{
			"fields": strings.Join(d.jsonFieldsR(), ","),
		},
		nil,
	}

	var respBody struct {
		Data R `json:"data"`
	}
	var empty R
	err := d.executeRequest(req, http.StatusOK, &respBody)
	if err != nil {
		return empty, fmt.Errorf("execute get by id request: %w", err)
	}
	return respBody.Data, nil
}

func (d API[R, W, PK]) Update(ctx context.Context, id PK, partials map[string]any) (R, error) {
	var empty R
	u := fmt.Sprintf("%s://%s/%s/items/%s/%v", d.Scheme, d.Host, d.Namespace, d.CollectionName, id)

	req := request{
		ctx,
		http.MethodPatch,
		u,
		map[string]string{
			"fields": strings.Join(d.jsonFieldsR(), ","),
		},
		partials,
	}

	var respBody struct {
		Data R `json:"data"`
	}
	err := d.executeRequest(req, http.StatusOK, &respBody)
	if err != nil {
		return empty, fmt.Errorf("execute update request: %w", err)
	}
	return respBody.Data, nil
}

func (d API[R, W, PK]) Set(ctx context.Context, id PK, item W) (R, error) {
	var empty R
	u := fmt.Sprintf("%s://%s/%s/items/%s/%v", d.Scheme, d.Host, d.Namespace, d.CollectionName, id)

	req := request{
		ctx,
		http.MethodPatch,
		u,
		map[string]string{
			"fields": strings.Join(d.jsonFieldsR(), ","),
		},
		item,
	}

	var respBody struct {
		Data R `json:"data"`
	}
	err := d.executeRequest(req, http.StatusOK, &respBody)
	if err != nil {
		return empty, fmt.Errorf("execute set request: %w", err)
	}
	return respBody.Data, nil
}

func (d API[R, W, PK]) Delete(ctx context.Context, id PK) error {
	u := fmt.Sprintf("%s://%s/%s/items/%s/%v", d.Scheme, d.Host, d.Namespace, d.CollectionName, id)
	req := request{
		ctx,
		http.MethodDelete,
		u,
		nil,
		nil,
	}

	err := d.executeRequest(req, http.StatusNoContent, nil)
	if err != nil {
		return fmt.Errorf("execute delete request: %w", err)
	}
	return nil
}

func (d API[R, W, PK]) Items(ctx context.Context, q query) ([]R, error) {
	u := fmt.Sprintf("%s://%s/%s/items/%s", d.Scheme, d.Host, d.Namespace, d.CollectionName)
	qv := q.asKeyValue()
	qv["fields"] = strings.Join(d.jsonFieldsR(), ",")

	req := request{
		ctx,
		http.MethodGet,
		u,
		qv,
		nil,
	}
	var respBody struct {
		Data []R `json:"data"`
	}
	err := d.executeRequest(req, http.StatusOK, &respBody)
	if err != nil {
		return nil, fmt.Errorf("execute items request: %w", err)
	}
	return respBody.Data, nil
}

var fieldsR []string

func (d API[R, W, PK]) jsonFieldsR() []string {
	if fieldsR == nil {
		var x R
		t := reflect.TypeOf(x)
		fieldsR = iterateFields(t, "")
	}
	return fieldsR
}

// iterateFields returns fields for all struct's fields
func iterateFields(t reflect.Type, prefix string) []string {
	fields := []string{}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fields = append(fields, structFields(f, prefix)...)
	}
	return fields
}

// structFields returns fields for a signle struct field
func structFields(f reflect.StructField, prefix string) []string {
	fields := []string{}
	tagVal := ""
	if v, ok := f.Tag.Lookup(tagName); ok {
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
	// case reflect.Pointer:
	// 	element := f.Type.Elem()

	case
		reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.String,
		reflect.Slice, reflect.Map, reflect.Pointer:
		// field is not nested
		v := tagVal
		if prefix != "" {
			v = prefix + "." + tagVal
		}
		fields = append(fields, v)
	default:
		panic(f.Type.Kind().String() + " not implemented")
	}
	return fields
}
