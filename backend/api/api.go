package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/shubmjagtap/goSpotify/backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "API is Running...")
	log.Printf("Home page accessed from IP: %s", r.RemoteAddr)
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Logging In User...")
	log.Printf("Home page accessed from IP: %s", r.RemoteAddr)
}

func CheckUserExistence(email string, client *mongo.Client) (bool, error) {

	goChatDB := client.Database("goChat")
	userCollection := goChatDB.Collection("userCollection")

	filter := bson.M{"email": email}

	var existingUser bson.M
	if err := userCollection.FindOne(context.Background(), filter).Decode(&existingUser); err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// RegisterUser inserts a new user into the MongoDB database.
func RegisterUser(newUser models.User, client *mongo.Client) error {
	goChatDB := client.Database("goChat")
	userCollection := goChatDB.Collection("userCollection")
	_, err := userCollection.InsertOne(context.Background(), newUser)
	if err != nil {
		return err
	}
	return nil
}

func SignUpUser(w http.ResponseWriter, r *http.Request) {

	var clientt *mongo.Client

	var newUser models.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&newUser); err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if newUser.Name == "" || newUser.Email == "" || newUser.Password == "" || newUser.Pic == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	userExists, err := CheckUserExistence(newUser.Email, clientt)
	if err != nil {
		http.Error(w, "Error checking user existence", http.StatusInternalServerError)
		return
	}
	if userExists {
		http.Error(w, "User with this email already exists", http.StatusConflict)
		return
	}

	registerationErr := RegisterUser(newUser, clientt)
	if registerationErr != nil {
		http.Error(w, "Error in registration of  user", http.StatusInternalServerError)
		return
	}
}

// if user is not present then make new entry for that user
// if user is sucessfully created then send that user else send unsuccessfull message
