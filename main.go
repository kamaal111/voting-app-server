package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"github.com/kamaal111/voting-app-server/config"
	"github.com/kamaal111/voting-app-server/handlers"
)

func main() {
	c := cors.New(cors.Options{
		// AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT"},
		AllowCredentials: true,
		Debug:            true,
	})

	handlers.ConnectToClient()

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/sessions", handlers.GetAllSessions).Methods("GET")
	router.HandleFunc("/sessions", handlers.CreateSession).Methods("POST")
	router.HandleFunc("/sessions/{id}", handlers.GetOneSession).Methods("GET")
	router.HandleFunc("/sessions/{id}/{vote}", handlers.UpdateVoteInToSession).Methods("PUT")

	fmt.Printf("Listening on %s\n", config.Config.Port)
	handler := c.Handler(router)
	log.Fatal(http.ListenAndServe(config.Config.Port, handler))
}
