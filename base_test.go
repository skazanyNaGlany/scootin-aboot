package main

import (
	"net/http"
	"os"
	"scootin-aboot/consts"
	"scootin-aboot/enums"
	"scootin-aboot/handlers"
	"scootin-aboot/models"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/rand"
)

// BaseTest represents a base test structure.
type BaseTest struct {
	staticAPIKey string
	wrappedAPI   humatest.TestAPI
	testScooters []*models.Scooter
	testEvents   []*models.Event
	testUsers    []*models.User
	isSet        bool
}

func (st *BaseTest) Test401UnauthorizedResponseMap(t *testing.T, responseMap map[string]any) {
	assert.Contains(t, responseMap, "$schema")
	assert.NotEmpty(t, responseMap["$schema"])
	assert.Contains(t, responseMap, "title")
	assert.NotEmpty(t, responseMap["title"])
	assert.Equal(t, "Unauthorized", responseMap["title"])
	assert.Contains(t, responseMap, "status")
	assert.NotEmpty(t, responseMap["status"])
	assert.Equal(t, http.StatusUnauthorized, int(responseMap["status"].(float64)))
	assert.Contains(t, responseMap, "detail")
	assert.NotEmpty(t, responseMap["detail"])
	assert.Contains(t, responseMap, "errors")
	assert.NotEmpty(t, responseMap["errors"])
	assert.Contains(t, responseMap["errors"].([]any)[0], "message")
	assert.NotEmpty(t, responseMap["errors"].([]any)[0].(map[string]any)["message"])
}

func (st *BaseTest) Test404NotFoundResponseMap(t *testing.T, responseMap map[string]any) {
	assert.Contains(t, responseMap, "$schema")
	assert.NotEmpty(t, responseMap["$schema"])
	assert.Contains(t, responseMap, "title")
	assert.NotEmpty(t, responseMap["title"])
	assert.Equal(t, "Not Found", responseMap["title"])
	assert.Contains(t, responseMap, "status")
	assert.NotEmpty(t, responseMap["status"])
	assert.Equal(t, http.StatusNotFound, int(responseMap["status"].(float64)))
	assert.Contains(t, responseMap, "detail")
	assert.NotEmpty(t, responseMap["detail"])
}

func (st *BaseTest) Test412PreconditionFailedResponseMap(t *testing.T, responseMap map[string]any) {
	assert.Contains(t, responseMap, "$schema")
	assert.NotEmpty(t, responseMap["$schema"])
	assert.Contains(t, responseMap, "title")
	assert.NotEmpty(t, responseMap["title"])
	assert.Equal(t, "Precondition Failed", responseMap["title"])
	assert.Contains(t, responseMap, "status")
	assert.NotEmpty(t, responseMap["status"])
	assert.Equal(t, http.StatusPreconditionFailed, int(responseMap["status"].(float64)))
	assert.Contains(t, responseMap, "detail")
	assert.NotEmpty(t, responseMap["detail"])
}

func (st *BaseTest) Test400BadRequestResponseMap(t *testing.T, responseMap map[string]any) {
	assert.Contains(t, responseMap, "$schema")
	assert.NotEmpty(t, responseMap["$schema"])
	assert.Contains(t, responseMap, "title")
	assert.NotEmpty(t, responseMap["title"])
	assert.Equal(t, "Bad Request", responseMap["title"])
	assert.Contains(t, responseMap, "status")
	assert.NotEmpty(t, responseMap["status"])
	assert.Equal(t, http.StatusBadRequest, int(responseMap["status"].(float64)))
	assert.Contains(t, responseMap, "detail")
	assert.NotEmpty(t, responseMap["detail"])
}

func (st *BaseTest) Test409ConflictResponseMap(t *testing.T, responseMap map[string]any) {
	assert.Contains(t, responseMap, "$schema")
	assert.NotEmpty(t, responseMap["$schema"])
	assert.Contains(t, responseMap, "title")
	assert.NotEmpty(t, responseMap["title"])
	assert.Equal(t, "Conflict", responseMap["title"])
	assert.Contains(t, responseMap, "status")
	assert.NotEmpty(t, responseMap["status"])
	assert.Equal(t, http.StatusConflict, int(responseMap["status"].(float64)))
	assert.Contains(t, responseMap, "detail")
	assert.NotEmpty(t, responseMap["detail"])
}

func (st *BaseTest) Test422UnprocessableEntityResponseMap(t *testing.T, responseMap map[string]any) {
	assert.Contains(t, responseMap, "$schema")
	assert.NotEmpty(t, responseMap["$schema"])
	assert.Contains(t, responseMap, "title")
	assert.NotEmpty(t, responseMap["title"])
	assert.Equal(t, "Unprocessable Entity", responseMap["title"])
	assert.Contains(t, responseMap, "status")
	assert.NotEmpty(t, responseMap["status"])
	assert.Equal(t, http.StatusUnprocessableEntity, int(responseMap["status"].(float64)))
	assert.Contains(t, responseMap, "detail")
	assert.NotEmpty(t, responseMap["detail"])
}

func (st *BaseTest) Test200OKResponseMapCollection(t *testing.T, responseMap map[string]any) {
	assert.Contains(t, responseMap, "$schema")
	assert.NotEmpty(t, responseMap["$schema"])
	assert.Contains(t, responseMap, "_embedded")
	assert.Contains(t, responseMap, "_links")
	assert.Contains(t, responseMap["_links"], "self")
	assert.Contains(t, responseMap["_links"].(map[string]any)["self"], "href")
	assert.NotEmpty(t, responseMap["_links"].(map[string]any)["self"].(map[string]any)["href"])
}

func (st *BaseTest) Test200OKResponseMapItem(t *testing.T, responseMap map[string]any) {
	assert.Contains(t, responseMap, "$schema")
	assert.NotEmpty(t, responseMap["$schema"])
	assert.Contains(t, responseMap, "id")
	assert.NotEmpty(t, responseMap["id"])
	assert.Contains(t, responseMap, "created_at")
	assert.NotEmpty(t, responseMap["created_at"])
	assert.Contains(t, responseMap, "updated_at")
	assert.NotEmpty(t, responseMap["updated_at"])
	assert.Contains(t, responseMap, "_links")
	assert.Contains(t, responseMap["_links"], "self")
	assert.Contains(t, responseMap["_links"].(map[string]any)["self"], "href")
	assert.NotEmpty(t, responseMap["_links"].(map[string]any)["self"].(map[string]any)["href"])

	if _, etagExists := responseMap["etag"]; etagExists {
		assert.NotEmpty(t, responseMap["etag"])
	}
}

func (st *BaseTest) TestScooterMap(t *testing.T, scooterMap map[string]any, withEmbeddedEvents bool) {
	assert.Contains(t, scooterMap, "id")
	assert.NotEmpty(t, scooterMap["id"])
	assert.Contains(t, scooterMap, "created_at")
	assert.NotEmpty(t, scooterMap["created_at"])
	assert.Contains(t, scooterMap, "updated_at")
	assert.NotEmpty(t, scooterMap["updated_at"])
	assert.Contains(t, scooterMap, "status")
	assert.NotEmpty(t, scooterMap["status"])
	assert.Contains(t, scooterMap, "user_id")
	assert.NotEmpty(t, scooterMap["user_id"])
	assert.Contains(t, scooterMap, "etag")
	assert.NotEmpty(t, scooterMap["etag"])

	if withEmbeddedEvents {
		assert.Contains(t, scooterMap, "_embedded")
		assert.Contains(t, scooterMap["_embedded"], "events")
		assert.NotEmpty(t, scooterMap["_embedded"].(map[string]any)["events"])

		st.TestEventList(t, scooterMap["_embedded"].(map[string]any)["events"].([]any))
	}
}

func (st *BaseTest) TestScooterList(t *testing.T, scooterList []any, withEmbeddedEvents bool) {
	for _, iscooter := range scooterList {
		st.TestScooterMap(t, iscooter.(map[string]any), withEmbeddedEvents)
	}
}

func (st *BaseTest) TestEventList(t *testing.T, eventList []any) {
	for _, ievent := range eventList {
		st.TestEventMap(t, ievent.(map[string]any))
	}
}

func (st *BaseTest) TestEventMap(t *testing.T, eventMap map[string]any) {
	assert.Contains(t, eventMap, "id")
	assert.NotEmpty(t, eventMap["id"])
	assert.Contains(t, eventMap, "created_at")
	assert.NotEmpty(t, eventMap["created_at"])
	assert.Contains(t, eventMap, "updated_at")
	assert.NotEmpty(t, eventMap["updated_at"])
	assert.Contains(t, eventMap, "scooter_id")
	assert.NotEmpty(t, eventMap["scooter_id"])
	assert.Contains(t, eventMap, "user_id")
	assert.NotEmpty(t, eventMap["user_id"])
	assert.Contains(t, eventMap, "event_type")
	assert.NotEmpty(t, eventMap["event_type"])
	assert.Contains(t, eventMap, "latitude")
	assert.Contains(t, eventMap, "longitude")
	assert.Contains(t, eventMap, "_links")
	assert.Contains(t, eventMap["_links"], "self")
	assert.Contains(t, eventMap["_links"].(map[string]any)["self"], "href")
	assert.NotEmpty(t, eventMap["_links"].(map[string]any)["self"].(map[string]any)["href"])
}

func (st *BaseTest) TestUserMap(t *testing.T, userMap map[string]any) {
	assert.Contains(t, userMap, "$schema")
	assert.NotEmpty(t, userMap["$schema"])
	assert.Contains(t, userMap, "id")
	assert.NotEmpty(t, userMap["id"])
	assert.Contains(t, userMap, "created_at")
	assert.NotEmpty(t, userMap["created_at"])
	assert.Contains(t, userMap, "updated_at")
	assert.NotEmpty(t, userMap["updated_at"])
	assert.Contains(t, userMap, "_links")
	assert.Contains(t, userMap["_links"], "self")
	assert.Contains(t, userMap["_links"].(map[string]any)["self"], "href")
	assert.NotEmpty(t, userMap["_links"].(map[string]any)["self"].(map[string]any)["href"])
}

func (st *BaseTest) HasEvent(t *testing.T, events []any, id int64) bool {
	for _, ievent := range events {
		if ievent.(map[string]any)["id"].(float64) == float64(id) {
			return true
		}
	}

	return false
}

// setup initializes the test environment for the BaseTest struct.
// It sets up the necessary dependencies, such as the API object, test scooters, test events, and test users.
// If any error occurs during the setup process, it will cause the test to fail.
func (st *BaseTest) setup(t *testing.T) {
	var err error

	if st.isSet {
		return
	}

	apiObj, _ := InitAPI()
	st.wrappedAPI = humatest.Wrap(t, apiObj)
	st.staticAPIKey = os.Getenv("STATIC_API_KEY")

	st.testScooters, st.testEvents, err = st.addTestScooters()

	if err != nil {
		t.Fatal(err)
	}

	st.testUsers, err = st.addTestUsers()

	if err != nil {
		t.Fatal(err)
	}

	st.isSet = true
}

// teardown is a helper function used to clean up the test environment after each test case.
// It deletes the test scooters and events, and also deletes the test users from the user repository.
// If any error occurs during the cleanup process, it will cause the test to fail.
func (st *BaseTest) teardown(t *testing.T) {
	if !st.isSet {
		return
	}

	if err := st.deleteScootersAndEvents(st.testScooters, st.testEvents); err != nil {
		t.Fatal(err)
	}

	if err := handlers.UserRepository.DeleteBatch(st.testUsers); err != nil {
		t.Fatal(err)
	}

	st.isSet = false
}

// addTestScooters adds test scooters and corresponding events to the system.
// It creates a specified number of test scooters with initial status as "free"
// and generates location update events for each scooter.
// The created scooters and events are then stored in the database.
// Returns the list of created scooters, list of created events, and any error encountered.
func (st *BaseTest) addTestScooters() ([]*models.Scooter, []*models.Event, error) {
	var scooters []*models.Scooter
	var events []*models.Event

	for i := 0; i < consts.INITIAL_TEST_SCOOTERS_COUNT; i++ {
		scooters = append(
			scooters,
			&models.Scooter{ID: uuid.New(), Status: string(enums.ScooterStatusFree), ETag: uuid.New()},
		)
	}

	for _, scooter := range scooters {
		events = append(
			events,
			&models.Event{
				ScooterID: scooter.ID,
				EventType: string(enums.EventTypeLocationUpdate),
				Latitude:  0,
				Longitude: 1,
			},
		)
	}

	if err := handlers.ScooterRepository.CreateBatch(scooters); err != nil {
		return nil, nil, err
	}

	if err := handlers.EventRepository.CreateBatch(events); err != nil {
		return nil, nil, err
	}

	return scooters, events, nil
}

// addTestUsers adds test users to the system.
// It creates a slice of User objects with randomly generated IDs and ETags.
// The number of test users created is determined by the INITIAL_TEST_USERS_COUNT constant.
// The created users are then stored in the database using the UserRepository.
// Returns the slice of created users and any error encountered during the creation process.
func (st *BaseTest) addTestUsers() ([]*models.User, error) {
	var users []*models.User

	for i := 0; i < consts.INITIAL_TEST_USERS_COUNT; i++ {
		users = append(
			users,
			&models.User{ID: uuid.New()},
		)
	}

	if err := handlers.UserRepository.CreateBatch(users); err != nil {
		return nil, err
	}

	return users, nil
}

// deleteScootersAndEvents deletes the given scooters and events from the database.
// It takes a slice of scooters and a slice of events as parameters and returns an error if any.
func (st *BaseTest) deleteScootersAndEvents(scooters []*models.Scooter, events []*models.Event) error {
	if err := handlers.ScooterRepository.DeleteBatch(scooters); err != nil {
		return err
	}

	if err := handlers.EventRepository.DeleteBatch(events); err != nil {
		return err
	}

	return nil
}

func (st *BaseTest) getRandomScooter() *models.Scooter {
	return st.testScooters[rand.Intn(len(st.testScooters))]
}

func (st *BaseTest) getRandomUser() *models.User {
	return st.testUsers[rand.Intn(len(st.testUsers))]
}

func (st *BaseTest) getRandomEvent() *models.Event {
	return st.testEvents[rand.Intn(len(st.testEvents))]
}
