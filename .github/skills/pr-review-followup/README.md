# pr-review-followup skill

Posts the Copilot reconciliation follow-up and suggested Rootstock lessons as
separate pull request conversation comments after a review that used
`.github/instructions/pull-request-review.instructions.md`.

## Required repository MCP settings

Configure these under **Settings → Copilot → MCP servers**. This JSON is not
stored in the repository.

<<<<<<< Updated upstream
`tools` is a strict allowlist. A single entry with `"tools": ["add_issue_comment"]`
removes every read tool the review procedure depends on. Register two server
entries instead: the stock read-only server, plus a narrow write server that
exposes nothing but the comment tool.

```json
{
  "mcpServers": {
    "github-mcp-server": {
      "type": "http",
      "url": "https://api.githubcopilot.com/mcp/readonly",
      "tools": ["*"],
      "headers": {
        "X-MCP-Toolsets": "repos,issues,users,pull_requests,code_security,secret_protection,actions,web_search"
      }
    },
    "github-mcp-comment": {
      "type": "http",
      "url": "https://api.githubcopilot.com/mcp/",
      "tools": ["add_issue_comment"],
      "headers": {
        "X-MCP-Toolsets": "issues"
      }
    }
  }
}
```

The `/readonly` endpoint strips write tools server-side, so `"tools": ["*"]` on
that entry stays read-only. The second entry is the only place write access is
granted, and it is limited to one tool.

If you would rather run a single server entry, drop `/readonly` and enumerate
every tool you need — the allowlist has no "all read tools plus this one" form:
=======
`tools` is a strict allowlist, and the `/readonly` endpoint strips write tools
server-side. There is no "keep all read tools and add this one write tool"
switch, so the entire default read-only tool set has to be enumerated alongside
`add_issue_comment` on the non-readonly endpoint.
>>>>>>> Stashed changes

```json
{
  "mcpServers": {
    "github-mcp-server": {
      "type": "http",
      "url": "https://api.githubcopilot.com/mcp/",
      "tools": [
<<<<<<< Updated upstream
        "pull_request_read",
        "list_pull_requests",
        "search_pull_requests",
        "issue_read",
        "list_issues",
        "search_issues",
        "get_file_contents",
        "get_commit",
        "list_commits",
        "list_branches",
        "search_code",
        "add_issue_comment"
      ],
      "headers": {
        "X-MCP-Toolsets": "repos,issues,pull_requests"
=======
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
>>>>>>> Stashed changes
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

<<<<<<< Updated upstream
Never point `"tools": ["*"]` at the non-readonly URL; that grants merge, push,
and branch-write tools. If `add_issue_comment` is unavailable or write calls
fail, the skill falls back to leaving the artifacts in the review session
=======
Never use `"tools": ["*"]` against the non-readonly URL; that also grants merge,
push, and branch-write tools. If `add_issue_comment` is unavailable or write
calls fail, the skill falls back to leaving the artifacts in the review session
>>>>>>> Stashed changes
instead of claiming they were posted.

## Related files

- Procedure: `.github/instructions/pull-request-review.instructions.md` (Step 7)
- Lesson schema: `references/lesson-schema.md`
