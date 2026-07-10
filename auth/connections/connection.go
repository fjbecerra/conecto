package connections


type Connection struct {
	ID string
	TenantID string
	Provider string
	ExternalID string //The identifier used by the external system.
	Metadata map[string]any //This lets each connector persist small pieces of configuration without changing the schema.
	Status string //Lifecycle of the connection.
}

