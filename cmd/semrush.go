package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

// Compile time check to ensure SEMRush type implements RankingProvider
var _ RankingProvider = (*SEMRush)(nil)

func init() {
	// Add SEMRush to the list of available ranking providers
	addRankingProvider(&SEMRush{})
}

// SEMRush is a RankingProvider that pull data from semrush.com
type SEMRush struct {
}

// Provider returns the name of the RankingProvider
func (s *SEMRush) Provider() string {
	return "SEMRush"
}

// TopSitesGlobal returns the list of URLs
// of the most visited sites in the world
func (s *SEMRush) TopSitesGlobal() ([]string, error) {
	return s.topSites("https://www.semrush.com/website/top")
}

// TopSitesCountry returns a map of lists containing URLS
// of the most visited websites by country indexed by the country name
func (s *SEMRush) TopSitesCountry() (map[string][]string, error) {

	// First, we need to get info which countries have data

	res, err := http.Get("https://www.semrush.com/website")
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Create doc from HTML response
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	// a map of lists containing URLS
	// of the most visited websites by country
	// indexed by the country name
	countries := make(map[string][]string)

	// Error that occurred when scraping
	// and evaulating scraped data
	var scrapeError error

	// The html element that contain the list of countries
	// are dynamically created based on Javascript data.
	// The data is in window.__INITIAL_STATE__ which is JSON string.

	// Loop through the script tags in the page...
	doc.Find("script").EachWithBreak(func(i int, sel *goquery.Selection) bool {

		// ... whose text contents...
		t := sel.Text()

		// contains the window.__INITIAL_STATE__ assignment...
		if !strings.Contains(t, "window.__INITIAL_STATE__") {
			return true
		}

		// We only need the JSON string
		t = t[strings.Index(t, "'{")+1 : strings.Index(t, "}'")+1]

		// The JSON string they use contains unicode escape sequences
		// We need them to be in their corresponding one-character string
		unquoted, _ := strconv.Unquote("\"" + t + "\"")

		// Anonymous struct to unmarshal the JSON string into.
		// We just need the part containing the countries list
		// and their corresponding URLs
		data := struct {
			Countries []struct {
				Title string
				Link  string
			}
		}{}
		err = json.Unmarshal([]byte(unquoted), &data)

		if err != nil {
			scrapeError = fmt.Errorf("failed to unmarshal JSON: %w", err)
			return true
		}

		// Channel to store extracted top website URLs for a country
		info := make(chan struct {
			Country string
			URLs    []string
		})

		var wg sync.WaitGroup
		var m sync.Mutex

		for _, c := range data.Countries {

			if c.Title == "World" {
				continue
			}

			// Go routine that crawls and scrapes data
			wg.Add(1)
			go func(country, url string) {
				urls, err := s.topSites("https://semrush.com" + url)

				// If we had an error crawling and scraping data
				// return immediately and update the workgroup
				if err != nil {
					wg.Done()
					return
				}

				info <- struct {
					Country string
					URLs    []string
				}{Country: country, URLs: urls}

			}(c.Title, c.Link)

			// Go routine that receives and stores
			// data scraped by another routine
			go func() {
				defer wg.Done()
				m.Lock()
				defer m.Unlock()
				n := <-info
				countries[n.Country] = n.URLs
			}()

		}

		wg.Wait()
		return false
	})

	return countries, scrapeError
}

// topSites() crawls the given URL of a page from SEMRush
// Scrapes top websites info from the page and returns the list of URLs
func (s *SEMRush) topSites(url string) ([]string, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	urls := []string{}
	doc.Find("span[itemprop=name]").Each(func(i int, s *goquery.Selection) {
		txt := s.Text()
		class, _ := s.Attr("class")

		if strings.Contains(class, "TopWebsitesTable") {
			urls = append(urls, txt)
		}

	})

	return urls, nil

}
