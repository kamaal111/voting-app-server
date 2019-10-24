package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DatabaseCollections defines the database collections structure
type DatabaseCollections struct {
	sessions string
}

// Config defines how the configuration structure
type Config struct {
	port                string
	databaseURL         string
	databaseName        string
	databaseCollections DatabaseCollections
}

var config = Config{
	port:         "127.0.0.1:8000",
	databaseURL:  "mongodb://127.0.0.1:27017",
	databaseName: "voting_app",
	databaseCollections: DatabaseCollections{
		sessions: "sessions",
	},
}

// Candidate struct represents a candidate
type Candidate struct {
	Name  string   `json:"name,omitempty" bson:"name,omitempty"`
	Votes []string `json:"votes,omitempty" bson:"votes,omitempty"`
}

// An Session represents a voting session
type Session struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name       string             `json:"name,omitempty" bson:"name,omitempty"`
	Candidates []Candidate        `json:"candidates,omitempty" bson:"candidates,omitempty"`
}

var client *mongo.Client

func createSession(response http.ResponseWriter, request *http.Request) {
	fmt.Printf("POST %s/sessions\n", config.port)
	response.Header().Add("content-type", "application/json")
	var session Session
	json.NewDecoder(request.Body).Decode(&session)
	collection := client.Database(config.databaseName).Collection(config.databaseCollections.sessions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, _ := collection.InsertOne(ctx, session)
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(result)
}

func getAllSessions(response http.ResponseWriter, request *http.Request) {
	fmt.Printf("GET %s/sessions\n", config.port)
	response.Header().Add("content-type", "application/json")
	var allSessions []Session
	collection := client.Database(config.databaseName).Collection(config.databaseCollections.sessions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var session Session
		cursor.Decode(&session)
		allSessions = append(allSessions, session)
	}
	if len(allSessions) == 0 {
		response.WriteHeader(http.StatusNotFound)
		response.Write([]byte(`{"message": "Not Found"}`))
		return
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(allSessions)
}

func getOneSession(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var session Session
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	fmt.Printf("GET %s/sessions/%s\n", config.port, params["id"])
	collection := client.Database(config.databaseName).Collection(config.databaseCollections.sessions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := collection.FindOne(ctx, Session{ID: id}).Decode(&session)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(response).Encode(session)
}

func main() {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
		Debug:            true,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI(config.databaseURL)
	client, _ = mongo.Connect(ctx, clientOptions)
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/sessions", getAllSessions).Methods("GET")
	router.HandleFunc("/sessions", createSession).Methods("POST")
	router.HandleFunc("/sessions/{id}", getOneSession).Methods("GET")

	fmt.Printf("Listening on %s\n", config.port)
	handler := c.Handler(router)
	log.Fatal(http.ListenAndServe(config.port, handler))
}
