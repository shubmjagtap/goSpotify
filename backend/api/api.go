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

func CheckUserExistenceForLogin(email string, client *mongo.Client) (bool, error) {

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

func CheckLogin(loginInfo models.UserLoginInformation, client *mongo.Client) (bool, error) {
	userCollection := client.Database("goChat").Collection("userCollection")

	var user models.User
	err := userCollection.FindOne(context.Background(), bson.M{"email": loginInfo.Email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		} else {
			return false, err
		}
	}

	if user.Password == loginInfo.Password {
		return true, nil
	} else {
		return false, nil
	}
}

func LoginUser(w http.ResponseWriter, r *http.Request, client *mongo.Client) {

	var loginInfo models.UserLoginInformation
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&loginInfo); err != nil {
		log.Printf("Error decoding JSON data: %v\n", err)
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if loginInfo.Email == "" || loginInfo.Password == "" {
		log.Println("Missing required fields")
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	userExists, err := CheckUserExistenceForLogin(loginInfo.Email, client)
	if err != nil {
		log.Printf("Error checking user existence: %v\n", err)
		http.Error(w, "Error checking user existence", http.StatusInternalServerError)
		return
	}
	if !userExists {
		log.Println("User with this email doesnot exist")
		http.Error(w, "User with this email doesnot exist", http.StatusConflict)
		return
	}

	passwordMatches, err := CheckLogin(loginInfo, client)
	if passwordMatches {
		// Passwords match, send a positive response indicating successful login.
		response := map[string]string{"message": "Successfully logged in"}
		jsonResponse, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	} else {
		// Passwords do not match, send a negative response indicating login failure.
		http.Error(w, "Failed to log in", http.StatusUnauthorized)
	}
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
