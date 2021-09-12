package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"unicode"
)

func main() {
	timeSeriesPtr := flag.Bool("timeseries", false, "To indicate that the user requires time series data")
	statPtr := flag.String("stat", "", "return data only of the specific type")
	locPtr := flag.String("loc", "", "return data only from the specified province or health region")
	datePtr := flag.String("date", "", "return date only from the specified date")
	beforePtr := flag.String("before", "", "return date on or after the specified date")
	afterPtr := flag.String("after", "", "return date on or before the specified date")

	flag.Parse()

	if *timeSeriesPtr {
		getTimeSeries(*statPtr, *locPtr, *datePtr, *beforePtr, *afterPtr)
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

	/** FLAG VALUE VALIDATION **/
	// Check if location flag value is all uppercase and is apart of the alphabet
	for _, i := range loc {
		if unicode.IsLower(i) {
			log.Fatal("'loc' flag value needs to be in all uppercase.")
		} else if !unicode.IsLetter(i) {
			log.Fatal("invalid value for 'loc' flag")
		}
	}

	for _, i := range stat {
		if !unicode.IsLetter(i) {
			log.Fatal("invalid value for 'stat' flag")
		}
	}

	re := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

	if date != "" {
		dateResult := re.MatchString(date)

		if !dateResult {
			log.Fatal("invalid value for 'date' flag")
		}
	}

	if before != "" {
		beforeResult := re.MatchString(before)

		if !beforeResult {
			log.Fatal("invalid value for 'before' flag")
		}
	}

	if after != "" {
		afterResult := re.MatchString(after)

		if !afterResult {
			log.Fatal("invalid value for 'after' flag")
		}
	}
	/** FLAG VALUE VALIDATION END **/

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
