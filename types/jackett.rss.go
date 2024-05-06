package types

import (
	"encoding/xml"

	"github.com/anacrolix/torrent"
)

type JackettRssReponse struct {
	XMLName xml.Name `xml:"rss"`
	Text    string   `xml:",chardata"`
	Version string   `xml:"version,attr"`
	Atom    string   `xml:"atom,attr"`
	Torznab string   `xml:"torznab,attr"`
	Channel struct {
		Text string `xml:",chardata"`
		Link struct {
			Text string `xml:",chardata"`
			Href string `xml:"href,attr"`
			Rel  string `xml:"rel,attr"`
			Type string `xml:"type,attr"`
		} `xml:"link"`
		Title       string `xml:"title"`
		Description string `xml:"description"`
		Language    string `xml:"language"`
		Category    string `xml:"category"`
		Item        []struct {
			Text           string `xml:",chardata"`
			Title          string `xml:"title"`
			Guid           string `xml:"guid"`
			Jackettindexer struct {
				Text string `xml:",chardata"`
				ID   string `xml:"id,attr"`
			} `xml:"jackettindexer"`
			Type        string   `xml:"type"`
			Comments    string   `xml:"comments"`
			PubDate     string   `xml:"pubDate"`
			Size        string   `xml:"size"`
			Grabs       string   `xml:"grabs"`
			Description string   `xml:"description"`
			Link        string   `xml:"link"`
			Category    []string `xml:"category"`
			Enclosure   struct {
				Text   string `xml:",chardata"`
				URL    string `xml:"url,attr"`
				Length string `xml:"length,attr"`
				Type   string `xml:"type,attr"`
			} `xml:"enclosure"`
			Attr []struct {
				Text  string `xml:",chardata"`
				Name  string `xml:"name,attr"`
				Value string `xml:"value,attr"`
			} `xml:"attr"`
		} `xml:"item"`
	} `xml:"channel"`
}

type ItemsParsed struct {
	Tracker     string         `json:"Tracker,omitempty"`
	Title       string         `json:"Title,omitempty"`
	Seeders     string         `json:"Seeders,omitempty"`
	Peers       string         `json:"Peers,omitempty"`
	Link        string         `json:"Link,omitempty"`
	MagnetURI   string         `json:"MagnetUri,omitempty"`
	TorrentData []torrent.File `json:"TorrentData,omitempty"`
}
