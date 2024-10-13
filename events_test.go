package main

import (
	"encoding/json"
	"net/http"
	"scootin-aboot/consts"
	"scootin-aboot/enums"
	"scootin-aboot/handlers"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var eventsTest = EventsTest{}

// EventsTest represents a test suite for events.
type EventsTest struct {
	BaseTest
}

// TestPostEventsUnauthorized tests the behavior of the PostEventsUnauthorized function.
// It sends a POST request to the /events endpoint without authorization and checks if the response is HTTP 401 Unauthorized.
// It also verifies the response body by unmarshaling it into a map[string]any and comparing it with the expected response map.
func (st *EventsTest) TestPostEventsUnauthorized(t *testing.T) {
	var responseMap map[string]any

	st.setup(t)
	defer st.teardown(t)

	response := st.wrappedAPI.Post(consts.EVENTS, "")

	assert.Equal(t, http.StatusUnauthorized, response.Code)
	json.Unmarshal(response.Body.Bytes(), &responseMap)

	st.Test401UnauthorizedResponseMap(t, responseMap)
}

// TestPostEventNoScooter tests the scenario where a scooter is not found when posting an event.
func (st *EventsTest) TestPostEventNoScooter(t *testing.T) {
	var responseMap map[string]any

	st.setup(t)
	defer st.teardown(t)

	user := st.getRandomUser()

	response := st.wrappedAPI.Post(consts.EVENTS, "Authorization: "+user.ID.String(), map[string]any{
		"scooter_id": uuid.New().String(),
		"event_type": string(enums.EventTypeStart),
		"latitude":   0,
		"longitude":  1,
	})

	assert.Equal(t, http.StatusNotFound, response.Code)
	json.Unmarshal(response.Body.Bytes(), &responseMap)

	st.Test404NotFoundResponseMap(t, responseMap)
}

// TestPostEventNoOccupiedScooter tests the scenario where a POST request is made to create an event
// when there is no occupied scooter available.
// It verifies that the response code is http.StatusBadRequest and the response body contains a
// specific error message.
func (st *EventsTest) TestPostEventNoOccupiedScooter(t *testing.T) {
	var responseMap map[string]any

	st.setup(t)
	defer st.teardown(t)

	scooter := st.getRandomScooter()
	user := st.getRandomUser()

	response := st.wrappedAPI.Post(consts.EVENTS, "Authorization: "+user.ID.String(), map[string]any{
		"scooter_id": scooter.ID.String(),
		"event_type": string(enums.EventTypeStart),
		"latitude":   0,
		"longitude":  1,
	})

	assert.Equal(t, http.StatusBadRequest, response.Code)
	json.Unmarshal(response.Body.Bytes(), &responseMap)

	st.Test400BadRequestResponseMap(t, responseMap)
}

// TestPostEventScooterOccupiedAnotherUser tests the scenario where a user tries to post an event for a scooter that is already occupied by another user.
// It verifies that the API returns a BadRequest status code and the appropriate error response.
func (st *EventsTest) TestPostEventScooterOccupiedAnotherUser(t *testing.T) {
	var responseMap map[string]any

	st.setup(t)
	defer st.teardown(t)

	scooter := st.getRandomScooter()
	user := st.testUsers[0]

	// occupy the scooter
	response := st.wrappedAPI.Patch(
		strings.ReplaceAll(consts.SCOOTERS_ITEM, "{id}", scooter.ID.String()),
		"Authorization: "+user.ID.String(),
		"If-Match: "+scooter.ETag.String(),
		map[string]any{
			"status": string(enums.ScooterStatusOccupied),
		},
	)

	assert.Equal(t, http.StatusOK, response.Code)
	json.Unmarshal(response.Body.Bytes(), &responseMap)

	st.TestScooterMap(t, responseMap, false)
	assert.Equal(t, string(enums.ScooterStatusOccupied), responseMap["status"])

	response = st.wrappedAPI.Post(consts.EVENTS, "Authorization: "+st.testUsers[1].ID.String(), map[string]any{
		"scooter_id": scooter.ID.String(),
		"event_type": string(enums.EventTypeStart),
		"latitude":   0,
		"longitude":  1,
	})

	assert.Equal(t, http.StatusConflict, response.Code)
	json.Unmarshal(response.Body.Bytes(), &responseMap)

	st.Test409ConflictResponseMap(t, responseMap)
}

// TestPostEvent tests the functionality of posting different types of events for a scooter.
// It verifies that the events are successfully added and can be retrieved from the API.
func (st *EventsTest) TestPostEvent(t *testing.T) {
	var responseMap map[string]any
	var addedIds []int64

	st.setup(t)
	defer st.teardown(t)

	scooter := st.getRandomScooter()
	user := st.getRandomUser()

	// occupy the scooter
	response := st.wrappedAPI.Patch(
		strings.ReplaceAll(consts.SCOOTERS_ITEM, "{id}", scooter.ID.String()),
		"Authorization: "+user.ID.String(),
		"If-Match: "+scooter.ETag.String(),
		map[string]any{
			"status": string(enums.ScooterStatusOccupied),
		},
	)

	assert.Equal(t, http.StatusOK, response.Code)
	json.Unmarshal(response.Body.Bytes(), &responseMap)

	st.TestScooterMap(t, responseMap, false)
	assert.Equal(t, string(enums.ScooterStatusOccupied), responseMap["status"])

	// add start event
	response = st.wrappedAPI.Post(consts.EVENTS, "Authorization: "+user.ID.String(), map[string]any{
		"scooter_id": scooter.ID.String(),
		"event_type": string(enums.EventTypeStart),
		"latitude":   0,
		"longitude":  1,
	})

	assert.Equal(t, http.StatusOK, response.Code)
	json.Unmarshal(response.Body.Bytes(), &responseMap)

	st.Test200OKResponseMapItem(t, responseMap)

	id := int64(responseMap["id"].(float64))
	addedIds = append(addedIds, id)

	// check if the event was added
	response = st.wrappedAPI.Get(consts.EVENTS, "Authorization: "+user.ID.String())

	assert.Equal(t, http.StatusOK, response.Code)
	json.Unmarshal(response.Body.Bytes(), &responseMap)

	st.Test200OKResponseMapCollection(t, responseMap)

	assert.Contains(t, responseMap["_embedded"], "events")
	assert.NotEmpty(t, responseMap["_embedded"].(map[string]any)["events"])

	events := responseMap["_embedded"].(map[string]any)["events"].([]any)

	assert.True(t, st.HasEvent(t, events, int64(id)), "Event not found in the list after POST")

	// add location_update event
	response = st.wrappedAPI.Post(consts.EVENTS, "Authorization: "+user.ID.String(), map[string]any{
		"scooter_id": scooter.ID.String(),
		"event_type": string(enums.EventTypeLocationUpdate),
		"latitude":   0,
		"longitude":  5,
	})

	assert.Equal(t, http.StatusOK, response.Code)
	json.Unmarshal(response.Body.Bytes(), &responseMap)

	st.Test200OKResponseMapItem(t, responseMap)

	id = int64(responseMap["id"].(float64))
	addedIds = append(addedIds, id)

	// check if the event was added
	response = st.wrappedAPI.Get(consts.EVENTS, "Authorization: "+user.ID.String())

	assert.Equal(t, http.StatusOK, response.Code)
	json.Unmarshal(response.Body.Bytes(), &responseMap)

	st.Test200OKResponseMapCollection(t, responseMap)

	assert.Contains(t, responseMap["_embedded"], "events")
	assert.NotEmpty(t, responseMap["_embedded"].(map[string]any)["events"])

	events = responseMap["_embedded"].(map[string]any)["events"].([]any)

	assert.True(t, st.HasEvent(t, events, int64(id)), "Event not found in the list after POST")

	// add stop event
	response = st.wrappedAPI.Post(consts.EVENTS, "Authorization: "+user.ID.String(), map[string]any{
		"scooter_id": scooter.ID.String(),
		"event_type": string(enums.EventTypeStop),
		"latitude":   0,
		"longitude":  10,
	})

	assert.Equal(t, http.StatusOK, response.Code)
	json.Unmarshal(response.Body.Bytes(), &responseMap)

	st.Test200OKResponseMapItem(t, responseMap)

	id = int64(responseMap["id"].(float64))
	addedIds = append(addedIds, id)

	// check if the event was added
	response = st.wrappedAPI.Get(consts.EVENTS, "Authorization: "+user.ID.String())

	assert.Equal(t, http.StatusOK, response.Code)
	json.Unmarshal(response.Body.Bytes(), &responseMap)

	st.Test200OKResponseMapCollection(t, responseMap)

	assert.Contains(t, responseMap["_embedded"], "events")
	assert.NotEmpty(t, responseMap["_embedded"].(map[string]any)["events"])

	events = responseMap["_embedded"].(map[string]any)["events"].([]any)

	assert.True(t, st.HasEvent(t, events, int64(id)), "Event not found in the list after POST")

	// remove added events
	err := handlers.EventRepository.DeleteBatchByIDs(addedIds)

	assert.Nil(t, err)
}

// TestGetEvents tests the GetEvents function.
// It sends a GET request to the /events endpoint and verifies the response.
func (st *EventsTest) TestGetEvents(t *testing.T) {
	var responseMap map[string]any

	st.setup(t)
	defer st.teardown(t)

	user := st.getRandomUser()

	response := st.wrappedAPI.Get(consts.EVENTS, "Authorization: "+user.ID.String())

	assert.Equal(t, http.StatusOK, response.Code)
	json.Unmarshal(response.Body.Bytes(), &responseMap)

	st.Test200OKResponseMapCollection(t, responseMap)

	assert.Contains(t, responseMap["_embedded"], "events")
	assert.NotEmpty(t, responseMap["_embedded"].(map[string]any)["events"])

	st.TestEventList(t, responseMap["_embedded"].(map[string]any)["events"].([]any))
}

func TestPostEventsUnauthorized(t *testing.T) {
	eventsTest.TestPostEventsUnauthorized(t)
}

func TestGetEvents(t *testing.T) {
	eventsTest.TestGetEvents(t)
}

func TestPostEventNoScooter(t *testing.T) {
	eventsTest.TestPostEventNoScooter(t)
}

func TestPostEventNoOccupiedScooter(t *testing.T) {
	eventsTest.TestPostEventNoOccupiedScooter(t)
}

func TestPostEventScooterOccupiedAnotherUser(t *testing.T) {
	eventsTest.TestPostEventScooterOccupiedAnotherUser(t)
}

func TestPostEvent(t *testing.T) {
	eventsTest.TestPostEvent(t)
}
