package rss

import "encoding/xml"

type RSS struct {
	Channel Channel `xml:"channel"`
}

type RDF struct {
	XMLName xml.Name   `xml:"RDF"`
	Channel RDFChannel `xml:"channel"`
	Items   []Item     `xml:"item"`
}

type RDFChannel struct {
	Title       string `xml:"title"`
	Hyperlink   string `xml:"link"`
	Description string `xml:"description"`
}

type Channel struct {
	Title       string `xml:"title"`
	Hyperlink   string `xml:"link"`
	Description string `xml:"description"`
	Published   string `xml:"pubDate"`
	Updated     string `xml:"lastBuildDate"`
	Items       []Item `xml:"item"`
}

type Item struct {
	Title     string `xml:"title"`
	Hyperlink string `xml:"link"`
	Published string `xml:"pubDate"`
}
