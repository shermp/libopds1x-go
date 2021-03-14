// Package opds1 provide parsing and generation method for an OPDS1.X feed
// https://github.com/opds-community/opds-revision/blob/master/opds-1.2.md
package opds1

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type imageRel string

const (
	stdImage   imageRel = "http://opds-spec.org/image"
	thumbImage imageRel = "http://opds-spec.org/image/thumbnail"
)

// Feed root element for acquisition or navigation feed
type Feed struct {
	ID               string    `xml:"id"`
	Title            string    `xml:"title"`
	Updated          time.Time `xml:"updated"`
	Entries          []Entry   `xml:"entry"`
	Links            []Link    `xml:"link"`
	TotalResults     int       `xml:"totalResults"`
	ItemsPerPage     int       `xml:"itemsPerPage"`
	feedTypeDetected bool
	isNavigation     bool
}

// Link link to different resources
type Link struct {
	Rel                 string                `xml:"rel,attr"`
	Href                string                `xml:"href,attr"`
	TypeLink            string                `xml:"type,attr"`
	Title               string                `xml:"title,attr"`
	FacetGroup          string                `xml:"facetGroup,attr"`
	Count               int                   `xml:"count,attr"`
	Price               Price                 `xml:"price"`
	IndirectAcquisition []IndirectAcquisition `xml:"indirectAcquisition"`
}

// Author represent the feed author or the entry author
type Author struct {
	Name string `xml:"name"`
	URI  string `xml:"uri"`
}

// Entry an atom entry in the feed
type Entry struct {
	Title      string     `xml:"title"`
	ID         string     `xml:"id"`
	Identifier string     `xml:"identifier"`
	Updated    *time.Time `xml:"updated"`
	Rights     string     `xml:"rights"`
	Publisher  string     `xml:"publisher"`
	Author     []Author   `xml:"author,omitempty"`
	Language   string     `xml:"language"`
	Issued     string     `xml:"issued"` // Check for format
	Published  *time.Time `xml:"published"`
	Category   []Category `xml:"category,omitempty"`
	Links      []Link     `xml:"link,omitempty"`
	Summary    Content    `xml:"summary"`
	Content    Content    `xml:"content"`
	Series     []Serie    `xml:"Series"`
	AppMeta    map[string]string
}

// Content content tag in an entry, the type will be html or text
type Content struct {
	Content     string `xml:",innerxml"`
	ContentType string `xml:"type,attr"`
}

// Category represent the book category with scheme and term to machine
// handling
type Category struct {
	Scheme string `xml:"scheme,attr"`
	Term   string `xml:"term,attr"`
	Label  string `xml:"label,attr"`
}

// Price represent the book price
type Price struct {
	CurrencyCode string  `xml:"currencycode,attr"`
	Value        float64 `xml:",cdata"`
}

// IndirectAcquisition represent the link mostly for buying or borrowing
// a book
type IndirectAcquisition struct {
	TypeAcquisition     string                `xml:"type,attr"`
	IndirectAcquisition []IndirectAcquisition `xml:"indirectAcquisition"`
}

// Serie store serie information from schema.org
type Serie struct {
	Name     string  `xml:"name,attr"`
	URL      string  `xml:"url,attr"`
	Position float32 `xml:"position,attr"`
}

// ParseURL take a url in entry and parse the feed
func ParseURL(url string) (*Feed, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	res, errReq := http.DefaultClient.Do(request)
	if errReq != nil {
		return nil, errReq
	}

	return ParseResponse(res)
}

// ParseResponse is useful if you need more control over the
// HTTP communications
func ParseResponse(r *http.Response) (*Feed, error) {
	var feed Feed
	buff, errRead := ioutil.ReadAll(r.Body)
	if errRead != nil {
		return nil, errRead
	}
	xml.Unmarshal(buff, &feed)
	return &feed, nil
}

func (f *Feed) detectFeedType() {
	if !f.feedTypeDetected {
		f.isNavigation = true
		for _, entry := range f.Entries {
			for _, link := range entry.Links {
				if strings.HasPrefix(link.Rel, "http://opds-spec.org/acquisition") ||
					link.TypeLink == "application/atom+xml;profile=opds-catalog;kind=acquisition" {
					f.isNavigation = false
				}
			}
		}
	}
}

// IsNavigation determines whether feed is a navigation feed.
func (f *Feed) IsNavigation() bool {
	f.detectFeedType()
	return f.isNavigation
}

// IsAcquisition determines whether feed is an acquisition feed
func (f *Feed) IsAcquisition() bool {
	f.detectFeedType()
	return !f.isNavigation
}

func (f *Feed) imageLink(rel imageRel) *Link {
	if !f.IsAcquisition() {
		return nil
	}
	for _, link := range f.Links {
		if link.Rel == string(rel) {
			return &link
		}
	}
	return nil
}

// ThumbnailLink gets the link for a thumbnail from an aquisition feed.
// nil will be returned if not an aquisition feed, or a thumbnail could not be found
func (f *Feed) ThumbnailLink() *Link {
	return f.imageLink(thumbImage)
}

// ImageLink gets the link for an image from an aquisition feed.
// nil will be returned if not an aquisition feed, or an image could not be found
func (f *Feed) ImageLink() *Link {
	return f.imageLink(stdImage)
}
