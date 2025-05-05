package github

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"github.com/shurcooL/githubv4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)


func TestListOrganizationProjects(t *testing.T) {
	tests := []struct {
		name        string
		input       *ListOrganizationProjectsInput
		mockHandler http.HandlerFunc
		wantErr     bool
		wantCount   int
	}{
		{
			name:    "missing organization",
			input:   &ListOrganizationProjectsInput{},
			wantErr: true,
		},
		{
			name:  "success",
			input: &ListOrganizationProjectsInput{Organization: "test-org"},
			mockHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				w.Write([]byte(`{"data":{"organization":{"projectsV2":{"nodes":[{"id":"1","number":1,"title":"Proj1","url":"http://example.com/p1"},{"id":"2","number":2,"title":"Proj2","url":"http://example.com/p2"}],"pageInfo":{"endCursor":"abc","hasNextPage":false}}}}}`))
			},
			wantErr:   false,
			wantCount: 2,
		},
		// Add more cases: API error, pagination, etc.
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var server *httptest.Server
			if tc.mockHandler != nil {
				server = httptest.NewServer(tc.mockHandler)
				defer server.Close()
			}

			httpClient := &http.Client{}
			if server != nil {
				httpClient = server.Client()
			}
			var ghClient *githubv4.Client
			if server != nil {
				ghClient = githubv4.NewEnterpriseClient(server.URL, httpClient)
			} else {
				ghClient = githubv4.NewClient(httpClient)
			}

			// Assume ListOrganizationProjects now accepts a githubv4.Client as a parameter
			out, err := ListOrganizationProjects(context.Background(), tc.input, ghClient)
			if tc.wantErr {
				assert.Error(t, err)
				assert.Nil(t, out)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, out)
				assert.Len(t, out.Projects, tc.wantCount)
			}
		})
	}
}

// Integration tests (real API) go in a separate section, skipped by default.


func TestListUserProjects(t *testing.T) {
	tests := []struct {
		name        string
		input       *ListUserProjectsInput
		mockHandler http.HandlerFunc
		wantErr     bool
		wantCount   int
	}{
		{
			name:    "missing user",
			input:   &ListUserProjectsInput{},
			wantErr: true,
		},
		{
			name:  "success",
			input: &ListUserProjectsInput{User: "test-user"},
			mockHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				w.Write([]byte(`{"data":{"user":{"projectsV2":{"nodes":[{"id":"1","number":1,"title":"Proj1","url":"http://example.com/p1"}],"pageInfo":{"endCursor":"abc","hasNextPage":false}}}}}`))
			},
			wantErr:   false,
			wantCount: 1,
		},
	}


	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var server *httptest.Server
			if tc.mockHandler != nil {
				server = httptest.NewServer(tc.mockHandler)
				defer server.Close()
			}
			httpClient := &http.Client{}
			if server != nil {
				httpClient = server.Client()
			}
			var ghClient *githubv4.Client
			if server != nil {
				ghClient = githubv4.NewEnterpriseClient(server.URL, httpClient)
			} else {
				ghClient = githubv4.NewClient(httpClient)
			}
			out, err := ListUserProjects(context.Background(), tc.input, ghClient)
			if tc.wantErr {
				assert.Error(t, err)
				assert.Nil(t, out)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, out)
				assert.Len(t, out.Projects, tc.wantCount)
			}
		})
	}
}


func TestGetProject(t *testing.T) {
	tests := []struct {
		name        string
		input       *GetProjectInput
		mockHandler http.HandlerFunc
		wantErr     bool
		wantID      string
	}{
		{
			name:    "missing owner or number",
			input:   &GetProjectInput{},
			wantErr: true,
		},
		{
			name:  "success",
			input: &GetProjectInput{Owner: "test-owner", Number: 123},
			mockHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				w.Write([]byte(`{"data":{"organization":{"projectV2":{"id":"proj123","title":"Test Project","number":123,"url":"http://example.com/project"}}}}`))
			},
			wantErr: false,
			wantID:  "proj123",
		},
		// Add more cases: API error, not found, etc.
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var server *httptest.Server
			if tc.mockHandler != nil {
				server = httptest.NewServer(tc.mockHandler)
				defer server.Close()
			}
			httpClient := &http.Client{}
			if server != nil {
				httpClient = server.Client()
			}
			var ghClient *githubv4.Client
			if server != nil {
				ghClient = githubv4.NewEnterpriseClient(server.URL, httpClient)
			} else {
				ghClient = githubv4.NewClient(httpClient)
			}
			out, err := GetProject(context.Background(), tc.input, ghClient)
			if tc.wantErr {
				assert.Error(t, err)
				assert.Nil(t, out)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, out)
				assert.Equal(t, tc.wantID, out.ID)
			}
		})
	}
}

func TestGetProjectItems(t *testing.T) {
	tests := []struct {
		name        string
		input       *GetProjectItemsInput
		mockHandler http.HandlerFunc
		wantErr     bool
		wantCount   int
	}{
		{
			name:    "missing project_id",
			input:   &GetProjectItemsInput{},
			wantErr: true,
		},
		{
			name:  "success",
			input: &GetProjectItemsInput{ProjectID: "proj123"},
			mockHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				w.Write([]byte(`{"data":{"node":{"items":{"nodes":[{"id":"item1","content":{"__typename":"Issue","id":"c1","title":"Issue1","url":"http://example.com/i1"}}],"pageInfo":{"endCursor":"abc","hasNextPage":false}}}}}`))
			},
			wantErr:   false,
			wantCount: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var server *httptest.Server
			if tc.mockHandler != nil {
				server = httptest.NewServer(tc.mockHandler)
				defer server.Close()
			}
			httpClient := &http.Client{}
			if server != nil {
				httpClient = server.Client()
			}
			var ghClient *githubv4.Client
			if server != nil {
				ghClient = githubv4.NewEnterpriseClient(server.URL, httpClient)
			} else {
				ghClient = githubv4.NewClient(httpClient)
			}
			out, err := GetProjectItems(context.Background(), tc.input, ghClient)
			if tc.wantErr {
				assert.Error(t, err)
				assert.Nil(t, out)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, out)
				assert.Len(t, out.Items, tc.wantCount)
			}
		})
	}
}

func TestOwnerResolutionInCreateProject(t *testing.T) {
	tests := []struct {
		name        string
		input       *CreateProjectInput
		mockHandler http.HandlerFunc
		wantErr     bool
		wantID      string
		wantErrMsg  string
	}{
		{
			name:  "owner is organization",
			input: &CreateProjectInput{Owner: "org-login", Title: "Project for Org"},
			mockHandler: func(w http.ResponseWriter, r *http.Request) {
				var buf bytes.Buffer
				_, _ = buf.ReadFrom(r.Body)
				body := buf.String()
				if strings.Contains(body, "organization") {
					w.WriteHeader(200)
					w.Write([]byte(`{"data":{"organization":{"id":"org123"}}}`))
				} else if strings.Contains(body, "createProjectV2") {
					w.WriteHeader(200)
					w.Write([]byte(`{"data":{"createProjectV2":{"projectV2":{"id":"projOrg","title":"Project for Org","number":1,"url":"http://example.com/orgproject"}}}}`))
				} else {
					w.WriteHeader(400)
					w.Write([]byte(`{"error":"unexpected request"}`))
				}
			},
			wantErr: false,
			wantID:  "projOrg",
		},
		{
			name:  "owner is user",
			input: &CreateProjectInput{Owner: "user-login", Title: "Project for User"},
			mockHandler: func(w http.ResponseWriter, r *http.Request) {
				var buf bytes.Buffer
				_, _ = buf.ReadFrom(r.Body)
				body := buf.String()
				if strings.Contains(body, "organization") {
					w.WriteHeader(200)
					w.Write([]byte(`{"data":{"organization":null}}`)) // Not an org
				} else if strings.Contains(body, "user") {
					w.WriteHeader(200)
					w.Write([]byte(`{"data":{"user":{"id":"user123"}}}`))
				} else if strings.Contains(body, "createProjectV2") {
					w.WriteHeader(200)
					w.Write([]byte(`{"data":{"createProjectV2":{"projectV2":{"id":"projUser","title":"Project for User","number":2,"url":"http://example.com/userproject"}}}}`))
				} else {
					w.WriteHeader(400)
					w.Write([]byte(`{"error":"unexpected request"}`))
				}
			},
			wantErr: false,
			wantID:  "projUser",
		},
		{
			name:  "owner is neither user nor org",
			input: &CreateProjectInput{Owner: "ghost-login", Title: "Project for Ghost"},
			mockHandler: func(w http.ResponseWriter, r *http.Request) {
				var buf bytes.Buffer
				_, _ = buf.ReadFrom(r.Body)
				body := buf.String()
				if strings.Contains(body, "organization") {
					w.WriteHeader(200)
					w.Write([]byte(`{"data":{"organization":null}}`))
				} else if strings.Contains(body, "user") {
					w.WriteHeader(200)
					w.Write([]byte(`{"data":{"user":null}}`))
				} else {
					w.WriteHeader(400)
					w.Write([]byte(`{"error":"unexpected request"}`))
				}
			},
			wantErr:    true,
			wantErrMsg: "owner not found",
		},
		{
			name:  "owner ambiguous (org and user both exist)",
			input: &CreateProjectInput{Owner: "ambiguous-login", Title: "Ambiguous Project"},
			mockHandler: func(w http.ResponseWriter, r *http.Request) {
				var buf bytes.Buffer
				_, _ = buf.ReadFrom(r.Body)
				body := buf.String()
				if strings.Contains(body, "organization") {
					w.WriteHeader(200)
					w.Write([]byte(`{"data":{"organization":{"id":"orgAmbig"}}}`))
				} else if strings.Contains(body, "user") {
					w.WriteHeader(200)
					w.Write([]byte(`{"data":{"user":{"id":"userAmbig"}}}`))
				} else if strings.Contains(body, "createProjectV2") {
					w.WriteHeader(200)
					w.Write([]byte(`{"data":{"createProjectV2":{"projectV2":{"id":"projAmbig","title":"Ambiguous Project","number":3,"url":"http://example.com/ambigproject"}}}}`))
				} else {
					w.WriteHeader(400)
					w.Write([]byte(`{"error":"unexpected request"}`))
				}
			},
			wantErr: false, // Should prefer org or document behavior
			wantID:  "projAmbig",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var server *httptest.Server
			if tc.mockHandler != nil {
				server = httptest.NewServer(tc.mockHandler)
				defer server.Close()
			}
			httpClient := &http.Client{}
			if server != nil {
				httpClient = server.Client()
			}
			var ghClient *githubv4.Client
			if server != nil {
				ghClient = githubv4.NewEnterpriseClient(server.URL, httpClient)
			} else {
				ghClient = githubv4.NewClient(httpClient)
			}
			out, err := CreateProject(context.Background(), tc.input, ghClient)
			if tc.wantErr {
				assert.Error(t, err)
				if tc.wantErrMsg != "" {
					assert.Contains(t, err.Error(), tc.wantErrMsg)
				}
				assert.Nil(t, out)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, out)
				assert.Equal(t, tc.wantID, out.ID)
			}
		})
	}
}

func TestCreateProject(t *testing.T) {
	tests := []struct {
		name        string
		input       *CreateProjectInput
		mockHandler http.HandlerFunc
		wantErr     bool
		wantID      string
	}{
		{
			name:    "missing owner/title",
			input:   &CreateProjectInput{},
			wantErr: true,
		},
		{
			name:  "success",
			input: &CreateProjectInput{Owner: "test-owner", Title: "Test Project"},
			mockHandler: func(w http.ResponseWriter, r *http.Request) {
				var buf bytes.Buffer
				_, _ = buf.ReadFrom(r.Body)
				body := buf.String()
				if strings.Contains(body, "organization") || strings.Contains(body, "user") {
					// Owner lookup query
					w.WriteHeader(200)
					w.Write([]byte(`{"data":{"organization":{"id":"owner123"}}}`))
				} else if strings.Contains(body, "createProjectV2") {
					w.WriteHeader(200)
					w.Write([]byte(`{"data":{"createProjectV2":{"projectV2":{"id":"proj456","title":"Test Project","number":456,"url":"http://example.com/project"}}}}`))
				} else {
					w.WriteHeader(400)
					w.Write([]byte(`{"error":"unexpected request"}`))
				}
			},
			wantErr: false,
			wantID:  "proj456",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var server *httptest.Server
			if tc.mockHandler != nil {
				server = httptest.NewServer(tc.mockHandler)
				defer server.Close()
			}
			httpClient := &http.Client{}
			if server != nil {
				httpClient = server.Client()
			}
			var ghClient *githubv4.Client
			if server != nil {
				ghClient = githubv4.NewEnterpriseClient(server.URL, httpClient)
			} else {
				ghClient = githubv4.NewClient(httpClient)
			}
			out, err := CreateProject(context.Background(), tc.input, ghClient)
			if tc.wantErr {
				assert.Error(t, err)
				assert.Nil(t, out)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, out)
				assert.Equal(t, tc.wantID, out.ID)
			}
		})
	}
}

func TestAddProjectItem(t *testing.T) {
	tests := []struct {
		name        string
		input       *AddProjectItemInput
		mockHandler http.HandlerFunc
		wantErr     bool
		wantID      string
	}{
		{
			name:    "missing project_id/content_id",
			input:   &AddProjectItemInput{},
			wantErr: true,
		},
		{
			name:  "success",
			input: &AddProjectItemInput{ProjectID: "proj123", ContentID: "c1"},
			mockHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				w.Write([]byte(`{"data":{"addProjectV2ItemById":{"item":{"id":"item2","content":{"__typename":"Issue","id":"c1","title":"Issue1","url":"http://example.com/i1"}}}}}`))
			},
			wantErr: false,
			wantID:  "item2",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var server *httptest.Server
			if tc.mockHandler != nil {
				server = httptest.NewServer(tc.mockHandler)
				defer server.Close()
			}
			httpClient := &http.Client{}
			if server != nil {
				httpClient = server.Client()
			}
			var ghClient *githubv4.Client
			if server != nil {
				ghClient = githubv4.NewEnterpriseClient(server.URL, httpClient)
			} else {
				ghClient = githubv4.NewClient(httpClient)
			}
			out, err := AddProjectItem(context.Background(), tc.input, ghClient)
			if tc.wantErr {
				assert.Error(t, err)
				assert.Nil(t, out)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, out)
				assert.Equal(t, tc.wantID, out.Item.ID)
			}
		})
	}
}

func TestUpdateProjectItemField(t *testing.T) {
	tests := []struct {
		name        string
		input       *UpdateProjectItemFieldInput
		mockHandler http.HandlerFunc
		wantErr     bool
		wantID      string
	}{
		{
			name:    "missing required fields",
			input:   &UpdateProjectItemFieldInput{},
			wantErr: true,
		},
		{
			name:  "success",
			input: &UpdateProjectItemFieldInput{ItemID: "item2", FieldID: "field1", Value: "new value"},
			mockHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				w.Write([]byte(`{"data":{"updateProjectV2ItemFieldValue":{"projectV2Item":{"id":"item2"}}}}`))
			},
			wantErr: false,
			wantID:  "item2",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var server *httptest.Server
			if tc.mockHandler != nil {
				server = httptest.NewServer(tc.mockHandler)
				defer server.Close()
			}
			httpClient := &http.Client{}
			if server != nil {
				httpClient = server.Client()
			}
			var ghClient *githubv4.Client
			if server != nil {
				ghClient = githubv4.NewEnterpriseClient(server.URL, httpClient)
			} else {
				ghClient = githubv4.NewClient(httpClient)
			}
			out, err := UpdateProjectItemField(context.Background(), tc.input, ghClient)
			if tc.wantErr {
				assert.Error(t, err)
				assert.Nil(t, out)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, out)
				assert.Equal(t, tc.wantID, out.Item.ID)
			}
		})
	}
}
