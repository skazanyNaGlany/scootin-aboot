package hal

// Links represents a collection of links related to a resource.
type Links struct {
	Self Self `json:"self" doc:"Link to this resource"`
}
