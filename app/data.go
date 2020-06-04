package main

import (
	"flag"
	"encoding/csv"
	"log"
	"os"
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