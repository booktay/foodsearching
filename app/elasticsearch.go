package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"
	"strings"
	"github.com/dustin/go-humanize"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

var ES01IP = flag.String("ES01IP", "http://172.21.0.2:9200", "ES01 IP Address")

func startElasticsearchConnection() {
	log.Print("Starting the Elasticsearch Client")
	var (
		r map[string]interface{}
	)

	cfg := elasticsearch.Config{
		Addresses: []string {
			*ES01IP,
		},
	}

	elasticClient, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	for {
		res, err := elasticClient.Info()
		if err == nil {
			json.NewDecoder(res.Body).Decode(&r)
			log.Println("Connected to Elasticsearch :", r["name"])
			log.Println("IP Address :", *ES01IP)
			res.Body.Close()
			break
		} else {
			log.Println("Waiting for Elasticsearch connection...")
			time.Sleep(5 * time.Second)
		}
	}

	insertBulkDocument(elasticClient)
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

func insertBulkDocument(es *elasticsearch.Client) {

	var (
		_ = fmt.Print
		count int
		batch int
	)

	rand.Seed(time.Now().UnixNano())
	
	type bulkResponse struct {
		Errors bool `json:"errors"`
		Items  []struct {
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

	indexName := "reviews"
	datas := reviewDatas
	// indexName := "food"
	// datas := foodKeyword

	flag.IntVar(&count, "count", len(datas), "Number of documents to generate")
	flag.IntVar(&batch, "batch", len(datas)/15+1 , "Number of documents to send in one batch")
	flag.Parse()
	
	log.SetFlags(0)

	log.Println("\x1b[1mBulk\x1b[0m: documents " + humanize.Comma(int64(count)) + " batch size " + humanize.Comma(int64(batch)))
	log.Println("→ Sending batch ")

	// Re-create the index
	if res, err = es.Indices.Delete([]string{indexName}); err != nil {
		log.Fatalf("Cannot delete index: %s", err)
	}
	res, err = es.Indices.Create(indexName)
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
		//
		meta := []byte(fmt.Sprintf(`{ "index" : { "_id" : "%d" } }%s`, a.ID, "\n"))
		// Prepare the data payload: encode to JSON
		//
		data, err := json.Marshal(&a)
		if err != nil {
			log.Fatalf("Cannot encode %d: %s", a.ID, err)
		}

		// Append newline to the data payload
		data = append(data, "\n"...) // <-- Comment out to trigger failure for batch

		// Append payloads to the buffer (ignoring write errors)
		buf.Grow(len(meta) + len(data))
		buf.Write(meta)
		buf.Write(data)

		// When a threshold is reached, execute the Bulk() request with body from buffer
		//
		if i > 0 && i%batch == 0 || i == count-1 {
			log.Printf(strings.Repeat("=", currBatch))

			res, err = es.Bulk(bytes.NewReader(buf.Bytes()), es.Bulk.WithIndex(indexName))
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