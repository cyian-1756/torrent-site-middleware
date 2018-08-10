package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/mux"
)

const (
	debug = true
	dURL  = "https://www.demonoid.pw/genlb.php?genid="
)

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

func main() {
	log.Println("Setting up routes")
	router := mux.NewRouter()
	router.HandleFunc("/", proxyRequest).Queries("url", "{url}", "use_regex", "{use_regex}", "to_replace", "{to_replace}", "replacement", "{replacement}").Methods("GET")
	router.HandleFunc("/extractor/{site}/{url_id}", extract).Methods("GET")
	log.Println("Routes set")
	log.Fatal(http.ListenAndServe(":5000", router))
}
