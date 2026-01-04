package app

// maybe move to main.go as separate func, and move all init and other setup logic there.
func Run() {

}

// This struct holds your entire running application
type App struct {
	// HomeUsecase    *home.Usecase
	// CatalogUsecase *catalog.Usecase
	// Add others here...
}

// New initializes everything in one place
// func New(apiHost string) *App {
// 1. Init the shared adapter once
// client := carapi.New(apiHost)

// 2. Init all use cases, injecting the client
// return &App{
// 	HomeUsecase:    home.New(client),
// 	CatalogUsecase: catalog.New(client),
// }
// }
