package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"./extractors"
	"./scrapers"
	"./shared"
	"github.com/gorilla/mux"
)

func magnetToTorrent(magnetLink string) []byte {
	return []byte("d10:magnet-uri" + strconv.Itoa(len(magnetLink)) + ":" + magnetLink + "e")
}

// Proxy a request and maybe use regex
func proxyRequest(rw http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	url := params["url"]
	shared.DebugPrint("Getting page " + url)
	page := string(shared.GetPage(url))
	toReplace := params["to_replace"]
	replacement := params["replacement"]
	if params["use_regex"] == "1" {
		page = strings.Replace(page, toReplace, replacement, -1)
	}
	shared.DebugPrint("Returning page")
	fmt.Fprintf(rw, page)

}

func siteToRss(rw http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	site := params["site"]
	url := params["url"]
	shared.DebugPrint("Using siteToRss for " + site)
	shared.DebugPrint("Turning " + url + " into a rss feed")
	feed := shared.GenerateFeed(url+" rss feed", "127.0.0.1/rss/"+site+"/"+url, "A rss feed generated from "+url)

	switch site {
	case "1337x":
		fmt.Fprintf(rw, scrapers.Scrape1337x(feed, url))
	default:
		// TODO find out if this should really throw a 400
		rw.WriteHeader(400)
		fmt.Fprintf(rw, "No rss generator for that site")
	}

}

func main() {
	log.Println("Setting up routes")
	router := mux.NewRouter()
	router.HandleFunc("/", proxyRequest).Queries("url", "{url}", "use_regex", "{use_regex}", "to_replace", "{to_replace}", "replacement", "{replacement}").Methods("GET")
	router.HandleFunc("/extractor/{site}/{url_id}", extractors.Extract).Methods("GET")
	router.HandleFunc("/rss/{site}", siteToRss).Queries("url", "{url}").Methods("GET")
	log.Println("Routes set")
	log.Fatal(http.ListenAndServe(":5000", router))
}
