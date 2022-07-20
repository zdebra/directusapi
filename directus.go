package directusapi

import (
	"context"
	"fmt"
	"net/http"
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

func (d API[R, W]) Insert(ctx context.Context, item W) (R, error) {
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

func (d API[R, W]) Create(ctx context.Context, partials map[string]any) (R, error) {
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

func (d API[R, W]) GetByID(ctx context.Context, id string) (R, error) {
	u := fmt.Sprintf("%s://%s/%s/items/%s/%s", d.Scheme, d.Host, d.Namespace, d.CollectionName, id)

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

func (d API[R, W]) Update(ctx context.Context, id string, partials map[string]any) (R, error) {
	var empty R
	u := fmt.Sprintf("%s://%s/%s/items/%s/%s", d.Scheme, d.Host, d.Namespace, d.CollectionName, id)

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

func (d API[R, W]) Set(ctx context.Context, id string, item W) (R, error) {
	var empty R
	u := fmt.Sprintf("%s://%s/%s/items/%s/%s", d.Scheme, d.Host, d.Namespace, d.CollectionName, id)

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

func (d API[R, W]) Delete(ctx context.Context, id string) error {
	u := fmt.Sprintf("%s://%s/%s/items/%s/%s", d.Scheme, d.Host, d.Namespace, d.CollectionName, id)
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

func (d API[R, W]) Items(ctx context.Context, q query) ([]R, error) {
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
