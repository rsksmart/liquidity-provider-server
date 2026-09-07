"""Classify prior findings against the current review and code changes."""

from __future__ import annotations

import re
from typing import Any

from .finding import (
    LINE_PROXIMITY,
    REPEAT_HINT_RE,
    WONT_FIX_RE,
    ChangeState,
    Classified,
    Finding,
)
from .parse import is_copilot_user


def normalize(text: str) -> str:
    return re.sub(r"\s+", " ", text).strip().lower()


def similar(a: Finding, b: Finding) -> bool:
    if a.path != b.path:
        return False
    if abs(a.line - b.line) > LINE_PROXIMITY:
        return False
    na, nb = normalize(a.body), normalize(b.body)
    if not na or not nb:
        return True
    # Cheap overlap: shared significant token prefix or substring.
    if na[:60] in nb or nb[:60] in na:
        return True
    a_tokens = set(re.findall(r"[a-z0-9_]{4,}", na))
    b_tokens = set(re.findall(r"[a-z0-9_]{4,}", nb))
    if not a_tokens or not b_tokens:
        return abs(a.line - b.line) <= 1
    overlap = len(a_tokens & b_tokens) / max(1, min(len(a_tokens), len(b_tokens)))
    return overlap >= 0.35


def find_decline(
    finding: Finding, comments: list[dict[str, Any]]
) -> tuple[str, str] | None:
    """Return (user, reason) if a Won't fix reply exists for this finding."""
    comments_by_id = {c.get("id"): c for c in comments}
    declines = [
        c
        for c in comments
        if WONT_FIX_RE.search(c.get("body") or "")
        and not is_copilot_user((c.get("user") or {}).get("login"))
    ]

    for decline in declines:
        reply_to = decline.get("in_reply_to_id")
        parent = comments_by_id.get(reply_to)
        direct_reply = bool(finding.comment_id and reply_to == finding.comment_id)
        top_level_match = (
            reply_to is None
            and finding.path in (decline.get("body") or "")
            and str(finding.line) in (decline.get("body") or "")
        )
        nearby_reply = parent_matches_finding(parent, finding)
        if direct_reply or top_level_match or nearby_reply:
            body = decline.get("body") or ""
            login = (decline.get("user") or {}).get("login") or "unknown"
            reason = WONT_FIX_RE.split(body, maxsplit=1)[-1].strip()
            return login, reason or body.strip()
    return None


def parent_matches_finding(
    parent: dict[str, Any] | None, finding: Finding
) -> bool:
    if not parent or parent.get("path") != finding.path:
        return False
    parent_line = parent.get("line") or parent.get("original_line") or 0
    return abs(int(parent_line) - finding.line) <= LINE_PROXIMITY


def classify(
    prior: list[Finding],
    current: list[Finding],
    comments: list[dict[str, Any]],
    state: ChangeState,
) -> Classified:
    result = Classified()
    for finding in prior:
        category, item = classify_finding(finding, current, comments, state)
        getattr(result, category).append(item)
    return result


def classify_finding(
    finding: Finding,
    current: list[Finding],
    comments: list[dict[str, Any]],
    state: ChangeState,
) -> tuple[str, tuple[Finding, str] | tuple[Finding, str, str]]:
    """Classify one prior finding, never silently assuming it was fixed.

    This pass cannot confirm a fix, so it never reports one. It reports a
    finding as open whenever it can point at evidence that the code did not
    change, and asks for confirmation otherwise.
    """
    decline = find_decline(finding, comments)
    if decline:
        return "declined", (finding, decline[0], decline[1])
    if any(similar(finding, candidate) for candidate in current):
        return "not_addressed", (finding, "raised again in the current review")
    unchanged = unchanged_evidence(finding, state)
    if unchanged:
        return "not_addressed", (finding, unchanged)
    return "unverified", (finding, needs_inspection_reason(state))


def unchanged_evidence(finding: Finding, state: ChangeState) -> str:
    """Evidence that the code behind this finding never changed, if any."""
    if state.touched is None:
        return ""
    short = state.prior_sha[:8]
    if not state.code_changed:
        return f"no commits since the prior review ({short})"
    if finding.path not in state.touched:
        return f"`{finding.path}` untouched by the commits since {short}"
    return ""


def needs_inspection_reason(state: ChangeState) -> str:
    if state.touched is None:
        return "could not determine what changed since the prior review"
    return "file changed since the prior review; needs code inspection"


def new_findings(prior: list[Finding], current: list[Finding]) -> list[Finding]:
    out: list[Finding] = []
    for cur in current:
        is_new = (
            not cur.is_repeat
            and not REPEAT_HINT_RE.search(cur.body)
            and not any(similar(cur, p) for p in prior)
        )
        if is_new:
            out.append(cur)
    # Prefer comment-sourced, higher signal; cap at 5.
    out.sort(key=lambda f: (0 if f.source == "comment" else 1, f.path, f.line))
    return out[:5]


def reconciliation_skip_reason(state: ChangeState) -> str:
    # A pass with no new commits is not a reason to stay quiet: it means every
    # prior finding is still open, which is exactly what has to be reported.
    if state.first_pass:
        return "first pass"
    return ""
