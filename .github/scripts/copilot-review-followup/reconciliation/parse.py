"""Load and parse Copilot review data into finding inventories."""

from __future__ import annotations

import re
import subprocess
from dataclasses import dataclass, field
from typing import Any

from .config import Config, log
from .finding import (
    COPILOT_LOGINS,
    DETAILS_SUPPRESSED_RE,
    HEADING_SUPPRESSED_RE,
    LINE_PROXIMITY,
    PREVIOUSLY_MISSED_RE,
    REPEAT_HINT_RE,
    REVIEW_STATS_RE,
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
    commit_id: str = ""
    findings: list[Finding] = field(default_factory=list)
    path: str | None = None
    line: int = 0
    body: list[str] = field(default_factory=list)

    def consume(self, body_line: str) -> None:
        header = SUPPRESSED_FINDING_HEADER_RE.match(body_line)
        if header:
            self.start(header.group("path"), int(header.group("line")))
            return
        if (
            body_line.strip() == "</details>"
            or REVIEW_STATS_RE.match(body_line)
        ):
            self.finish()
            return
        if PREVIOUSLY_MISSED_RE.match(body_line):
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
                commit_id=self.commit_id,
                is_repeat=bool(REPEAT_HINT_RE.search(text)),
            )
        )


def is_copilot_user(login: str | None) -> bool:
    if not login:
        return False
    return login in COPILOT_LOGINS


def parse_suppressed(body: str, review_id: int, commit_id: str = "") -> list[Finding]:
    section = suppressed_section(body)
    if not section:
        return []
    parser = SuppressedParser(review_id, commit_id=commit_id)
    for body_line in section.splitlines():
        parser.consume(body_line)
    parser.finish()
    return parser.findings


def suppressed_section(body: str) -> str:
    """Slice the suppressed-comments block out of a review body.

    Copilot has used two layouts: a ``<summary>Suppressed comments`` details
    block, and a Markdown heading (often nested under Review details). The
    details form ends at ``</details>``. The heading form also ends at the next
    heading of the same or higher level, so a later ``**path:line**`` section
    cannot be ingested as another suppressed finding.
    """
    details = DETAILS_SUPPRESSED_RE.search(body)
    heading = HEADING_SUPPRESSED_RE.search(body)
    if details and (heading is None or details.start() <= heading.start()):
        return body[details.end() :].partition("</details>")[0]
    if heading is None:
        return ""
    after_marker = body[heading.end() :]
    level = len(heading.group("hashes"))
    peer = re.compile(rf"^#{{1,{level}}}\s", re.MULTILINE)
    stops = [len(after_marker)]
    details_end = after_marker.find("</details>")
    if details_end != -1:
        stops.append(details_end)
    peer_match = peer.search(after_marker)
    if peer_match:
        stops.append(peer_match.start())
    return after_marker[: min(stops)]


def suppressed_summary(lines: list[str]) -> str:
    raw = "\n".join(lines).strip()
    summary = re.split(r"\n```", raw, maxsplit=1)[0].strip()
    return re.sub(r"^\*\s+", "", summary)


def findings_from_comments(
    comments: list[dict[str, Any]], review_id: int, commit_id: str = ""
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
                    commit_id=commit_id,
                    is_repeat=bool(REPEAT_HINT_RE.search(comment_body)),
                )
            )
    return out


def inventory_for_review(
    review: dict[str, Any], comments: list[dict[str, Any]]
) -> list[Finding]:
    rid = int(review["id"])
    commit_id = review.get("commit_id") or ""
    findings = findings_from_comments(comments, rid, commit_id)
    findings.extend(parse_suppressed(review.get("body") or "", rid, commit_id))
    return dedupe_findings(findings)


def dedupe_findings(findings: list[Finding]) -> list[Finding]:
    """Deduplicate the same nearby finding, preferring inline comments."""
    best: list[Finding] = []
    for f in findings:
        match_idx = next(
            (
                i
                for i, other in enumerate(best)
                if is_duplicate_finding(other, f)
            ),
            None,
        )
        if match_idx is None:
            best.append(f)
        else:
            best[match_idx] = preferred_finding(best[match_idx], f)
    return best


TOKEN_OVERLAP_THRESHOLD = 0.35


def is_duplicate_finding(previous: Finding, candidate: Finding) -> bool:
    if (
        previous.comment_id is not None
        and previous.comment_id == candidate.comment_id
    ):
        return True
    if previous.path != candidate.path:
        return False
    if abs(previous.line - candidate.line) > LINE_PROXIMITY:
        return False
    previous_tokens = significant_tokens(previous.body)
    candidate_tokens = significant_tokens(candidate.body)
    if not previous_tokens or not candidate_tokens:
        return False
    overlap = len(previous_tokens & candidate_tokens)
    return (
        overlap / min(len(previous_tokens), len(candidate_tokens))
        >= TOKEN_OVERLAP_THRESHOLD
    )


def significant_tokens(body: str) -> set[str]:
    return set(re.findall(r"[a-z0-9_]{4,}", body.lower()))


def preferred_finding(previous: Finding, candidate: Finding) -> Finding:
    prefer_candidate = (
        previous.source == "suppressed" and candidate.source == "comment"
    ) or (previous.is_repeat and not candidate.is_repeat)
    return candidate if prefer_candidate else previous


def load_review_data(config: Config) -> ReviewData | None:
    reviews = gh_paginated(f"repos/{config.repo}/pulls/{config.pr}/reviews")
    review_comments = gh_paginated(
        f"repos/{config.repo}/pulls/{config.pr}/comments"
    )
    issue_comments = gh_paginated(
        f"repos/{config.repo}/issues/{config.pr}/comments"
    )
    comments = review_comments + issue_comments
    copilot_reviews = sorted_copilot_reviews(reviews)
    current = find_current_review(copilot_reviews, config.review_id)
    if current is None:
        log(f"Copilot review {config.review_id} not found on PR {config.pr}")
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
    copilot_reviews: list[dict[str, Any]],
    review_id: int,
) -> dict[str, Any] | None:
    return next(
        (review for review in copilot_reviews if int(review["id"]) == review_id),
        None,
    )


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
    return Inventories(
        current=current,
        prior=all_prior,
        all_prior=all_prior,
    )


def build_change_state(config: Config, data: ReviewData) -> ChangeState:
    first_pass = not data.prior
    prior_sha = (data.prior[-1].get("commit_id") if data.prior else "") or ""
    current_sha = data.current.get("commit_id") or ""
    code_changed = bool(prior_sha and current_sha and prior_sha != current_sha)
    origin_shas: list[str] = []
    for review in data.prior:
        sha = review.get("commit_id") or ""
        if sha and sha not in origin_shas:
            origin_shas.append(sha)
    touched_since = {
        sha: files_touched_since(config.repo, sha, current_sha) for sha in origin_shas
    }
    return ChangeState(
        first_pass=first_pass,
        code_changed=code_changed,
        prior_sha=prior_sha,
        current_sha=current_sha,
        touched_since=touched_since,
    )


COMPARE_FILE_CAP = 300


def files_touched_since(
    repo: str, prior_sha: str, current_sha: str
) -> set[str] | None:
    """Files changed between the two commits, or None if that is unknowable.

    An empty set means nothing changed, which is evidence a finding is still
    open. A failed or truncated compare means we know nothing, so it must not
    be mistaken for an untouched file; those cases return None.
    """
    if not prior_sha or not current_sha:
        return None
    if prior_sha == current_sha:
        return set()
    try:
        compare = gh_json([f"repos/{repo}/compare/{prior_sha}...{current_sha}"])
    except subprocess.CalledProcessError as exc:
        log(f"compare failed: {exc.stderr}")
        return None
    if not isinstance(compare, dict):
        return None
    files = compare.get("files") or []
    # GitHub does not paginate compare `files`, it caps the array at 300 and sets truncated
    if compare.get("truncated") or len(files) >= COMPARE_FILE_CAP:
        log(
            f"compare {prior_sha[:8]}...{current_sha[:8]} file list truncated "
            f"({len(files)} files); treating change set as unknown"
        )
        return None
    return compare_filenames(files)


def compare_filenames(files: list[Any]) -> set[str]:
    names: set[str] = set()
    for entry in files:
        for key in ("filename", "previous_filename"):
            name = entry.get(key)
            if name:
                names.add(name)
    return names


def log_state(state: ChangeState, inventories: Inventories) -> None:
    origins = len(state.touched_since)
    log(
        f"first_pass={state.first_pass} code_changed={state.code_changed} "
        f"origin_shas={origins} "
        f"prior_findings={len(inventories.prior)} "
        f"current_findings={len(inventories.current)} "
        f"prior_sha={state.prior_sha[:8]} current_sha={state.current_sha[:8]}"
    )
