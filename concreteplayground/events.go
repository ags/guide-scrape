package main

import (
	"encoding/csv"
	"fmt"
	"html"
	"os"
)

func searchEvents(client *Client, csv *csv.Writer, region string) error {
	header := []string{
		"type",
		"name",
		"description",
		"start date",
		"end date",
		"location",
		"address",
		"latitude",
		"longitude",
		"url",
		"permalink",
	}
	_ = csv.Write(header)

	for page := 1; ; page++ {
		fmt.Fprintf(os.Stderr, "region=%s type=event page=%d\n", region, page)
		res, err := client.SearchEvents(region, page)
		if err != nil {
			return err
		}
		for _, r := range res.Results {
			row := make([]string, len(header))
			row[0] = r.StructuredData.Type
			row[1] = html.UnescapeString(r.PostTitle)
			row[2] = html.UnescapeString(r.PostExcerpt)
			if r.StructuredData.Type != "" {
				row[3] = r.StructuredData.StartDate
				row[4] = r.StructuredData.EndDate
				row[5] = r.StructuredData.Location.Name
				row[6] = r.StructuredData.Location.Address.Name
			}
			if r.Location.Latitude != 0 && r.Location.Longitude != 0 {
				row[7] = fmt.Sprint(r.Location.Latitude)
				row[8] = fmt.Sprint(r.Location.Longitude)
			}
			if r.StructuredData.Type != "" {
				row[9] = r.StructuredData.URL
			}
			row[10] = r.Permalink

			if err := csv.Write(row); err != nil {
				return err
			}
		}
		if len(res.Results) < 20 {
			break
		}
	}

	return nil
}
