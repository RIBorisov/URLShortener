package models

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	Result string `json:"result"`
}

type BatchRequest struct {
	CorrelationId string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type BatchResponse struct {
	CorrelationId string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

type BatchIn []BatchRequest

type BatchOut []BatchResponse
