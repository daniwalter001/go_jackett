package types

type StreamMeta struct {
	Streams []TorrentStreams `json:"streams,omitempty"`
}

type TorrentStreams struct {
	Name          string        `json:"name,omitempty"`
	Type          string        `json:"type,omitempty"`
	InfoHash      string        `json:"infoHash,omitempty"`
	FileIdx       int           `json:"fileIdx,omitempty"`
	Sources       []string      `json:"sources,omitempty"`
	Title         string        `json:"title,omitempty"`
	URL           string        `json:"url,omitempty"`
	BehaviorHints BehaviorHints `json:"behaviorHints,omitempty"`
}

type BehaviorHints struct {
	BingeGroup  string `json:"bingeGroup,omitempty"`
	NotWebReady bool   `json:"notWebReady,omitempty"`
}
