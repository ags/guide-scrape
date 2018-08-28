package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

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

const apiURL = "https://concreteplayground.com/ajax.php"

func (c *Client) Search(
	region string,
	placeType string,
	pageNumber int,
) (Response, error) {
	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return Response{}, nil
	}
	q := req.URL.Query()
	q.Set("post_type", "places")
	q.Set("place_type", placeType)
	q.Set("region", region)
	q.Set("sort", "all")
	q.Set("paged", strconv.Itoa(pageNumber))
	q.Set("action", "directory_search")
	req.URL.RawQuery = q.Encode()

	res, err := c.http.Do(req)
	if err != nil {
		return Response{}, nil
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

type Response struct {
	FoundPosts int      `json:"found_posts"`
	Results    []Result `json:"results"`
}

type Result struct {
	Permalink      string `json:"permalink"`
	StructuredData struct {
		Type    string `json:"@type"`
		Address struct {
			Locality      string `json:"addressLocality"`
			StreetAddress string `json:"streetaddress"`
		} `json:"address"`
		Description string `json:"description"`
		Geo         struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		} `json:"geo"`
		Name          string   `json:"name"`
		PriceRange    string   `json:"priceRange"`
		ServesCuisine []string `json:"servesCuisine"`
		URL           string   `json:"url"`
	} `json:"structured_data"`
}
