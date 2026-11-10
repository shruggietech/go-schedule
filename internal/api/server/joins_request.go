package server

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// Join writes share the existing portable configuration byte budget. This limit
// belongs to the request body, not to any individual accepted name.
const maxJoinDefinitionBody = 2 << 20

func decodeJoinDefinitionRequest(w http.ResponseWriter, r *http.Request, target any) bool {
	if r.ContentLength > maxJoinDefinitionBody {
		writeError(w, http.StatusRequestEntityTooLarge, CodeValidation, "body", "join definition request exceeds the 2 MiB configuration body limit")
		return false
	}
	r.Body = http.MaxBytesReader(w, r.Body, maxJoinDefinitionBody)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(target)
	if err == nil {
		var extra any
		err = decoder.Decode(&extra)
		if errors.Is(err, io.EOF) {
			return true
		}
		if err == nil {
			err = errors.New("request body must contain one JSON value")
		}
	}
	var oversized *http.MaxBytesError
	if errors.As(err, &oversized) {
		writeError(w, http.StatusRequestEntityTooLarge, CodeValidation, "body", "join definition request exceeds the 2 MiB configuration body limit")
	} else {
		writeError(w, http.StatusBadRequest, CodeValidation, "body", "invalid join JSON; provide one value containing only writable fields")
	}
	return false
}
