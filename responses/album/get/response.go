package albumget

import "github.com/trevorstarick/qobuz-sync/responses"

type Response struct {
	*responses.Album

	Tracks struct {
		Offset int               `json:"offset"`
		Limit  int               `json:"limit"`
		Total  int               `json:"total"`
		Items  []responses.Track `json:"items"`
	} `json:"tracks"`
}
