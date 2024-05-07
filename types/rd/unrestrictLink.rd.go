package rd

type UnrestrictLinkResponse struct {
	ID         string `json:"id,omitempty"`
	Filename   string `json:"filename,omitempty"`
	MimeType   string `json:"mimeType,omitempty"`
	Filesize   int64  `json:"filesize,omitempty"`
	Link       string `json:"link,omitempty"`
	Host       string `json:"host,omitempty"`
	HostIcon   string `json:"host_icon,omitempty"`
	Chunks     int    `json:"chunks,omitempty"`
	Crc        int    `json:"crc,omitempty"`
	Download   string `json:"download,omitempty"`
	Streamable int    `json:"streamable,omitempty"`
}
