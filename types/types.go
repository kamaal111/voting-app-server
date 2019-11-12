package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DatabaseCollections defines the database collections structure
type DatabaseCollections struct {
	Sessions string
}

// Config defines how the configuration structure
type Config struct {
	Port                string
	DatabaseURL         string
	DatabaseName        string
	DatabaseCollections DatabaseCollections
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

// An Response represents the response from the server
type Response struct {
	Session Session `json:"session,omitempty" bson:"session,omitempty"`
	Message string  `json:"message,omitempty" bson:"message,omitempty"`
}
