package meta

type SearchResult[T any] struct {
	Hits hits[T] `json:"hits"`
}

type hits[T any] struct {
	Hits  []HitsRow[T] `json:"hits"`
	Total *struct {
		Value int `json:"value"`
	} `json:"total"`
}

// HitsRow .
type HitsRow[T any] struct {
	Id        string     `json:"_id"`
	Score     float64    `json:"_score"`
	Source    T          `json:"_source"`
	Highlight *Highlight `json:"highlight"`
}

// Highlight .
type Highlight struct{}
