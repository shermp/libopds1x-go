package opds1

import (
	"reflect"
	"testing"
)

func TestCalSeriesRegexp(t *testing.T) {
	tests := []struct {
		name    string
		content string
		result  [][]string
	}{
		{
			name: "Test Series",
			content: `<div xmlns="http://www.w3.org/1999/xhtml">RATING: ★★★★<br/>
			TAGS: Fantasy<br/>
			SERIES: Series Name [3.50]<br/>
			<div><p>Summary Content Here</p>
			</div></div>`,
			result: [][]string{
				{"RATING: ★★★★<br/>", "RATING", "★★★★"},
				{"TAGS: Fantasy<br/>", "TAGS", "Fantasy"},
				{"SERIES: Series Name [3.50]<br/>", "SERIES", "Series Name [3.50]"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calibreMetaRegexp.FindAllStringSubmatch(tt.content, -1)
			if got == nil || len(got) != 3 {
				t.Errorf("Len mismatch: Wanted %v\n got %v\n", tt.result, got)
			}
			for i, res := range got {
				if !reflect.DeepEqual(res, tt.result[i]) {
					t.Errorf("Wanted: \n%v got: \n%v\n", tt.result, got)
				}
			}
		})
	}
}

func TestFeed_ParseCalibreMetadata(t *testing.T) {
	tests := []struct {
		name   string
		f      *Feed
		result map[string]string
	}{
		{
			name: "Test Parse calibre metadata",
			f: &Feed{
				Entries: []Entry{
					{
						Content: Content{
							Content: `<div xmlns="http://www.w3.org/1999/xhtml">RATING: ★★★★<br/>
							TAGS: Fantasy<br/>
							SERIES: Series Name [3.50]<br/>
							<div><p>Summary Content Here</p>
							</div></div>`,
						},
					},
				},
				feedTypeDetected: true,
				isNavigation:     false,
			},
			result: map[string]string{"RATING": "★★★★", "TAGS": "Fantasy", "SERIES": "Series Name [3.50]"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.f.ParseCalibreMetadata()
			if !reflect.DeepEqual(tt.f.Entries[0].AppMeta, tt.result) {
				t.Errorf("Got:\n%v, expected:\n%v\n", tt.f.Entries[0].AppMeta, tt.result)
			}
		})
	}
}
