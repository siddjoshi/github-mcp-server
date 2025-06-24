package github

import (
	"testing"

	"github.com/github/github-mcp-server/internal/toolsnaps"
	"github.com/github/github-mcp-server/pkg/translations"
	"github.com/shurcooL/githubv4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ListProjectV2(t *testing.T) {
	// Verify tool definition once
	mockClient := githubv4.NewClient(nil)
	tool, _ := ListProjectV2(stubGetGQLClientFn(mockClient), translations.NullTranslationHelper)
	require.NoError(t, toolsnaps.Test(tool.Name, tool))

	assert.Equal(t, "list_projects_v2", tool.Name)
	assert.NotEmpty(t, tool.Description)
	assert.Contains(t, tool.InputSchema.Properties, "owner")
	assert.Contains(t, tool.InputSchema.Properties, "owner_type")
	assert.ElementsMatch(t, tool.InputSchema.Required, []string{"owner", "owner_type"})
}

func Test_GetProjectV2(t *testing.T) {
	// Verify tool definition once
	mockClient := githubv4.NewClient(nil)
	tool, _ := GetProjectV2(stubGetGQLClientFn(mockClient), translations.NullTranslationHelper)
	require.NoError(t, toolsnaps.Test(tool.Name, tool))

	assert.Equal(t, "get_project_v2", tool.Name)
	assert.NotEmpty(t, tool.Description)
	assert.Contains(t, tool.InputSchema.Properties, "project_id")
	assert.ElementsMatch(t, tool.InputSchema.Required, []string{"project_id"})
}

func Test_GetIssueNodeId(t *testing.T) {
	// Verify tool definition once
	mockClient := githubv4.NewClient(nil)
	tool, _ := GetIssueNodeId(stubGetGQLClientFn(mockClient), translations.NullTranslationHelper)
	require.NoError(t, toolsnaps.Test(tool.Name, tool))

	assert.Equal(t, "get_issue_node_id", tool.Name)
	assert.NotEmpty(t, tool.Description)
	assert.Contains(t, tool.InputSchema.Properties, "owner")
	assert.Contains(t, tool.InputSchema.Properties, "repo")
	assert.Contains(t, tool.InputSchema.Properties, "issue_number")
	assert.ElementsMatch(t, tool.InputSchema.Required, []string{"owner", "repo", "issue_number"})
}

func Test_AddIssueToProjectV2(t *testing.T) {
	// Verify tool definition once
	mockClient := githubv4.NewClient(nil)
	tool, _ := AddIssueToProjectV2(stubGetGQLClientFn(mockClient), translations.NullTranslationHelper)
	require.NoError(t, toolsnaps.Test(tool.Name, tool))

	assert.Equal(t, "add_issue_to_project_v2", tool.Name)
	assert.NotEmpty(t, tool.Description)
	assert.Contains(t, tool.InputSchema.Properties, "project_id")
	assert.Contains(t, tool.InputSchema.Properties, "issue_id")
	assert.ElementsMatch(t, tool.InputSchema.Required, []string{"project_id", "issue_id"})
}

