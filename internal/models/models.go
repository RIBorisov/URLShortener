package models

// ShortenRequest shorten request model.
type ShortenRequest struct {
	URL string `json:"url"`
}

// ShortenResponse shorten response model.
type ShortenResponse struct {
	Result string `json:"result"`
}

// BatchRequest shorten request model.
type BatchRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// BatchResponse shorten response model.
type BatchResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// Batch shorten model.
type Batch struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
	OriginalURL   string `json:"original_url"`
}

// BatchArray shorten type.
type BatchArray []Batch

// BatchResponseArray shorten type.
type BatchResponseArray []BatchResponse

// URL model.
type URL struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// UserURLs model.
type UserURLs []URL

// User model.
type User struct {
	ID string `json:"user_id"`
}

// BaseRow model.
type BaseRow struct {
	Short string `json:"short"`
	Long  string `json:"long"`
}

// DeleteURLs model.
type DeleteURLs []string

type key int

// CtxUserIDKey context userID key.
const CtxUserIDKey key = iota
