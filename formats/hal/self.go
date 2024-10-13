package hal

// Self represents a link to the current resource.
type Self struct {
	Href string `json:"href" doc:"Link URL"`
}
