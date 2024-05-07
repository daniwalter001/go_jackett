package rd

type AvailabilityFileInfo struct {
	Filename string `json:"filename,omitempty"`
	Filesize int64  `json:"filesize,omitempty"`
}
type AvailabilityHoster map[string][]map[string]AvailabilityFileInfo
type AvailabilityResponse map[string]AvailabilityHoster
