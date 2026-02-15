package auth

import "regexp"

var bearerRe = regexp.MustCompile(`(?i)bearer\s+[A-Za-z0-9._\-]+`)

func RedactToken(input string) string {
	return bearerRe.ReplaceAllString(input, "Bearer [REDACTED]")
}
