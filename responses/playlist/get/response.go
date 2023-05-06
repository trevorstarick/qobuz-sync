//nolint:tagliatelle
package playlistget

import "github.com/trevorstarick/qobuz-sync/responses"

type Response struct {
	ImageRectangleMini []string `json:"image_rectangle_mini"`
	FeaturedArtists    []any    `json:"featured_artists"`
	Description        string   `json:"description"`
	CreatedAt          int      `json:"created_at"`
	TimestampPosition  int      `json:"timestamp_position"`
	Images300          []string `json:"images300"`
	Duration           int      `json:"duration"`
	UpdatedAt          int      `json:"updated_at"`
	Genres             []any    `json:"genres"`
	ImageRectangle     []string `json:"image_rectangle"`
	ID                 int      `json:"id"`
	Slug               string   `json:"slug"`
	Owner              struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"owner"`
	UsersCount      int      `json:"users_count"`
	Images150       []string `json:"images150"`
	Images          []string `json:"images"`
	IsCollaborative bool     `json:"is_collaborative"`
	Stores          []string `json:"stores"`
	TracksCount     int      `json:"tracks_count"`
	PublicAt        int      `json:"public_at"`
	Name            string   `json:"name"`
	IsPublic        bool     `json:"is_public"`
	IsFeatured      bool     `json:"is_featured"`
	Tracks          struct {
		Offset int               `json:"offset"`
		Limit  int               `json:"limit"`
		Total  int               `json:"total"`
		Items  []responses.Track `json:"items"`
	} `json:"tracks"`
}
