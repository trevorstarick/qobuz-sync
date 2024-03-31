package responses

type Artist struct {
	Image       any      `json:"image"`
	Name        string   `json:"name"`
	ID          int      `json:"id"`
	AlbumsCount int      `json:"albums_count"`
	Slug        string   `json:"slug"`
	Picture     any      `json:"picture"`
	FavoritedAt int      `json:"favorited_at"`
	Roles       []string `json:"roles"`
	SupplierID  int      `json:"supplier_id"`
}
