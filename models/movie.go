package models

type Movie struct {
	ID          int      `json:"id"`
	TMDB_ID     int      `json:"tmdb_id,omitempty"`
	Title       string   `json:"title"`
	Tagline     *string  `json:"tagline,omitempty"`
	ReleaseYear int      `json:"release_year"`
	Genres      []Genre  `json:"genres"`
	Overview    *string  `json:"overview,omitempty"`
	Score       *float32 `json:"score,omitempty"`
	Popularity  *float32 `json:"popularity,omitempty"`
	Keywords    []string `json:"keywords"`
	Language    *string  `json:"language,omitempty"`
	PosterURL   *string  `json:"poster_url,omitempty"`
	TrailerURL  *string  `json:"trailer_url,omitempty"`
	Casting     []Actor  `json:"casting"`
}
