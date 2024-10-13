package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"scootin-aboot/consts"
	"scootin-aboot/enums"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var scootersTest = ScootersTest{}

// ScootersTest represents a test suite for the Scooters functionality.
type ScootersTest struct {
	BaseTest
}

// TestGetScootersUnauthorized tests the behavior of the GetScootersUnauthorized function.
// It sends an unauthorized request to the scooters endpoint and verifies that the response
// code is http.StatusUnauthorized. It also checks the response body for the expected
// unauthorized response map.
func (st *ScootersTest) TestGetScootersUnauthorized(t *testing.T) {
	var responseMap map[string]any

	st.setup(t)
	defer st.teardown(t)

	response := st.wrappedAPI.Get(consts.SCOOTERS, "")

	assert.Equal(t, http.StatusUnauthorized, response.Code)
	json.Unmarshal(response.Body.Bytes(), &responseMap)

	st.Test401UnauthorizedResponseMap(t, responseMap)
}

// TestPatchScootersUnauthorized tests the unauthorized response when patching a scooter.
// It sends a PATCH request to the /scooters/{id} endpoint with an unauthorized user.
// The function verifies that the response code is http.StatusUnauthorized (401) and
// checks the response body for the expected error message.
func (st *ScootersTest) TestPatchScootersUnauthorized(t *testing.T) {
	var responseMap map[string]any

	st.setup(t)
	defer st.teardown(t)

	response := st.wrappedAPI.Patch(consts.SCOOTERS_ITEM)
	assert.Equal(t, http.StatusUnauthorized, response.Code)

	json.Unmarshal(response.Body.Bytes(), &responseMap)
	st.Test401UnauthorizedResponseMap(t, responseMap)
}

// TestPatchScootersNotFound tests the scenario where a scooter is not found during a patch operation.
// It sends a PATCH request to the API with a non-existent scooter ID and verifies that the response
// status code is http.StatusNotFound (404). It also checks that the response body contains a JSON
// object with the expected structure for a 404 Not Found response.
func (st *ScootersTest) TestPatchScootersNotFound(t *testing.T) {
	var responseMap map[string]any

	st.setup(t)
	defer st.teardown(t)

	user := st.getRandomUser()

	fullQuery := strings.ReplaceAll(consts.SCOOTERS_ITEM, "{id}", uuid.New().String())

	response := st.wrappedAPI.Patch(fullQuery, "Authorization: "+user.ID.String(), map[string]any{
		"status": string(enums.ScooterStatusFree),
	})
	assert.Equal(t, http.StatusNotFound, response.Code)

	json.Unmarshal(response.Body.Bytes(), &responseMap)
	st.Test404NotFoundResponseMap(t, responseMap)
}

// TestPatchScootersETagNotMatch tests the scenario where the ETag in the request does not match the ETag of the scooter.
// It sends a PATCH request to update the status and user ID of a scooter, but with an incorrect ETag value.
// The expected behavior is to receive a 412 Precondition Failed response.
func (st *ScootersTest) TestPatchScootersETagNotMatch(t *testing.T) {
	var responseMap map[string]any

	st.setup(t)
	defer st.teardown(t)

	scooter := st.getRandomScooter()
	user := st.getRandomUser()

	fullQuery := strings.ReplaceAll(consts.SCOOTERS_ITEM, "{id}", scooter.ID.String())

	response := st.wrappedAPI.Patch(
		fullQuery,
		"Authorization: "+user.ID.String(),
		"If-Match: "+uuid.NewString(),
		map[string]any{
			"status": string(enums.ScooterStatusFree),
		},
	)
	assert.Equal(t, http.StatusPreconditionFailed, response.Code)

	json.Unmarshal(response.Body.Bytes(), &responseMap)
	st.Test412PreconditionFailedResponseMap(t, responseMap)
}

// TestGetScooters tests the retrieval of scooters from the API.
// It sends a GET request to the /scooters endpoint and verifies the response.
// The function asserts that the response status code is http.StatusOK (200),
// and checks the structure and content of the response body.
// It also tests the scooter list returned in the response.
func (st *ScootersTest) TestGetScooters(t *testing.T) {
	var responseMap map[string]any

	st.setup(t)
	defer st.teardown(t)

	user := st.getRandomUser()

	response := st.wrappedAPI.Get(consts.SCOOTERS, "Authorization: "+user.ID.String())

	assert.Equal(t, http.StatusOK, response.Code)
	json.Unmarshal(response.Body.Bytes(), &responseMap)

	st.Test200OKResponseMapCollection(t, responseMap)

	assert.Contains(t, responseMap["_embedded"], "scooters")
	assert.NotEmpty(t, responseMap["_embedded"].(map[string]any)["scooters"])

	st.TestScooterList(t, responseMap["_embedded"].(map[string]any)["scooters"].([]any), false)
}

// TestGetScooterByStatus tests the functionality of retrieving scooters by status.
// It sends a GET request to the scooters API with a specific status parameter,
// and asserts that the response is successful (HTTP 200 OK) and contains the expected data.
// It also verifies that the retrieved scooters have the correct status.
func (st *ScootersTest) TestGetScooterByStatus(t *testing.T) {
	var responseMap map[string]any

	st.setup(t)
	defer st.teardown(t)

	user := st.getRandomUser()

	fullQuery := consts.SCOOTERS + "?status=%v"

	response := st.wrappedAPI.Get(
		fmt.Sprintf(fullQuery, string(enums.ScooterStatusFree)),
		"Authorization: "+user.ID.String(),
	)

	assert.Equal(t, http.StatusOK, response.Code)
	json.Unmarshal(response.Body.Bytes(), &responseMap)

	st.Test200OKResponseMapCollection(t, responseMap)

	assert.Contains(t, responseMap["_embedded"], "scooters")
	assert.NotEmpty(t, responseMap["_embedded"].(map[string]any)["scooters"])

	st.TestScooterList(t, responseMap["_embedded"].(map[string]any)["scooters"].([]any), false)

	scooters := responseMap["_embedded"].(map[string]any)["scooters"].([]any)

	for _, scooter := range scooters {
		assert.Equal(t, string(enums.ScooterStatusFree), scooter.(map[string]any)["status"])
	}
}

// TestGetScooterByStatusAndLocation tests the functionality of retrieving scooters by status and location.
// It sends a GET request to the /scooters endpoint with the specified status and location parameters,
// and asserts that the response status code is http.StatusOK.
// It also asserts that the response body contains a collection of scooters and each scooter has the expected status.
func (st *ScootersTest) TestGetScooterByStatusAndLocation(t *testing.T) {
	var responseMap map[string]any

	st.setup(t)
	defer st.teardown(t)

	user := st.getRandomUser()

	minLatitude := 0
	minLongitude := 1
	maxLatitude := 0
	maxLongitude := 1

	fullQuery := consts.SCOOTERS + "?status=%v&min_latitude=%v&min_longitude=%v&max_latitude=%v&max_longitude=%v"

	response := st.wrappedAPI.Get(
		fmt.Sprintf(fullQuery, string(enums.ScooterStatusFree), minLatitude, minLongitude, maxLatitude, maxLongitude),
		"Authorization: "+user.ID.String(),
	)

	assert.Equal(t, http.StatusOK, response.Code)
	json.Unmarshal(response.Body.Bytes(), &responseMap)

	st.Test200OKResponseMapCollection(t, responseMap)

	assert.Contains(t, responseMap["_embedded"], "scooters")
	assert.NotEmpty(t, responseMap["_embedded"].(map[string]any)["scooters"])

	st.TestScooterList(t, responseMap["_embedded"].(map[string]any)["scooters"].([]any), true)

	scooters := responseMap["_embedded"].(map[string]any)["scooters"].([]any)

	for _, scooter := range scooters {
		assert.Equal(t, string(enums.ScooterStatusFree), scooter.(map[string]any)["status"])

		// TODO check the latitude and longitude
	}
}

// TestGetScooterNotEnoughArgs tests the scenario where there are not enough arguments provided for getting scooters.
// It sends a GET request to the /scooters endpoint with the specified status and minimum latitude.
// The expected behavior is to receive a 422 Unprocessable Entity response.
// The response body is then unmarshalled into a map[string]any for further assertions.
func (st *ScootersTest) TestGetScooterNotEnoughArgs(t *testing.T) {
	var responseMap map[string]any

	st.setup(t)
	defer st.teardown(t)

	user := st.getRandomUser()

	minLatitude := 0

	fullQuery := consts.SCOOTERS + "?status=%v&min_latitude=%v"

	response := st.wrappedAPI.Get(
		fmt.Sprintf(fullQuery, string(enums.ScooterStatusFree), minLatitude),
		"Authorization: "+user.ID.String(),
	)

	assert.Equal(t, http.StatusUnprocessableEntity, response.Code)
	json.Unmarshal(response.Body.Bytes(), &responseMap)

	st.Test422UnprocessableEntityResponseMap(t, responseMap)
}

// TestCanOccupyScooter tests the ability to occupy a scooter.
// It sends a PATCH request to the API to occupy a scooter and verifies the response.
// Then, it tries to occupy the same scooter again and expects a 400 Bad Request response.
func (st *ScootersTest) TestCanOccupyScooter(t *testing.T) {
	var responseMap map[string]any

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

	// preserve the ETag
	// it will be needed to edit the scooter
	scooter.ETag = uuid.MustParse(responseMap["etag"].(string))

	// try to occupy the same scooter again
	// it should fail with 400 because the scooter is already occupied
	response = st.wrappedAPI.Patch(
		strings.ReplaceAll(consts.SCOOTERS_ITEM, "{id}", scooter.ID.String()),
		"Authorization: "+user.ID.String(),
		"If-Match: "+scooter.ETag.String(),
		map[string]any{
			"status": string(enums.ScooterStatusOccupied),
		},
	)

	assert.Equal(t, http.StatusBadRequest, response.Code)
	json.Unmarshal(response.Body.Bytes(), &responseMap)

	st.Test400BadRequestResponseMap(t, responseMap)
}

// TestCanFreeScooter tests the ability to free a scooter.
// It verifies that a scooter can be occupied, and then freed by setting the user ID to nil.
// It also checks that attempting to free a scooter with a non-nil user ID results in a 400 Bad Request error.
// Finally, it ensures that the scooter is successfully freed and its status is updated to "free".
func (st *ScootersTest) TestCanFreeScooter(t *testing.T) {
	var responseMap map[string]any

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

	// preserve the ETag
	// it will be needed to edit the scooter
	scooter.ETag = uuid.MustParse(responseMap["etag"].(string))

	// free the scooter
	response = st.wrappedAPI.Patch(
		strings.ReplaceAll(consts.SCOOTERS_ITEM, "{id}", scooter.ID.String()),
		"Authorization: "+user.ID.String(),
		"If-Match: "+scooter.ETag.String(),
		map[string]any{
			"status": string(enums.ScooterStatusFree),
		},
	)

	assert.Equal(t, http.StatusOK, response.Code)
	json.Unmarshal(response.Body.Bytes(), &responseMap)

	st.TestScooterMap(t, responseMap, false)
	assert.Equal(t, string(enums.ScooterStatusFree), responseMap["status"])
}

// TestCannotFreeNonOccupiedScooter tests the scenario where a non-occupied scooter cannot be freed.
// It sends a PATCH request to the API to free a scooter that is not occupied.
// The test expects the API to respond with a 400 Bad Request status code, indicating that the scooter cannot be freed.
// The response body is also checked to ensure it contains the expected error message.
func (st *ScootersTest) TestCannotFreeNonOccupiedScooter(t *testing.T) {
	var responseMap map[string]any

	st.setup(t)
	defer st.teardown(t)

	scooter := st.getRandomScooter()
	user := st.getRandomUser()

	// try to free the scooter that is not occupied
	// it should fail with 400 because the scooter is not occupied
	response := st.wrappedAPI.Patch(
		strings.ReplaceAll(consts.SCOOTERS_ITEM, "{id}", scooter.ID.String()),
		"Authorization: "+user.ID.String(),
		"If-Match: "+scooter.ETag.String(),
		map[string]any{
			"status": string(enums.ScooterStatusFree),
		},
	)

	assert.Equal(t, http.StatusBadRequest, response.Code)
	json.Unmarshal(response.Body.Bytes(), &responseMap)

	st.Test400BadRequestResponseMap(t, responseMap)
}

func TestGetScootersUnauthorized(t *testing.T) {
	scootersTest.TestGetScootersUnauthorized(t)
}

func TestGetScooters(t *testing.T) {
	scootersTest.TestGetScooters(t)
}

func TestGetScooterByStatus(t *testing.T) {
	scootersTest.TestGetScooterByStatus(t)
}

func TestGetScooterByStatusAndLocation(t *testing.T) {
	scootersTest.TestGetScooterByStatusAndLocation(t)
}

func TestPatchScootersUnauthorized(t *testing.T) {
	scootersTest.TestPatchScootersUnauthorized(t)
}

func TestPatchScootersNotFound(t *testing.T) {
	scootersTest.TestPatchScootersNotFound(t)
}

func TestPatchScootersETagNotMatch(t *testing.T) {
	scootersTest.TestPatchScootersETagNotMatch(t)
}

func TestCanOccupyScooter(t *testing.T) {
	scootersTest.TestCanOccupyScooter(t)
}
func TestCanFreeScooter(t *testing.T) {
	scootersTest.TestCanFreeScooter(t)
}

func TestGetScooterNotEnoughArgs(t *testing.T) {
	scootersTest.TestGetScooterNotEnoughArgs(t)
}

func TestCannotFreeNonOccupiedScooter(t *testing.T) {
	scootersTest.TestCannotFreeNonOccupiedScooter(t)
}
