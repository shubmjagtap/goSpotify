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

	// Setup the database and collection
	userCollection, err := database.SetupDatabase(client)
	if err != nil {
		fmt.Println("Failed to setup database:", err)
		return
	}
	fmt.Println(userCollection)

	// Create a router
	router := server.CreateRouter()

	// Create routes for the router
	server.CreateRoutes(router, client)

	// Setup CORS middleware
	corsRouter := server.SetupCorsMiddleware(router)

	// Start the server
	port := "8081"
	server.StartServer(port, corsRouter)
}
