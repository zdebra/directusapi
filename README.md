# Directus API client

This is generics golang client for [Directus](https://directus.io/) v8 CMS. Never write the same API client again.
Just define your collection model and use strongly typed methods.

---

## Example usage

```go
// define Read and Write model
type FruitR struct {
    ID          int             `json:"id"`
    Name        string          `json:"name"`
    Weight      float64         `json:"weight"`
}

type FruitW struct {
    Name        string          `json:"name"`
    Weight      float64         `json:"weight"`
}

// initialize generic API client
api := API[FruitR, FruitW, int]{
    Scheme:         "http",
    Host:           "localhost:8080",
    Namespace:      "_",
    CollectionName: "fruits",
    HTTPClient:     http.DefaultClient,
    BearerToken:    "1a2bd3db-8026-4494-ad36-9873ee46c0af"
}

// use typed methods
// - insert
watermelon, err := api.Insert(ctx, FruitW{
	Name:         "watermelon",
	Weight:       20.3,
})
// watermelon's type is FruitR

// - retrieve collection of items
fruits, err := api.Items(ctx, None())
// fruits's type is []FruitR

// update (set) item
passionfruit, err := api.Set(ctx, watermelonID, FruitW{
	Name:         "passionfruit",
	Weight:       3.3,
})
// passionfruit's type is FruitR

```

Go to the documentation to see all available methods.

## Features

- strongly-typed API methods based on [directus reference](https://v8.docs.directus.io/api/reference.html)
- different models for reads and writes
- collection querying support: filtering, sorting, limit, offset, fulltext search
- custom `directusapi.Time` to support Directus API time format
- custom `directusapi.Optional` to support optional fields

## What is Directus?

[Directus](https://directus.io/) is open sourced Content Management System, it has UI and exposed API for dynamicly created collections.

## Setup

- todo

## Limitations

- directus v9 is not supported at this moment; this library was developed for directus v8
- pointers are not allowed in your Read and Write models, `directusapi.Optional` should be used for optional fields
- `directusapi.Time` has to be used instead of `time.Time`

## License

> You can check out the full license [here](https://github.com/zdebra/directusapi/blob/master/LICENSE)

This project is licensed under the terms of the **MIT** license.

## Buy me a coffee

Whether you use this project, have learned something from it, or just like it, please consider supporting it by buying me a coffee, so I can dedicate more time on open-source projects like this :)

<a href="https://www.buymeacoffee.com/zdebra" target="_blank"><img src="https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png" alt="Buy Me A Coffee" style="height: auto !important;width: auto !important;" ></a>
