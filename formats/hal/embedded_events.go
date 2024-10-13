package hal

// EmbeddedEvents represents a collection of embedded events in HAL format.
type EmbeddedEvents struct {
	Events []Event `json:"events" doc:"List of events"`
}
