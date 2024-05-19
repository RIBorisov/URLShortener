package models

import (
	"errors"
)

type Request struct {
	Request SimpleRequest `json:"request"`
}

type Response struct {
	Result string `json:"result"`
}

type SimpleRequest struct {
	URL string `json:"url"`
}

func (s *SimpleRequest) Validate() error {
	const minURLLength = 3

	if len(s.URL) < minURLLength {
		return errors.New("URL should be at least 3 characters long")
	}
	return nil
}
