package main

func main() {
	// Create a router
	router := createRouter()

	// Create routes for the router
	createRoutes(router)

	// Setup CORS middleware
	corsRouter := setupCorsMiddleware(router)

	// Start the server
	port := "8081"
	startServer(port, corsRouter)
}
