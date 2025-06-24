package github

import (
	"context"
	"fmt"

	"github.com/github/github-mcp-server/pkg/translations"
	"github.com/go-viper/mapstructure/v2"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/shurcooL/githubv4"
)

// ListProjectV2 creates a tool to list Projects v2 using GraphQL
func ListProjectV2(getGQLClient GetGQLClientFn, t translations.TranslationHelperFunc) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("list_projects_v2",
			mcp.WithDescription(t("TOOL_LIST_PROJECTS_V2_DESCRIPTION", "List Projects v2 for an organization or user.")),
			mcp.WithToolAnnotation(mcp.ToolAnnotation{
				Title:        t("TOOL_LIST_PROJECTS_V2_USER_TITLE", "List Projects v2"),
				ReadOnlyHint: ToBoolPtr(true),
			}),
			mcp.WithString("owner",
				mcp.Required(),
				mcp.Description("Organization or user login"),
			),
			mcp.WithString("owner_type",
				mcp.Required(),
				mcp.Description("Owner type: organization or user"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var params struct {
				Owner     string `mapstructure:"owner"`
				OwnerType string `mapstructure:"owner_type"`
			}
			if err := mapstructure.Decode(request.Params.Arguments, &params); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			client, err := getGQLClient(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to get GitHub GraphQL client: %w", err)
			}

			var query struct {
				Organization struct {
					ProjectsV2 struct {
						Nodes []struct {
							ID     githubv4.ID
							Title  githubv4.String
							Number githubv4.Int
							URL    githubv4.String
							State  githubv4.String
						}
					} `graphql:"projectsV2(first: 20)"`
				} `graphql:"organization(login: $login)" json:",omitempty"`
				User struct {
					ProjectsV2 struct {
						Nodes []struct {
							ID     githubv4.ID
							Title  githubv4.String
							Number githubv4.Int
							URL    githubv4.String
							State  githubv4.String
						}
					} `graphql:"projectsV2(first: 20)"`
				} `graphql:"user(login: $login)" json:",omitempty"`
			}

			variables := map[string]interface{}{
				"login": githubv4.String(params.Owner),
			}

			err = client.Query(ctx, &query, variables)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to query projects v2: %v", err)), nil
			}

			if params.OwnerType == "organization" {
				return mcp.NewToolResultText(formatProjectsV2Response(query.Organization.ProjectsV2.Nodes)), nil
			} else {
				return mcp.NewToolResultText(formatProjectsV2Response(query.User.ProjectsV2.Nodes)), nil
			}
		}
}

// GetProjectV2 creates a tool to get details of a specific Projects v2
func GetProjectV2(getGQLClient GetGQLClientFn, t translations.TranslationHelperFunc) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("get_project_v2",
			mcp.WithDescription(t("TOOL_GET_PROJECT_V2_DESCRIPTION", "Get details of a specific Projects v2.")),
			mcp.WithToolAnnotation(mcp.ToolAnnotation{
				Title:        t("TOOL_GET_PROJECT_V2_USER_TITLE", "Get Projects v2 details"),
				ReadOnlyHint: ToBoolPtr(true),
			}),
			mcp.WithString("project_id",
				mcp.Required(),
				mcp.Description("Projects v2 ID"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var params struct {
				ProjectID string `mapstructure:"project_id"`
			}
			if err := mapstructure.Decode(request.Params.Arguments, &params); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			client, err := getGQLClient(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to get GitHub GraphQL client: %w", err)
			}

			var query struct {
				Node struct {
					ProjectV2 struct {
						ID     githubv4.ID
						Title  githubv4.String
						Number githubv4.Int
						URL    githubv4.String
						State  githubv4.String
						Public githubv4.Boolean
						Fields struct {
							Nodes []struct {
								ID   githubv4.ID
								Name githubv4.String
							}
						} `graphql:"fields(first: 20)"`
					} `graphql:"... on ProjectV2"`
				} `graphql:"node(id: $id)"`
			}

			variables := map[string]interface{}{
				"id": githubv4.ID(params.ProjectID),
			}

			err = client.Query(ctx, &query, variables)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to query project v2: %v", err)), nil
			}

			return mcp.NewToolResultText(formatSingleProjectV2Response(query.Node.ProjectV2)), nil
		}
}

// AddIssueToProjectV2 creates a tool to add an issue to a Projects v2
func AddIssueToProjectV2(getGQLClient GetGQLClientFn, t translations.TranslationHelperFunc) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("add_issue_to_project_v2",
			mcp.WithDescription(t("TOOL_ADD_ISSUE_TO_PROJECT_V2_DESCRIPTION", "Add an issue to a Projects v2.")),
			mcp.WithToolAnnotation(mcp.ToolAnnotation{
				Title:        t("TOOL_ADD_ISSUE_TO_PROJECT_V2_USER_TITLE", "Add issue to Projects v2"),
				ReadOnlyHint: ToBoolPtr(false),
			}),
			mcp.WithString("project_id",
				mcp.Required(),
				mcp.Description("Projects v2 ID"),
			),
			mcp.WithString("issue_id",
				mcp.Required(),
				mcp.Description("Issue ID"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var params struct {
				ProjectID string `mapstructure:"project_id"`
				IssueID   string `mapstructure:"issue_id"`
			}
			if err := mapstructure.Decode(request.Params.Arguments, &params); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			client, err := getGQLClient(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to get GitHub GraphQL client: %w", err)
			}

			var mutation struct {
				AddProjectV2ItemById struct {
					Item struct {
						ID githubv4.ID
					}
				} `graphql:"addProjectV2ItemById(input: $input)"`
			}

			input := map[string]interface{}{
				"projectId": githubv4.ID(params.ProjectID),
				"contentId": githubv4.ID(params.IssueID),
			}

			err = client.Mutate(ctx, &mutation, input, nil)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to add issue to project v2: %v", err)), nil
			}

			return mcp.NewToolResultText(fmt.Sprintf("Successfully added issue to project v2. Item ID: %s", mutation.AddProjectV2ItemById.Item.ID)), nil
		}
}

// GetIssueNodeId creates a tool to get the node ID of an issue (needed for Projects v2)
func GetIssueNodeId(getGQLClient GetGQLClientFn, t translations.TranslationHelperFunc) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("get_issue_node_id",
			mcp.WithDescription(t("TOOL_GET_ISSUE_NODE_ID_DESCRIPTION", "Get the node ID of an issue for use with Projects v2.")),
			mcp.WithToolAnnotation(mcp.ToolAnnotation{
				Title:        t("TOOL_GET_ISSUE_NODE_ID_USER_TITLE", "Get issue node ID"),
				ReadOnlyHint: ToBoolPtr(true),
			}),
			mcp.WithString("owner",
				mcp.Required(),
				mcp.Description("Repository owner"),
			),
			mcp.WithString("repo",
				mcp.Required(),
				mcp.Description("Repository name"),
			),
			mcp.WithNumber("issue_number",
				mcp.Required(),
				mcp.Description("Issue number"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var params struct {
				Owner       string `mapstructure:"owner"`
				Repo        string `mapstructure:"repo"`
				IssueNumber int32  `mapstructure:"issue_number"`
			}
			if err := mapstructure.Decode(request.Params.Arguments, &params); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			client, err := getGQLClient(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to get GitHub GraphQL client: %w", err)
			}

			var query struct {
				Repository struct {
					Issue struct {
						ID     githubv4.ID
						Number githubv4.Int
						Title  githubv4.String
					} `graphql:"issue(number: $issue_number)"`
				} `graphql:"repository(owner: $owner, name: $name)"`
			}

			variables := map[string]interface{}{
				"owner":        githubv4.String(params.Owner),
				"name":         githubv4.String(params.Repo),
				"issue_number": githubv4.Int(params.IssueNumber),
			}

			err = client.Query(ctx, &query, variables)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to get issue node ID: %v", err)), nil
			}

			result := fmt.Sprintf("Issue #%d: %s\nNode ID: %s", 
				query.Repository.Issue.Number, 
				query.Repository.Issue.Title, 
				query.Repository.Issue.ID)

			return mcp.NewToolResultText(result), nil
		}
}

// Helper functions for formatting responses
func formatProjectsV2Response(projects []struct {
	ID     githubv4.ID
	Title  githubv4.String
	Number githubv4.Int
	URL    githubv4.String
	State  githubv4.String
}) string {
	if len(projects) == 0 {
		return "No Projects v2 found."
	}

	result := fmt.Sprintf("Found %d Projects v2:\n\n", len(projects))
	for _, project := range projects {
		result += fmt.Sprintf("- **%s** (#%d)\n", project.Title, project.Number)
		result += fmt.Sprintf("  - ID: %s\n", project.ID)
		result += fmt.Sprintf("  - State: %s\n", project.State)
		result += fmt.Sprintf("  - URL: %s\n\n", project.URL)
	}
	return result
}

func formatSingleProjectV2Response(project struct {
	ID     githubv4.ID
	Title  githubv4.String
	Number githubv4.Int
	URL    githubv4.String
	State  githubv4.String
	Public githubv4.Boolean
	Fields struct {
		Nodes []struct {
			ID   githubv4.ID
			Name githubv4.String
		}
	} `graphql:"fields(first: 20)"`
}) string {
	result := fmt.Sprintf("**%s** (#%d)\n", project.Title, project.Number)
	result += fmt.Sprintf("- ID: %s\n", project.ID)
	result += fmt.Sprintf("- State: %s\n", project.State)
	result += fmt.Sprintf("- Public: %t\n", project.Public)
	result += fmt.Sprintf("- URL: %s\n", project.URL)
	
	if len(project.Fields.Nodes) > 0 {
		result += "\n**Custom Fields:**\n"
		for _, field := range project.Fields.Nodes {
			result += fmt.Sprintf("- %s (ID: %s)\n", field.Name, field.ID)
		}
	}
	
	return result
}