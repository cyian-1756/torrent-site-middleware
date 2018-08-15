package shared

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/feeds"
)

const (
	debug = true
)

// GetPage Get a page as a []byte
func GetPage(url string) []byte {
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

func GetDoc(url string) *goquery.Document {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	g, err := goquery.NewDocumentFromReader(resp.Body)
	HandleError(err)
	return g
}

func HandleError(err error) {
	if err != nil {
		panic(err)
	}
}

func DebugPrint(line string) {
	if debug {
		log.Println(line)
	}
}

func GenerateFeed(title string, link string, des string) *feeds.Feed {
	feed := &feeds.Feed{
		Title:       title,
		Link:        &feeds.Link{Href: link},
		Description: des,
		Author:      &feeds.Author{Name: "middleware"},
		Created:     time.Now(),
	}
	return feed
}

func GenerateItem(title string, link string) *feeds.Item {
	return &feeds.Item{
		Title:       title,
		Link:        &feeds.Link{Href: link},
		Description: "",
		Author:      &feeds.Author{Name: "", Email: ""},
		Created:     time.Now(),
	}
}
