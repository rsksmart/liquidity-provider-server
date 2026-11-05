# How a developer declines a finding

A finding is **Declined** when a human replies in its thread with a message beginning `Won't fix:`
followed by a reason:

```text
Won't fix: the nil case is unreachable here, the caller validates the pointer before the call.
```

A thread reply binds to that comment only. Nearby findings on the same file are not declined.

A suppressed finding has no thread to reply in, so it may also be declined from a top-level pull
request comment that names the location:

```text
Won't fix: path/to/file.go:96 — the nil case is unreachable, the caller validates the pointer
before the call.
```

Backticks around the path are ignored. After `path:line`, a dash, colon, period, comma, or a blank
space may separate the reason:

```text
Won't fix: path/to/file.go:96 is for a separate task
```

- Read top-level declines from `get_comments` and match them by file path and line, allowing for
  small line drift caused by later edits.
- If a top-level decline names an exact `path:line`, decline every finding at that location.
- If it only lands near a line and more than one finding could match, do not guess: leave the
  findings unchanged and report the decline as unmatched.
- Match the marker case-insensitively and tolerate a missing apostrophe, so `Won't fix:`,
  `won't fix:`, and `Wont fix:` all count.
- Ignore the marker in comments authored by Copilot or any other bot. Honor a human decline only
  when GitHub's `author_association` is `OWNER`, `MEMBER`, `COLLABORATOR`, `CONTRIBUTOR`,
  `FIRST_TIME_CONTRIBUTOR`, or `FIRST_TIMER` — that includes the PR author. Ignore `NONE`.
- A reason is required. Without one, keep the finding **Not addressed** and say that a decline was
  found but not honored because no reason was given.
- Resolving a conversation is not a decline, and neither is a thumbs-down.
- If the code satisfies the comment, classify it **Addressed** even when a decline reply is present.
  Evidence in the code wins.
- A decline covers the point that was raised, not the location forever.
