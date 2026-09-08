"""Classify prior findings against the current review and code changes."""

from __future__ import annotations

import re
from typing import Any

from .finding import (
    DECLINE_ASSOCIATIONS,
    LINE_PROXIMITY,
    REPEAT_HINT_RE,
    WONT_FIX_RE,
    ChangeState,
    Classified,
    Finding,
)
from .parse import is_copilot_user

TOP_LEVEL_DECLINE_RE = re.compile(
    r"^\s*Won'?t\s+fix\s*:\s*(?P<path>.+):(?P<line>\d+)"
    r"\s*(?:—|–|-)\s*(?P<reason>\S.*)$",
    re.IGNORECASE,
)


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
    finding: Finding,
    comments: list[dict[str, Any]],
    findings: list[Finding],
) -> tuple[str, str] | None:
    """Return (user, reason) if a Won't fix reply exists for this finding."""
    comments_by_id = {c.get("id"): c for c in comments}
    for decline in human_declines(comments):
        matched = decline_for_finding(
            finding, decline, comments_by_id, findings
        )
        if matched:
            return matched
    return None


def human_declines(comments: list[dict[str, Any]]) -> list[dict[str, Any]]:
    return [
        comment
        for comment in comments
        if WONT_FIX_RE.search(comment.get("body") or "")
        and may_decline(comment)
    ]


def may_decline(comment: dict[str, Any]) -> bool:
    if (comment.get("user") or {}).get("type") == "Bot":
        return False
    if is_copilot_user((comment.get("user") or {}).get("login")):
        return False
    association = (comment.get("author_association") or "").upper()
    return association in DECLINE_ASSOCIATIONS


def decline_for_finding(
    finding: Finding,
    decline: dict[str, Any],
    comments_by_id: dict[Any, dict[str, Any]],
    findings: list[Finding],
) -> tuple[str, str] | None:
    body = decline.get("body") or ""
    marker_reason = WONT_FIX_RE.split(body, maxsplit=1)[-1].strip()
    if not marker_reason:
        return None
    reply_to = decline.get("in_reply_to_id")
    parent = comments_by_id.get(reply_to)
    if thread_decline_matches(finding, reply_to, parent):
        return comment_author(decline), marker_reason
    return top_level_decline_for_finding(finding, decline, findings)


def thread_decline_matches(
    finding: Finding,
    reply_to: Any,
    parent: dict[str, Any] | None,
) -> bool:
    direct_reply = bool(finding.comment_id and reply_to == finding.comment_id)
    return direct_reply or parent_matches_finding(parent, finding)


def top_level_decline_for_finding(
    finding: Finding,
    decline: dict[str, Any],
    findings: list[Finding],
) -> tuple[str, str] | None:
    if decline.get("in_reply_to_id") is not None:
        return None
    parsed = TOP_LEVEL_DECLINE_RE.match(decline.get("body") or "")
    if parsed is None:
        return None
    path = parsed.group("path").strip()
    line = int(parsed.group("line"))
    matches = [
        candidate
        for candidate in findings
        if candidate.path == path
        and abs(candidate.line - line) <= LINE_PROXIMITY
    ]
    if len(matches) != 1 or matches[0] != finding:
        return None
    return comment_author(decline), parsed.group("reason").strip()


def comment_author(comment: dict[str, Any]) -> str:
    return (comment.get("user") or {}).get("login") or "unknown"


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
        category, item = classify_finding(
            finding, prior, current, comments, state
        )
        getattr(result, category).append(item)
    return result


def classify_finding(
    finding: Finding,
    prior: list[Finding],
    current: list[Finding],
    comments: list[dict[str, Any]],
    state: ChangeState,
) -> tuple[str, tuple[Finding, str] | tuple[Finding, str, str]]:
    """Classify one prior finding, never silently assuming it was fixed.

    This pass cannot confirm a fix, so it never reports one. It reports a
    finding as open whenever it can point at evidence that the code did not
    change, and asks for confirmation otherwise.
    """
    decline = find_decline(finding, comments, prior)
    if decline:
        return "declined", (finding, decline[0], decline[1])
    if any(similar(finding, candidate) for candidate in current):
        return "not_addressed", (finding, "raised again in the current review")
    unchanged = unchanged_evidence(finding, state)
    if unchanged:
        return "not_addressed", (finding, unchanged)
    return "unverified", (finding, needs_inspection_reason(finding, state))


def unchanged_evidence(finding: Finding, state: ChangeState) -> str:
    """Evidence that the code behind this finding never changed, if any.

    Measured from the review that originally raised the finding, not from the
    most recent prior review. A later review that left the file alone is not
    evidence the earlier finding is still open.
    """
    origin = finding.commit_id
    if not origin or not state.current_sha:
        return ""
    if origin == state.current_sha:
        return f"no commits since this finding's review ({origin[:8]})"
    if origin not in state.touched_since:
        return ""
    touched = state.touched_since[origin]
    if touched is None:
        return ""
    if finding.path not in touched:
        return f"`{finding.path}` untouched since this finding's review ({origin[:8]})"
    return ""


def needs_inspection_reason(finding: Finding, state: ChangeState) -> str:
    origin = finding.commit_id
    if not origin or origin not in state.touched_since:
        return "could not determine what changed since this finding's review"
    if state.touched_since[origin] is None:
        return "could not determine what changed since this finding's review"
    return "file changed since this finding's review; needs code inspection"


def new_findings(prior: list[Finding], current: list[Finding]) -> list[Finding]:
    """Posted comments that are new on this pass. Suppressed findings never qualify."""
    out: list[Finding] = []
    for cur in current:
        if cur.source != "comment":
            continue
        is_new = (
            not cur.is_repeat
            and not REPEAT_HINT_RE.search(cur.body)
            and not any(similar(cur, p) for p in prior)
        )
        if is_new:
            out.append(cur)
    out.sort(key=lambda f: (f.path, f.line))
    return out[:5]


def reconciliation_skip_reason(state: ChangeState) -> str:
    # A pass with no new commits is not a reason to stay quiet: it means every
    # prior finding is still open, which is exactly what has to be reported.
    if state.first_pass:
        return "first pass"
    return ""
