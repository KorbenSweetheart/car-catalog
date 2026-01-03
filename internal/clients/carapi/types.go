package carapi

// TODO: Maybe it should be in Usecase????

// type Fetcher interface {
// 	GetModels(ctx context.Context, number int) ([]domain.Car, error)
// 	GetCarModel(ctx context.Context, id int) (domain.Car, error)
// 	GetManufacturers(ctx context.Context) ([]domain.Manufacturer, error)
// 	GetCarManufacturer(ctx context.Context, id int) (domain.Manufacturer, error)
// 	GetCarCategories(ctx context.Context) ([]domain.Manufacturer, error)
// 	GetCategory(ctx context.Context, id int) (domain.Manufacturer, error)
// }

/*

The api exposes the following endpoints:
```
GET /api/models
GET /api/models/{id}
GET /api/manufacturers
GET /api/manufacturers/{id}
GET /api/categories
GET /api/categories/{id}
```
*/

// Not used old structs
type UpdatesResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Update struct {
	ID      int    `json:"update_id"`
	Message string `json:"message"`
}
