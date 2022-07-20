package directusapi

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Category string

const (
	Red   Category = "red"
	Blue  Category = "blue"
	Green Category = "green"
)

type FruitR struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	Weight       int      `json:"weight"`
	Status       string   `json:"status"`
	Category     Category `json:"category"`
	Enabled      bool     `json:"enabled"`
	Price        *float64 `json:"price"`
	DiscoveredAt *Time    `json:"discovered_at"`
	Area         []string `json:"area"`
}

type FruitW struct {
	Name         string    `json:"name"`
	Weight       int       `json:"weight"`
	Status       string    `json:"status"`
	Category     Category  `json:"category"`
	Enabled      bool      `json:"enabled"`
	Price        *float64  `json:"price"`
	DiscoveredAt time.Time `json:"discovered_at"`
	Area         []string  `json:"area"`
}

func TestFlow(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFn()
	api := API[FruitR, FruitW, int]{
		Scheme:         "http",
		Host:           "192.168.64.3:8080",
		Namespace:      "_",
		CollectionName: "fruits",
		HTTPClient:     http.DefaultClient,
		debug:          false,
	}

	email := "zdenek@zdebra.com"
	password := "hovno"
	token, err := api.CreateToken(ctx, email, password)
	require.NoError(t, err)
	api.BearerToken = token

	// todo: create collection first

	watermelonID := 0
	t1 := time.Date(2022, 5, 5, 10, 30, 0, 0, time.UTC)
	t.Run("insert", func(t *testing.T) {
		price := 120.36
		melon, err := api.Insert(ctx, FruitW{
			Name:         "watermelon",
			Weight:       20,
			Status:       "published",
			Category:     Green,
			Enabled:      true,
			Price:        &price,
			DiscoveredAt: t1,
			Area:         []string{"europe", "africa"},
		})
		require.NoError(t, err)
		assert.NotEmpty(t, melon.ID)
		assert.Equal(t, melon.Name, "watermelon")
		assert.Equal(t, melon.Weight, 20)
		assert.Equal(t, melon.Status, "published")
		assert.Equal(t, melon.Enabled, true)
		assert.Equal(t, &price, melon.Price)
		assert.Equal(t, t1.Unix(), melon.DiscoveredAt.Unix())
		assert.Equal(t, []string{"europe", "africa"}, melon.Area)
		watermelonID = melon.ID
	})

	t.Run("get by id", func(t *testing.T) {
		melon, err := api.GetByID(ctx, watermelonID)
		require.NoError(t, err)
		assert.Equal(t, melon.ID, watermelonID)
		assert.Equal(t, melon.Name, "watermelon")
		assert.Equal(t, melon.Weight, 20)
		assert.Equal(t, melon.Status, "published")
		assert.Equal(t, melon.Category, Green)
		assert.Equal(t, melon.Enabled, true)
	})

	t.Run("set item", func(t *testing.T) {
		melonRepl := FruitW{
			Name:   "pasionfruit",
			Weight: 10,
			Status: "published",
		}
		pasionfruit, err := api.Set(ctx, watermelonID, melonRepl)
		require.NoError(t, err)
		assert.Equal(t, pasionfruit.ID, watermelonID)
		assert.Equal(t, pasionfruit.Name, "pasionfruit")
		assert.Equal(t, pasionfruit.Weight, 10)
		assert.Equal(t, pasionfruit.Status, "published")
		expectedPrice := float64(0)
		assert.Equal(t, &expectedPrice, pasionfruit.Price)
	})

	t.Run("update partials", func(t *testing.T) {
		pasionfruit, err := api.Update(ctx, watermelonID, map[string]any{
			"status":   "draft",
			"category": Blue,
		})
		require.NoError(t, err)
		assert.Equal(t, pasionfruit.ID, watermelonID)
		assert.Equal(t, pasionfruit.Name, "pasionfruit")
		assert.Equal(t, pasionfruit.Weight, 10)
		assert.Equal(t, pasionfruit.Status, "draft")
		assert.Equal(t, pasionfruit.Category, Blue)
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
		err := api.Delete(ctx, watermelonID)
		require.NoError(t, err)
	})

	// todo: drop collection

}
