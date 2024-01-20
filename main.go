package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	PHOTO_API = "https://api.pexels/v1"
	VIDEO_API = "https://api.pexels/videos"
)

type Client struct {
	Token          string
	hc             http.Client
	RemainingTimes int32
}

func NewClient(token string) *Client {
	c := http.Client{}

	return &Client{Token: token, hc: c}
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env: %v", err)
	}

	var TOKEN = os.Getenv("PEXELS_API_KEY")

	var c = NewClient(TOKEN)

	result, err := c.SearchPhotos("waves")
	if err != nil {
		fmt.Errorf("Search Error: %v", err)
	}

	if result.Page == 0 {
		fmt.Errorf("Search Result Wrong")
	}

	fmt.Println(result)
}
