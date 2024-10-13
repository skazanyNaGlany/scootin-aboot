package params

import "github.com/google/uuid"

// AuthorizationParam represents the parameters required for authorization.
type AuthorizationParam struct {
	Authorization uuid.UUID `doc:"API key (user ID)" header:"Authorization"`
}
