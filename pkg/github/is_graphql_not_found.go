package github

import "strings"

// isGraphQLNotFound returns true if the error is a GraphQL 'not found' error.
func isGraphQLNotFound(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "could not resolve to a User") ||
		strings.Contains(msg, "could not resolve to an Organization") ||
		strings.Contains(msg, "non-200 OK status code: 400") ||
		strings.Contains(msg, "non-200 OK status code: 404")
}
