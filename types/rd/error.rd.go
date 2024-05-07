package rd

type RdError struct {
	Error        string `json:"error,omitempty"`
	ErrorCode    int    `json:"error_code,omitempty"`
	ErrorDetails string `json:"error_details,omitempty"`
}
