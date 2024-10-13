package params

import "github.com/google/uuid"

// ETagUUIDParam represents the ETag and UUID parameters for a resource.
type ETagUUIDParam struct {
	ETag uuid.UUID `doc:"ETag of the resource, for concurrency control" header:"If-Match"`
}
