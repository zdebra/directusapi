package directusapi

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type FruitR struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Weight int    `json:"weight"`
	Status string `json:"status"`
}

type FruitW struct {
	Name   string `json:"name"`
	Weight int    `json:"weight"`
	Status string `json:"status"`
}

func TestFlow(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFn()
	api := API[FruitR, FruitW]{
		Scheme:         "http",
		Host:           "192.168.64.3:8080",
		Namespace:      "_",
		CollectionName: "fruits",
		HTTPClient:     http.DefaultClient,
	}

	email := "zdenek@zdebra.com"
	password := "hovno"
	token, err := api.CreateToken(ctx, email, password)
	require.NoError(t, err)
	api.BearerToken = token

	// todo: create collection first

	watermelonID := 0
	t.Run("insert", func(t *testing.T) {
		melon, err := api.Insert(ctx, FruitW{
			Name:   "watermelon",
			Weight: 20,
			Status: "published",
		})
		require.NoError(t, err)
		assert.NotEmpty(t, melon.ID)
		assert.Equal(t, melon.Name, "watermelon")
		assert.Equal(t, melon.Weight, 20)
		assert.Equal(t, melon.Status, "published")
		watermelonID = melon.ID
	})

	t.Run("get by id", func(t *testing.T) {
		melon, err := api.GetByID(ctx, fmt.Sprint(watermelonID))
		require.NoError(t, err)
		assert.Equal(t, melon.ID, watermelonID)
		assert.Equal(t, melon.Name, "watermelon")
		assert.Equal(t, melon.Weight, 20)
		assert.Equal(t, melon.Status, "published")
	})

	t.Run("set item", func(t *testing.T) {
		melonRepl := FruitW{
			Name:   "pasionfruit",
			Weight: 10,
			Status: "published",
		}
		pasionfruit, err := api.Set(ctx, fmt.Sprint(watermelonID), melonRepl)
		require.NoError(t, err)
		assert.Equal(t, pasionfruit.ID, watermelonID)
		assert.Equal(t, pasionfruit.Name, "pasionfruit")
		assert.Equal(t, pasionfruit.Weight, 10)
		assert.Equal(t, pasionfruit.Status, "published")
	})

	t.Run("update partials", func(t *testing.T) {
		pasionfruit, err := api.Update(ctx, fmt.Sprint(watermelonID), map[string]any{
			"status": "draft",
		})
		require.NoError(t, err)
		assert.Equal(t, pasionfruit.ID, watermelonID)
		assert.Equal(t, pasionfruit.Name, "pasionfruit")
		assert.Equal(t, pasionfruit.Weight, 10)
		assert.Equal(t, pasionfruit.Status, "draft")
	})

	t.Run("read items", func(t *testing.T) {
		fruits, err := api.Items(ctx, None())
		require.NoError(t, err)

		assert.True(t, len(fruits) > 0)
	})

	t.Run("create partials", func(t *testing.T) {
		peach, err := api.Create(ctx, map[string]any{
			"name": "peach",
		})
		require.NoError(t, err)
		assert.NotEmpty(t, peach.ID)
		assert.Equal(t, "peach", peach.Name)
		assert.Empty(t, peach.Weight)
	})

	t.Run("delete item", func(t *testing.T) {
		err := api.Delete(ctx, fmt.Sprint(watermelonID))
		require.NoError(t, err)
	})

	// todo: drop collection

}