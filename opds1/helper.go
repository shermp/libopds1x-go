package opds1

import (
	"strings"
)

func (f *Feed) getFeedType() feedType {
	for _, entry := range f.Entries {
		for _, link := range entry.Links {
			if strings.HasPrefix(link.Rel, "http://opds-spec.org/acquisition") ||
				link.TypeLink == "application/atom+xml;profile=opds-catalog;kind=acquisition" {
				return acquisition
			}
		}
	}
	return navigation
}

// IsNavigation determines whether feed is a navigation feed.
func (f *Feed) IsNavigation() bool {
	return f.getFeedType() == navigation
}

// IsAcquisition determines whether feed is an acquisition feed
func (f *Feed) IsAcquisition() bool {
	return f.getFeedType() == acquisition
}
