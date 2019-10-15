package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

const port = "127.0.0.1:8000"

// An Session represents a voting session
type Session struct {
	ID   int    `json:"id"`   // Session ID
	Name string `json:"name"` // Session Name
}

func homeLink(writer http.ResponseWriter, request *http.Request) {

	fmt.Printf("%s/\n", port)
	fmt.Fprintf(writer, "Hello and welcome")
}

func getAllSessions(writer http.ResponseWriter, request *http.Request) {
	// AllSessions represents all available sessions
	type AllSessions []Session

	sessions := AllSessions{
		{
			ID:   1,
			Name: "Kamaal",
		},
	}

	fmt.Printf("GET %s/sessions\n", port)
	json.NewEncoder(writer).Encode(sessions)
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", homeLink)
	router.HandleFunc("/sessions", getAllSessions).Methods("GET")

	fmt.Printf("Listening on %s\n", port)
	http.ListenAndServe(port, router)
}
