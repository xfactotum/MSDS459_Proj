package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"net/http"
	"time"

	huggingface "github.com/hupe1980/go-huggingface"
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

func hf_main() {
	ic := huggingface.NewInferenceClient(os.Getenv("HUGGINGFACEHUB_API_TOKEN"))

	res1, err1 := ic.FeatureExtraction(context.Background(), &huggingface.FeatureExtractionRequest{
		Inputs: []string{"Hello World"},
		Options: huggingface.Options{
			WaitForModel: huggingface.PTR(true),
		},
	})
	if err1 != nil {
		log.Fatal(err1)
	}

	fmt.Println("FeatureExtraction:")
	fmt.Println(res1[0])
	fmt.Println()

	res2, err2 := ic.FeatureExtractionWithAutomaticReduction(context.Background(), &huggingface.FeatureExtractionRequest{
		Inputs: []string{"Hello World"},
		Model:  "sentence-transformers/all-mpnet-base-v2",
		// Model:  "intfloat/e5-small-v2",
		Options: huggingface.Options{
			WaitForModel: huggingface.PTR(true),
		},
	})
	if err2 != nil {
		log.Fatal(err2)
	}

	fmt.Println("FeatureExtractionWithAutomaticReduction:")
	fmt.Println(res2[0])
	fmt.Println()
}
