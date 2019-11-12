package config

import (
	"github.com/kamaal111/voting-app-server/types"
)

// Config ...
var Config = types.Config{
	Port:         "127.0.0.1:8000",
	DatabaseURL:  "mongodb://127.0.0.1:27017",
	DatabaseName: "voting_app",
	DatabaseCollections: types.DatabaseCollections{
		Sessions: "sessions",
	},
}
