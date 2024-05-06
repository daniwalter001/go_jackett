package types

type StreamMeta struct {
	Streams []TorrentStreams `json:"streams"`
}

type Streams struct {
	Title         string `json:"title"`
	Name          string `json:"name"`
	URL           string `json:"url"`
	BehaviorHints struct {
		NotWebReady  bool `json:"notWebReady"`
		ProxyHeaders struct {
			Request struct {
				Accept        string `json:"Accept"`
				Authorization string `json:"Authorization"`
			} `json:"request"`
		} `json:"proxyHeaders"`
	} `json:"behaviorHints"`
}

type TorrentStreams struct {
	Name          string        `json:"name"`
	Type          string        `json:"type"`
	InfoHash      string        `json:"infoHash"`
	FileIdx       int           `json:"fileIdx"`
	Sources       []string      `json:"sources"`
	Title         string        `json:"title"`
	BehaviorHints BehaviorHints `json:"behaviorHints"`
}

type BehaviorHints struct {
	BingeGroup  string `json:"bingeGroup"`
	NotWebReady bool   `json:"notWebReady"`
}
