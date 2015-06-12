package main

import (
	"encoding/xml"
	"github.com/PuerkitoBio/goquery"
	"net/http"
)

type Rss struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel *Channel
}

type Channel struct {
	XMLName     xml.Name `xml:"channel"`
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	Description string   `xml:"description"`
	Items       []Item
}

type Item struct {
	XMLName     xml.Name `xml:"item"`
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	Description string   `xml:"description"`
}

func NewRss() Rss {
	return Rss{
		Version: "2.0",
		Channel: &Channel{
			Title:       "Berlin Startup Jobs",
			Link:        "http://berlinstartupjobs.com/engineering/",
			Description: "New Engineering Jobs from Berlins Start Up Scene",
			Items:       make([]Item, 0)}}
}

func createFeed() (Rss, error) {
	doc, err := goquery.NewDocument("http://berlinstartupjobs.com/engineering/")
	if err != nil {
		return NewRss(), err
	}

	feed := NewRss()
	doc.Find(".post").Each(func(i int, s *goquery.Selection) {
		link := s.Find("a").First()
		linkURL, _ := link.Attr("href")
		title := s.Find("h2 a").First().Text()
		postContent := s.Find(".post-content")
		postText := postContent.Find("p").First().Text()

		item := Item{
			Title:       title,
			Link:        linkURL,
			Description: postText}

		feed.Channel.Items = append(feed.Channel.Items, item)
	})
	return feed, nil
}

func serveRSS(w http.ResponseWriter, r *http.Request) {
	rss, err := createFeed()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rssAsXml, err := xml.Marshal(rss)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	w.Write(rssAsXml)
}

func main() {
	http.HandleFunc("/jobs", serveRSS)
	http.ListenAndServe(":8000", nil)
}
