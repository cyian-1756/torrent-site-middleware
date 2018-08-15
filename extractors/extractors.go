package extractors

import (
	"net/http"
	"strings"

	"../shared"
	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/mux"
)

const (
	dURL     = "https://www.demonoid.pw/genlb.php?genid="
	leetxURL = "https://1337x.to"
)

// Extract a torrent from a page
func Extract(rw http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	shared.DebugPrint("Using extractor " + params["site"])
	shared.DebugPrint("URL id " + params["url_id"])
	site := params["site"]
	switch site {

	case "demonoid":
		shared.DebugPrint("Running " + params["site"] + " extractor")
		extractDemonoid(params["url_id"], rw)
	case "1337x":
		shared.DebugPrint("Running " + params["site"] + " extractor")
		extract1337x(params["url_id"], rw)
	default:
		// TODO return error instead of just printing to log
		shared.DebugPrint("No extractor for " + params["site"])
		rw.Write([]byte("No extractor for " + params["site"]))

	}
}

func extractDemonoid(id string, rw http.ResponseWriter) {
	doc := shared.GetDoc(dURL + id)
	doc.Find(".ctable_content div a").Each(func(i int, s *goquery.Selection) {
		e, _ := s.Attr("href")
		shared.DebugPrint(e)
		// Make sure the link points to Demonoids cache
		if strings.Contains(e, "www.hypercache.pw") {
			toReturn := []byte(shared.GetPage(e))
			rw.Write(toReturn)

		}
	})
}

func extract1337x(id string, rw http.ResponseWriter) {
	// Because of 1227xs bugging url routing titleless urls need to end with an extra /
	doc := shared.GetDoc("https://1337x.to/torrent/" + id + "//")
	sel := doc.Find("a")
	for i := range sel.Nodes {
		shared.DebugPrint("Looping over torrent links")
		s := sel.Eq(i)
		e, _ := s.Attr("href")
		if !strings.Contains(e, "magnet:") && strings.Contains(e, "/torrent/") {
			shared.DebugPrint("Downloading page " + e)
			rw.Write([]byte(shared.GetPage(e)))
			return

		}
	}
}
