package opds1

import (
	"strings"
)

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
}
