package handlers

import (
	"fmt"
	"net/http"
)

func ErrorMethodHandler(w http.ResponseWriter, allowedMethods []string) {
	allowed := fmt.Sprintf("Only %s methods allowed", allowedMethods)
	http.Error(w, allowed, http.StatusMethodNotAllowed)
}
