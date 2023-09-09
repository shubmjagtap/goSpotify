package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/shubmjagtap/goSpotify/backend/api"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateRouter() *mux.Router {
	router := mux.NewRouter()
	return router
}

func CreateRoutes(router *mux.Router, client *mongo.Client) {
	router.HandleFunc("/user/api/login", api.LoginUser).Methods("POST")
	router.HandleFunc("/user/api/signup", func(w http.ResponseWriter, r *http.Request) {
		api.SignUpUser(w, r, client)
	}).Methods("POST")
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		api.HomeHandler(w, r, client)
	}).Methods("GET")
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
