---
applyTo: "**"
excludeAgent: "cloud-agent"
---

# Pull request review procedure

## Purpose and scope

- These instructions define *how* to review across repeated passes on the same pull request.
- They do not define *what* to look for. Apply any repository review checklist or coding standards
  alongside this procedure.
- On a repeat pass, run the phases in this order: reconcile prior comments (Step 3), then review
  the new diff (Step 4). Do not interleave them.

## Step 1 — Read the pull request with the GitHub MCP server

Read the current state of the pull request before writing any comment.

### Resolve the coordinates

`pull_request_read` requires `owner`, `repo`, and `pullNumber`. Missing arguments are never a valid
reason to report the context as unavailable.

- Derive `owner` and `repo` from `git remote get-url origin`: take the last two path segments and
  strip a trailing `.git`. For example `git@github.com:org/repo.git` yields owner `org` and repo
  `repo`.
- Use `origin` specifically when more than one remote is configured; other remotes may point at
  forks whose names differ from the repository under review.
- If `origin` is missing or unparseable, use owner `rsksmart` and the checkout's root directory name
  as the repo.
- Use `pullNumber` from the review context when available. Otherwise:
  1. Run `git rev-parse HEAD` and `git rev-parse --abbrev-ref HEAD`.
  2. Call `list_pull_requests` with `state: "open"` and `head: "OWNER:BRANCH"`.
  3. Match the returned head SHA against the local head commit.
  4. If the branch filter returns nothing, call `list_pull_requests` with `state: "open"` and no head
     filter, then match on head SHA or branch name.
- Report the context as unavailable only after all four steps fail, and name what each returned.

### Read these methods

- `get` — title, body, base and head refs, draft state, linked issues.
- `get_files` and `get_diff` — the lines under review.
- `get_reviews` — the full review history. Copilot's reviews are authored by
  `copilot-pull-request-reviewer[bot]`. Read each review `body`: withheld findings are listed there
  and nowhere else.
- `get_review_comments` — threads with their `isResolved`, `isOutdated`, and `isCollapsed` metadata.
  Read every reply, not just the first comment; declines live in replies.
- `get_comments` — the pull request conversation.
- `get_commits` — the commits that landed after the most recent Copilot review.

Paginate until the review history is complete.

If a call fails, say so, name the tool and method and the error it returned, and note that the review
ran without prior-review context. Do not skip reconciliation silently.

## Step 2 — First pass or repeat pass

- **First pass** — no prior Copilot review exists. Review normally against the repository's review
  standards. Skip Step 3 and omit the follow-up record entirely.
- **Repeat pass** — at least one prior Copilot review exists. Do Steps 3, 4, and 5.
- Track only Copilot's own prior comments. Human review comments inform the review but are never
  counted as addressed or unaddressed, and never make a finding a repeat.

## Step 3 — Reconcile prior Copilot comments

Complete this step before looking for anything new. Here you may read any file at any path, whether
or not the latest push touched it.

Inventory every prior Copilot finding. They come from two sources and both count:

- **Posted comments** — the threads from `get_review_comments`.
- **Suppressed findings** — those listed in each prior Copilot review `body`, usually under a heading
  such as "Comments suppressed due to low confidence". Being withheld does not make a finding new. A
  suppressed finding has no thread, so it carries no thread metadata and no comment URL; link the
  review that contains it instead.

Classify each finding as:

- **Addressed** — the code now does what the comment asked.
- **Partially addressed** — some of the comment was acted on, some was not.
- **Not addressed** — the code is unchanged in the relevant respect.
- **Declined** — a human deliberately declined it, per the convention below.
- **No longer applicable** — the code was deleted or reworked and the concern no longer exists.
- **Unverified** — the evidence was inconclusive.

Classify from code evidence in the head commit, not from thread state. Resolved, collapsed, and
outdated threads, dismissed reviews, and reactions are hints only. Read the file at the comment's
path in the current head before calling anything addressed, and never assert that something was
addressed without checking.

Record per item: the file and line, a one-sentence restatement of the original point, the
classification, and the evidence — the commit that changed it, or the current code showing it
unchanged. For **Unverified**, record what blocked the determination.

### How a developer declines a finding

A finding is **Declined** when a human replies in its thread with a message beginning `Won't fix:`
followed by a reason:

```text
Won't fix: the nil case is unreachable here, the caller validates the pointer before the call.
```

A suppressed finding has no thread to reply in, so it may also be declined from a top-level pull
request comment that names the location:

```text
Won't fix: path/to/file.go:96 — the nil case is unreachable, the caller validates the pointer
before the call.
```

- Read top-level declines from `get_comments` and match them by file path and line, allowing for line
  drift caused by later edits.
- If a top-level decline could match more than one finding in that file, do not guess: leave the
  findings unchanged and say the decline could not be matched to a single location.
- Match the marker case-insensitively and tolerate a missing apostrophe, so `Won't fix:`,
  `won't fix:`, and `Wont fix:` all count.
- Ignore the marker in comments authored by Copilot or any other bot. Any human commenter may
  decline, including the author.
- A reason is required. Without one, keep the finding **Not addressed** and use its inline comment to
  say that a decline was found, that it was not honored because no reason was given, and that adding
  a reason will retire the comment.
- Resolving a conversation is not a decline, and neither is a thumbs-down.
- If the code satisfies the comment, classify it **Addressed** even when a decline reply is present.
  Evidence in the code wins.
- A decline covers the point that was raised, not the location forever. A different problem
  introduced on those lines by a later push is a new finding.

## Step 4 — Review the new changes, scoped to the diff

A finding is new only if it is absent from the Step 3 inventory. Deriving it in this step does not
make it new, and neither does its appearing in the diff.

Raise new findings only on lines that appear in the diff under review. Unchanged surrounding code is
context for understanding the change, not review surface. Leave problems in untouched code alone even
when they violate the repository's review standards.

- On a first pass, the diff under review is the full pull request diff from `get_diff`.
- On a repeat pass, narrow it to what changed since the most recent Copilot review. Take that
  review's `commit_id` from `get_reviews`, then:
  1. Confirm the commit exists locally with `git cat-file -e <commit_id>^{commit}`. Never assume it
     does; the checkout may be shallow, and the branch may have been rebased or force-pushed.
  2. If it is present, use `git diff <commit_id>...HEAD`.
  3. If it is missing, rebuild the change set through the API: call `get_commits`, take every commit
     made after the review's `submitted_at`, and call `get_commit` on each one. The union of those
     changes is the incremental diff.
  4. Raise new findings only on lines in that change set.
- If the review carries no `commit_id`, go straight to the API path using `submitted_at` as the
  cutoff.
- Fall back to the full pull request diff only when neither git nor the API can produce a change set,
  and say plainly that you did. When the pull request adds a file, the full diff contains every line
  of it, so every earlier finding in that file will look like it is in scope. The inventory check
  below is what stops that from turning old findings into new ones.
- A file that was reviewed in an earlier pass and has not changed since produces no new comments.
  Comment on it only for a repeat finding under Step 6, or when a change elsewhere in this push broke
  it.
- Do not reopen reconciliation here. The Step 3 classifications stand as recorded, including for
  files outside this diff.

### Check every finding against the inventory before posting it

Match each finding you are about to post against the Step 3 inventory. Do this for every finding,
however you arrived at it and whether or not it fell inside the diff.

- If it matches an inventory entry, it is not new. Report it under Step 6 with the repeat marker,
  even though you derived it here.
- Match on the substance of the problem and the code it concerns, not on line numbers. Lines drift
  between passes, so the same finding rarely sits at the same line twice.
- A finding on a line this push did not change is almost always an inventory entry. If you cannot
  match it to one, say so in the comment rather than presenting it as new.
- Report a finding as new only after this check fails to match it.

## Step 5 — Report the reconciliation

Inline comments are the channel that reliably reaches the reader, so nothing actionable may exist
only in the reconciliation record. The pull request overview is generated by GitHub; do not treat it
as where the reconciliation lands.

- Post an inline comment for every **Not addressed** and **Partially addressed** finding, following
  Step 6.
- Post an inline comment for every **Unverified** finding, phrased as a question: say the point was
  raised earlier, that you could not determine whether it was addressed, and what blocked the
  determination, then ask the author to confirm. Do not assert that the code is unfixed and do not
  give it an occurrence count.
- Post no inline comment for **Addressed**, **Declined**, or **No longer applicable** findings.

Where you summarize the review, group the findings under a `Previous Copilot review follow-up`
heading with one subsection per classification — Addressed, Partially addressed, Not addressed,
Declined, No longer applicable, Unverified — omitting any that have no entries. Give each entry as
`path:line`, the restatement, and the evidence:

```markdown
### Not addressed
- `path/to/file.go:210` — handler still returns the domain entity instead of a DTO.

### Declined
- `path/to/file.go:96` — declined by @someuser: the nil case is unreachable because the caller
  validates the pointer first.
```

- Name who declined each declined finding and the reason they gave.
- If a decline rests on a factual mistake, say so once, in one sentence, next to that entry. Do not
  re-post the inline comment, do not repeat the disagreement on later passes, and do not reclassify
  the finding away from **Declined**.
- If none of the prior comments were addressed, say so directly.
- Keep the record factual. Report what is and is not addressed without editorializing about the
  author.

## Step 6 — Mark repeat findings

When a finding was already raised in a prior Copilot review on this pull request, say so in the
comment and link the original. Link the review instead when the earlier finding was suppressed rather
than posted.

This covers every finding that matched the inventory check in Step 4, not only the ones carried over
from the reconciliation. A repeat that reaches the reader through the diff review is still a repeat.

For example: "This was already raised earlier in this pull request (link), and the code has not
changed since."

- Restate the issue and why it still matters. A bare link is not enough.
- If the same point has come up more than twice, say how many times.
- Report a repeat at the same severity as the first time. A blocking problem stays blocking.
- For an **Unverified** item, state the history but give no verdict; you are asking whether it was
  handled, not reporting that it was not.
- Never post an inline comment for a **Declined** finding, on this pass or any later one.

Do not call a finding a repeat when:

- The code at that location changed and this is a genuinely different problem, even if it is in the
  same file or of the same category.
- The earlier mention came from a human reviewer rather than Copilot.
- The original comment cannot be located. Treat the finding as new instead of guessing.

## Step 7 — Publish follow-up artifacts

After Steps 5 and 6 are complete, if `.github/skills/code-review/SKILL.md` exists,
follow that skill to post the reconciliation and lessons comments via the GitHub
MCP server. If the skill is missing, do nothing extra.
