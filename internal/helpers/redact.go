package helpers

import "strings"

func Redact(token string, visibleChars int) string {
	if len(token) <= 2*visibleChars {
		return token
	}
	return token[:visibleChars] + strings.Repeat("*", len(token)-2*visibleChars) + token[len(token)-visibleChars:]
}
