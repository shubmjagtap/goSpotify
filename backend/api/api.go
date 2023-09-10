package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/shubmjagtap/goSpotify/backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetAllUsersHandler(w http.ResponseWriter, r *http.Request, client *mongo.Client) {
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

func CheckUserExistenceForLogin(email string, client *mongo.Client) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	userCollection := client.Database("goChat").Collection("userCollection")
	filter := bson.M{"email": email}
	var user models.User
	err := userCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.User{}, nil // User not found, return empty user
		}
		return models.User{}, err
	}
	return user, nil
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
		log.Println("Missing required fields in LoginUser")
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	user, err := CheckUserExistenceForLogin(loginInfo.Email, client)

	if err != nil {
		log.Printf("Error checking user existence: %v\n", err)
		http.Error(w, "Error checking user existence", http.StatusInternalServerError)
		return
	}

	if user.Email == "" {
		log.Println("User with this email does not exist")
		http.Error(w, "User with this email does not exist", http.StatusConflict)
		return
	}

	// User with the given email exists, now check the password
	if user.Password == loginInfo.Password {
		// Passwords match, return positive response with user data
		jsonResponse, _ := json.Marshal(user)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	} else {
		// Passwords do not match, return negative response with user (nil)
		jsonResponse, _ := json.Marshal(nil)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized) // Use 401 for authentication failure
		w.Write(jsonResponse)
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

func RegisterUser(newUser models.User, client *mongo.Client) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	goChatDB := client.Database("goChat")
	userCollection := goChatDB.Collection("userCollection")
	filter := bson.M{"email": newUser.Email}
	_, err := userCollection.InsertOne(context.Background(), newUser)
	if err != nil {
		fmt.Println("Insertion Error")
		return models.User{}, err
	}
	var user models.User
	errr := userCollection.FindOne(ctx, filter).Decode(&user)
	if errr != nil {
		if err == mongo.ErrNoDocuments {
			return models.User{}, nil // User not found, return empty user
		}
		return models.User{}, err
	}
	return user, nil
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

	if newUser.Name == "" || newUser.Email == "" || newUser.Password == "" {
		log.Println("Missing required fields in SignUp function")
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

	newUser, registrationErr := RegisterUser(newUser, client)
	if registrationErr != nil {
		log.Printf("Error in registration of user: %v\n", registrationErr)
		http.Error(w, "Error in registration of user", http.StatusInternalServerError)
		return
	}
	jsonResponse, _ := json.Marshal(newUser)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func AccessChat(w http.ResponseWriter, r *http.Request, client *mongo.Client) {
	fmt.Println("AccessChat function called")
	// You can add any specific logic here if needed
	w.WriteHeader(http.StatusOK)
}

func FetchChats(w http.ResponseWriter, r *http.Request, client *mongo.Client) {
	fmt.Println("FetchChats function called")
	// You can add any specific logic here if needed
	w.WriteHeader(http.StatusOK)
}

func CreateGroupChat(w http.ResponseWriter, r *http.Request, client *mongo.Client) {
	fmt.Println("CreateGroupChat function called")
	// You can add any specific logic here if needed
	w.WriteHeader(http.StatusOK)
}

func RenameGroup(w http.ResponseWriter, r *http.Request, client *mongo.Client) {
	fmt.Println("RenameGroup function called")
	// You can add any specific logic here if needed
	w.WriteHeader(http.StatusOK)
}

func RemoveFromGroup(w http.ResponseWriter, r *http.Request, client *mongo.Client) {
	fmt.Println("RemoveFromGroup function called")
	// You can add any specific logic here if needed
	w.WriteHeader(http.StatusOK)
}

func AddToGroup(w http.ResponseWriter, r *http.Request, client *mongo.Client) {
	fmt.Println("AddToGroup function called")
	// You can add any specific logic here if needed
	w.WriteHeader(http.StatusOK)
}
