package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/mux"
	"github.com/gorilla/feeds"
)

const (
	debug = true
	dURL  = "https://www.demonoid.pw/genlb.php?genid="
	leetxURL = "https://1337x.to"
)

func magnetToTorrent(magnetLink string) []byte {
	return []byte("d10:magnet-uri" + strconv.Itoa(len(magnetLink)) + ":" + magnetLink + "e")
}

func getPage(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return body
}

func getDoc(url string) *goquery.Document {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	g, err := goquery.NewDocumentFromReader(resp.Body)
	handleError(err)
	return g
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

func debugPrint(line string) {
	if debug {
		log.Println(line)
	}
}

func proxyRequest(rw http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	url := params["url"]
	debugPrint("Getting page " + url)
	page := string(getPage(url))
	toReplace := params["to_replace"]
	replacement := params["replacement"]
	if params["use_regex"] == "1" {
		page = strings.Replace(page, toReplace, replacement, -1)
	}
	debugPrint("Returning page")
	fmt.Fprintf(rw, page)

}

func extract(rw http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	debugPrint("Using extractor " + params["site"])
	debugPrint("URL id " + params["url_id"])
	site := params["site"]
	switch site {

	case "demonoid":
		debugPrint("Running " + params["site"] + " extractor")
		extractDemonoid(params["url_id"], rw)
	case "1337x":
		debugPrint("Running " + params["site"] + " extractor")
		extract1337x(params["url_id"], rw)
	default:
		debugPrint("No extractor for " + params["site"])

	}
}

func extractDemonoid(id string, rw http.ResponseWriter) {
	doc := getDoc(dURL + id)
	doc.Find(".ctable_content div a").Each(func(i int, s *goquery.Selection) {
		e, _ := s.Attr("href")
		debugPrint(e)
		if strings.Contains(e, "www.hypercache.pw") {
			toReturn := []byte(getPage(e))
			rw.Write(toReturn)

		}
	})
}

func extract1337x(id string, rw http.ResponseWriter) {
	debugPrint("Starting")
	doc := getDoc("https://1337x.to/torrent/" + id + "//")
	debugPrint("Finding torrent links")
	sel := doc.Find("a")
	debugPrint("got torrent links")
	for i := range sel.Nodes {
		debugPrint("Looping over torrent links")
		s := sel.Eq(i)
		e, _ := s.Attr("href")
		if !strings.Contains(e, "magnet:") && strings.Contains(e, "/torrent/") {
			debugPrint("Downloading page " + e)
			rw.Write([]byte(getPage(e)))
			return

		}
	}
		
}

func generateFeed(title string, link string, des string) *feeds.Feed {
	feed := &feeds.Feed{
        Title:       title,
        Link:        &feeds.Link{Href: link},
        Description: des,
        Author:      &feeds.Author{Name: "test"},
        Created:     time.Now(),
	}
	return feed
}

func generateItem(title string, link string) *feeds.Item {
	return &feeds.Item{
		Title:       title,
		Link:        &feeds.Link{Href: link},
		Description: "",
		Author:      &feeds.Author{Name: "", Email: ""},
		Created:     time.Now(),
	}
}

func scrape1337x(feed *feeds.Feed, urlToScrape string) string {
	doc := getDoc(urlToScrape)
	doc.Find("tbody > tr > td > a").Each(func(i int, s *goquery.Selection) {
		e, _ := s.Attr("href")
		debugPrint(e)
		if strings.Contains(e, "/torrent/") {
			feed.Add(generateItem(s.Text(), "http://127.0.0.1:5000/extractor/1337x/" + strings.Split(e, "/")[2]))
		}

	})
	e, err := feed.ToRss()
	handleError(err)
	return e

}

func siteToRss(rw http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	site := params["site"]
	url := params["url"]
	debugPrint("Using siteToRss for " + site)
	debugPrint("Turning " + url + " into a rss feed")
	feed := generateFeed(url + " rss feed", "127.0.0.1/rss/" + site + "/" + url, "A rss feed generated from " + url)


	switch site {
	case "1337x":
		fmt.Fprintf(rw, scrape1337x(feed, url))
	default:
		rw.WriteHeader(400)
		fmt.Fprintf(rw, "No rss generator for that site")
	}
	

}

func main() {
	log.Println("Setting up routes")
	router := mux.NewRouter()
	router.HandleFunc("/", proxyRequest).Queries("url", "{url}", "use_regex", "{use_regex}", "to_replace", "{to_replace}", "replacement", "{replacement}").Methods("GET")
	router.HandleFunc("/extractor/{site}/{url_id}", extract).Methods("GET")
	router.HandleFunc("/rss/{site}", siteToRss).Queries("url", "{url}").Methods("GET")
	log.Println("Routes set")
	log.Fatal(http.ListenAndServe(":5000", router))
}
