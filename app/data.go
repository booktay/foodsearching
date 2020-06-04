package main

import (
	"flag"
	"encoding/csv"
	"log"
	"os"
	"fmt"
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

    // Open the CSV file
    csvfile, err := os.Open(*reviewDirectory)
    if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
        return []FoodReview{}, err
    }
    defer csvfile.Close()

    // Read all lines in file
	file := csv.NewReader(csvfile)
	file.Comma = ';'
	datas, err := file.ReadAll()
    if err != nil {
		log.Fatal(err)
        return []FoodReview{}, err
	}

	fmt.Println(datas)
	
	reviewsData := []FoodReview{}
	return reviewsData, nil
}