package scrapers

import (
	"strings"

	"../shared"
	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/feeds"
)

func Scrape1337x(feed *feeds.Feed, urlToScrape string) string {
	doc := shared.GetDoc(urlToScrape)
	doc.Find("tbody > tr > td > a").Each(func(i int, s *goquery.Selection) {
		e, _ := s.Attr("href")
		shared.DebugPrint(e)
		if strings.Contains(e, "/torrent/") {
			// Set the link to the url of the 1337x extractor
			// TODO dehardcode url
			feed.Add(shared.GenerateItem(s.Text(), "http://127.0.0.1:5000/extractor/1337x/"+strings.Split(e, "/")[2]))
		}

	})
	e, err := feed.ToRss()
	shared.HandleError(err)
	return e

}
