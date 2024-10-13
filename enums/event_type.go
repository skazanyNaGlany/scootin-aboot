package enums

type EventType string

const (
	EventTypeStart          EventType = "start"
	EventTypeStop           EventType = "stop"
	EventTypeLocationUpdate EventType = "location_update"
)
