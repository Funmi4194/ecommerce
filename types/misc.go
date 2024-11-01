package types

type Health struct {
	Name    string `json:"name"`
	Status  bool   `json:"status"`
	Version string `json:"version"`
}

type M map[string]interface{}
