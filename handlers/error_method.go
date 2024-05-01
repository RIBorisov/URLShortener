package handlers

import "net/http"

func ErrorMethodHandler(w http.ResponseWriter, allowedMethods []string) {
	methodsString := ""
	for _, method := range allowedMethods {
		methodsString = methodsString + ", " + method
	}
	http.Error(w, "Only [%s] methods are allowed", http.StatusMethodNotAllowed)
}
