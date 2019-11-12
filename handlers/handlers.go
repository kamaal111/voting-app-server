package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/kamaal111/voting-app-server/config"
	"github.com/kamaal111/voting-app-server/types"
)

var client *mongo.Client

// ConnectToClient connects to MongoDB client
func ConnectToClient() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(config.Config.DatabaseURL)
	client, _ = mongo.Connect(ctx, clientOptions)

	fmt.Println("Connected to MongoDB!")
}

func externalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil || ip.IsLoopback() {
				continue
			}

			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}

			return ip.String(), nil
		}
	}

	return "", errors.New("are you connected to the network?")
}

// CreateSession is a POST handler to create one session
func CreateSession(response http.ResponseWriter, request *http.Request) {
	fmt.Printf("POST %s/sessions\n", config.Config.Port)
	response.Header().Add("content-type", "application/json")
	var session types.Session
	json.NewDecoder(request.Body).Decode(&session)
	collection := client.Database(config.Config.DatabaseName).Collection(config.Config.DatabaseCollections.Sessions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, _ := collection.InsertOne(ctx, session)
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(result)
}

// GetAllSessions is a GET handler to read all sessions
func GetAllSessions(response http.ResponseWriter, request *http.Request) {
	fmt.Printf("GET %s/sessions\n", config.Config.Port)
	response.Header().Add("content-type", "application/json")
	var allSessions []types.Session
	collection := client.Database(config.Config.DatabaseName).Collection(config.Config.DatabaseCollections.Sessions)
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
		var session types.Session
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

// GetOneSession is a GET handler to read one session
func GetOneSession(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var session types.Session
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	fmt.Printf("GET %s/sessions/%s\n", config.Config.Port, params["id"])
	collection := client.Database(config.Config.DatabaseName).Collection(config.Config.DatabaseCollections.Sessions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := collection.FindOne(ctx, types.Session{ID: id}).Decode(&session)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(response).Encode(session)
}

// UpdateVoteInToSession is a PUT handler to update one session with one vote
func UpdateVoteInToSession(response http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	vote, _ := params["vote"]
	fmt.Printf("PUT %s/sessions/%s\n", config.Config.Port, params["id"])
	response.Header().Add("content-type", "application/json")
	var sessionResponse types.Response
	collection := client.Database(config.Config.DatabaseName).Collection(config.Config.DatabaseCollections.Sessions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := collection.FindOne(ctx, types.Session{ID: id}).Decode(&sessionResponse.Session)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	ip, err := externalIP()
	if err != nil {
		fmt.Println(err)
	}
	sessionResponse.Message = vote
	for index := 0; index < len(sessionResponse.Session.Candidates); index++ {
		if sessionResponse.Session.Candidates[index].Name == vote {
			sessionResponse.Session.Candidates[index].Votes = append(sessionResponse.Session.Candidates[index].Votes, ip)
		}
	}
	_, err = collection.UpdateOne(ctx, types.Session{ID: id}, sessionResponse.Session.Candidates)
	if err != nil {
		fmt.Println(err)
	}
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(sessionResponse)
}
