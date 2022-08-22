package directusapi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJsonFields(t *testing.T) {
	api := API[FruitR, FruitW, int]{}

	expected := []string{
		"id",
		"name",
		"weight",
		"status",
		"category",
		"enabled",
		"price",
		"discovered_at",
		"area",
		"favorites",
		"lefield.id",
		"lefield.email",
	}
	jsonFields := api.jsonFieldsR()
	assert.Equal(t, expected, jsonFields)
}
