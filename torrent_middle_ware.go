package main

import (
    "github.com/gorilla/mux"
    "log"
	"net/http"
	"io/ioutil"
	"strings"
	"fmt"
	"github.com/PuerkitoBio/goquery"
)

const (
	debug = true
	dURL = "https://www.demonoid.pw/genlb.php?genid="
)

func getPageAsString(url string) string {
	if debug {
		log.Println("URL: " + url)
	}
	resp, err := http.Get(url)
	if err != nil {
    	panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
    	panic(err)
	}
	return string(body)
}

func getPageAsBytes(url string) []byte {
	if debug {
		log.Println("URL: " + url)
	}
	resp, err := http.Get(url)
	if err != nil {
    	panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
    	panic(err)
	}
	return []byte(body)
}

func getPage(url string) *goquery.Document {
	if debug {
		log.Println("URL: " + url)
	}
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
	page := getPageAsString(params["url"])
	toReplace := params["to_replace"]
	replacement := params["replacement"]
	if params["use_regex"] == "1" {
		debugPrint("Using regex")
		page = strings.Replace(page, toReplace, replacement, -1)
	} else {
		debugPrint("Not using regex")
	}
	fmt.Fprintf(rw, page)
    

}

func extract(rw http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	debugPrint("Using extractor " + params["site"])
	debugPrint("URL id " + params["url_id"])
	switch site = params["site"] {
		case "demonoid":
			debugPrint("Running " +  params["site"] + " extractor")
			extractDemonoid(params["url_id"], rw)
		default:
			debugPrint("No extractor for " +  params["site"])

	}
}

func extractDemonoid(id string, rw http.ResponseWriter) {
	doc := getPage(dURL + id)
	doc.Find(".ctable_content div a").Each(func(i int, s *goquery.Selection) {
		e, _ := s.Attr("href")
		debugPrint(e)
		if strings.Contains(e, "www.hypercache.pw") {
			debugPrint("Returning file from: " + e)
			toReturn := getPageAsBytes(e)
			rw.Write(toReturn)
			
		}
	  })
}

// main function to boot up everything
func main() {
    log.Println("Setting up routes")
    router := mux.NewRouter()
	router.HandleFunc("/", proxyRequest).Queries("url","{url}", "use_regex", "{use_regex}", "to_replace", "{to_replace}", "replacement", "{replacement}").Methods("GET")
	router.HandleFunc("/extractor/{site}/{url_id}", extract).Methods("GET")
    log.Println("Routes set")
    log.Fatal(http.ListenAndServe(":5000", router))
}