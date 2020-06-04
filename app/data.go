package main

import (
	"flag"
	"encoding/csv"
	"log"
	"os"
	"time"
	"strconv"
)

var foodReviewsDir = flag.String("foodReviewsDir", "data/test_file.csv", "Food Reviews Directory")
var foodKeywordDir = flag.String("foodKeywordDir", "data/food_20k_keyword.txt", "Food Keyword Directory")

type FoodReview struct {
	ID int `json:"reviewid"`
	ReviewText string `json:"reviewtext"`
	CreatedTime int64 `json:"created"`
	ModifiedTime int64 `json:"modified"`
}

type FoodKeyword struct {
	ID int `json:"keywordid"`
	Keyword string `json:"keyword"`
}

func getReviewData() ([]FoodReview, error) {
	log.Print("Reading Review Data...")

    // Open the CSV file from Food Reviews Directory
    csvfile, err := os.Open(*foodReviewsDir)
    if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
        return []FoodReview{}, err
    }
    defer csvfile.Close()

    // Read all lines in CSV file
	file := csv.NewReader(csvfile)
	file.Comma = ';' // Declared seperate character
	datas, err := file.ReadAll()
    if err != nil {
		log.Fatal(err)
        return []FoodReview{}, err
	}

	// Get Present Time in Unix format
	timeNow := time.Now().UnixNano()

	// Transform CSV format to FoodReview struct
	reviewsData := []FoodReview{}
	for _, line := range datas[1:] {
		id, _ := strconv.Atoi(line[0])
        reviewsData = append(reviewsData, FoodReview{
			ID: id,
			ReviewText: line[1],
			CreatedTime: timeNow,
			ModifiedTime: timeNow,
		})
	}

	log.Print("Reading Completed")
	return reviewsData, nil
}