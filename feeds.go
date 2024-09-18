package feeds

import (
	"bytes"
	"encoding/xml"
	"github.com/IucassacuI/feeds/atom"
	"github.com/IucassacuI/feeds/rss"
	"strings"
	"time"
)

type Item struct {
	Hyperlink string
	Title     string
	Published string
	Updated   string
}

type Feed struct {
	Format      string
	Hyperlink   string
	Title       string
	Description string
	Published   string
	Updated     string
	Author      string
	Items       []Item
}

func rfc(datetime string) string {
	t, _ := time.Parse(time.DateTime, datetime)
	return t.Format(time.RFC3339)
}

func datetime(rfcdate string) string {
	t, _ := time.Parse(time.RFC3339, rfcdate)
	return t.Format(time.DateTime)
}

func ParseRSS(feed []byte) Feed {
	var doc rss.RSS
	err := xml.Unmarshal(feed, &doc)
	if err != nil {
		return Feed{}
	}

	var final Feed
	final.Format = "rss"
	final.Hyperlink = doc.Channel.Hyperlink
	final.Title = doc.Channel.Title
	final.Description = doc.Channel.Description
	final.Published = doc.Channel.Published
	final.Updated = doc.Channel.Updated
	final.Author = "N/A"

	for _, field := range []*string{&final.Hyperlink, &final.Title, &final.Description, &final.Published, &final.Updated} {
		if *field == "" {
			*field = "N/A"
		} else {
			*field = strings.TrimSpace(*field)
		}
	}

	for _, item := range doc.Channel.Items {
		for _, field := range []*string{&item.Hyperlink, &item.Title, &item.Published} {
			if *field == "" {
				*field = "N/A"
			} else {
				*field = strings.TrimSpace(*field)
			}
		}

		final.Items = append(final.Items, Item{Hyperlink: item.Hyperlink, Title: item.Title, Published: item.Published, Updated: "N/A"})
	}

	return final
}

func ParseRDF(feed []byte) Feed {
	var doc rss.RDF
	err := xml.Unmarshal(feed, &doc)

	if err != nil {
		return Feed{}
	}

	var final Feed
	final.Format = "rdf"
	final.Hyperlink = doc.Channel.Hyperlink
	final.Title = doc.Channel.Title
	final.Description = doc.Channel.Description
	final.Published = "N/A"
	final.Updated = "N/A"
	final.Author = "N/A"

	for _, field := range []*string{&final.Hyperlink, &final.Title, &final.Description} {
		if *field == "" {
			*field = "N/A"
		} else {
			*field = strings.TrimSpace(*field)
		}
	}

	for _, item := range doc.Items {
		for _, field := range []*string{&item.Hyperlink, &item.Title, &item.Published} {
			if *field == "" {
				*field = "N/A"
			} else {
				*field = strings.TrimSpace(*field)
			}
		}

		final.Items = append(final.Items, Item{Hyperlink: item.Hyperlink, Title: item.Title, Published: item.Published, Updated: "N/A"})
	}

	return final
}

func ParseAtom(feed []byte) Feed {
	var doc atom.Atom
	err := xml.Unmarshal(feed, &doc)
	if err != nil {
		return Feed{}
	}

	var final Feed
	final.Format = "atom"
	final.Hyperlink = doc.Hyperlink.Href
	final.Title = doc.Title
	final.Description = "N/A"
	final.Published = doc.Published
	final.Updated = doc.Updated
	final.Author = doc.Author

	for _, field := range []*string{&final.Hyperlink, &final.Title, &final.Published, &final.Updated, &final.Author} {
		if *field == "" {
			*field = "N/A"
		} else {
			*field = strings.TrimSpace(*field)
		}
	}

	if final.Published != "N/A" {
		final.Published = datetime(final.Published)
	}

	if final.Updated != "N/A" {
		final.Updated = datetime(final.Updated)
	}

	for _, entry := range doc.Entries {
		for _, field := range []*string{&entry.Hyperlink.Href, &entry.Title, &entry.Published} {
			if *field == "" {
				*field = "N/A"
			} else {
				*field = strings.TrimSpace(*field)
			}
		}

		final.Items = append(final.Items, Item{
			Title:     entry.Title,
			Hyperlink: entry.Hyperlink.Href,
			Published: datetime(entry.Published),
			Updated:   datetime(entry.Updated),
		})
	}

	return final
}

func Marshal(feed Feed) []byte {
	var data []byte
	var err error

	if feed.Format == "rss" {
		doc := rss.RSS{
			Channel: rss.Channel{
				Title:       feed.Title,
				Description: feed.Description,
				Hyperlink:   feed.Hyperlink,
				Published:   feed.Published,
				Updated:     feed.Updated,
			},
		}

		for _, item := range feed.Items {
			doc.Channel.Items = append(doc.Channel.Items, rss.Item{Title: item.Title, Hyperlink: item.Hyperlink, Published: item.Published})
		}

		data, err = xml.Marshal(doc)
		if err == nil {
			data = bytes.Replace(data, []byte("RSS>"), []byte("rss>"), 2)
		}
	} else if feed.Format == "rdf" {
		doc := rss.RDF{
			Channel: rss.RDFChannel{
				Title:       feed.Title,
				Hyperlink:   feed.Hyperlink,
				Description: feed.Description,
			},
		}

		for _, item := range feed.Items {
			doc.Items = append(doc.Items, rss.Item{Title: item.Title, Hyperlink: item.Hyperlink, Published: item.Published})
		}

		data, err = xml.Marshal(doc)
		if err == nil {
			data = bytes.Replace(data, []byte("RDF>"), []byte("rdf:RDF>"), 2)
		}
	} else if feed.Format == "atom" {
		doc := atom.Atom{
			Title:     feed.Title,
			Author:    feed.Author,
			Published: rfc(feed.Published),
			Updated:   rfc(feed.Updated),
			Hyperlink: atom.Hyperlink{Href: feed.Hyperlink},
		}

		for _, item := range feed.Items {
			doc.Entries = append(doc.Entries, atom.Entry{
				Title:     item.Title,
				Published: rfc(item.Published),
				Updated:   rfc(item.Updated),
				Hyperlink: atom.Hyperlink{Href: item.Hyperlink},
			})
		}

		data, err = xml.Marshal(doc)
		if err == nil {
			data = bytes.Replace(data, []byte("Atom>"), []byte("feed>"), 2)
		}
	}

	if err != nil {
		print(err)
	}

	return data
}

func Parse(feed []byte) Feed {
	if bytes.Contains(feed, []byte("<feed")) {
		feed := feed[bytes.Index(feed, []byte("<feed")):]
		return ParseAtom(feed)
	} else if bytes.Contains(feed, []byte("<rss")) {
		feed := feed[bytes.Index(feed, []byte("<rss")):]
		return ParseRSS(feed)
	} else if bytes.Contains(feed, []byte("<rdf:RDF")) {
		feed := feed[bytes.Index(feed, []byte("<rdf:RDF")):]
		return ParseRDF(feed)
	}

	return Feed{}
}
