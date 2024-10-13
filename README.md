# Scootin' Aboot

An example of REST-like API written in Go powered by HUMA and GORM frameworks.

## Overview

Let's say a company called Scootin' Aboot will deploy electric scooters in Ottawa and
Montreal. We need to design and implement a backend service that exposes a REST-like
API intended for scooter event collecting and reporting to users.

1. The scooters report an event when a trip begins, report an event when the
   trip ends, and send in periodic updates on their location. After beginning a
   trip, the scooter is considered occupied. After a trip ends the scooter
   becomes free for use. A location update must contain the time, and
   geographical coordinates.
2. Users can query scooter locations and statuses in any rectangular
   location (e.g. two pair of coordinates), and filter them by status.
3. For the sake of simplicity, users will authenticate with the server
   using ID. Of course real-world API should use authentication and authorization
   using for example JWT token.

The API will expose a few REST-like endpoints to manage scooters and their timeline:

- `GET /scooters` - to get list of available scooters, also by the status and the coordinates (free, occupied)
- `POST /scooters` - to create a scooter
- `PATCH /scooters/{id}` - to edit a scooter
- `POST /users` - to create an user
- `GET /events` - to get list of all events across all scooters
- `POST /events` - to create an event for the scooter

## Frameworks and technologies used

- [HUMA](https://huma.rocks/) - a modern, simple, fast & flexible micro framework for building HTTP REST/RPC APIs in Golang backed by OpenAPI 3 and JSON Schema
- [GORM](https://gorm.io/) - the fantastic ORM library for Golang
- [PostgreSQL](https://www.postgresql.org/) - the World's Most Advanced Open Source Relational Database
- [pgAdmin](https://www.pgadmin.org/) - pgAdmin is the most popular and feature rich Open Source administration and development platform for PostgreSQL, the most advanced Open Source database in the world. *I added it only for checking and viewing DB contents, of course it should not be included in the production environment.*

## Implementation

The API is implemented as REST-like in the HAL format. Each item is an resource and have its `_links` object.

The implementation consists of 2 components:
- `api` - which is the implemented API exposed by the `main.go` on the `localhost:80` whithin the container (`localhost:8080` from the host system)
- `tests` - test suites for testing the API

The implementation does not need any migration scripts and database because it is using builtin GORM auto-migration feature.

The project is using `docker` with `docker-compose` to easily deploy and run the API, as well to connect to running API container and pgAdmin.

## Running the API

First clone the repository
```bash
$ git clone https://github.com/skazanyNaGlany/scootin-aboot.git
```

Change your current working directory to that cloned repository
```bash
$ cd scootin-aboot
```

Build the project using `docker-compose`
```bash
$ docker-compose build
```

Finally run the project
```bash
$ docker-compose up
```

After a while the API will be running and accessible on the `http://localhost:8080` URL. Use `CTRL+C` to stop the API.

## Accessing running pgAdmin instance

To access a running pgAdmin instance type `http://localhost:8888/browser/` in your web browser. Hit `ENTER` when prompt for a password because there is no a password set. Remember to always set a very-strong password on the production environment. I set it to empty only for testing purposes.

## Documentation

Thanks to HUMA framework the documentation is generated automatically from the code. You can access it at http://localhost:8080/docs or http://localhost:8080/openapi.json for the OpenAPI spec in the JSON format or http://localhost:8080/openapi.yaml in the YAML format.

## Testing

Before running the test suite make sure the project is running using `docker-compose up`.

Access `api` container shell
```bash
$ cd scootin-aboot
$ ./bin/shell.sh api
```

In the container's shell
```bash
root@api:/var/www# cd tests
root@api:/var/www/tests# go test
```

You can also run `./coverage.sh` to see tests coverage in percent and view generated `coverage.html` in your web-browser.

```bash
root@api:/var/www/tests# ./coverage.sh
```

You can access the generated `coverage.html` file from your host system in the `scootin-aboot/src/tests` directory.


## Using the API

To use the API you need to generate user ID which will be used as the access-token for calling other endpoints. There is no any other authorization system implemented, of course real-world project should use authentication with authorization and access-tokens for example in the JWT format.

Make sure the project is running before accessing the endpoints.

`POST /users` is only one endpoint which is not requiring any authorization.

Generate a user to get the user ID
```bash
$ curl --request POST \
  --url http://localhost:8080/users \
  --header 'Content-Type: application/json' \
  --data '{}'
```

The response should look like this
```json
{
	"$schema": "http://localhost:8080/schemas/User.json",
	"id": "6d962a89-e9ec-4b1f-8e93-24b9fb56e40c",
	"created_at": "2024-09-26T10:45:47.632435486Z",
	"updated_at": "2024-09-26T10:45:47.632435486Z",
	"_links": {
		"self": {
			"href": "/users/6d962a89-e9ec-4b1f-8e93-24b9fb56e40c"
		}
	}
}
```

Please note the user's ID, it will be used for all other endpoints.

Create some free scooters
```bash
curl --request POST \
  --url http://localhost:8080/scooters \
  --header 'Authorization: 6d962a89-e9ec-4b1f-8e93-24b9fb56e40c' \
  --header 'Content-Type: application/json' \
  --data '{
}
'
```

Example response
```json
{
	"$schema": "http://localhost:8080/schemas/Scooter.json",
	"id": "61f06bc7-c356-4f4b-ad92-6692e861e96f",
	"created_at": "2024-10-13T20:56:15.74682651Z",
	"updated_at": "2024-10-13T20:56:15.74682651Z",
	"status": "free",
	"user_id": "00000000-0000-0000-0000-000000000000",
	"etag": "7a65dd50-1e77-46bf-943e-31ab474ed8f2",
	"_links": {
		"self": {
			"href": "/scooters/61f06bc7-c356-4f4b-ad92-6692e861e96f"
		}
	}
}
```

Find a free scooter (or use that one created in the previous step)
```bash
curl --request GET \
  --url 'http://localhost:8080/scooters?status=free' \
  --header 'Authorization: 6d962a89-e9ec-4b1f-8e93-24b9fb56e40c'
```

Example response
```json
{
	"$schema": "http://localhost:8080/schemas/GET_Scooters_OutputBody.json",
	"_embedded": {
		"scooters": [
			{
				"id": "2410b744-5e15-4aef-8c93-be0f751ab254",
				"created_at": "2024-09-26T10:50:17.541101Z",
				"updated_at": "2024-09-26T10:50:17.541101Z",
				"status": "free",
				"user_id": "00000000-0000-0000-0000-000000000000",
				"etag": "cf7ddf9a-3088-4adf-9570-43a3f67be2f0",
				"_links": {
					"self": {
						"href": "/scooters/2410b744-5e15-4aef-8c93-be0f751ab254"
					}
				}
			},
```

Please note the scooter's `id` and `etag`, the etag is very important thing, it is some kind of the password for editing the scooter resource. It is using for optimistic locking and will avoid race-condition issues. Without proper etag you will **not** be able to edit the scooter resource.

`PATCH /scooters/{id}` is only one endpoint which is requiring etag, because only scooter resource can be patched.

Occupy the scooter using saved `id` and `etag`
```bash
$ curl --request PATCH \
  --url http://localhost:8080/scooters/2410b744-5e15-4aef-8c93-be0f751ab254 \
  --header 'Authorization: 6d962a89-e9ec-4b1f-8e93-24b9fb56e40c' \
  --header 'Content-Type: application/json' \
  --header 'If-Match: cf7ddf9a-3088-4adf-9570-43a3f67be2f0' \
  --data '{
	"status": "occupied"}'
```

Response
```json
{
	"$schema": "http://localhost:8080/schemas/Scooter.json",
	"id": "2410b744-5e15-4aef-8c93-be0f751ab254",
	"created_at": "2024-09-26T10:50:17.541101Z",
	"updated_at": "2024-09-26T10:57:33.936261074Z",
	"status": "occupied",
	"user_id": "6d962a89-e9ec-4b1f-8e93-24b9fb56e40c",
	"etag": "68ddbf51-f1a7-4309-bc16-b683aceec01f",
	"_links": {
		"self": {
			"href": "/scooters/2410b744-5e15-4aef-8c93-be0f751ab254"
		}
	}
}
```

That call will set scooter's status to `occupied` so from now on your user will be occuping the scooter, and you can create some events for it. Please not that the scooter's etag needs to be passed in the `If-Match` header anytime you want to edit the scooter.

Remember to note scooter's etag after editing it because the etag is rotated on any `PATCH` call. To edit the scooter you always need the fresh etag.

Create `start` event with some latitude and longitude coordinates

```bash
$ curl --request POST \
  --url http://localhost:8080/events \
  --header 'Authorization: 6d962a89-e9ec-4b1f-8e93-24b9fb56e40c' \
  --header 'Content-Type: application/json' \
  --data '{
	"scooter_id": "2410b744-5e15-4aef-8c93-be0f751ab254",
	"event_type": "start",
	"latitude": 8,
	"longitude": -79
}
'
```

From now on your scooter will be marked as traveling.

Create a `location_update` event
```bash
$ curl --request POST \
  --url http://localhost:8080/events \
  --header 'Authorization: 6d962a89-e9ec-4b1f-8e93-24b9fb56e40c' \
  --header 'Content-Type: application/json' \
  --data '{
	"scooter_id": "2410b744-5e15-4aef-8c93-be0f751ab254",
	"event_type": "location_update",
	"latitude": 51,
	"longitude": 19
}
'
```

Stop traveling using your scooter
```bash
curl --request POST \
  --url http://localhost:8080/events \
  --header 'Authorization: 6d962a89-e9ec-4b1f-8e93-24b9fb56e40c' \
  --header 'Content-Type: application/json' \
  --data '{
	"scooter_id": "2410b744-5e15-4aef-8c93-be0f751ab254",
	"event_type": "stop",
	"latitude": 51,
	"longitude": 19
}
'
```

Finally release the scooter so another user will be able to use it (only one user can occupy the scooter at a time)

```bash
$ curl --request PATCH \
  --url http://localhost:8080/scooters/2410b744-5e15-4aef-8c93-be0f751ab254 \
  --header 'Authorization: 6d962a89-e9ec-4b1f-8e93-24b9fb56e40c' \
  --header 'Content-Type: application/json' \
  --header 'If-Match: 68ddbf51-f1a7-4309-bc16-b683aceec01f' \
  --data '{
	"status": "free"
}
'
```

You can also query for scooters status and location
```bash
$ curl --request GET \
  --url 'http://localhost:8080/scooters?status=free&min_latitude=50&max_latitude=52&min_longitude=18&max_longitude=20' \
  --header 'Authorization: 6d962a89-e9ec-4b1f-8e93-24b9fb56e40c'
```

Response
```json
{
	"$schema": "http://localhost:8080/schemas/GET_Scooters_OutputBody.json",
	"_embedded": {
		"scooters": [
			{
				"id": "2410b744-5e15-4aef-8c93-be0f751ab254",
				"created_at": "2024-09-26T10:50:17.541101Z",
				"updated_at": "2024-09-26T11:10:22.367805Z",
				"status": "free",
				"user_id": "6d962a89-e9ec-4b1f-8e93-24b9fb56e40c",
				"etag": "e62e60ee-87f2-426d-86c8-155bebef77e1",
				"_embedded": {
					"events": [
						{
							"id": 562,
							"created_at": "2024-09-26T11:09:12.037076Z",
							"updated_at": "2024-09-26T11:09:12.037076Z",
							"scooter_id": "2410b744-5e15-4aef-8c93-be0f751ab254",
							"user_id": "6d962a89-e9ec-4b1f-8e93-24b9fb56e40c",
							"event_type": "stop",
							"latitude": 51,
							"longitude": 19,
							"_links": {
								"self": {
									"href": "/events/562"
								}
							}
						}
					]
				},
				"_links": {
					"self": {
						"href": "/scooters/2410b744-5e15-4aef-8c93-be0f751ab254"
					}
				}
			}
		]
	},
	"_links": {
		"self": {
			"href": "/scooters"
		}
	}
}
```

Get all events
```bash
$ curl --request GET \
  --url http://localhost:8080/events \
  --header 'Authorization: 6d962a89-e9ec-4b1f-8e93-24b9fb56e40c' \
  --header 'Content-Type: application/json'
```

Response
```json
{
	"$schema": "http://localhost:8080/schemas/GET_Events_OutputBody.json",
	"_embedded": {
		"events": [
			{
				"id": 550,
				"created_at": "2024-09-26T10:50:17.556717Z",
				"updated_at": "2024-09-26T10:50:17.556717Z",
				"scooter_id": "2410b744-5e15-4aef-8c93-be0f751ab254",
				"user_id": "00000000-0000-0000-0000-000000000000",
				"event_type": "location_update",
				"latitude": 0,
				"longitude": 1,
				"_links": {
					"self": {
						"href": "/events/550"
					}
				}
			},
```

## Known issues

- `$schema` is generated by the HUMA framework and currently the link goes to `404 Not Found`
- there is no pagination implemented, of course real-world API should have a pagination implemented to return resources only on the desired page
