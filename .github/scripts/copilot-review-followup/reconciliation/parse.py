"""Load and parse Copilot review data into finding inventories."""

from __future__ import annotations

import re
import subprocess
from dataclasses import dataclass, field
from typing import Any

from .config import Config, log
from .finding import (
    COPILOT_LOGINS,
    LINE_PROXIMITY,
    REPEAT_HINT_RE,
    SUPPRESSED_FINDING_HEADER_RE,
    ChangeState,
    Finding,
    Inventories,
    ReviewData,
)
from .github import gh_json, gh_paginated


@dataclass
class SuppressedParser:
    review_id: int
    findings: list[Finding] = field(default_factory=list)
    path: str | None = None
    line: int = 0
    body: list[str] = field(default_factory=list)

    def consume(self, body_line: str) -> None:
        header = SUPPRESSED_FINDING_HEADER_RE.match(body_line)
        if header:
            self.start(header.group("path"), int(header.group("line")))
            return
        if body_line == "</details>":
            self.finish()
            return
        if self.path is not None:
            self.body.append(body_line)

    def start(self, path: str, line: int) -> None:
        self.append_current()
        self.path = path.strip()
        self.line = line
        self.body = []

    def finish(self) -> None:
        self.append_current()
        self.path = None
        self.line = 0
        self.body = []

    def append_current(self) -> None:
        if self.path is None:
            return
        text = suppressed_summary(self.body)
        self.findings.append(
            Finding(
                path=self.path,
                line=self.line,
                body=text,
                source="suppressed",
                review_id=self.review_id,
                is_repeat=bool(REPEAT_HINT_RE.search(text)),
            )
        )


def is_copilot_user(login: str | None) -> bool:
    if not login:
        return False
    return login in COPILOT_LOGINS or login.lower().startswith("copilot")


def parse_suppressed(body: str, review_id: int) -> list[Finding]:
    section = suppressed_section(body)
    if not section:
        return []
    parser = SuppressedParser(review_id)
    for body_line in section.splitlines():
        parser.consume(body_line)
    parser.finish()
    return parser.findings


def suppressed_section(body: str) -> str:
    marker = "<summary>Suppressed comments"
    if marker not in body:
        return ""
    after_marker = body.partition(marker)[2]
    return after_marker.partition("</details>")[0]


def suppressed_summary(lines: list[str]) -> str:
    raw = "\n".join(lines).strip()
    summary = re.split(r"\n```", raw, maxsplit=1)[0].strip()
    return re.sub(r"^\*\s+", "", summary)


def findings_from_comments(
    comments: list[dict[str, Any]], review_id: int
) -> list[Finding]:
    out: list[Finding] = []
    for c in comments:
        is_review_finding = (
            c.get("pull_request_review_id") == review_id
            and not c.get("in_reply_to_id")
            and is_copilot_user((c.get("user") or {}).get("login"))
        )
        if is_review_finding:
            path = c.get("path") or ""
            line = c.get("line") or c.get("original_line") or 0
            comment_body = c.get("body") or ""
            out.append(
                Finding(
                    path=path,
                    line=int(line),
                    body=comment_body,
                    source="comment",
                    comment_id=c.get("id"),
                    review_id=review_id,
                    is_repeat=bool(REPEAT_HINT_RE.search(comment_body)),
                )
            )
    return out


def inventory_for_review(
    review: dict[str, Any], comments: list[dict[str, Any]]
) -> list[Finding]:
    rid = int(review["id"])
    findings = findings_from_comments(comments, rid)
    findings.extend(parse_suppressed(review.get("body") or "", rid))
    return dedupe_findings(findings)


def dedupe_findings(findings: list[Finding]) -> list[Finding]:
    """Keep one finding per path + nearby line, preferring inline comments."""
    best: list[Finding] = []
    for f in findings:
        match_idx = next(
            (
                i
                for i, other in enumerate(best)
                if other.path == f.path and abs(other.line - f.line) <= LINE_PROXIMITY
            ),
            None,
        )
        if match_idx is None:
            best.append(f)
        else:
            best[match_idx] = preferred_finding(best[match_idx], f)
    return best


def preferred_finding(previous: Finding, candidate: Finding) -> Finding:
    prefer_candidate = (
        previous.source == "suppressed" and candidate.source == "comment"
    ) or (previous.is_repeat and not candidate.is_repeat)
    return candidate if prefer_candidate else previous


def load_review_data(config: Config) -> ReviewData | None:
    reviews = gh_paginated(f"repos/{config.repo}/pulls/{config.pr}/reviews")
    comments = gh_paginated(f"repos/{config.repo}/pulls/{config.pr}/comments")
    copilot_reviews = sorted_copilot_reviews(reviews)
    current = find_current_review(reviews, copilot_reviews, config.review_id)
    if current is None:
        log(f"review {config.review_id} not found on PR {config.pr}")
        return None
    prior = reviews_before(copilot_reviews, current)
    return ReviewData(current=current, prior=prior, comments=comments)


def sorted_copilot_reviews(
    reviews: list[dict[str, Any]],
) -> list[dict[str, Any]]:
    allowed_states = {"COMMENTED", "APPROVED", "CHANGES_REQUESTED", "DISMISSED"}
    copilot_reviews = [
        review
        for review in reviews
        if is_copilot_user((review.get("user") or {}).get("login"))
        and (review.get("state") or "").upper() in allowed_states
    ]
    return sorted(copilot_reviews, key=review_submitted_at)


def review_submitted_at(review: dict[str, Any]) -> str:
    return review.get("submitted_at") or ""


def find_current_review(
    reviews: list[dict[str, Any]],
    copilot_reviews: list[dict[str, Any]],
    review_id: int,
) -> dict[str, Any] | None:
    current = next(
        (review for review in copilot_reviews if int(review["id"]) == review_id),
        None,
    )
    if current is None:
        current = next(
            (review for review in reviews if int(review["id"]) == review_id),
            None,
        )
    return current


def reviews_before(
    reviews: list[dict[str, Any]], current: dict[str, Any]
) -> list[dict[str, Any]]:
    current_submitted = review_submitted_at(current)
    prior = [
        review
        for review in reviews
        if int(review["id"]) != int(current["id"])
        and review_submitted_at(review) < current_submitted
    ]
    return sorted(prior, key=review_submitted_at)


def build_inventories(data: ReviewData) -> Inventories:
    current = inventory_for_review(data.current, data.comments)
    all_prior_findings: list[Finding] = []
    for review in data.prior:
        all_prior_findings.extend(inventory_for_review(review, data.comments))
    all_prior = dedupe_findings(all_prior_findings)
    original_prior = [finding for finding in all_prior if not finding.is_repeat]
    return Inventories(
        current=current,
        prior=original_prior or all_prior,
        all_prior=all_prior,
    )


def build_change_state(config: Config, data: ReviewData) -> ChangeState:
    first_pass = not data.prior
    prior_sha = (data.prior[-1].get("commit_id") if data.prior else "") or ""
    current_sha = data.current.get("commit_id") or ""
    code_changed = bool(prior_sha and current_sha and prior_sha != current_sha)
    touched = changed_files(config.repo, prior_sha, current_sha, code_changed)
    return ChangeState(
        first_pass=first_pass,
        code_changed=code_changed,
        prior_sha=prior_sha,
        current_sha=current_sha,
        touched=touched,
    )


def changed_files(
    repo: str, prior_sha: str, current_sha: str, code_changed: bool
) -> set[str]:
    return files_touched_since(repo, prior_sha, current_sha) if code_changed else set()


def files_touched_since(
    repo: str, prior_sha: str, current_sha: str
) -> set[str]:
    if not prior_sha or not current_sha or prior_sha == current_sha:
        return set()
    try:
        compare = gh_json([f"repos/{repo}/compare/{prior_sha}...{current_sha}"])
    except subprocess.CalledProcessError as exc:
        log(f"compare failed: {exc.stderr}")
        return set()
    files = compare.get("files") or []
    return {f.get("filename") for f in files if f.get("filename")}


def log_state(state: ChangeState, inventories: Inventories) -> None:
    log(
        f"first_pass={state.first_pass} code_changed={state.code_changed} "
        f"prior_findings={len(inventories.prior)} "
        f"current_findings={len(inventories.current)} "
        f"prior_sha={state.prior_sha[:8]} current_sha={state.current_sha[:8]}"
    )
