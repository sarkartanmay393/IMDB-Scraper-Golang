package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

// Star holds a celebrety's name and other imdb information.
type Star struct {
	Name      string  `json:"name"`
	Photo     string  `json:"photo"`
	JobTitle  string  `json:"job_title"`
	Birthdate string  `json:"birthdate"`
	Bio       string  `json:"bio"`
	TopMovies []Movie `json:"top_movies"`
}

// Movie holds basic information about a movie.
type Movie struct {
	Title   string `json:"title"`
	Release string `json:"release"`
}

func main() {

	// Printing prompt for user to input values.
	var input string
	// User should type just like this. e.g.. 12-23
	fmt.Println("Type a date and month in (MM-DD) format:	")
	_, err := fmt.Scanf("%s", &input)
	if err != nil {
		log.Fatalln("Unable to read input because	", err)
	}

	// Splitting the input into month and day.
	sliceOfInput := strings.Split(input, "-")

	// Converting the month into integer.
	month, err := strconv.Atoi(sliceOfInput[0])
	if err != nil {
		log.Fatalln("Unable to convert month into int because	", err)
	}
	// Converting the month into integer.
	day, err := strconv.Atoi(sliceOfInput[1])
	if err != nil {
		log.Fatalln("Unable to convert date into int because	", err)
	}

	// Crawling the website with the given month and day and stores the data into a slice of Star.
	stars := Crawl(month, day)

	// To print all the scraped data in terminal.
	// encoder := json.NewEncoder(os.Stderr)
	// encoder.SetIndent("", "   ")
	// err := encoder.Encode(stars)

	// Converting Star struct into JSON.
	data, err := json.MarshalIndent(stars, "", "	")
	if err != nil {
		log.Println("Error while marshalling JSON: ", err)
	}

	// Writing JSON into a file called 'mm-dd-output.json'.
	fileName := fmt.Sprintf("%s-output.json", input)
	err = os.WriteFile(fileName, data, 0644)
	if err != nil {
		log.Println("Error while writing JSON: ", err)
	}

}

// Crawl crawls the imdb website for a given page and returns a list of stars.
func Crawl(month int, day int) []Star {

	// stars will hold all the scraped data.
	var stars []Star

	// Main collector initialized.
	c := colly.NewCollector(
		// Colly collector will only visit the given domains.
		colly.AllowedDomains("www.imdb.com", "imdb.com"),
		// Colly collector stores cache and uses from here.
		colly.CacheDir("./.imdb_cache"),
	)

	// Info Collector is made using cloing.
	ic := c.Clone()

	// Where crawler sees 'mode-detail' class attribute, it will call callback function.
	c.OnHTML(".mode-detail", func(e *colly.HTMLElement) {
		// Getting the profile url from the href attribute, it is goquery selector string.
		// 'div.lister-item-image > a' is path of profile url in HTML.
		// At fist there is a Div tag with class 'lister-item-image' and
		// inside that there is a tag 'a' where href attribute is profile url.
		profileUrl := e.ChildAttr("div.lister-item-image > a", "href")
		// But the profile url is relative, so we need to add the base url.
		// profileUrl is "/name/nm0000123/" corrently.
		profileUrl = e.Request.AbsoluteURL(profileUrl)
		// Now profileUrl is "https://www.imdb.com/name/nm0000123/".

		// Asking info collector to visit the profile url.
		ic.Visit(profileUrl)
	})

	// This crawler function gets into next page, if there is one.
	c.OnHTML("a.lister-page-next", func(e *colly.HTMLElement) {
		// Getting the next page url from the href attribute.
		nextPageUrl := e.Attr("href")
		nextPageUrl = e.Request.AbsoluteURL(nextPageUrl)

		// Asking main collector to visit the next page url.
		c.Visit(nextPageUrl)
	})

	// This info collector crawler function gets the information of the celebrety inside the 'profileURL' page.
	ic.OnHTML("#content-2-wide", func(e *colly.HTMLElement) {

		// Getting all details using info collector and storing them in a Star struct.
		temporaryStar := Star{
			Name:      e.ChildText("h1.header > span.itemprop"),
			Photo:     e.ChildAttr("#name-poster", "src"),
			JobTitle:  e.ChildText("#name-job-categories > a > span.itemprop"),
			Birthdate: e.ChildAttr("#name-born-info > time", "datetime"),
			Bio:       strings.TrimSpace(e.ChildText("#name-bio-text > div.name-trivia-bio-text > div.inline")),
			TopMovies: []Movie{},
		}

		// Now iterating over all the top movies of the profile url page.
		e.ForEach("div.knownfor-title", func(_ int, el *colly.HTMLElement) {
			temporaryStar.TopMovies = append(temporaryStar.TopMovies, Movie{
				Title:   el.ChildText("div.knownfor-title-role > a.knownfor-ellipsis"),
				Release: el.ChildText("div.knownfor-title > div.knownfor-year > span.knownfor-ellipsis"),
			})
		})

		// Now appending the temporaryStar to the stars slice.
		stars = append(stars, temporaryStar)
	})

	// Printing text of every request made on info collector.
	ic.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting:	", r.URL.String())
	})

	// Our link to crawl.
	startUrl := fmt.Sprintf("https://www.imdb.com/search/name/?birth_monthday=%d-%d", month, day)
	fmt.Println("Starting crawling into	", startUrl)

	// Starting the main collector.
	c.Visit(startUrl)

	// Returning all the scraped data.
	return stars
}
