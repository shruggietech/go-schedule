package mcpobserve

import (
	"encoding/base64"
	"encoding/json"
	"errors"
)

var errInvalidCursor = errors.New("invalid_cursor: continuation cursor is invalid")

type cursor struct {
	Version int `json:"v"`
	Offset  int `json:"o"`
}

func encodeCursor(offset int) string {
	encoded, _ := json.Marshal(cursor{Version: 1, Offset: offset})
	return base64.RawURLEncoding.EncodeToString(encoded)
}

func decodeCursor(value string) (int, error) {
	raw, err := base64.RawURLEncoding.DecodeString(value)
	if err != nil {
		return 0, errInvalidCursor
	}
	var decoded cursor
	if err := json.Unmarshal(raw, &decoded); err != nil || decoded.Version != 1 || decoded.Offset < 0 {
		return 0, errInvalidCursor
	}
	return decoded.Offset, nil
}
