package common

import "time"

type Metadata struct {
	Album       string
	AlbumArtist string
	Artist      string
	Comment     string
	Composer    string
	Genre       string
	Date        time.Time
	Disc        int
	DiscTotal   int
	Track       int
	TrackTotal  int
	Title       string
}
