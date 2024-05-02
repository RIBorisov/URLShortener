package handlers

import (
	"fmt"
	"net/http"
)

func GenerateURL(r *http.Request, shortLink string) string {
	resultString := fmt.Sprintf("http://%s/%s", r.Host, shortLink)
	return resultString
}
