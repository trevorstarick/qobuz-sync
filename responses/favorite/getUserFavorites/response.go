package favoritesgetuserfavorites

import "github.com/trevorstarick/qobuz-sync/responses"

type Response struct {
	Albums  AlbumsRes  `json:"albums"`
	Tracks  TracksRes  `json:"tracks"`
	Artists ArtistsRes `json:"artists"`
	User    User       `json:"user"`
}
type TracksRes struct {
	Offset int               `json:"offset"`
	Limit  int               `json:"limit"`
	Total  int               `json:"total"`
	Items  []responses.Track `json:"items"`
}
type AlbumsRes struct {
	Offset int               `json:"offset"`
	Limit  int               `json:"limit"`
	Total  int               `json:"total"`
	Items  []responses.Album `json:"items"`
}
type ArtistsRes struct {
	Offset int                `json:"offset"`
	Limit  int                `json:"limit"`
	Total  int                `json:"total"`
	Items  []responses.Artist `json:"items"`
}
type User struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
}
