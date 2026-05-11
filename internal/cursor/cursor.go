package cursor

import (
	"encoding/base64"
	"encoding/json"
	"time"
)

type token struct {
	T int64  `json:"t"`
	I string `json:"i"`
}

// Encode encodes a (timestamp, id) pair as base64url(JSON{"t":unix_nano,"i":id}).
func Encode(t time.Time, id string) string {
	b, _ := json.Marshal(token{T: t.UnixNano(), I: id})
	return base64.RawURLEncoding.EncodeToString(b)
}

// Decode decodes a cursor produced by Encode.
// Returns an error when s is not valid base64url or not valid JSON.
func Decode(s string) (time.Time, string, error) {
	b, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return time.Time{}, "", err
	}
	var tok token
	if err := json.Unmarshal(b, &tok); err != nil {
		return time.Time{}, "", err
	}
	return time.Unix(0, tok.T).UTC(), tok.I, nil
}
