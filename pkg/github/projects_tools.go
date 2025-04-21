package github

import (
	"context"
	"encoding/json"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/github/github-mcp-server/pkg/translations"
)

// MCP tool factory for listing organization projects
func ListOrganizationProjectsTool(getClient GetGraphQLClientFn, t translations.TranslationHelperFunc) (mcp.Tool, server.ToolHandlerFunc) {
	tool := mcp.NewTool(
		"list_organization_projects",
		mcp.WithDescription("List Projects for an organization"),
		mcp.WithString("organization", mcp.Required(), mcp.Description("The organization login")),
		mcp.WithNumber("first", mcp.Description("Max number of projects to return")),
		mcp.WithString("after", mcp.Description("Cursor for pagination")),
	)
	handler := func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		client, err := getClient(ctx)
		if err != nil {
			return nil, err
		}

		organization, err := requiredParam[string](req, "organization")
		if err != nil {
			return nil, err
		}
		first, _ := requiredParam[float64](req, "first") // optional
		after, _ := requiredParam[string](req, "after") // optional
		input := &ListOrganizationProjectsInput{
			Organization: organization,
			First:        int(first),
			After:        after,
		}
		out, err := ListOrganizationProjects(ctx, input, client)
		if err != nil {
			return nil, err
		}
		b, _ := json.Marshal(out)
		return mcp.NewToolResultText(string(b)), nil
	}
	return tool, handler
}

// MCP tool factory for listing user projects
func ListUserProjectsTool(getClient GetGraphQLClientFn, t translations.TranslationHelperFunc) (mcp.Tool, server.ToolHandlerFunc) {
	tool := mcp.NewTool(
		"list_user_projects",
		mcp.WithDescription("List Projects for a user"),
		mcp.WithString("user", mcp.Required(), mcp.Description("The user login")),
		mcp.WithNumber("first", mcp.Description("Max number of projects to return")),
		mcp.WithString("after", mcp.Description("Cursor for pagination")),
	)
	handler := func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		client, err := getClient(ctx)
		if err != nil {
			return nil, err
		}

		user, err := requiredParam[string](req, "user")
		if err != nil {
			return nil, err
		}
		first, _ := requiredParam[float64](req, "first") // optional
		after, _ := requiredParam[string](req, "after") // optional
		input := &ListUserProjectsInput{
			User:  user,
			First: int(first),
			After: after,
		}
		out, err := ListUserProjects(ctx, input, client)
		if err != nil {
			return nil, err
		}
		b, _ := json.Marshal(out)
		return mcp.NewToolResultText(string(b)), nil
	}
	return tool, handler
}

// MCP tool factory for getting a project
func GetProjectTool(getClient GetGraphQLClientFn, t translations.TranslationHelperFunc) (mcp.Tool, server.ToolHandlerFunc) {
	tool := mcp.NewTool(
		"get_project",
		mcp.WithDescription("Get a project by owner and number"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("The organization or user login")),
		mcp.WithNumber("number", mcp.Required(), mcp.Description("Project number")),
	)
	handler := func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		client, err := getClient(ctx)
		if err != nil {
			return nil, err
		}

		owner, err := requiredParam[string](req, "owner")
		if err != nil {
			return nil, err
		}
		number, err := requiredParam[float64](req, "number")
		if err != nil {
			return nil, err
		}
		input := &GetProjectInput{
			Owner:  owner,
			Number: int(number),
		}
		out, err := GetProject(ctx, input, client)
		if err != nil {
			return nil, err
		}
		b, _ := json.Marshal(out)
		return mcp.NewToolResultText(string(b)), nil
	}
	return tool, handler
}

// MCP tool factory for getting project items
func GetProjectItemsTool(getClient GetGraphQLClientFn, t translations.TranslationHelperFunc) (mcp.Tool, server.ToolHandlerFunc) {
	tool := mcp.NewTool(
		"get_project_items",
		mcp.WithDescription("Get items for a project"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("Project node ID")),
		mcp.WithNumber("first", mcp.Description("Max number of items to return")),
		mcp.WithString("after", mcp.Description("Cursor for pagination")),
	)
	handler := func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		client, err := getClient(ctx)
		if err != nil {
			return nil, err
		}

		projectID, err := requiredParam[string](req, "project_id")
		if err != nil {
			return nil, err
		}
		first, _ := requiredParam[float64](req, "first") // optional
		after, _ := requiredParam[string](req, "after") // optional
		input := &GetProjectItemsInput{
			ProjectID: projectID,
			First:     int(first),
			After:     after,
		}
		out, err := GetProjectItems(ctx, input, client)
		if err != nil {
			return nil, err
		}
		b, _ := json.Marshal(out)
		return mcp.NewToolResultText(string(b)), nil
	}
	return tool, handler
}

// MCP tool factory for creating a project
func CreateProjectTool(getClient GetGraphQLClientFn, t translations.TranslationHelperFunc) (mcp.Tool, server.ToolHandlerFunc) {
	tool := mcp.NewTool(
		"create_project",
		mcp.WithDescription("Create a new project"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("The organization or user login")),
		mcp.WithString("title", mcp.Required(), mcp.Description("Project title")),
		mcp.WithString("description", mcp.Description("Project description")),
	)
	handler := func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		client, err := getClient(ctx)
		if err != nil {
			return nil, err
		}

		owner, err := requiredParam[string](req, "owner")
		if err != nil {
			return nil, err
		}
		title, err := requiredParam[string](req, "title")
		if err != nil {
			return nil, err
		}
		description, _ := requiredParam[string](req, "description") // optional
		input := &CreateProjectInput{
			Owner:       owner,
			Title:       title,
			Description: description,
		}
		out, err := CreateProject(ctx, input, client)
		if err != nil {
			return nil, err
		}
		b, _ := json.Marshal(out)
		return mcp.NewToolResultText(string(b)), nil
	}
	return tool, handler
}

// MCP tool factory for adding a project item
func AddProjectItemTool(getClient GetGraphQLClientFn, t translations.TranslationHelperFunc) (mcp.Tool, server.ToolHandlerFunc) {
	tool := mcp.NewTool(
		"add_project_item",
		mcp.WithDescription("Add an item to a project"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("Project node ID")),
		mcp.WithString("content_id", mcp.Required(), mcp.Description("Content node ID (issue, PR, etc)")),
	)
	handler := func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		client, err := getClient(ctx)
		if err != nil {
			return nil, err
		}

		projectID, err := requiredParam[string](req, "project_id")
		if err != nil {
			return nil, err
		}
		contentID, err := requiredParam[string](req, "content_id")
		if err != nil {
			return nil, err
		}
		input := &AddProjectItemInput{
			ProjectID: projectID,
			ContentID: contentID,
		}
		out, err := AddProjectItem(ctx, input, client)
		if err != nil {
			return nil, err
		}
		b, _ := json.Marshal(out)
		return mcp.NewToolResultText(string(b)), nil
	}
	return tool, handler
}

// MCP tool factory for updating a project item field
func UpdateProjectItemFieldTool(getClient GetGraphQLClientFn, t translations.TranslationHelperFunc) (mcp.Tool, server.ToolHandlerFunc) {
	tool := mcp.NewTool(
		"update_project_item_field",
		mcp.WithDescription("Update a field on a project item"),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("Project node ID")),
		mcp.WithString("item_id", mcp.Required(), mcp.Description("Item node ID")),
		mcp.WithString("field_id", mcp.Required(), mcp.Description("Field node ID")),
		mcp.WithString("value", mcp.Required(), mcp.Description("New value for the field")),
	)
	handler := func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		client, err := getClient(ctx)
		if err != nil {
			return nil, err
		}

		projectID, err := requiredParam[string](req, "project_id")
		if err != nil {
			return nil, err
		}
		itemID, err := requiredParam[string](req, "item_id")
		if err != nil {
			return nil, err
		}
		fieldID, err := requiredParam[string](req, "field_id")
		if err != nil {
			return nil, err
		}
		value, err := requiredParam[string](req, "value")
		if err != nil {
			return nil, err
		}
		input := &UpdateProjectItemFieldInput{
			ProjectID: projectID,
			ItemID:    itemID,
			FieldID:   fieldID,
			Value:     value,
		}
		out, err := UpdateProjectItemField(ctx, input, client)
		if err != nil {
			return nil, err
		}
		b, _ := json.Marshal(out)
		return mcp.NewToolResultText(string(b)), nil
	}
	return tool, handler
}
