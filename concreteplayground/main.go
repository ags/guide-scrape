package main

import (
	"encoding/csv"
	"fmt"
	"html"
	"os"
	"strings"
)

func main() {
	var (
		region    = os.Args[1]
		placeType = os.Args[2]
	)

	client := NewClient()
	csv := csv.NewWriter(os.Stdout)
	page := 1

	_ = csv.Write([]string{
		"type",
		"name",
		"description",
		"latitude",
		"longitude",
		"locality",
		"address",
		"price range",
		"serves cuisine",
		"permalink",
		"url",
	})

	for {
		fmt.Fprintf(os.Stderr, "request region=%s type=%s page=%d\n", region, placeType, page)
		res, err := client.Search(region, placeType, page)
		if err != nil {
			fmt.Fprintf(os.Stderr, "request region=%s type=%s page=%d err=%v\n", region, placeType, page, err)
			os.Exit(1)
		}
		for _, r := range res.Results {
			err := csv.Write([]string{
				r.StructuredData.Type,
				html.UnescapeString(r.StructuredData.Name),
				html.UnescapeString(r.StructuredData.Description),
				fmt.Sprint(r.StructuredData.Geo.Latitude),
				fmt.Sprint(r.StructuredData.Geo.Longitude),
				r.StructuredData.Address.Locality,
				r.StructuredData.Address.StreetAddress,
				r.StructuredData.PriceRange,
				strings.Join(r.StructuredData.ServesCuisine, ","),
				r.Permalink,
				r.StructuredData.URL,
			})
			if err != nil {
				fmt.Fprintf(os.Stderr, "request region=%s type=%s page=%d err=%v\n", region, placeType, page, err)
				os.Exit(1)
			}
		}
		if len(res.Results) < 20 {
			break
		}
		page++
	}
}
