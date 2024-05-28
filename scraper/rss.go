package scraper

import (
	"encoding/xml"
	"io"
	"net/http"
	"time"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Language    string    `xml:"language"`
		Items       []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func URLToFeed(url string) (RSSFeed, error) {
	// 1. create HTTP client + FETCH data
	httpClient := http.Client{
		Timeout: time.Second * 3,
	}

	resp, err := httpClient.Get(url)

	if err != nil {
		return RSSFeed{}, err
	}

	// 2. read data
	defer resp.Body.Close()

	dat, err := io.ReadAll(resp.Body)

	if err != nil {
		return RSSFeed{}, err
	}

	// 3. unmarshal data = parse XML
	rssFeed := RSSFeed{}

	err = xml.Unmarshal(dat, &rssFeed)

	if err != nil {
		return RSSFeed{}, err
	}

	return rssFeed, nil
}
