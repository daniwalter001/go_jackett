package types

type HttpClientError struct {
	Op  string `json:"Op,omitempty"`
	URL string `json:"URL,omitempty"`
	Err any    `json:"Err,omitempty"`
}
