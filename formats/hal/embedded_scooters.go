package hal

// EmbeddedScooters represents a collection of scooters embedded in a HAL response.
type EmbeddedScooters struct {
	Scooters []Scooter `json:"scooters" doc:"List of scooters"`
}
