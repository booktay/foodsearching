package main

import (
	"flag"
)

var reviewDirectory = flag.String("reviewDirectory", "data/test_file.csv", "Review Directory")
var foodKeywordDirectory = flag.String("foodKeywordDirectory", "data/food_20k_keyword.txt", "Food Keyword Directory")

type FoodReview struct {
	ID int `json:"reviewid"`
	ReviewText string `json:"reviewtext"`
	CreatedTime int64 `json:"created"`
	ModifiedTime int64 `json:"modified"`
}

type FoodDictionary struct {
	ID int `json:"keywordid"`
	Keyword string `json:"keyword"`
}

