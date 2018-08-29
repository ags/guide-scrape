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

func (c *Client) SearchEvents(region string, pageNumber int) (EventResponse, error) {
	res, err := c.search(region, "event", pageNumber)
	if err != nil {
		return EventResponse{}, nil
	}
	defer res.Body.Close()

	var r EventResponse
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return EventResponse{}, err
	}
	return r, nil
}

func (c *Client) SearchPlaces(region, placeType string, pageNumber int) (PlaceResponse, error) {
	res, err := c.search(region, placeType, pageNumber)
	if err != nil {
		return PlaceResponse{}, nil
	}
	defer res.Body.Close()

	var r PlaceResponse
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return PlaceResponse{}, err
	}
	return r, nil
}

func (c Client) search(region, placeType string, pageNumber int) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	postType := "places"
	if placeType == "event" {
		postType = "tribe_events"
	}
	q.Set("post_type", postType)
	q.Set("place_type", placeType)
	q.Set("region", region)
	q.Set("sort", "all")
	q.Set("paged", strconv.Itoa(pageNumber))
	q.Set("action", "directory_search")
	req.URL.RawQuery = q.Encode()

	res, err := c.http.Do(req)
	if err != nil {
		return res, err
	}
	if res.StatusCode != 200 {
		return res, fmt.Errorf("expected 200, got %d", res.StatusCode)
	}
	return res, nil
}

type PlaceResponse struct {
	FoundPosts int           `json:"found_posts"`
	Results    []PlaceResult `json:"results"`
}

type PlaceResult struct {
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

type EventResponse struct {
	FoundPosts int           `json:"found_posts"`
	Results    []EventResult `json:"results"`
}

type EventResult struct {
	PostTitle   string `json:"post_title"`
	PostExcerpt string `json:"post_excerpt"`
	Permalink   string `json:"permalink"`
	Location    struct {
		Latitude  float64 `json:"lat"`
		Longitude float64 `json:"lon"`
	} `json:"location"`
	StructuredData EventStructuredData `json:"structured_data"`
}

type EventStructuredData struct {
	Type        string `json:"@type"`
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
	StartDate   string `json:"startDate"`
	EndDate     string `json:"endDate"`
	Location    struct {
		Name    string `json:"name"`
		Address struct {
			Name string `json:"name"`
		} `json:"address"`
	} `json:"location"`
}

func (d *EventStructuredData) UnmarshalJSON(b []byte) error {
	if string(b) == "false" {
		return nil
	}
	var tmp struct {
		Type        string `json:"@type"`
		Name        string `json:"name"`
		Description string `json:"description"`
		URL         string `json:"url"`
		StartDate   string `json:"startDate"`
		EndDate     string `json:"endDate"`
		Location    struct {
			Name    string `json:"name"`
			Address struct {
				Name string `json:"name"`
			} `json:"address"`
		} `json:"location"`
	}
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}
	*d = tmp
	return nil
}
