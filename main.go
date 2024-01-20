package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
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

type SearchResult struct {
	Page         int32   `json:"page"`
	PerPage      int32   `json:"per_page"`
	TotalResults int32   `json:"total_result"`
	NextPage     string  `json:"next_page"`
	Photos       []Photo `json:"photos"`
}

type Photo struct {
	Id              int32       `json:"id"`
	Width           int32       `json:"width"`
	Height          int32       `json:"height"`
	Url             string      `json:"url"`
	Photographer    string      `json:"photographer"`
	PhotographerUrl string      `json:"photographer_url"`
	Src             PhotoSource `json:"src"`
}

type PhotoSource struct {
	Original  string `json:"original"`
	Large     string `json:"large"`
	Large2x   string `json:"large2x"`
	Medium    string `json:"medium"`
	Small     string `json:"small"`
	Potrait   string `json:"potrait"`
	Square    string `json:"square"`
	Landscape string `json:"landscape"`
	Tiny      string `json:"tiny"`
}

type curatedResult struct {
	Page     int32   `json:"page"`
	PerPage  int32   `json:"per_page"`
	NextPage string  `json:"next_page"`
	Photos   []Photo `json:"photos"`
}

func (c *Client) SearchPhotos(query string, perPage, page int) (*SearchResult, error) {
	url := fmt.Sprint(PHOTO_API+"/search?query=%s&per_page=%d&page=%d", query, perPage, page)

	res, err := c.requestDoWithAuth("GET", url)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var result SearchResult

	err = json.Unmarshal(data, &result)

	return &result, err
}

func (c *Client) curatedPhotos(perPage, page int) (*curatedResult, error) {
	url := fmt.Sprintf(PHOTO_API+"/curated?per_page=%d&page=%d", perPage, page)

	res, err := c.requestDoWithAuth("GET", url)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var result curatedResult

	err = json.Unmarshal(data, &result)

	return &result, err
}

func (c *Client) requestDoWithAuth(method, url string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", c.Token)

	res, err := c.hc.Do(req)
	if err != nil {
		return res, err
	}

	times, err := strconv.Atoi(res.Header.Get("X-Ratelimit-Remaining"))
	if err != nil {
		return res, nil
	} else {
		c.RemainingTimes = int32(times)
	}

	return res, nil
}

func (c *Client) GetPhotoById(id int32) (*Photo, error) {
	url := fmt.Sprintf(PHOTO_API+"/photos/%d", id)

	res, err := c.requestDoWithAuth("GET", url)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var result Photo

	err = json.Unmarshal(data, &result)

	return &result, err
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env: %v", err)
	}

	var TOKEN = os.Getenv("PEXELS_API_KEY")

	var c = NewClient(TOKEN)

	result, err := c.SearchPhotos("waves", 15, 1)
	if err != nil {
		fmt.Errorf("Search Error: %v", err)
	}

	if result.Page == 0 {
		fmt.Errorf("Search Result Wrong")
	}

	fmt.Println(result)
}
