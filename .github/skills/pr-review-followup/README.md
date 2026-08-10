# pr-review-followup skill

Posts the Copilot reconciliation follow-up and suggested Rootstock lessons as
separate pull request conversation comments after a review that used
`.github/instructions/pull-request-review.instructions.md`.

## Required repository MCP settings

Configure these under **Settings → Copilot → MCP servers**. This JSON is not
stored in the repository.

Widen the GitHub MCP endpoint past the read-only URL, and allowlist only the
comment tool:

```json
{
  "mcpServers": {
    "github-mcp-server": {
      "type": "http",
      "url": "https://api.githubcopilot.com/mcp/",
      "tools": ["add_issue_comment"],
      "headers": {
        "X-MCP-Toolsets": "issues,pull_requests"
      }
    }
  }
}
```

Also keep **Allow Copilot to use MCP tools when reviewing pull requests**
enabled for the repository.

Do not use `"tools": ["*"]` unless you intentionally want broader write access.
If `add_issue_comment` is unavailable or write calls fail, the skill falls back
to leaving the artifacts in the review session instead of claiming they were
posted.

## Related files

- Procedure: `.github/instructions/pull-request-review.instructions.md` (Step 7)
- Lesson schema: `references/lesson-schema.md`
