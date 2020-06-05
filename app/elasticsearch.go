package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"
	"flag"
	"strings"
	"os/exec"
	"github.com/dustin/go-humanize"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

var (
	elasticClient *elasticsearch.Client
)

func startElasticsearchConnection() {
	log.Print("Starting the Elasticsearch Client")
	var (
		r map[string]interface{}
	)

	ES01IP := getlocalIPAddress()
	cfg := elasticsearch.Config{
		Addresses: []string {
			ES01IP,
		},
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	for {
		res, err := es.Info()
		if err == nil {
			json.NewDecoder(res.Body).Decode(&r)
			log.Println("Connected to Elasticsearch :", r["name"])
			log.Println("IP Address :", ES01IP)
			res.Body.Close()
			break
		} else {
			log.Println("Waiting for Elasticsearch connection...")
			time.Sleep(5 * time.Second)
		}
	}

	// Passs es variable to elasticClient variable
	es_temp := &elasticClient
	*es_temp = es

	// Run First time only to create elasticsearch db
	//
	// insertBulkDocument()

}

func getlocalIPAddress() string {
	output, err := exec.Command("hostname", "-i").Output()
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(5 * time.Second)
	output[len(output)-2] -= 1
	return "http://" + string(output[:len(output)-1]) + ":9200"
}

func loadReviewsAndKeyword() ([]FoodReview, []FoodKeyword, error){
	log.Print("Loading data...")
	
	// Get reviewDatas from getReviewData()
	reviewsData, err := getReviewData()
	if err != nil {
		log.Fatal(err)
	}

	// Get foodKeywords from getFoodKeyword()
	foodKeywords, err := getFoodKeyword()
	if err != nil {
		log.Fatal(err)
	}
	
	// log.Print(foodKeyword)
	log.Print("Loading Completed")
	return reviewsData, foodKeywords, nil
}

func insertBulkDocument() {

	var (
		_ = fmt.Print
		count int
		batch int
	)

	rand.Seed(time.Now().UnixNano())
	
	type bulkResponse struct {
		Errors bool `json:"errors"`
		Items []struct {
			Index struct {
				ID     string `json:"_id"`
				Result string `json:"result"`
				Status int    `json:"status"`
				Error  struct {
					Type   string `json:"type"`
					Reason string `json:"reason"`
					Cause  struct {
						Type   string `json:"type"`
						Reason string `json:"reason"`
					} `json:"caused_by"`
				} `json:"error"`
			} `json:"index"`
		} `json:"items"`
	}

	var (
		buf bytes.Buffer
		res *esapi.Response
		err error
		raw map[string]interface{}
		blk *bulkResponse

		numItems   int
		numErrors  int
		numIndexed int
		currBatch  int
	)

	// reviewDatas, foodKeyword, err := loadReviewsAndKeyword()
	reviewDatas, _, err := loadReviewsAndKeyword()
	// _, foodKeyword, err := loadReviewsAndKeyword()
	if err != nil {
		log.Fatalf("Error loading data")
	}

	// Use to create Food reviews index
	//
	indexName := "reviews"
	datas := reviewDatas

	// Use to create Food keywords index
	//
	// indexName := "food"
	// datas := foodKeyword

	flag.IntVar(&count, "count", len(datas), "Number of documents to generate")
	flag.IntVar(&batch, "batch", len(datas)/15+1 , "Number of documents to send in one batch")
	flag.Parse()
	
	log.SetFlags(0)

	log.Println("\x1b[1mBulk\x1b[0m: documents " + humanize.Comma(int64(count)) + " batch size " + humanize.Comma(int64(batch)))
	log.Println("→ Sending batch ")

	// Re-create the index
	if res, err = elasticClient.Indices.Delete([]string{indexName}); err != nil {
		log.Fatalf("Cannot delete index: %s", err)
	}
	res, err = elasticClient.Indices.Create(indexName)
	if err != nil {
		log.Fatalf("Cannot create index: %s", err)
	}
	if res.IsError() {
		log.Fatalf("Cannot create index: %s", res)
	}

	start := time.Now().UTC()

	// Loop over the collection
	for i, a := range datas {
		numItems++

		currBatch = i / batch
		if i == count-1 {
			currBatch++
		}

		// Prepare the metadata payload
		meta := []byte(fmt.Sprintf(`{ "index" : { "_id" : "%d" } }%s`, a.ID, "\n"))

		// Prepare the data payload: encode to JSON
		data, err := json.Marshal(&a)
		if err != nil {
			log.Fatalf("Cannot encode %d: %s", a.ID, err)
		}

		// Append newline to the data payload
		data = append(data, "\n"...)

		// Append payloads to the buffer (ignoring write errors)
		buf.Grow(len(meta) + len(data))
		buf.Write(meta)
		buf.Write(data)

		// When a threshold is reached, execute the Bulk() request with body from buffer
		if i > 0 && i%batch == 0 || i == count-1 {
			log.Printf(strings.Repeat("=", currBatch))

			res, err = elasticClient.Bulk(bytes.NewReader(buf.Bytes()), elasticClient.Bulk.WithIndex(indexName))
			if err != nil {
				log.Fatalf("Failure indexing batch %d: %s", currBatch, err)
			}
			// If the whole request failed, print error and mark all documents as failed
			if res.IsError() {
				numErrors += numItems
				if err := json.NewDecoder(res.Body).Decode(&raw); err != nil {
					log.Fatalf("Failure to to parse response body: %s", err)
				} else {
					log.Printf("  Error: [%d] %s: %s",
						res.StatusCode,
						raw["error"].(map[string]interface{})["type"],
						raw["error"].(map[string]interface{})["reason"],
					)
				}
				// A successful response might still contain errors for particular documents...
				//
			} else {
				if err := json.NewDecoder(res.Body).Decode(&blk); err != nil {
					log.Fatalf("Failure to to parse response body: %s", err)
				} else {
					for _, d := range blk.Items {
						// ... so for any HTTP status above 201 ...
						if d.Index.Status > 201 {
							// ... increment the error counter ...
							numErrors++

							// ... and print the response status and error information ...
							log.Printf("  Error: [%d]: %s: %s: %s: %s",
								d.Index.Status,
								d.Index.Error.Type,
								d.Index.Error.Reason,
								d.Index.Error.Cause.Type,
								d.Index.Error.Cause.Reason,
							)
						} else {
							// ... otherwise increase the success counter.
							numIndexed++
						}
					}
				}
			}

			// Close the response body, to prevent reaching the limit for goroutines or file handles
			res.Body.Close()

			// Reset the buffer and items counter
			buf.Reset()
			numItems = 0
		}
	}

	dur := time.Since(start)

	if numErrors > 0 {
		log.Fatalf(
			"Indexed [%s] documents with [%s] errors in %s (%s docs/sec)",
			humanize.Comma(int64(numIndexed)),
			humanize.Comma(int64(numErrors)),
			dur.Truncate(time.Millisecond),
			humanize.Comma(int64(1000.0/float64(dur/time.Millisecond)*float64(numIndexed))),
		)
	} else {
		log.Printf(
			"Sucessfuly indexed [%s] documents in %s (%s docs/sec)",
			humanize.Comma(int64(numIndexed)),
			dur.Truncate(time.Millisecond),
			humanize.Comma(int64(1000.0/float64(dur/time.Millisecond)*float64(numIndexed))),
		)
	}
}

func searchByMatchKeyword(field string, keyword string) map[string]interface{} {
	query := `{"query":{"match":{"` + field + `":"` + keyword + `"}},"highlight": {"order":"score"}}`
	
	query = checkValidJson(query)

	// Build a new string from JSON query
	var b strings.Builder
	b.WriteString(query)
	
	// Instantiate a *strings.Reader object from string
	read := strings.NewReader(b.String())

	var mapResp map[string]interface{}
	var buf bytes.Buffer

	// Attempt to encode the JSON query and look for errors
	json.NewEncoder(&buf).Encode(read)
	// Pass the JSON query to the Golang client's Search() method
	res, err := elasticClient.Search(
		elasticClient.Search.WithIndex("reviews"),
		elasticClient.Search.WithBody(read),
		elasticClient.Search.WithPretty(),
	)

	// Check for any errors returned by API call to Elasticsearch
	if err != nil {
		log.Fatalf("Elasticsearch Search() API ERROR:", err)
	} else {
		// Close the result body when the function call is complete
		defer res.Body.Close()

		// Decode the JSON response and using a pointer
		json.NewDecoder(res.Body).Decode(&mapResp)
		// fmt.Println(`mapResp["_shards"] :`, mapResp["_shards"])
		// fmt.Println(`mapResp["hits"] :`, mapResp["hits"])
	}
	return mapResp
}

func checkValidJson(text string) string {
	// Check for JSON errors
	if json.Valid([]byte(text)) {
		return text
	} else {
		return "{}"
	}
}