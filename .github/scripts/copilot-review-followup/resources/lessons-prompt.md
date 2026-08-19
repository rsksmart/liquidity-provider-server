# Author Rootstock lesson cards from Copilot code-review findings

You are authoring lesson YAML cards for the Rootstock lessons DB from **new**
findings on a pull request. Output **only** the pull request comment markdown
described below — no preamble, no analysis outside that comment.

## Hard rules

- Use only the findings in `FINDINGS_JSON`. Do not invent extra findings.
- Cap at 5 lessons. Prefer high-severity / repo-pattern issues over style nits.
- Skip anything that looks like a repeat, decline, or formatting-only nit.
- Each lesson MUST match the schema in SCHEMA_MD exactly (no extra top-level fields).
- Do not commit files. Do not call tools. Emit the comment body only.
- `frequency` must be `1`.
- Choose a category that fits the lesson. Examples: `go`, `reliability`,
  `security`, `testing`, `process`, `observability`, `sql`, `rpc`. Other
  categories are fine when they better match the finding.
- Lessons are exported into a durable lessons DB. Every field — especially `why` —
  must be understandable without any pull-request context.
- Never write PR-relative phrasing such as "this PR", "the current pull request",
  "this change", "here", "in this review", or "on this pass". Name the durable
  symbols instead (use case, function, type, contract, layer).
- `why` must explain the lasting technical reason the pattern matters (contract,
  failure mode, layer boundary), not where or when it was found.

## Comment shape (emit exactly this structure)

```markdown
<!-- ccr-followup:lessons:REVIEW_ID -->
## Suggested lessons for the Rootstock lessons DB

Copy each block into `skills/lessons/<category>/<id>.yaml`, then add a row to
`skills/lessons/README.md`.

<details>
<summary>category / lesson-id (severity)</summary>

```yaml
id: lesson-id
title: Human-readable title
tags:
  - go
  - reliability
severity: high
frequency: 1
trigger_patterns:
  - pattern
problem: |
  What goes wrong.
bad: |
  Bad pattern.
good: |
  Safer pattern.
why: |
  Lasting technical reason this pattern matters, naming durable symbols
  (use case, function, contract). No PR-relative phrasing.
```

</details>
```

Replace `REVIEW_ID` with the numeric review id provided below. Repeat `<details>`
blocks once per lesson. If there are no suitable findings, emit exactly:

```markdown
<!-- ccr-followup:lessons:REVIEW_ID -->
## Suggested lessons for the Rootstock lessons DB

No eligible new findings for lesson cards on this pass.
```

## SCHEMA_MD

<<<SCHEMA>>>

## REVIEW_ID

<<<REVIEW_ID>>>

## FINDINGS_JSON

<<<FINDINGS>>>
