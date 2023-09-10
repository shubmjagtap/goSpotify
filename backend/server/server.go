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

	router.HandleFunc("/user/api/signup", func(w http.ResponseWriter, r *http.Request) {
		api.SignUpUser(w, r, client)
	}).Methods("POST")

	router.HandleFunc("/user/api/allusers", func(w http.ResponseWriter, r *http.Request) {
		api.GetAllUsersHandler(w, r, client)
	}).Methods("GET")

	router.HandleFunc("/user/api/login", func(w http.ResponseWriter, r *http.Request) {
		api.LoginUser(w, r, client)
	}).Methods("POST")

	router.HandleFunc("/user/api/chat", func(w http.ResponseWriter, r *http.Request) {
		api.AccessChat(w, r, client) // Pass client to AccessChat
	}).Methods("POST")

	router.HandleFunc("/user/api/chat", func(w http.ResponseWriter, r *http.Request) {
		api.FetchChats(w, r, client) // Pass client to FetchChats
	}).Methods("GET")

	router.HandleFunc("/user/api/chat/group", func(w http.ResponseWriter, r *http.Request) {
		api.CreateGroupChat(w, r, client) // Pass client to CreateGroupChat
	}).Methods("POST")

	router.HandleFunc("/user/api/chat/rename", func(w http.ResponseWriter, r *http.Request) {
		api.RenameGroup(w, r, client) // Pass client to RenameGroup
	}).Methods("PUT")

	router.HandleFunc("/user/api/chat/groupremove", func(w http.ResponseWriter, r *http.Request) {
		api.RemoveFromGroup(w, r, client) // Pass client to RemoveFromGroup
	}).Methods("PUT")

	router.HandleFunc("/user/api/chat/groupadd", func(w http.ResponseWriter, r *http.Request) {
		api.AddToGroup(w, r, client) // Pass client to AddToGroup
	}).Methods("PUT")
}

func SetupCorsMiddleware(router *mux.Router) http.Handler {
	corsMiddleware := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "OPTIONS", "POST"}),
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
