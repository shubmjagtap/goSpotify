package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", api.homeHandler).Methods("GET")

	router.HandleFunc("/api/data/{id:[0-9a-zA-Z]+}", api.apiDataIdHandler).Methods("GET")

	router.HandleFunc("/api/data", api.apiDataHandler).Methods("GET")

	// Create CORS middleware with desired options
	corsMiddleware := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),              // Allow all origins (you can specify specific origins if needed)
		handlers.AllowedMethods([]string{"GET", "OPTIONS"}), // Allow only GET and OPTIONS methods
		handlers.AllowedHeaders([]string{"Content-Type"}),   // Allow only the Content-Type header
	)

	// Wrap the router with the CORS middleware
	corsRouter := corsMiddleware(router)

	port := "8081"
	fmt.Printf("Server started on :%s\n", port)
	err := http.ListenAndServe(":"+port, corsRouter)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
