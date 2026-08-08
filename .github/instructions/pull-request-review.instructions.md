---
applyTo: "**"
excludeAgent: "cloud-agent"
---

# Pull request review procedure

These instructions govern *how* to conduct a review across repeated passes on the same pull
request. They do not replace the review checklist in `.github/prompts/review-code.prompt.md`,
which defines *what* to look for. Apply both.

A repeat pass has two phases and they run in this order: reconcile the prior Copilot comments
first (Step 3), then review what changed since the last pass (Step 4). Do not interleave them.
Reconciliation may read anywhere in the codebase; the search for new findings may not.

## Step 1 — Establish pull request context with the GitHub MCP server

Before writing a single comment, use the GitHub MCP server to read the current state of this
pull request. Do not rely on recollection of earlier passes and do not assume this is a fresh
pull request.

### Resolving the pull request coordinates

The `pull_request_read` tool requires `owner`, `repo`, and `pullNumber`. These are always
obtainable — never skip this step on the grounds that they were not supplied.

- Read `owner` and `repo` from the `origin` remote rather than assuming a name. Run
  `git remote get-url origin`, take the last two path segments, and strip any trailing `.git`.
  For example `git@github.com:rsksmart/liquidity-provider-server.git` yields owner `rsksmart` and
  repo `liquidity-provider-server`. Use `origin` specifically — this repository also has remotes
  pointing at security advisory forks whose names differ from the real repository.
- If `origin` is missing or unparseable, use `rsksmart` as the owner and take the repository name
  from the root directory of the checkout.
- `pullNumber` is the number of the pull request currently being reviewed. Use it directly if it
  is available in the review context.

If the pull request number is not directly available, resolve it in this order:

1. Get the head commit with `git rev-parse HEAD` and the current branch with
   `git rev-parse --abbrev-ref HEAD`.
2. Call `list_pull_requests` with the resolved `owner` and `repo`, `state: "open"`, and
   `head: "OWNER:BRANCH"` using the owner and branch from the previous steps.
3. Match the returned pull request's head SHA against the local head commit and use that pull
   request's number.
4. If the branch filter returns nothing, call `list_pull_requests` with `state: "open"` and no
   head filter, then match on head SHA or branch name.

Only report the context as unavailable after all four steps have failed, and when you do, name
the steps you attempted and what each returned.

Once the coordinates are resolved, call `pull_request_read` with them, using these methods:

- `get` — title, body, base and head refs, draft state, and linked issues.
- `get_files` and `get_diff` — the changed lines that are actually under review.
- `get_reviews` — the full review history. Identify reviews authored by Copilot
  (`copilot-pull-request-reviewer[bot]`, or any review whose author is the Copilot reviewer bot).
- `get_review_comments` — review threads, including their `isResolved`, `isOutdated`, and
  `isCollapsed` metadata, along with the comments in each thread. Read every comment in a thread,
  not only the first one: a developer's decision to decline a comment lives in a reply.
- `get_comments` — the pull request conversation, where an author's explanation of how they
  handled earlier feedback often lives.
- `get_commits` — which commits landed after the most recent Copilot review.

Paginate until the review history is complete. A truncated list of prior comments produces a
wrong "addressed" accounting, which is worse than no accounting at all.

If the GitHub MCP server is genuinely unavailable, or a `pull_request_read` call fails, say so
explicitly rather than staying silent, and note that the review ran without prior-review context.
Name the specific tool and method that failed and the error it returned. Missing arguments are not
a valid reason to report the context as unavailable; resolve them as described above. Do not
silently skip the reconciliation below.

## Step 2 — Determine whether this is a first pass or a repeat pass

**First pass** — no prior review on this pull request was authored by Copilot. Perform the review
normally against the checklist. There is nothing to reconcile, so skip Step 3 and omit the
follow-up record entirely.

**Repeat pass** — at least one prior Copilot review exists. Reconcile the prior comments per Step 3,
then review the new changes per Step 4, then report both per Step 5.

Only Copilot's own prior comments are tracked in the addressed accounting. Comments from human
reviewers are useful context and may inform the review, but they are never listed as addressed or
unaddressed, and they never make a finding a "repeat".

## Step 3 — First, reconcile prior Copilot comments against the current code

Do this before looking for anything new. Reconciliation and the new review are separate phases and
must not be interleaved: finish this step completely, then move to Step 4.

In this step you may read any file at any path, whether or not the latest push touched it. Answering
"was this prior comment addressed?" requires reading the current state of that code wherever it
lives.

Build an inventory of every comment from every prior Copilot review on this pull request.
Classify each one as:

- **Addressed** — the code now does what the comment asked.
- **Partially addressed** — some of the comment was acted on, some was not.
- **Not addressed** — the code is unchanged in the relevant respect.
- **Declined** — a human deliberately decided not to act on it, following the convention below.
- **No longer applicable** — the code the comment referred to was deleted or reworked such that
  the concern no longer exists.
- **Unverified** — the evidence was inconclusive and you could not determine whether it was
  addressed.

Classify based on evidence in the head commit, not on thread state. A resolved, collapsed, or
outdated thread is a hint only. GitHub's own documentation notes that resolving a conversation
does not mean the underlying issue was fixed, so read the file at the comment's path in the
current head and confirm the change is really there before calling anything addressed.

### How a developer declines a comment

A comment is **Declined** when someone replied in its thread with a message beginning `Won't fix:`
followed by a reason. For example:

```text
Won't fix: the nil case is unreachable here, the caller validates the pointer before the call.
```

Honor a decline under these conditions:

- Match the marker case-insensitively and tolerate a missing apostrophe, so `Won't fix:`,
  `won't fix:`, and `Wont fix:` all count.
- The reply must come from a human. Ignore the marker in replies authored by Copilot or any other
  bot.
- Any human commenter on the pull request can decline, including the author.
- A reason is required. If the marker appears with nothing substantive after the colon, do **not**
  honor it: keep the comment classified as **Not addressed**. That classification produces an
  inline comment, so say it there rather than only in the record — state that a decline was found,
  that it was not honored because no reason was given, and that adding a reason will retire the
  comment. The developer who declined has to be able to see why it had no effect without opening
  the session log.

Resolving the conversation is not a decline on its own, and neither is a thumbs-down. Only the
marker with a reason counts.

If the code actually satisfies the comment, classify it as **Addressed** even when a decline reply
is present — evidence in the code wins. **Declined** applies only where the code still does not
satisfy the comment and a human said so deliberately.

A decline covers the specific point that was raised, not the location forever. If a later push
changes those lines and introduces a different problem, that is a new finding under the normal
rules.

Record for each item: the file and line, a one-sentence restatement of the original point, the
classification, and the evidence (the commit that changed it, or the current code that shows it
unchanged).

If a prior comment cannot be verified either way, classify it as **Unverified** and say why.
Never assert that something was addressed without having checked.

This step is complete when every prior Copilot comment carries a classification and its supporting
evidence. Only then continue.

## Step 4 — Then review the new changes, scoped to the diff

Reconciliation is finished by the time you reach this step. Everything from here on is about finding
new problems, and it is strictly limited to the diff.

Raise new findings only on lines that appear in the diff under review. Unchanged surrounding code is
context for understanding the change, not review surface. When a problem sits in code this pull
request did not touch, leave it alone even if it violates the checklist.

On a first pass, the diff under review is the full pull request diff from `get_diff`.

On a repeat pass, narrow it to what changed since the last Copilot review:

1. Take the `commit_id` of the most recent Copilot review from `get_reviews`.
2. Diff that commit against the current head with `git diff <commit_id>...HEAD`.
3. Raise new findings only on lines that appear in that incremental diff.

If the review carries no `commit_id`, call `get_commits` and use the last commit whose timestamp is
at or before that review's `submitted_at`. If the baseline still cannot be determined, fall back to
the full pull request diff and say so.

A file that was reviewed in an earlier pass and has not changed since should produce no new
comments. If you are about to comment on such a file, the only valid reasons are that it is a
repeat finding under Step 6, or that a change elsewhere in this push broke it.

Do not reopen reconciliation here, and do not let this narrower scope walk back anything you
verified in Step 3. The classifications from that step stand as recorded, including for files that
fall outside this diff.

## Step 5 — Record the reconciliation

The pull request overview comment is written by GitHub, not by you. Do not attempt to change its
wording, structure, or content, and do not treat it as the place where the reconciliation lands.

Everything a reviewer needs to act on must therefore reach the pull request as inline comments.
Every prior comment classified **Not addressed** or **Partially addressed** gets an inline comment
that says it has come up before, following Step 6. That channel is the one that reliably reaches
the reader, so nothing actionable may exist only in the reconciliation record.

A prior comment classified **Unverified** also gets an inline comment, but phrased as a question
rather than a finding. Say that the point was raised earlier, that you could not determine whether
it was addressed, and what specifically blocked the determination, then ask the author to confirm.
Do not assert that the code is unfixed, and do not give it an occurrence count — you do not know
that it recurred. Silence is the worst outcome here, because an unverified item is the one most
likely to be a real problem nobody looked at.

Do not post inline comments for prior comments classified **Addressed**, **Declined**, or **No
longer applicable**. Those belong in the record only.

Where you summarize the review, organize the reconciliation with this structure, omitting any
subsection that has no entries:

```markdown
## Previous Copilot review follow-up

### Addressed
- `internal/usecases/pegout/send_pegout.go:142` — nil check missing before dereferencing the
  quote pointer. Fixed in abc1234; the pointer is now validated before use.

### Partially addressed
- `internal/adapters/dataproviders/rootstock/common.go:88` — repeated string literal moved to a
  constant, but two of the four occurrences still inline the literal.

### Not addressed
- `pkg/liquidity_provider.go:210` — handler still returns the domain entity instead of a DTO.

### Declined
- `internal/usecases/pegin/call_for_user.go:96` — declined by @someuser: the nil case is
  unreachable because the caller validates the pointer first.

### No longer applicable
- `internal/usecases/reports/get_assets_report.go:55` — the function this referred to was removed.

### Unverified
- `internal/adapters/entrypoints/watcher/pegout_rsk_watcher.go:71` — could not determine whether
  the goroutine nil check was added; the surrounding code was restructured.
```

List each declined comment with who declined it and the reason they gave, so the record is visible
without opening every thread.

If you believe a decline rests on a factual mistake, you may say so once, in one sentence, next to
that entry. Do not re-post the inline comment, do not repeat the disagreement on later passes, and
do not reclassify the comment away from **Declined**.

If none of the prior comments were addressed, say so directly.

## Step 6 — Say when a finding is a repeat

When a finding was already raised in a prior Copilot review on this pull request, state that fact
in the opening sentence of the comment and link the original. This is a requirement about content,
not styling: what matters is that the reader learns this is not the first time it has come up.

For example: "This was already raised earlier in this pull request (link), and the code has not
changed since."

If the same point has been raised more than twice, say how many times it has come up.

Then restate the issue and why it still matters. A bare link is not enough; the reader should not
have to open the original comment to understand the problem.

A repeated finding is never less important than it was the first time. If it was reported as a
blocking problem before, report it as a blocking problem again.

The exception is an **Unverified** item. State its history, but not a verdict: you are asking
whether it was handled, not reporting that it was not.

Do **not** describe a finding as a repeat when:

- The code at that location changed and this is a genuinely different problem, even if it is in
  the same file or of the same category.
- The earlier mention came from a human reviewer rather than Copilot.
- The original comment cannot be located. Treat the finding as new instead of guessing.

Never post an inline comment at all for a finding classified as **Declined**. Repeating a declined
comment inline is the specific behavior the decline convention exists to prevent.

## Step 7 — Constraints

- Never try to change the pull request overview comment; it is generated by GitHub.
- Never leave an actionable prior comment visible only in the reconciliation record. If it still
  needs work, or you could not tell whether it does, it needs an inline comment.
- Never report an **Unverified** item as though you had determined it was unfixed.
- Never re-post a prior comment without saying that it has come up before.
- Never re-raise a comment that was declined with a reason, on this pass or any later one.
- Never silently ignore a decline that lacked a reason; say that it was not honored and why.
- Never treat a thumbs-down reaction, a dismissed review, or a resolved conversation as evidence
  that the code was fixed, or as a decline.
- Never include the follow-up section on a first pass.
- Never begin looking for new findings before reconciliation is complete.
- Do not raise new findings on code outside the diff scope from Step 4, and do not re-audit files
  that this push did not touch.
- Keep the follow-up section factual and neutral. Report what is and is not addressed without
  editorializing about the author.
