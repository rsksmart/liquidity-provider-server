---
applyTo: "**"
excludeAgent: "cloud-agent"
---

# Pull request review procedure

These instructions govern *how* to conduct a review across repeated passes on the same pull
request. They do not replace the review checklist in `.github/prompts/review-code.prompt.md`,
which defines *what* to look for. Apply both.

## Step 1 — Establish pull request context with the GitHub MCP server

Before writing a single comment, use the GitHub MCP server to read the current state of this
pull request. Do not rely on recollection of earlier passes and do not assume this is a fresh
pull request.

Call `pull_request_read` for the current owner, repo, and pull request number, using these
methods:

- `get` — title, body, base and head refs, draft state, and linked issues.
- `get_files` and `get_diff` — the changed lines that are actually under review.
- `get_reviews` — the full review history. Identify reviews authored by Copilot
  (`copilot-pull-request-reviewer[bot]`, or any review whose author is the Copilot reviewer bot).
- `get_review_comments` — review threads, including their `isResolved`, `isOutdated`, and
  `isCollapsed` metadata, along with the comments in each thread.
- `get_comments` — the pull request conversation, where an author's explanation of how they
  handled earlier feedback often lives.
- `get_commits` — which commits landed after the most recent Copilot review.

Paginate until the review history is complete. A truncated list of prior comments produces a
wrong "addressed" accounting, which is worse than no accounting at all.

If the GitHub MCP server is unavailable or a call fails, state that at the top of the pull
request summary and note that the review ran without prior-review context. Do not silently skip
the reconciliation described below.

## Step 2 — Determine whether this is a first pass or a repeat pass

**First pass** — no prior review on this pull request was authored by Copilot. Perform the review
normally against the checklist. In the summary, include a single line stating that this is the
first Copilot review of this pull request, and omit the follow-up section entirely.

**Repeat pass** — at least one prior Copilot review exists. Perform the normal review *and* the
reconciliation in steps 3 and 4.

Only Copilot's own prior comments are tracked in the addressed accounting. Comments from human
reviewers are useful context and may inform the review, but they are never listed as addressed or
unaddressed, and they never make a finding a "repeat".

## Step 3 — Reconcile prior Copilot comments against the current code

Build an inventory of every comment from every prior Copilot review on this pull request.
Classify each one as:

- **Addressed** — the code now does what the comment asked.
- **Partially addressed** — some of the comment was acted on, some was not.
- **Not addressed** — the code is unchanged in the relevant respect.
- **No longer applicable** — the code the comment referred to was deleted or reworked such that
  the concern no longer exists.

Classify based on evidence in the head commit, not on thread state. A resolved, collapsed, or
outdated thread is a hint only. GitHub's own documentation notes that resolving a conversation
does not mean the underlying issue was fixed, so read the file at the comment's path in the
current head and confirm the change is really there before calling anything addressed.

Record for each item: the file and line, a one-sentence restatement of the original point, the
classification, and the evidence (the commit that changed it, or the current code that shows it
unchanged).

If a prior comment cannot be verified either way, classify it as **Unverified** and say why.
Never assert that something was addressed without having checked.

## Step 4 — Report the reconciliation in the pull request summary

On a repeat pass, open the summary with a section titled **Previous Copilot review follow-up**,
placed before the summary of the new changes. Use this shape, and omit any subsection that has no
entries:

```markdown
## Previous Copilot review follow-up

N of M comments from previous Copilot reviews on this pull request are addressed.

### Addressed
- `internal/usecases/pegout/send_pegout.go:142` — nil check missing before dereferencing the
  quote pointer. Fixed in abc1234; the pointer is now validated before use.

### Partially addressed
- `internal/adapters/dataproviders/rootstock/common.go:88` — repeated string literal moved to a
  constant, but two of the four occurrences still inline the literal.

### Not addressed
- `pkg/liquidity_provider.go:210` — handler still returns the domain entity instead of a DTO.

### No longer applicable
- `internal/usecases/reports/get_assets_report.go:55` — the function this referred to was removed.

### Unverified
- `internal/adapters/entrypoints/watcher/pegout_rsk_watcher.go:71` — could not determine whether
  the goroutine nil check was added; the surrounding code was restructured.
```

State the counts plainly. If none of the prior comments were addressed, say so directly. Follow
this section with the normal review summary covering what changed since the last pass.

## Step 5 — Mark repeat findings in inline comments

When a finding was already raised in a prior Copilot review on this pull request, say so in the
comment itself. Begin the comment with a marker that links to the original:

```markdown
**Repeat finding — this was already raised in this pull request** ([original comment](URL))
```

If the same point has been raised more than twice, include the count, for example
`**Repeat finding (3rd time) — ...**`.

After the marker, restate the issue and why it still matters. A bare link is not enough; the
reader should not have to open the original comment to understand the problem.

Severity does not decrease on repetition. A finding raised before is at minimum a 🟡 suggestion,
and anything previously marked 🔴 critical stays 🔴 critical.

Do **not** apply the repeat marker when:

- The code at that location changed and this is a genuinely different problem, even if it is in
  the same file or of the same category.
- The earlier mention came from a human reviewer rather than Copilot.
- The original comment cannot be located. Treat the finding as new instead of guessing.

## Step 6 — Constraints

- Never re-post a prior comment verbatim without the repeat marker.
- Never treat a thumbs-down reaction, a dismissed review, or a resolved conversation as evidence
  that the code was fixed.
- Never include the follow-up section on a first pass.
- Do not re-review unchanged code that no prior comment touched; the follow-up section is about
  prior comments, not a full re-audit.
- Keep the follow-up section factual and neutral. Report what is and is not addressed without
  editorializing about the author.
