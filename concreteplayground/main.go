package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

func main() {
	var (
		region    = os.Args[1]
		placeType = os.Args[2]
	)

	client := NewClient()
	csv := csv.NewWriter(os.Stdout)

	search := func() error {
		if placeType == "event" {
			return searchEvents(client, csv, region)
		}
		return searchPlaces(client, csv, region, placeType)
	}

	if err := search(); err != nil {
		fmt.Fprintf(os.Stderr, "region=%s type=%s err=%v\n", region, placeType, err)
		os.Exit(1)
	}

	csv.Flush()
}
