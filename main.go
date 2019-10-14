package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

const port = "127.0.0.1:8000"

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "Hello world")
	})

	fmt.Printf("Listening on %s\n", port)
	http.ListenAndServe(port, router)
}
