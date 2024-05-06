package types

import (
	"encoding/json"
	"time"
)

func UnmarshalMeta(data []byte) (IMDBMeta, error) {
	var r IMDBMeta
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *IMDBMeta) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type IMDBMeta struct {
	Meta struct {
		Awards       *string            `json:"awards,omitempty"`
		Cast         []string           `json:"cast,omitempty"`
		Country      *string            `json:"country,omitempty"`
		Description  *string            `json:"description,omitempty"`
		Director     interface{}        `json:"director"`
		DVDRelease   interface{}        `json:"dvdRelease"`
		Genre        []string           `json:"genre,omitempty"`
		ImdbRating   *string            `json:"imdbRating,omitempty"`
		ImdbID       *string            `json:"imdb_id,omitempty"`
		Name         *string            `json:"name,omitempty"`
		Popularity   *float64           `json:"popularity,omitempty"`
		Poster       *string            `json:"poster,omitempty"`
		Released     *time.Time         `json:"released,omitempty"`
		Runtime      *string            `json:"runtime,omitempty"`
		Status       *string            `json:"status,omitempty"`
		TvdbID       *int64             `json:"tvdb_id,omitempty"`
		Type         *string            `json:"type,omitempty"`
		Writer       []string           `json:"writer,omitempty"`
		Year         *string            `json:"year,omitempty"`
		Background   *string            `json:"background,omitempty"`
		Logo         *string            `json:"logo,omitempty"`
		Popularities map[string]float64 `json:"popularities,omitempty"`
		MoviedbID    *int64             `json:"moviedb_id,omitempty"`
		Slug         *string            `json:"slug,omitempty"`
		Trailers     []struct {
			Source *string `json:"source,omitempty"`
			Type   *string `json:"type,omitempty"`
		} `json:"trailers,omitempty"`
		ID          *string  `json:"id,omitempty"`
		Genres      []string `json:"genres,omitempty"`
		ReleaseInfo *string  `json:"releaseInfo,omitempty"`
		Videos      []struct {
			Name        *string    `json:"name,omitempty"`
			Season      *int64     `json:"season,omitempty"`
			Number      *int64     `json:"number,omitempty"`
			FirstAired  *time.Time `json:"firstAired,omitempty"`
			TvdbID      *int64     `json:"tvdb_id,omitempty"`
			Rating      *string    `json:"rating,omitempty"`
			Overview    *string    `json:"overview,omitempty"`
			Thumbnail   *string    `json:"thumbnail,omitempty"`
			ID          *string    `json:"id,omitempty"`
			Released    *time.Time `json:"released,omitempty"`
			Episode     *int64     `json:"episode,omitempty"`
			Description *string    `json:"description,omitempty"`
		} `json:"videos,omitempty"`
		TrailerStreams []struct {
			Title *string `json:"title,omitempty"`
			YtID  *string `json:"ytId,omitempty"`
		} `json:"trailerStreams,omitempty"`
		Links []struct {
			Name     *string `json:"name,omitempty"`
			Category *string `json:"category,omitempty"`
			URL      *string `json:"url,omitempty"`
		} `json:"links,omitempty"`
		BehaviorHints struct {
			DefaultVideoID     interface{} `json:"defaultVideoId"`
			HasScheduledVideos *bool       `json:"hasScheduledVideos,omitempty"`
		} `json:"behaviorHints,omitempty"`
	} `json:"meta,omitempty"`
}
