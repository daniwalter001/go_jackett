package rd

type TorrentInfoResponse struct {
	ID               string   `json:"id,omitempty"`
	Filename         string   `json:"filename,omitempty"`
	OriginalFilename string   `json:"original_filename,omitempty"`
	Hash             string   `json:"hash,omitempty"`
	Bytes            int64    `json:"bytes,omitempty"`
	OriginalBytes    int64    `json:"original_bytes,omitempty"`
	Host             string   `json:"host,omitempty"`
	Split            int      `json:"split,omitempty"`
	Progress         int      `json:"progress,omitempty"`
	Status           string   `json:"status,omitempty"`
	Added            string   `json:"added,omitempty"`
	Files            []Files  `json:"files,omitempty"`
	Links            []string `json:"links,omitempty"`
}

type Files struct {
	ID       int    `json:"id,omitempty"`
	Path     string `json:"path,omitempty"`
	Bytes    int64  `json:"bytes,omitempty"`
	Selected int    `json:"selected,omitempty"`
}
