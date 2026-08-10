# code-review skill

Posts the Copilot reconciliation follow-up and suggested Rootstock lessons as
separate pull request conversation comments after a review that used
`.github/instructions/pull-request-review.instructions.md`.

The directory is named `code-review` because Copilot code review is documented to
favour review-focused skill directory names when deciding which skills apply.

## Required repository MCP settings

Configure these under **Settings → Copilot → MCP servers**. This JSON is not
stored in the repository.

`tools` is a strict allowlist, and the `/readonly` endpoint strips write tools
server-side. There is no "keep all read tools and add this one write tool"
switch, so the entire default read-only tool set has to be enumerated alongside
`add_issue_comment` on the non-readonly endpoint.

```json
{
  "mcpServers": {
    "github-mcp-server": {
      "type": "http",
      "url": "https://api.githubcopilot.com/mcp/",
      "tools": [
        "actions_get",
        "actions_list",
        "get_code_scanning_alert",
        "get_commit",
        "get_discussion",
        "get_discussion_comments",
        "get_file_contents",
        "get_job_logs",
        "get_label",
        "get_latest_release",
        "get_release_by_tag",
        "get_secret_scanning_alert",
        "get_tag",
        "issue_read",
        "list_branches",
        "list_code_scanning_alerts",
        "list_commits",
        "list_discussion_categories",
        "list_discussions",
        "list_issue_fields",
        "list_issue_types",
        "list_issues",
        "list_label",
        "list_pull_requests",
        "list_releases",
        "list_repository_collaborators",
        "list_secret_scanning_alerts",
        "list_tags",
        "pull_request_read",
        "search_code",
        "search_commits",
        "search_issues",
        "search_pull_requests",
        "search_repositories",
        "search_users",
        "add_issue_comment"
      ],
      "headers": {
        "X-MCP-Toolsets": "actions,code_security,discussions,issues,labels,pull_requests,repos,secret_protection,users"
      }
    }
  }
}
```

That list is the default read-only surface plus exactly one addition,
`add_issue_comment`. A tool is only available when its toolset is enabled *and*
it appears in `tools`, so the toolsets header has to stay wide enough to cover
every allowlisted tool.

Also keep **Allow Copilot to use MCP tools when reviewing pull requests**
enabled for the repository.

Never use `"tools": ["*"]` against the non-readonly URL; that also grants merge,
push, and branch-write tools. If `add_issue_comment` is unavailable or write
calls fail, the skill falls back to leaving the artifacts in the review session
instead of claiming they were posted.

## Related files

- Procedure: `.github/instructions/pull-request-review.instructions.md` (Step 7)
- Lesson schema: `references/lesson-schema.md`
