package directusapi

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	_ "embed"

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
	ID           int               `json:"id"`
	Name         string            `json:"name"`
	Weight       int               `json:"weight"`
	Status       string            `json:"status"`
	Category     Category          `json:"category"`
	Enabled      bool              `json:"enabled"`
	Price        Optional[float64] `json:"price"`
	DiscoveredAt Optional[Time]    `json:"discovered_at"`
	Area         []string          `json:"area"`
	Favorites    map[string]string `json:"favorites"`
	LeField      UserR             `json:"lefield"`
	// Poc          *UserR            `json:"poc"`
}

type UserR struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

type FruitW struct {
	Name         string            `json:"name"`
	Weight       int               `json:"weight"`
	Status       string            `json:"status"`
	Category     Category          `json:"category"`
	Enabled      bool              `json:"enabled"`
	Price        Optional[float64] `json:"price"`
	DiscoveredAt Time              `json:"discovered_at"`
	Area         []string          `json:"area"`
	Favorites    map[string]string `json:"favorites"`
	LeFieldRef   int               `json:"lefield"`
	// PocID        *int              `json:"poc"`
}

func TestFlow(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFn()
	api := API[FruitR, FruitW, int]{
		Scheme:         "http",
		Host:           "localhost:8080",
		Namespace:      "_",
		CollectionName: "fruits",
		HTTPClient:     http.DefaultClient,
		debug:          true,
	}

	email := "email@example.com"
	password := "d1r3ctu5"
	token, err := api.CreateToken(ctx, email, password)
	require.NoError(t, err)
	api.BearerToken = token

	// cleanup db before start
	_ = dropCollection(token, "http", "localhost:8080", "_", "fruits")

	err = createCollection(token, "http", "localhost:8080", "_")
	if err != nil {
		panic(err)
	}

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
			Price:        SetOptional(price),
			DiscoveredAt: Time{t1},
			Area:         []string{"europe", "africa"},
			Favorites: map[string]string{
				"josef": "10",
			},
			LeFieldRef: 1,
			// PocID:      lo.ToPtr(1),
		})
		require.NoError(t, err)
		assert.NotEmpty(t, melon.ID)
		assert.Equal(t, melon.Name, "watermelon")
		assert.Equal(t, melon.Weight, 20)
		assert.Equal(t, melon.Status, "published")
		assert.Equal(t, melon.Enabled, true)
		assert.Equal(t, SetOptional(price), melon.Price)
		assert.Equal(t, t1.Unix(), melon.DiscoveredAt.ValueOrZero().Unix())
		assert.Equal(t, []string{"europe", "africa"}, melon.Area)
		assert.Equal(t, map[string]string{
			"josef": "10",
		}, melon.Favorites)
		assert.Equal(t, UserR{
			ID:    1,
			Email: email,
		}, melon.LeField)
		// assert.Equal(t, &UserR{
		// 	ID:    1,
		// 	Email: email,
		// }, melon.Poc)
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
			Name:       "pasionfruit",
			Weight:     10,
			Status:     "published",
			LeFieldRef: 1,
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
			"name":    "peach",
			"lefield": 1,
		})
		require.NoError(t, err)
		assert.NotEmpty(t, peach.ID)
		assert.Equal(t, "peach", peach.Name)
		assert.Empty(t, peach.Weight)
		assert.Equal(t, UserR{
			ID:    1,
			Email: email,
		}, peach.LeField)
	})

	t.Run("delete item", func(t *testing.T) {
		err := api.Delete(ctx, watermelonID)
		require.NoError(t, err)
	})
}

//go:embed test_collection.json
var createCollBody string

func createCollection(apiToken, scheme, hostname, project string) error {
	u := fmt.Sprintf("%s://%s/%s/collections", scheme, hostname, project)

	req, _ := http.NewRequest(http.MethodPost, u, strings.NewReader(createCollBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("invalid status code %s: %q", resp.Status, string(b))
	}

	return nil
}

func dropCollection(apiToken, scheme, hostname, project, collectionName string) error {
	u := fmt.Sprintf("%s://%s/%s/collections/%s", scheme, hostname, project, collectionName)

	req, _ := http.NewRequest(http.MethodDelete, u, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("invalid status code %s", resp.Status)
	}
	return nil
}
