package opds1

import (
	"regexp"
	"strings"
)

//var calSeriesRegexp = regexp.MustCompile(`(SERIES: *)([^\[]+)\[([^\]]+)]<br *\/>`)
var calibreMetaRegexp = regexp.MustCompile(`(SERIES|TAGS|RATING): *([^<]+)<br *\/>`)

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

// ParseCalibreMetadata attempts to parse Calibre metadata embedded into the content tag
func (f *Feed) ParseCalibreMetadata() {
	if f.IsNavigation() {
		return
	}
	for i := range f.Entries {
		// Calibre stores extra metadata in the 'content' tag, along with the book description
		meta := calibreMetaRegexp.FindAllStringSubmatch(f.Entries[i].Content.Content, -1)
		if meta != nil {
			if f.Entries[i].AppMeta == nil {
				f.Entries[i].AppMeta = make(map[string]string)
			}
			for _, md := range meta {
				f.Entries[i].AppMeta[md[1]] = md[2]
			}
		}
	}
}
