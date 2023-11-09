package main

import (
	"context"
	"fmt"
	"io"
	"encoding/json"
	"os"
	"net/http"
	"strings"

	huggingface "github.com/hupe1980/go-huggingface"
	"github.com/edgedb/edgedb-go"
)

var results []byte
var userwish = ""
var embedding []float64

func arrayToString(a []float64, delim string) string {
    return strings.Trim(strings.Replace(fmt.Sprint(a), " ", delim, -1), "[]")
}

func findReview() {
    opts := edgedb.Options{Concurrency: 4}
    ctx := context.Background()
    db, err := edgedb.CreateClient(ctx, opts)
    if err != nil {
        panic(err)
    }
    defer db.Close()

	query := fmt.Sprintf(
				`select Reviews {
					breed := (select Breed {name} filter Reviews.BreedReview.id = .id)
				} 
				order by ext::pgvector::cosine_distance(Reviews.embedding, <TxEmbedding>[%s])
				limit 5`, 
				arrayToString(embedding, ", "))
    err = db.QueryJSON(ctx, query, &results)
    if err != nil {
        panic(err)
    }
}

func hfModel() {
	ic := huggingface.NewInferenceClient(os.Getenv("HUGGINGFACEHUB_API_TOKEN"))

	feat, err := ic.FeatureExtractionWithAutomaticReduction(context.Background(), &huggingface.FeatureExtractionRequest{
		Inputs: []string{userwish},
		Model:  "sentence-transformers/all-mpnet-base-v2",
		Options: huggingface.Options{
			WaitForModel: huggingface.PTR(true),
		},
	})
	if err != nil {
		panic(err)
	}

	embedding = feat[0]

	fmt.Println("FeatureExtractionWithAutomaticReduction:")
	fmt.Println(embedding)
	fmt.Println()
}

func parseMap(aMap map[string]interface{}) {
    for key, val := range aMap {
        switch concreteVal := val.(type) {
        case map[string]interface{}:
            fmt.Println(key)
            parseMap(val.(map[string]interface{}))
        case []interface{}:
            fmt.Println(key)
            parseArray(val.([]interface{}))
        default:
            fmt.Println(key, ":", concreteVal)
			switch key {
				case "userwish":
					userwish = fmt.Sprintf("%v", concreteVal)
					hfModel()
					findReview()
			}
        }
    }
}

func parseArray(anArray []interface{}) {
    for i, val := range anArray {
        switch concreteVal := val.(type) {
        case map[string]interface{}:
            fmt.Println("Index:", i)
            parseMap(val.(map[string]interface{}))
        case []interface{}:
            fmt.Println("Index:", i)
            parseArray(val.([]interface{}))
        default:
            fmt.Println("Index", i, ":", concreteVal)

        }
    }
}

// POST /recommenderBreed
func recommendBreed(writer http.ResponseWriter, req *http.Request) {
	var jsonBody map[string]interface{}

	body, _ := io.ReadAll(req.Body)
	json.Unmarshal([]byte(body), &jsonBody)
	parseMap(jsonBody)

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(string(results[:]))
    fmt.Println("JSON Body:", string(results[:]))

}
