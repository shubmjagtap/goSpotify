package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/shubmjagtap/goSpotify/backend/api"
)

func CreateRouter() *mux.Router {
	router := mux.NewRouter()
	return router
}

func CreateRoutes(router *mux.Router) {
	router.HandleFunc("/", api.HomeHandler).Methods("GET")
	router.HandleFunc("/user/api/login", api.LoginUser).Methods("POST")
	router.HandleFunc("/user/api/signup", api.SignUpUser).Methods("POST")
}

func SetupCorsMiddleware(router *mux.Router) http.Handler {
	corsMiddleware := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
	)
	corsRouter := corsMiddleware(router)
	return corsRouter
}

func StartServer(port string, handler http.Handler) {
	fmt.Printf("Server started on :%s\n", port)
	err := http.ListenAndServe(":"+port, handler)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
