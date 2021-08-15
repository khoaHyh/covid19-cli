package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"unicode"
)

func main() {
	timeSeriesPtr := flag.String("timeseries", "", "To indicate that the user requires time series data")
	stat := flag.String("stat", "", "return data only of the specific type")
	loc := flag.String("loc", "", "return data only from the specified province or health region")
	date := flag.String("date", "", "return date only from the specified date")
	before := flag.String("before", "", "return date on or after the specified date")
	after := flag.String("after", "", "return date on or before the specified date")

	flag.Parse()

	if *timeSeriesPtr != "" {
		getTimeSeries(*stat, *loc, *date, *before, *after)
		return
	}

	getMostRecentAvailableDayCanada()
}

func prettyPrintJsonFromApi(url string) {
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Fatal(err)
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Fatal(err)
	}

	var prettyJSON bytes.Buffer
	error := json.Indent(&prettyJSON, body, "", "\t")

	if error != nil {
		log.Fatal(error)
	}

	log.Println(prettyJSON.String())
}

func getMostRecentAvailableDayCanada() {
	url := "https://api.opencovid.ca"

	prettyPrintJsonFromApi(url)
}

func getTimeSeries(stat string, loc string, date string, before string, after string) {
	if stat == "" && loc == "" && date == "" && before == "" && after == "" {
		log.Fatal("Need to provide flags for tempseries query")
	}

	// Check if location flag value is all uppercase
	for _, i := range loc {
		if unicode.IsLower(i) && unicode.IsLetter(i) {
			log.Fatal("Location flag value needs to be in all uppercase.")
		}
	}

	url := "https://api.opencovid.ca/timeseries?"

	params := 0

	if stat != "" {
		if params > 0 {
			url += "&"
		}
		url += "stat=" + stat
		params++
	}
	if loc != "" {
		if params > 0 {
			url += "&"
		}
		url += "loc=" + loc
		params++
	}
	if date != "" {
		if params > 0 {
			url += "&"
		}
		url += "date=" + date
		params++
	}
	if before != "" {
		if params > 0 {
			url += "&"
		}
		url += "before=" + before
		params++
	}
	if after != "" {
		if params > 0 {
			url += "&"
		}
		url += "after=" + after
		params++
	}

	prettyPrintJsonFromApi(url)
}
