package github

import (
	"context"
	"errors"
	"fmt"

	ghv4 "github.com/shurcooL/githubv4"
)

// resolveOwnerID resolves an owner login (org or user) to a GraphQL ID, preferring org if both exist.
// Returns the ID or an error ("owner not found" if neither found).
type GraphQLClient interface {
	Query(ctx context.Context, q interface{}, vars map[string]interface{}) error
}

func resolveOwnerID(ctx context.Context, client GraphQLClient, owner string) (ghv4.ID, error) {
	var orgQ struct {
		Organization *struct{ ID ghv4.ID } `graphql:"organization(login: $login)"`
	}
	orgVars := map[string]interface{}{"login": ghv4.String(owner)}
	orgErr := client.Query(ctx, &orgQ, orgVars)
	orgNotFound := orgErr != nil && isGraphQLNotFound(orgErr)
	if orgErr != nil && !orgNotFound {
		return "", fmt.Errorf("organization lookup failed: %w", orgErr)
	}
	if orgQ.Organization != nil {
		return orgQ.Organization.ID, nil
	}

	var userQ struct {
		User *struct{ ID ghv4.ID } `graphql:"user(login: $login)"`
	}
	userVars := map[string]interface{}{"login": ghv4.String(owner)}
	userErr := client.Query(ctx, &userQ, userVars)
	userNotFound := userErr != nil && isGraphQLNotFound(userErr)
	if userErr != nil && !userNotFound {
		return "", fmt.Errorf("user lookup failed: %w", userErr)
	}
	if userQ.User != nil {
		return userQ.User.ID, nil
	}
	if orgNotFound && userNotFound {
		return "", errors.New("owner not found")
	}
	return "", errors.New("owner not found") // Defensive fallback
}
