---
description: Reconciles prior Copilot findings and reviews only the new pull request diff on each repeat pass. Use for Copilot code review, pull request review, repeat-pass review, and when reconciling prior findings or honoring Won't fix declines.
metadata:
    github-path: skills/code-review
    github-pinned: b7b583d6d7b3cbee347991c970d1ffebde489636
    github-ref: b7b583d6d7b3cbee347991c970d1ffebde489636
    github-repo: https://github.com/rsksmart/copilot-review-template
    github-tree-sha: 47e6496c6f2b044c5890b435939867ddd312083f
name: code-review
---
# Pull request review procedure

## Purpose and scope

- These instructions define *how* to review across repeated passes on the same pull request.
- They do not define *what* to look for. Apply any repository review checklist or coding standards
  alongside this procedure.
- On a repeat pass, run the phases in this order: reconcile prior comments (Step 3), then review
  the new diff (Step 4). Do not interleave them.

## Do not review vendored framework copies

In an **integration repository** (any repository that installed this skill, other than the
framework that publishes it), do not raise findings on these paths only:

- `.github/skills/code-review/**`
- `.github/workflows/ccr-dispatch.yml`

Those two trees are vendored copies of this framework. Review them in the framework repository,
not in consumer pull requests. That is what stops an infinite loop of new findings on the
procedure itself.

Leave every other skill, instruction file, and workflow in scope. Do not skip
`.github/skills/**` as a whole.

In the framework repository itself, these paths stay fully in scope.

## Step 1 — Read the pull request with the GitHub MCP server

Read the current state of the pull request before writing any comment.

### Resolve the coordinates

`pull_request_read` requires `owner`, `repo`, and `pullNumber`. Missing arguments are never a valid
reason to report the context as unavailable.

- Prefer `owner`, `repo`, and `pullNumber` from the Copilot code review context when they are present.
- Otherwise derive `owner` and `repo` from `git remote get-url origin`: take the last two path
  segments and strip a trailing `.git`. Use `origin` specifically when more than one remote is
  configured.
- If `origin` is missing or unparseable, use the GitHub organization of the repository under review
  if you can see it, otherwise the checkout's root directory name as the repo.
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
  such as "Suppressed comments". Being withheld does not make a finding new. A suppressed finding
  has no thread; link the review that contains it instead.

Classify each finding as Addressed, Partially addressed, Not addressed, Declined, No longer
applicable, or Unverified. Classify from code evidence in the head commit, not from thread state.

Every finding in the inventory takes exactly one classification and none may be dropped. When the
evidence is thin, report the finding as open rather than closing it.

Decline matching, markers, and author rules are in [references/declines.md](references/declines.md).

## Step 4 — Review the new changes, scoped to the diff

A finding is new only if it is absent from the Step 3 inventory.

Raise new findings only on lines that appear in the diff under review. Leave problems in untouched
code alone.

- On a first pass, the diff under review is the full pull request diff from `get_diff`.
- On a repeat pass, narrow it to what changed since the most recent Copilot review. Take that
  review's `commit_id` from `get_reviews`, then:
  1. Confirm the commit exists locally with `git cat-file -e <commit_id>^{commit}`.
  2. If it is present, use `git diff <commit_id>...HEAD`.
  3. If it is missing, rebuild the change set through the API: call `get_commits`, take every commit
     made after the review's `submitted_at`, and call `get_commit` on each one.
  4. Raise new findings only on lines in that change set.
- If the review carries no `commit_id`, go straight to the API path using `submitted_at` as the
  cutoff.
- Fall back to the full pull request diff only when neither git nor the API can produce a change set,
  and say plainly that you did.
- A file that was reviewed in an earlier pass and has not changed since produces no new comments
  except a repeat finding under Step 6.

### Check every finding against the inventory before posting it

Match each finding you are about to post against the Step 3 inventory on the substance of the
problem and the code it concerns, not on line numbers. If it matches, it is a repeat (Step 6).
Report a finding as new only after this check fails to match it.

## Step 5 — Report the reconciliation

Post an inline comment for every **Not addressed**, **Partially addressed**, and **Unverified**
finding. Post no inline comment for **Addressed**, **Declined**, or **No longer applicable**.

Where you summarize the review, group the findings under a `Previous Copilot review follow-up`
heading with one subsection per classification. Give each entry as `path:line`, the restatement,
and the evidence. Name who declined each declined finding and the reason they gave.

## Step 6 — Mark repeat findings

When a finding was already raised in a prior Copilot review on this pull request, say so in the
comment and link the original. Restate the issue and why it still matters. Report a repeat at the
same severity as the first time. Never post an inline comment for a **Declined** finding.

Do not call a finding a repeat when the code changed and this is a different problem, when the
earlier mention came from a human, or when the original comment cannot be located.

## Additional resources

- Decline convention and author rules: [references/declines.md](references/declines.md)
- Full original procedure notes: [references/procedure.md](references/procedure.md)
