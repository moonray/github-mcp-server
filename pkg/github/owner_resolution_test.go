package github

import (
	"context"
	"errors"
	"testing"

	ghv4 "github.com/shurcooL/githubv4"
	"github.com/stretchr/testify/assert"
)


type fakeGraphQLClient struct {
	orgID      ghv4.ID
	userID     ghv4.ID
	orgErr     error
	userErr    error
}

func (f *fakeGraphQLClient) Query(ctx context.Context, q interface{}, vars map[string]interface{}) error {
	if oq, ok := q.(*struct {
		Organization *struct{ ID ghv4.ID } `graphql:"organization(login: $login)"`
	}); ok {
		if f.orgErr != nil {
			return f.orgErr
		}
		if f.orgID != "" {
			*oq = struct {
				Organization *struct{ ID ghv4.ID } `graphql:"organization(login: $login)"`
			}{Organization: &struct{ ID ghv4.ID }{ID: f.orgID}}
		}
		return nil
	}
	if uq, ok := q.(*struct {
		User *struct{ ID ghv4.ID } `graphql:"user(login: $login)"`
	}); ok {
		if f.userErr != nil {
			return f.userErr
		}
		if f.userID != "" {
			*uq = struct {
				User *struct{ ID ghv4.ID } `graphql:"user(login: $login)"`
			}{User: &struct{ ID ghv4.ID }{ID: f.userID}}
		}
		return nil
	}
	return errors.New("unexpected query type")
}

func (f *fakeGraphQLClient) Mutate(ctx context.Context, m interface{}, input interface{}, v map[string]interface{}) error {
	return errors.New("not implemented")
}

func TestResolveOwnerID(t *testing.T) {
	ctx := context.Background()
	owner := "testowner"

	tests := []struct {
		name      string
		orgID     ghv4.ID
		userID    ghv4.ID
		orgErr    error
		userErr   error
		expectID  ghv4.ID
		expectErr string
	}{
		{
			name:     "org exists",
			orgID:    "ORGID",
			userID:   "",
			expectID:  "ORGID",
		},
		{
			name:     "user exists",
			orgID:    "",
			userID:   "USERID",
			expectID:  "USERID",
		},
		{
			name:     "both org and user exist (prefer org)",
			orgID:    "ORGID",
			userID:   "USERID",
			expectID:  "ORGID",
		},
		{
			name:     "neither org nor user exist",
			orgID:    "",
			userID:   "",
			orgErr:   errors.New("non-200 OK status code: 404"),
			userErr:  errors.New("non-200 OK status code: 404"),
			expectErr: "owner not found",
		},
		{
			name:     "org fatal error",
			orgID:    "",
			userID:   "USERID",
			orgErr:   errors.New("fatal org error"),
			expectErr: "organization lookup failed",
		},
		{
			name:     "user fatal error",
			orgID:    "",
			userID:   "",
			orgErr:   errors.New("non-200 OK status code: 404"),
			userErr:  errors.New("fatal user error"),
			expectErr: "user lookup failed",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client := &fakeGraphQLClient{
				orgID:   tc.orgID,
				userID:  tc.userID,
				orgErr:  tc.orgErr,
				userErr: tc.userErr,
			}
			id, err := resolveOwnerID(ctx, client, owner)
			if tc.expectErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectID, id)
			}
		})
	}
}
