"""Finding models and shared constants."""

from __future__ import annotations

import re
from dataclasses import dataclass, field
from typing import Any

COPILOT_LOGINS = {
    "copilot-pull-request-reviewer[bot]",
    "Copilot",
    "copilot-pull-request-reviewer",
}

WONT_FIX_RE = re.compile(r"Won'?t\s+fix\s*:", re.IGNORECASE)
SUPPRESSED_FINDING_HEADER_RE = re.compile(
    r"^\*\*(?P<path>.+):(?P<line>\d+)\*\*$"
)
SUPPRESSED_SECTION_RE = re.compile(
    r"(?:<summary>\s*)?#{0,6}\s*Suppressed comments",
    re.IGNORECASE,
)
REVIEW_STATS_RE = re.compile(
    r"^- \*\*(Files reviewed|Comments generated|Review effort)",
    re.IGNORECASE,
)
PREVIOUSLY_MISSED_RE = re.compile(r"^\*\*Previously missed", re.IGNORECASE)
REPEAT_HINT_RE = re.compile(
    r"(already raised|repeat from earlier|was already raised|previously missed|previously raised)",
    re.IGNORECASE,
)
MARKER_RECONCILIATION = "<!-- ccr-followup:reconciliation:{review_id} -->"
MARKER_LESSONS = "<!-- ccr-followup:lessons:{review_id} -->"
LINE_PROXIMITY = 5


@dataclass
class Finding:
    path: str
    line: int
    body: str
    source: str  # "comment" | "suppressed"
    comment_id: int | None = None
    review_id: int | None = None
    commit_id: str = ""
    is_repeat: bool = False

    @property
    def key(self) -> str:
        return f"{self.path}:{self.line}"

    def summary(self, limit: int = 180) -> str:
        text = " ".join(self.body.split())
        if len(text) > limit:
            return text[: limit - 1] + "…"
        return text


@dataclass
class Classified:
    """Classifications this pass can reach without reading the code.

    There is no "addressed" bucket: confirming a fix requires inspecting the
    file at the head commit, which this pass never does. Every prior finding
    lands in one of the buckets below, so an unfixed finding is never dropped
    for lack of proof — it is reported as open or flagged for confirmation.

    Entries carry the evidence behind the classification.
    """

    not_addressed: list[tuple[Finding, str]] = field(default_factory=list)
    declined: list[tuple[Finding, str, str]] = field(default_factory=list)
    unverified: list[tuple[Finding, str]] = field(default_factory=list)

    def has_findings(self) -> bool:
        return any(
            [
                self.not_addressed,
                self.declined,
                self.unverified,
            ]
        )


@dataclass(frozen=True)
class ReviewData:
    current: dict[str, Any]
    prior: list[dict[str, Any]]
    comments: list[dict[str, Any]]


@dataclass(frozen=True)
class Inventories:
    current: list[Finding]
    prior: list[Finding]
    all_prior: list[Finding]


@dataclass(frozen=True)
class ChangeState:
    first_pass: bool
    code_changed: bool
    prior_sha: str
    current_sha: str
    # Origin review SHA -> files changed between that SHA and current_sha.
    # Missing key or None value means the change set is unknowable.
    touched_since: dict[str, set[str] | None]
