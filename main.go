package conecto

func main() {

	// Database
	db := postgres.Connect(...)

	// Stores
	connectionStore := connections.NewPostgresStore(db)
	credentialStore := credentials.NewAESStore(...)
	stateStore := oauth2.NewPostgresStateStore(db)

	// Connectors
	shopify := shopify.New(...)
	meta := meta.New(...)
	google := googleads.New(...)

	connectors := map[string]connector.Connector{
		shopify.Name(): shopify,
		meta.Name(): meta,
		google.Name(): google,
	}

	// OAuth service
	oauthService := oauth2.NewService(
		connectionStore,
		credentialStore,
		stateStore,
		connectors,
	)

	// HTTP handlers
	oauthHandler := oauth2.NewHandler(oauthService)

	// Router
	router := NewRouter(oauthHandler)

	http.ListenAndServe(":8080", router)
}