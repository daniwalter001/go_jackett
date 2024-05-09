package types

type StreamManifest struct {
	Description string   `json:"description"`
	ID          string   `json:"id"`
	Logo        string   `json:"logo"`
	Name        string   `json:"name"`
	Resources   []string `json:"resources"`
	IdPrefixes  []string `json:"idPrefixes"`
	Types       []string `json:"types"`
	Version     string   `json:"version"`
	Catalogs    []string `json:"catalogs"`
}
