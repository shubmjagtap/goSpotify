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

func HomeHandler(w http.ResponseWriter, r *http.Request, client *mongo.Client) {
	// Get a handle to the "goChat" database and the "userCollection" collection
	database := client.Database("goChat")
	collection := database.Collection("userCollection")

	// Define a filter to match all documents (if you want all data)
	filter := bson.M{}

	// Perform a find operation to retrieve all documents in the collection
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Printf("Error querying database: %v\n", err)
		http.Error(w, "Error querying database", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	// Create a slice to hold the decoded documents
	var results []bson.M

	// Iterate over the cursor and decode documents into the results slice
	for cursor.Next(context.TODO()) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			log.Printf("Error decoding document: %v\n", err)
			http.Error(w, "Error decoding document", http.StatusInternalServerError)
			return
		}
		results = append(results, result)
	}

	// If there are no results, return a message
	if len(results) == 0 {
		fmt.Fprintln(w, "No data found")
		return
	}

	// Convert the results slice to JSON
	jsonData, err := json.Marshal(results)
	if err != nil {
		log.Printf("Error encoding JSON: %v\n", err)
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}

	// Set the response content type and write the JSON data to the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
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

func RegisterUser(newUser models.User, client *mongo.Client) error {
	goChatDB := client.Database("goChat")
	userCollection := goChatDB.Collection("userCollection")
	_, err := userCollection.InsertOne(context.Background(), newUser)
	if err != nil {
		fmt.Println("Insertion Error")
		return err
	}
	return nil
}

// function for handling signup endpoint
func SignUpUser(w http.ResponseWriter, r *http.Request, client *mongo.Client) {
	var newUser models.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&newUser); err != nil {
		log.Printf("Error decoding JSON data: %v\n", err)
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if newUser.Name == "" || newUser.Email == "" || newUser.Password == "" || newUser.Pic == "" {
		log.Println("Missing required fields")
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	userExists, err := CheckUserExistence(newUser.Email, client)
	if err != nil {
		log.Printf("Error checking user existence: %v\n", err)
		http.Error(w, "Error checking user existence", http.StatusInternalServerError)
		return
	}
	if userExists {
		log.Println("User with this email already exists")
		http.Error(w, "User with this email already exists", http.StatusConflict)
		return
	}

	registrationErr := RegisterUser(newUser, client)
	if registrationErr != nil {
		log.Printf("Error in registration of user: %v\n", registrationErr)
		http.Error(w, "Error in registration of user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User registered successfully"))
}
