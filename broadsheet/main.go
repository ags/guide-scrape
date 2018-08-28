package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

func main() {
	scraper := &scraper{}
	var (
		city     = "sydney"
		category = "restaurant"
		suburbs  = []string{}
	)
	for _, s := range sydneySuburbs {
		suburbs = append(suburbs, strings.ToLower(s))
	}
	//suburbs = []string{"prahran", "south yarra", "collingwood"}

	results, err := scraper.scrape(city, category, suburbs)
	if err != nil {
		log.Fatal(err)
	}

	csv := csv.NewWriter(os.Stdout)

	if err := csv.Write([]string{
		"city",
		"suburb",
		"category",
		"name",
		"url",
		"description",
		"latitude",
		"longitude",
		"features",
	}); err != nil {
		log.Fatal(err)
	}

	for _, r := range results {
		err := csv.Write([]string{
			city,
			r.Suburb,
			r.Category,
			r.Title,
			r.URL,
			strings.TrimRight(r.Description, "\n"),
			fmt.Sprint(r.PrimaryAddress.Latitude),
			fmt.Sprint(r.PrimaryAddress.Longitude),
			strings.Join(r.Features, ","),
		})
		if err != nil {
			log.Fatal(err)
		}
	}
	csv.Flush()
}

type scraper struct{}

func (s *scraper) scrape(city, category string, suburbs []string) ([]Result, error) {
	sch := make(chan string, 10)
	outch := make(chan output, 10)
	var wg sync.WaitGroup

	for w := 1; w <= 10; w++ {
		wg.Add(1)

		go func(id int) {
			client := NewClient()
			for s := range sch {
				fmt.Fprintf(os.Stderr, "id=%d suburb=%s state=%s\n", id, s, "start")
				r, err := client.Search(Query{
					City:       city,
					Suburb:     s,
					Categories: []string{category},
					PageNumber: 1,
				})
				outch <- output{suburb: s, response: r, err: err}
				fmt.Fprintf(os.Stderr, "id=%d suburb=%s state=%s err=%v\n", id, s, "done", err)
			}
			wg.Done()
		}(w)
	}

	go func() {
		for _, s := range suburbs {
			sch <- s
		}
		close(sch)

		wg.Wait()
		close(outch)
	}()

	results := []Result{}
	for m := range outch {
		if m.err != nil {
			return []Result{}, m.err
		}
		for _, r := range m.response.Results {
			results = append(results, r)
		}
	}

	wg.Wait()

	return results, nil
}

type output struct {
	suburb   string
	response Response
	err      error
}
