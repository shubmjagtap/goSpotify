package api

import (
	"fmt"
	"log"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "API is Running...")
	log.Printf("Home page accessed from IP: %s", r.RemoteAddr)
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Logging In User...")
	log.Printf("Home page accessed from IP: %s", r.RemoteAddr)
}

func SignUpUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Signing up user...")
	log.Printf("Home page accessed from IP: %s", r.RemoteAddr)
}
