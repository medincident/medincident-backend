// Package urlutil provides URL manipulation helpers.
package urlutil

import "net/url"

// redactedPlaceholder is the replacement text used for masked userinfo
// components.
const redactedPlaceholder = "***"

// Redact parses raw as a URL and replaces any embedded userinfo
// (username / password) with masked placeholders. If parsing fails the
// function returns the redactedPlaceholder so that no credentials can
// leak through a malformed URL.
func Redact(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return redactedPlaceholder
	}
	if u.User != nil {
		u.User = url.UserPassword(redactedPlaceholder, redactedPlaceholder)
	}
	return u.String()
}
