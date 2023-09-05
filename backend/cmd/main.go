package main

import (
	"fmt"

	"github.com/shubmjagtap/goSpotify/backend/database"
	"github.com/shubmjagtap/goSpotify/backend/server"
)

func main() {
	// Initialize the database
	client, ctx, cancel, err := database.InitDb()
	if err != nil {
		fmt.Println("Failed to initialize MongoDB:", err)
		return
	}
	defer database.Close(client, ctx, cancel)

	// Create a router
	router := server.CreateRouter()

	// Create routes for the router
	server.CreateRoutes(router)

	// Setup CORS middleware
	corsRouter := server.SetupCorsMiddleware(router)

	// Start the server
	port := "8081"
	server.StartServer(port, corsRouter)
}
