package cursor_test

import (
	"testing"
	"time"

	"github.com/medincident/medincident-backend/internal/cursor"
)

func TestEncodeDecodeRoundTrip(t *testing.T) {
	ts := time.Unix(0, 1746950000000000000).UTC()
	id := "550e8400-e29b-41d4-a716-446655440000"

	encoded := cursor.Encode(ts, id)
	if encoded == "" {
		t.Fatal("Encode returned empty string")
	}

	gotT, gotID, err := cursor.Decode(encoded)
	if err != nil {
		t.Fatalf("Decode error: %v", err)
	}
	if !gotT.Equal(ts) {
		t.Errorf("time mismatch: got %v want %v", gotT, ts)
	}
	if gotID != id {
		t.Errorf("id mismatch: got %q want %q", gotID, id)
	}
}

func TestDecodeInvalidBase64(t *testing.T) {
	_, _, err := cursor.Decode("not-valid-base64!!!")
	if err == nil {
		t.Fatal("expected error for invalid base64")
	}
}

func TestDecodeInvalidJSON(t *testing.T) {
	// base64url of "notjson"
	_, _, err := cursor.Decode("bm90anNvbg")
	if err == nil {
		t.Fatal("expected error for invalid JSON inside base64")
	}
}

func TestEncodeUsesRawURLEncoding(t *testing.T) {
	ts := time.Now()
	id := "test-id"
	s := cursor.Encode(ts, id)
	// RawURLEncoding has no padding '=' and no '+' or '/'
	for _, c := range s {
		if c == '=' || c == '+' || c == '/' {
			t.Errorf("unexpected character %q in cursor (expected RawURLEncoding)", c)
		}
	}
}
