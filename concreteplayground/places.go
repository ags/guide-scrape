package main

import (
	"encoding/csv"
	"fmt"
	"html"
	"os"
	"strings"
)

func searchPlaces(client *Client, csv *csv.Writer, region, placeType string) error {
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

	for page := 1; ; page++ {
		fmt.Fprintf(os.Stderr, "region=%s type=%s page=%d\n", region, placeType, page)
		res, err := client.SearchPlaces(region, placeType, page)
		if err != nil {
			return err
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
				return err
			}
		}
		if len(res.Results) < 20 {
			break
		}
	}

	return nil
}
