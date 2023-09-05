package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/shubmjagtap/goSpotify/backend/models"
)

var chats = []models.Chat{
	{
		IsGroupChat: false,
		Users: []models.User{
			{
				Name:  "John Doe",
				Email: "john@example.com",
			},
			{
				Name:  "Piyush",
				Email: "piyush@example.com",
			},
		},
		ID:       "617a077e18c25468bc7c4dd4",
		ChatName: "John Doe",
	},
	{
		IsGroupChat: false,
		Users: []models.User{
			{
				Name:  "Guest User",
				Email: "guest@example.com",
			},
			{
				Name:  "Piyush",
				Email: "piyush@example.com",
			},
		},
		ID:       "617a077e18c25468b27c4dd4",
		ChatName: "Guest User",
	},
	{
		IsGroupChat: false,
		Users: []models.User{
			{
				Name:  "Anthony",
				Email: "anthony@example.com",
			},
			{
				Name:  "Piyush",
				Email: "piyush@example.com",
			},
		},
		ID:       "617a077e18c2d468bc7c4dd4",
		ChatName: "Anthony",
	},
	{
		IsGroupChat: true,
		Users: []models.User{
			{
				Name:  "John Doe",
				Email: "jon@example.com",
			},
			{
				Name:  "Piyush",
				Email: "piyush@example.com",
			},
			{
				Name:  "Guest User",
				Email: "guest@example.com",
			},
		},
		ID:       "617a518c4081150716472c78",
		ChatName: "Friends",
		GroupAdmin: &models.User{
			Name:  "Guest User",
			Email: "guest@example.com",
		},
	},
	{
		IsGroupChat: false,
		Users: []models.User{
			{
				Name:  "Jane Doe",
				Email: "jane@example.com",
			},
			{
				Name:  "Piyush",
				Email: "piyush@example.com",
			},
		},
		ID:       "617a077e18c25468bc7cfdd4",
		ChatName: "Jane Doe",
	},
	{
		IsGroupChat: true,
		Users: []models.User{
			{
				Name:  "John Doe",
				Email: "jon@example.com",
			},
			{
				Name:  "Piyush",
				Email: "piyush@example.com",
			},
			{
				Name:  "Guest User",
				Email: "guest@example.com",
			},
		},
		ID:       "617a518c4081150016472c78",
		ChatName: "Chill Zone",
		GroupAdmin: &models.User{
			Name:  "Guest User",
			Email: "guest@example.com",
		},
	},
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "API is Running...")
	log.Printf("Home page accessed from IP: %s", r.RemoteAddr)
}

func ApiDataHandler(w http.ResponseWriter, r *http.Request) {
	chatsJSON, err := json.Marshal(chats)
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(chatsJSON)
	if err != nil {
		http.Error(w, "Failed to write JSON response", http.StatusInternalServerError)
	}
	log.Printf("API data served to IP: %s", r.RemoteAddr)
}

func ApiDataIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Find the chat with the matching ID
	var foundChat models.Chat
	for _, chat := range chats {
		if chat.ID == id {
			foundChat = chat
			break
		}
	}

	// Check if a chat with the specified ID was found
	if foundChat.ID == "" {
		http.Error(w, "Chat not found", http.StatusNotFound)
		return
	}

	// Marshal the found chat to JSON
	chatJSON, err := json.Marshal(foundChat)
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(chatJSON)
	if err != nil {
		http.Error(w, "Failed to write JSON response", http.StatusInternalServerError)
	}
	log.Printf("API data served to IP: %s for chat ID: %s", r.RemoteAddr, id)
}
