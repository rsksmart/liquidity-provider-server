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
REPEAT_HINT_RE = re.compile(
    r"(already raised|repeat from earlier|was already raised)",
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
    addressed: list[Finding] = field(default_factory=list)
    not_addressed: list[Finding] = field(default_factory=list)
    declined: list[tuple[Finding, str, str]] = field(default_factory=list)
    unverified: list[Finding] = field(default_factory=list)

    def has_findings(self) -> bool:
        return any(
            [
                self.addressed,
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
    touched: set[str]
