package main

import (
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"
)

func inSlice[T comparable](slice []T, element T) bool {
	for k := range slice {
		if slice[k] == element {
			return true
		}
	}

	return false
}

func requireMethod(w http.ResponseWriter, r *http.Request, allowedMethods ...string) (cont bool) {
	if !inSlice(allowedMethods, r.Method) {
		log.Error("Bad Method Request", zap.String("Method", r.Method), zap.String("Allowed Methods", fmt.Sprint(allowedMethods)), zap.String("From", fmtip(r.RemoteAddr)))

		w.Header().Add("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "allowed methods: "+fmt.Sprint(allowedMethods))

		return false
	}

	return true
}
