package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"net/http"
	"time"
)


func main() {
	// handle static assets
	mux := http.NewServeMux()
	files := http.FileServer(http.Dir("public"))

	//
	// all route patterns matched here
	// route handler functions defined in other files
	//

	// defined in route_auth.go
	mux.Handle("/", files)
  
	mux.HandleFunc("/recommendBreed", recommendBreed)
  
	mux.HandleFunc("/getBreedDetail", getBreedDetail)

	// starting up the server
	server := &http.Server{
		Addr:           "127.0.0.1:8989",
		Handler:        mux,
		ReadTimeout:    time.Duration(10 * int64(time.Second)),
		WriteTimeout:   time.Duration(600 * int64(time.Second)),
		MaxHeaderBytes: 1 << 20,
	}
	server.ListenAndServe()
}