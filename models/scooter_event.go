package models

// ScooterEvent represents an event related to a scooter.
// It is virual and does not have a corresponding table in the database.
type ScooterEvent struct {
	Scooter *Scooter
	Event   *Event
}
