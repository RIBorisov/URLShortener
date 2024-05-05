package handlers

import (
	"fmt"
	c "shortener/internal/config"
)

func GenerateURL(shortLink string) string {
	cfg := c.LoadConfig()
	resultString := fmt.Sprintf("%s/%s", cfg.Server.BaseURL, shortLink)
	return resultString
}
