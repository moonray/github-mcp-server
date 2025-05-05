package github

import (
	"errors"
	"testing"
)

func TestIsGraphQLNotFound(t *testing.T) {
	cases := []struct {
		name    string
		err     error
		expects bool
	}{
		{"nil error", nil, false},
		{"empty error", errors.New(""), false},
		{"user not found", errors.New("could not resolve to a User"), true},
		{"org not found", errors.New("could not resolve to an Organization"), true},
		{"400 code", errors.New("non-200 OK status code: 400"), true},
		{"404 code", errors.New("non-200 OK status code: 404"), true},
		{"other error", errors.New("some other error"), false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := isGraphQLNotFound(c.err)
			if result != c.expects {
				t.Errorf("expected %v, got %v for input %v", c.expects, result, c.err)
			}
		})
	}
}
