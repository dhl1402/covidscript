package svc

import (
	"context"
	"encoding/json"
	"net/http"
)

// EncodeError encodes an error to http response
func EncodeError(ctx context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	data := map[string]interface{}{
		"error": err.Error(),
	}
	status := http.StatusBadRequest
	if e, ok := err.(Error); ok {
		data["code"] = e.Code()
		data["message"] = e.Message()
		switch e.Code() {
		case ErrCodeUnauthorized:
			status = http.StatusUnauthorized
		case ErrCodePermissionDenied:
			status = http.StatusForbidden
		case ErrCodeSystem:
			status = http.StatusInternalServerError
		}
	}
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
