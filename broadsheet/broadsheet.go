package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Query struct {
	City       string
	Suburb     string
	Categories []string
	PageNumber int
}

type Client struct {
	http *http.Client
}

func NewClient() *Client {
	return &Client{
		http: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

type Response struct {
	Count   int      `json:"count"`
	Results []Result `json:"results"`
}

type Result struct {
	URL            string   `json:"url"`
	Suburb         string   `json:"suburb"`
	Title          string   `json:"title"`
	Category       string   `json:"category"`
	Description    string   `json:"description"`
	Features       []string `json:"features"`
	PrimaryAddress struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"primary_address"`
}

type searchOption struct {
	Categories []string `json:"categories"`
}

const (
	apiURL = "https://www.broadsheet.com.au/api"
)

func (c *Client) Search(q Query) (Response, error) {
	opt, err := json.Marshal(searchOption{
		Categories: q.Categories,
	})
	if err != nil {
		return Response{}, err
	}

	// TODO URL encoding suburb?
	url := fmt.Sprintf(
		"%s/%s/search/%s?page=%d&o=%s",
		apiURL, q.City, q.Suburb, q.PageNumber, opt,
	)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return Response{}, err
	}

	req.Header.Add("Accept", "application/json")

	res, err := c.http.Get(url)
	if err != nil {
		return Response{}, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return Response{}, fmt.Errorf("expected 200, got %d", res.StatusCode)
	}

	var r Response
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return Response{}, err
	}
	return r, nil
}
