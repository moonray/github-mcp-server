package github

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	ghv4 "github.com/shurcooL/githubv4"
)

// --- Struct definitions (colocated, per codebase convention) ---

type ListOrganizationProjectsInput struct {
	Organization string `json:"organization"`
	First        int    `json:"first,omitempty"`
	After        string `json:"after,omitempty"`
}

type Project struct {
	ID     string `json:"id"`
	Number int    `json:"number"`
	Title  string `json:"title"`
	URL    string `json:"url"`
}

type ListOrganizationProjectsOutput struct {
	Projects    []Project `json:"projects"`
	EndCursor   string    `json:"end_cursor,omitempty"`
	HasNextPage bool      `json:"has_next_page"`
}

type ListUserProjectsInput struct {
	User  string `json:"user"`
	First int    `json:"first,omitempty"`
	After string `json:"after,omitempty"`
}

type GetProjectInput struct {
	Owner  string `json:"owner"`
	Number int    `json:"number"`
}

type GetProjectItemsInput struct {
	ProjectID string `json:"project_id"`
	First     int    `json:"first,omitempty"`
	After     string `json:"after,omitempty"`
}

type ProjectItem struct {
	ID          string `json:"id"`
	ContentID   string `json:"content_id"`
	ContentType string `json:"content_type"`
	Title       string `json:"title"`
	State       string `json:"state"`
	URL         string `json:"url"`
}

type GetProjectItemsOutput struct {
	Items       []ProjectItem `json:"items"`
	EndCursor   string        `json:"end_cursor,omitempty"`
	HasNextPage bool          `json:"has_next_page"`
}

type CreateProjectInput struct {
	Owner       string `json:"owner"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
}

type AddProjectItemInput struct {
	ProjectID string `json:"project_id"`
	ContentID string `json:"content_id"`
}

type AddProjectItemOutput struct {
	Item ProjectItem `json:"item"`
}

type UpdateProjectItemFieldInput struct {
	ProjectID string `json:"project_id"`
	ItemID    string `json:"item_id"`
	FieldID   string `json:"field_id"`
	Value     string `json:"value"`
}

type UpdateProjectItemFieldOutput struct {
	Item ProjectItem `json:"item"`
}

// --- Handler scaffolds (not implemented yet; return errors) ---

// ListOrganizationProjects lists projects for an organization using the provided githubv4.Client.
// If client is nil, a default client is created using GITHUB_TOKEN from environment.
func ListOrganizationProjects(ctx context.Context, in *ListOrganizationProjectsInput, client *ghv4.Client) (*ListOrganizationProjectsOutput, error) {
	if in.Organization == "" {
		return nil, errors.New("organization is required")
	}

	if client == nil {
		token := os.Getenv("GITHUB_PERSONAL_ACCESS_TOKEN")
		if token == "" {
			return nil, errors.New("GITHUB_PERSONAL_ACCESS_TOKEN not set")
		}
		client = ghv4.NewClient(&http.Client{Transport: &authTransport{token: token}})
	}

	var q struct {
		Organization struct {
			ProjectsV2 struct {
				Nodes []struct {
					ID          ghv4.ID
					Number      ghv4.Int
					Title       ghv4.String
					URL         ghv4.URI
				} `graphql:"nodes"`
				PageInfo struct {
					EndCursor   ghv4.String
					HasNextPage bool
				}
			} `graphql:"projectsV2(first: $first, after: $after)"`
		} `graphql:"organization(login: $org)"`
	}
	vars := map[string]interface{}{
		"org":   ghv4.String(in.Organization),
		"first": ghv4.Int(in.First),
		"after": ghv4.String(in.After),
	}

	err := client.Query(ctx, &q, vars)
	if err != nil {
		return nil, fmt.Errorf("github graphql error: %w", err)
	}

	out := &ListOrganizationProjectsOutput{
		Projects:    []Project{},
		EndCursor:   string(q.Organization.ProjectsV2.PageInfo.EndCursor),
		HasNextPage: q.Organization.ProjectsV2.PageInfo.HasNextPage,
	}
	for _, n := range q.Organization.ProjectsV2.Nodes {
		out.Projects = append(out.Projects, Project{
			ID:          fmt.Sprint(n.ID),
			Number:      int(n.Number),
			Title:       string(n.Title),
			URL:         n.URL.String(),
		})
	}
	return out, nil
}

// authTransport is a simple http.RoundTripper that injects the GitHub token
// (matches patterns used in other MCP Go codebases)
type authTransport struct {
	token string
}

func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+t.token)
	return http.DefaultTransport.RoundTrip(req)
}


// ListUserProjects lists projects for a user using the provided githubv4.Client.
// If client is nil, a default client is created using GITHUB_TOKEN from environment.
func ListUserProjects(ctx context.Context, in *ListUserProjectsInput, client *ghv4.Client) (*ListOrganizationProjectsOutput, error) {
	if in.User == "" {
		return nil, errors.New("user is required")
	}

	if client == nil {
		token := os.Getenv("GITHUB_PERSONAL_ACCESS_TOKEN")
		if token == "" {
			return nil, errors.New("GITHUB_PERSONAL_ACCESS_TOKEN not set")
		}
		client = ghv4.NewClient(&http.Client{Transport: &authTransport{token: token}})
	}

	var q struct {
		User struct {
			ProjectsV2 struct {
				Nodes []struct {
					ID          ghv4.ID
					Number      ghv4.Int
					Title       ghv4.String
					URL         ghv4.URI
				} `graphql:"nodes"`
				PageInfo struct {
					EndCursor   ghv4.String
					HasNextPage bool
				}
			} `graphql:"projectsV2(first: $first, after: $after)"`
		} `graphql:"user(login: $login)"`
	}
	vars := map[string]interface{}{
		"login": ghv4.String(in.User),
		"first": ghv4.Int(in.First),
		"after": ghv4.String(in.After),
	}

	err := client.Query(ctx, &q, vars)
	if err != nil {
		return nil, fmt.Errorf("github graphql error: %w", err)
	}

	out := &ListOrganizationProjectsOutput{
		Projects:    []Project{},
		EndCursor:   string(q.User.ProjectsV2.PageInfo.EndCursor),
		HasNextPage: q.User.ProjectsV2.PageInfo.HasNextPage,
	}
	for _, n := range q.User.ProjectsV2.Nodes {
		out.Projects = append(out.Projects, Project{
			ID:          fmt.Sprint(n.ID),
			Number:      int(n.Number),
			Title:       string(n.Title),
			URL:         n.URL.String(),
		})
	}
	return out, nil
}


// GetProject fetches a project by owner and number using the provided githubv4.Client.
// If client is nil, a default client is created using GITHUB_TOKEN from environment.
func GetProject(ctx context.Context, in *GetProjectInput, client *ghv4.Client) (*Project, error) {
	if in.Owner == "" || in.Number == 0 {
		return nil, errors.New("owner and number are required")
	}

	if client == nil {
		token := os.Getenv("GITHUB_PERSONAL_ACCESS_TOKEN")
		if token == "" {
			return nil, errors.New("GITHUB_PERSONAL_ACCESS_TOKEN not set")
		}
		client = ghv4.NewClient(&http.Client{Transport: &authTransport{token: token}})
	}

	var q struct {
		Organization *struct {
			ProjectV2 *struct {
				ID          ghv4.ID
				Number      ghv4.Int
				Title       ghv4.String
				URL         ghv4.URI
			} `graphql:"projectV2(number: $number)"`
		} `graphql:"organization(login: $owner)"`
		User *struct {
			ProjectV2 *struct {
				ID          ghv4.ID
				Number      ghv4.Int
				Title       ghv4.String
				URL         ghv4.URI
			} `graphql:"projectV2(number: $number)"`
		} `graphql:"user(login: $owner)"`
	}
	vars := map[string]interface{}{
		"owner":  ghv4.String(in.Owner),
		"number": ghv4.Int(in.Number),
	}

	err := client.Query(ctx, &q, vars)
	if err != nil {
		return nil, fmt.Errorf("github graphql error: %w", err)
	}

	var p *struct {
		ID     ghv4.ID
		Number ghv4.Int
		Title  ghv4.String
		URL    ghv4.URI
	}
	if q.Organization != nil && q.Organization.ProjectV2 != nil {
		p = q.Organization.ProjectV2
	} else if q.User != nil && q.User.ProjectV2 != nil {
		p = q.User.ProjectV2
	} else {
		return nil, errors.New("project not found")
	}

	return &Project{
		ID:     fmt.Sprint(p.ID),
		Number: int(p.Number),
		Title:  string(p.Title),
		URL:    p.URL.String(),
	}, nil
}


// GetProjectItems fetches project items using the provided githubv4.Client.
// If client is nil, a default client is created using GITHUB_TOKEN from environment.
func GetProjectItems(ctx context.Context, in *GetProjectItemsInput, client *ghv4.Client) (*GetProjectItemsOutput, error) {
	if in.ProjectID == "" {
		return nil, errors.New("projectID is required")
	}

	if client == nil {
		token := os.Getenv("GITHUB_PERSONAL_ACCESS_TOKEN")
		if token == "" {
			return nil, errors.New("GITHUB_PERSONAL_ACCESS_TOKEN not set")
		}
		client = ghv4.NewClient(&http.Client{Transport: &authTransport{token: token}})
	}

	var q struct {
		Node struct {
			ProjectV2 struct {
				Items struct {
					Nodes []struct {
						ID      ghv4.ID
						Content struct {
							Typename string     `graphql:"__typename"`
							ID       ghv4.ID    `graphql:"id"`
							Title    ghv4.String `graphql:"title"`
							State    ghv4.String `graphql:"state"`
							URL      ghv4.URI   `graphql:"url"`
						} `graphql:"content"`
					} `graphql:"nodes"`
					PageInfo struct {
						EndCursor   ghv4.String
						HasNextPage bool
					}
				} `graphql:"items(first: $first, after: $after)"`
			} `graphql:"... on ProjectV2"`
		} `graphql:"node(id: $id)"`
	}
	vars := map[string]interface{}{
		"id":    ghv4.ID(in.ProjectID),
		"first": ghv4.Int(in.First),
		"after": ghv4.String(in.After),
	}

	err := client.Query(ctx, &q, vars)
	if err != nil {
		return nil, fmt.Errorf("github graphql error: %w", err)
	}

	out := &GetProjectItemsOutput{
		Items:      []ProjectItem{},
		EndCursor:  string(q.Node.ProjectV2.Items.PageInfo.EndCursor),
		HasNextPage: q.Node.ProjectV2.Items.PageInfo.HasNextPage,
	}
	for _, n := range q.Node.ProjectV2.Items.Nodes {
		out.Items = append(out.Items, ProjectItem{
			ID:          fmt.Sprint(n.ID),
			ContentID:   fmt.Sprint(n.Content.ID),
			ContentType: n.Content.Typename,
			Title:       string(n.Content.Title),
			URL:         n.Content.URL.String(),
		})
	}
	return out, nil
}


// CreateProject creates a new project using the provided githubv4.Client.
// If client is nil, a default client is created using GITHUB_TOKEN from environment.
func CreateProject(ctx context.Context, in *CreateProjectInput, client *ghv4.Client) (*Project, error) {
	if in.Owner == "" || in.Title == "" {
		return nil, errors.New("owner and title are required")
	}

	if client == nil {
		token := os.Getenv("GITHUB_PERSONAL_ACCESS_TOKEN")
		if token == "" {
			return nil, errors.New("GITHUB_PERSONAL_ACCESS_TOKEN not set")
		}
		client = ghv4.NewClient(&http.Client{Transport: &authTransport{token: token}})
	}

	// Lookup owner ID (org or user)
	var ownerQ struct {
		Organization *struct{ ID ghv4.ID } `graphql:"organization(login: $login)"`
		User         *struct{ ID ghv4.ID } `graphql:"user(login: $login)"`
	}
	ownerVars := map[string]interface{}{"login": ghv4.String(in.Owner)}
	if err := client.Query(ctx, &ownerQ, ownerVars); err != nil {
		return nil, fmt.Errorf("owner lookup failed: %w", err)
	}
	var ownerID ghv4.ID
	if ownerQ.Organization != nil {
		ownerID = ownerQ.Organization.ID
	} else if ownerQ.User != nil {
		ownerID = ownerQ.User.ID
	} else {
		return nil, errors.New("owner not found")
	}

	type createProjectInput struct {
		OwnerID          ghv4.ID     `json:"ownerId"`
		Title            ghv4.String `json:"title"`
		ShortDescription ghv4.String `json:"shortDescription,omitempty"`
	}
	input := createProjectInput{
		OwnerID: ownerID,
		Title:   ghv4.String(in.Title),
	}
	if in.Description != "" {
		input.ShortDescription = ghv4.String(in.Description)
	}

	var m struct {
		CreateProjectV2 struct {
			ProjectV2 struct {
				ID     ghv4.ID
				Number ghv4.Int
				Title  ghv4.String
				URL    ghv4.URI
			}
		} `graphql:"createProjectV2(input: $input)"`
	}
	if err := client.Mutate(ctx, &m, input, nil); err != nil {
		return nil, fmt.Errorf("github graphql error: %w", err)
	}
	p := m.CreateProjectV2.ProjectV2
	return &Project{
		ID:     fmt.Sprint(p.ID),
		Number: int(p.Number),
		Title:  string(p.Title),
		URL:    p.URL.String(),
	}, nil
}


// AddProjectItem adds an item to a project using the provided githubv4.Client.
// If client is nil, a default client is created using GITHUB_TOKEN from environment.
func AddProjectItem(ctx context.Context, in *AddProjectItemInput, client *ghv4.Client) (*AddProjectItemOutput, error) {
	if in.ProjectID == "" || in.ContentID == "" {
		return nil, errors.New("projectID and contentID are required")
	}

	if client == nil {
		token := os.Getenv("GITHUB_PERSONAL_ACCESS_TOKEN")
		if token == "" {
			return nil, errors.New("GITHUB_PERSONAL_ACCESS_TOKEN not set")
		}
		client = ghv4.NewClient(&http.Client{Transport: &authTransport{token: token}})
	}

	type addItemInput struct {
		ProjectID ghv4.ID `json:"projectId"`
		ContentID ghv4.ID `json:"contentId"`
	}
	input := addItemInput{
		ProjectID: ghv4.ID(in.ProjectID),
		ContentID: ghv4.ID(in.ContentID),
	}

	var m struct {
		AddProjectV2ItemById struct {
			Item struct {
				ID      ghv4.ID
				Content *struct {
					Typename string     `graphql:"__typename"`
					ID       ghv4.ID    `graphql:"id"`
					Title    ghv4.String `graphql:"title"`
					State    ghv4.String `graphql:"state"`
					URL      ghv4.URI   `graphql:"url"`
				} `graphql:"content"`
			}
		} `graphql:"addProjectV2ItemById(input: $input)"`
	}
	if err := client.Mutate(ctx, &m, input, nil); err != nil {
		return nil, fmt.Errorf("github graphql error: %w", err)
	}

	item := ProjectItem{
		ID: fmt.Sprint(m.AddProjectV2ItemById.Item.ID),
	}
	if m.AddProjectV2ItemById.Item.Content != nil {
		item.ContentID = fmt.Sprint(m.AddProjectV2ItemById.Item.Content.ID)
		item.ContentType = m.AddProjectV2ItemById.Item.Content.Typename
		item.Title = string(m.AddProjectV2ItemById.Item.Content.Title)
		item.State = string(m.AddProjectV2ItemById.Item.Content.State)
		item.URL = m.AddProjectV2ItemById.Item.Content.URL.String()
	}
	return &AddProjectItemOutput{Item: item}, nil
}


// UpdateProjectItemField updates a project item field using the provided githubv4.Client.
// If client is nil, a default client is created using GITHUB_TOKEN from environment.
func UpdateProjectItemField(ctx context.Context, in *UpdateProjectItemFieldInput, client *ghv4.Client) (*UpdateProjectItemFieldOutput, error) {
	if in.ItemID == "" || in.FieldID == "" || in.Value == "" {
		return nil, errors.New("itemID, fieldID, and value are required")
	}

	if client == nil {
		token := os.Getenv("GITHUB_PERSONAL_ACCESS_TOKEN")
		if token == "" {
			return nil, errors.New("GITHUB_PERSONAL_ACCESS_TOKEN not set")
		}
		client = ghv4.NewClient(&http.Client{Transport: &authTransport{token: token}})
	}

	type updateFieldInput struct {
		ProjectID ghv4.ID     `json:"projectId"`
		ItemID    ghv4.ID     `json:"itemId"`
		FieldID   ghv4.ID     `json:"fieldId"`
		Value     ghv4.String `json:"value"`
	}
	input := updateFieldInput{
		ProjectID: ghv4.ID(in.ProjectID),
		ItemID:    ghv4.ID(in.ItemID),
		FieldID:   ghv4.ID(in.FieldID),
		Value:     ghv4.String(in.Value),
	}

	var m struct {
		UpdateProjectV2ItemFieldValue struct {
			ProjectV2Item struct {
				ID ghv4.ID
			}
		} `graphql:"updateProjectV2ItemFieldValue(input: $input)"`
	}
	if err := client.Mutate(ctx, &m, input, nil); err != nil {
		return nil, fmt.Errorf("github graphql error: %w", err)
	}

	item := ProjectItem{ID: fmt.Sprint(m.UpdateProjectV2ItemFieldValue.ProjectV2Item.ID)}
	return &UpdateProjectItemFieldOutput{Item: item}, nil
}

