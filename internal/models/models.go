package models

import (
	"fmt"
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
	if len(s.URL) < 3 {
		return fmt.Errorf("URL should be at least 3 characters long")
	}
	return nil
}
