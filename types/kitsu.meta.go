package types

import (
	"time"
)

type KitsuMeta struct {
	Meta struct {
		ID          string   `json:"id,omitempty"`
		KitsuID     string   `json:"kitsu_id,omitempty"`
		Type        string   `json:"type,omitempty"`
		AnimeType   string   `json:"animeType,omitempty"`
		Name        string   `json:"name,omitempty"`
		Slug        string   `json:"slug,omitempty"`
		Aliases     []string `json:"aliases,omitempty"`
		Genres      []string `json:"genres,omitempty"`
		Logo        string   `json:"logo,omitempty"`
		Poster      string   `json:"poster,omitempty"`
		Background  string   `json:"background,omitempty"`
		Description string   `json:"description,omitempty"`
		ReleaseInfo string   `json:"releaseInfo,omitempty"`
		Year        string   `json:"year,omitempty"`
		ImdbRating  string   `json:"imdbRating,omitempty"`
		UserCount   int      `json:"userCount,omitempty"`
		Status      string   `json:"status,omitempty"`
		Runtime     string   `json:"runtime,omitempty"`
		Trailers    []struct {
			Source string `json:"source,omitempty"`
			Type   string `json:"type,omitempty"`
		} `json:"trailers,omitempty"`
		Videos []Videos `json:"videos,omitempty"`
		Links  []struct {
			Name     string `json:"name,omitempty"`
			Category string `json:"category,omitempty"`
			URL      string `json:"url,omitempty"`
		} `json:"links,omitempty"`
		ImdbID string `json:"imdb_id,omitempty"`
	} `json:"meta,omitempty"`
	CacheMaxAge int `json:"cacheMaxAge,omitempty"`
}

type Videos struct {
	ID          string    `json:"id,omitempty"`
	Title       string    `json:"title,omitempty"`
	Released    time.Time `json:"released,omitempty"`
	Season      int       `json:"season,omitempty"`
	Episode     int       `json:"episode,omitempty"`
	Thumbnail   string    `json:"thumbnail,omitempty"`
	ImdbID      string    `json:"imdb_id,omitempty"`
	ImdbSeason  int       `json:"imdbSeason,omitempty"`
	ImdbEpisode int       `json:"imdbEpisode,omitempty"`
}
