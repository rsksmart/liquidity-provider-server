# pr-review-followup skill

Posts the Copilot reconciliation follow-up and suggested Rootstock lessons as
separate pull request conversation comments after a review that used
`.github/instructions/pull-request-review.instructions.md`.

## Required repository MCP settings

Configure these under **Settings → Copilot → MCP servers**. This JSON is not
stored in the repository.

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

```json
{
  "mcpServers": {
    "github-mcp-server": {
      "type": "http",
      "url": "https://api.githubcopilot.com/mcp/",
      "tools": [
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
      }
    }
  }
}
```

Also keep **Allow Copilot to use MCP tools when reviewing pull requests**
enabled for the repository.

Never point `"tools": ["*"]` at the non-readonly URL; that grants merge, push,
and branch-write tools. If `add_issue_comment` is unavailable or write calls
fail, the skill falls back to leaving the artifacts in the review session
instead of claiming they were posted.

## Related files

- Procedure: `.github/instructions/pull-request-review.instructions.md` (Step 7)
- Lesson schema: `references/lesson-schema.md`
