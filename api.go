package main

import (
	"log"
	"scootin-aboot/consts"
	"scootin-aboot/handlers"
	"scootin-aboot/middlewares"
	"scootin-aboot/models"
	"scootin-aboot/repositories"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// initDB initializes the database connection and performs necessary migrations.
// It returns a pointer to the gorm.DB instance.
func initDB() *gorm.DB {
	db, err := gorm.Open(postgres.Open("host=db user=postgres dbname=postgres sslmode=disable"), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&models.Scooter{}, &models.Event{}, &models.User{})

	return db
}

// initAPI initializes the API and returns the API instance and the router.
func initAPI() (huma.API, *chi.Mux) {
	router := chi.NewMux()
	api := humachi.New(router, huma.DefaultConfig("Scootin' Aboot API", "1.0.0"))

	return api, router
}

// initRespositories initializes the repositories used by the API handlers.
// It sets the DB instance for each repository.
func initRespositories(db *gorm.DB) {
	handlers.ScooterRepository = &repositories.ScooterRepository{DB: db}
	handlers.EventRepository = &repositories.EventRepository{DB: db}
	handlers.UserRepository = &repositories.UserRepository{DB: db}
}

// initRoutes initializes the routes for the API.
// It sets up the HTTP methods and their corresponding handlers for each route.
// The routes include listing scooters, creating users, creating events, listing events, and updating scooters.
// Each route is associated with a summary, description, and tags for documentation purposes.
func initRoutes(api huma.API) {
	huma.Post(api, consts.SCOOTERS, handlers.POST_Scooters, func(o *huma.Operation) {
		o.Summary = "Create scooter"
		o.Description = `Create a new scooter.
		It requires proper API key (user ID) to be provided in the Authorization header.`
		o.Tags = []string{"Scooters"}
	})

	// Route for listing scooters
	huma.Get(
		api,
		consts.SCOOTERS,
		handlers.GET_Scooters,
		func(o *huma.Operation) {
			o.Summary = "List scooters"
			o.Description = `List all scooters.
			It requires proper API key (user ID) to be provided in the Authorization header.
			You can query scooters by status or even latitude/longitude coordinates.

			Some examples:
			/scooters?status=free
			/scooters?status=free&min_latitude=0&min_longitude=1&max_latitude=0&max_longitude=1`
			o.Tags = []string{"Scooters"}
		},
	)

	// Route for creating users
	huma.Post(api, consts.USERS, handlers.POST_Users, func(o *huma.Operation) {
		o.Summary = "Create an user"
		o.Description = "Create a new user. It not requires any parameters or API key. It returns the user ID which is used for authorization for the all other endpoints."
		o.Tags = []string{"Users"}
	})

	// Route for creating events
	huma.Post(api, consts.EVENTS, handlers.POST_Events, func(o *huma.Operation) {
		o.Summary = "Create event"
		o.Description = `Create a new event.
		It requires API key (user ID) in the Authorization header, scooter ID, event type (start, stop, location_update), latitude, and longitude.
		It returns Event object. Returns "404 Not Found" when the scooter is not found, "400 Bad Request" when the scooter is not occupied and "409 Conflict" when the scooter is occupied by another user.`
		o.Tags = []string{"Events"}
	})

	// Route for listing events
	huma.Get(api, consts.EVENTS, handlers.GET_Events, func(o *huma.Operation) {
		o.Summary = "List events"
		o.Description = "List all events. It requires proper API key (user ID) to be provided in the Authorization header."
		o.Tags = []string{"Events"}
	})

	// Route for updating scooters
	huma.Patch(api, consts.SCOOTERS_ITEM, handlers.PATCH_Scooters, func(o *huma.Operation) {
		o.Summary = "Update scooter"
		o.Description = `Update a scooter.
		It requires proper API key (user ID) to be provided in the Authorization header and the current etag in the If-Match header.
		Returns "404 Not Found" when the scooter is not found, "412 Precondition Failed" when the etag does not match, and "400 Bad Request" when the scooter is already occupied or want to free not occupied scooter.`
		o.Tags = []string{"Scooters"}
	})
}

// InitAPI initializes the API and returns an instance of the API and the router.
func InitAPI() (huma.API, *chi.Mux) {
	db := initDB()
	api, router := initAPI()

	initRespositories(db)

	api.UseMiddleware(middlewares.RequestLogMiddleware)
	api.UseMiddleware(middlewares.NewAuthorizationMiddleware(api).Middleware)

	initRoutes(api)

	return api, router
}
