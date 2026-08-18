#!/usr/bin/env python3
"""Post Copilot code-review reconciliation and emit lessons-input.json.

Uses `gh` for GitHub API access. Env:
  REPO          owner/repo (required)
  PR_NUMBER     pull request number (required)
  REVIEW_ID     triggering Copilot review id (required)
  GH_TOKEN      token with pull-requests:write (optional; gh uses its own auth)
  OUT_DIR       directory for outputs (default: .)
  DRY_RUN       if "1", do not post comments
"""

from __future__ import annotations

import json
import os
import re
import subprocess
import sys
from dataclasses import dataclass, field
from pathlib import Path
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
MARKER_RECON = "<!-- ccr-followup:recon:{review_id} -->"
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
class Config:
    repo: str
    pr: int
    review_id: int
    out_dir: Path
    dry_run: bool


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


def log(msg: str) -> None:
    print(msg, file=sys.stderr)


def gh_json(args: list[str]) -> Any:
    cmd = ["gh", "api", *args]
    result = subprocess.run(cmd, check=True, capture_output=True, text=True)
    if not result.stdout.strip():
        return None
    return json.loads(result.stdout)


def gh_paginated(path: str) -> list[Any]:
    items: list[Any] = []
    page = 1
    while True:
        batch = gh_json([f"{path}?per_page=100&page={page}"])
        if not batch:
            break
        if not isinstance(batch, list):
            raise RuntimeError(f"expected list from {path}, got {type(batch)}")
        items.extend(batch)
        if len(batch) < 100:
            break
        page += 1
    return items


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


def classify(
    prior: list[Finding],
    current: list[Finding],
    comments: list[dict[str, Any]],
    touched: set[str],
) -> Classified:
    result = Classified()
    for finding in prior:
        category, item = classify_finding(finding, current, comments, touched)
        getattr(result, category).append(item)
    return result


def classify_finding(
    finding: Finding,
    current: list[Finding],
    comments: list[dict[str, Any]],
    touched: set[str],
) -> tuple[str, Finding | tuple[Finding, str, str]]:
    decline = find_decline(finding, comments)
    if decline:
        return "declined", (finding, decline[0], decline[1])
    if any(similar(finding, candidate) for candidate in current):
        return "not_addressed", finding
    if finding.path in touched:
        return "addressed", finding
    return "unverified", finding


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


def render_recon(classified: Classified, review_id: int) -> str:
    sections: list[str] = [
        MARKER_RECON.format(review_id=review_id),
        "## Previous Copilot review follow-up",
        "",
    ]

    def add(title: str, lines: list[str]) -> None:
        if not lines:
            return
        sections.append(f"### {title}")
        sections.extend(lines)
        sections.append("")

    add(
        "Addressed",
        [f"- `{f.key}` — {f.summary()}" for f in classified.addressed],
    )
    add(
        "Not addressed",
        [f"- `{f.key}` — {f.summary()}" for f in classified.not_addressed],
    )
    add(
        "Declined",
        [
            f"- `{f.key}` — declined by @{user}: {reason}"
            for f, user, reason in classified.declined
        ],
    )
    add(
        "Unverified",
        [f"- `{f.key}` — {f.summary()}" for f in classified.unverified],
    )
    return "\n".join(sections).rstrip() + "\n"


def comment_exists(repo: str, pr: int, marker: str) -> bool:
    comments = gh_paginated(f"repos/{repo}/issues/{pr}/comments")
    return any(marker in (c.get("body") or "") for c in comments)


def post_comment(repo: str, pr: int, body: str, dry_run: bool) -> None:
    if dry_run:
        log("DRY_RUN: would post comment:")
        print(body)
        return
    payload = json.dumps({"body": body})
    proc = subprocess.run(
        [
            "gh",
            "api",
            f"repos/{repo}/issues/{pr}/comments",
            "-X",
            "POST",
            "--input",
            "-",
        ],
        input=payload,
        capture_output=True,
        text=True,
    )
    if proc.returncode != 0:
        raise RuntimeError(f"failed to post comment: {proc.stderr}")
    log("posted reconciliation comment")


def write_github_output(name: str, value: str) -> None:
    path = os.environ.get("GITHUB_OUTPUT")
    if not path:
        return
    with open(path, "a", encoding="utf-8") as fh:
        if "\n" in value:
            fh.write(f"{name}<<EOF\n{value}\nEOF\n")
        else:
            fh.write(f"{name}={value}\n")


def load_config() -> Config | None:
    repo = os.environ.get("REPO") or os.environ.get("GITHUB_REPOSITORY")
    if not repo:
        log("REPO / GITHUB_REPOSITORY is required")
        return None
    return Config(
        repo=repo,
        pr=int(os.environ["PR_NUMBER"]),
        review_id=int(os.environ["REVIEW_ID"]),
        out_dir=Path(os.environ.get("OUT_DIR") or "."),
        dry_run=os.environ.get("DRY_RUN") == "1",
    )


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


def log_state(state: ChangeState, inventories: Inventories) -> None:
    log(
        f"first_pass={state.first_pass} code_changed={state.code_changed} "
        f"prior_findings={len(inventories.prior)} "
        f"current_findings={len(inventories.current)} "
        f"prior_sha={state.prior_sha[:8]} current_sha={state.current_sha[:8]}"
    )


def reconcile(
    config: Config,
    data: ReviewData,
    inventories: Inventories,
    state: ChangeState,
) -> bool:
    reason = reconciliation_skip_reason(state)
    if reason:
        log(f"skipping reconciliation ({reason})")
        return False
    marker = MARKER_RECON.format(review_id=config.review_id)
    if comment_exists(config.repo, config.pr, marker):
        log(f"reconciliation already posted for review {config.review_id}")
        return False
    classified = classify(
        inventories.prior,
        inventories.current,
        data.comments,
        state.touched,
    )
    if not classified.has_findings():
        log("no classified findings; skipping reconciliation post")
        return False
    body = render_recon(classified, config.review_id)
    (config.out_dir / "recon.md").write_text(body, encoding="utf-8")
    post_comment(config.repo, config.pr, body, config.dry_run)
    return True


def reconciliation_skip_reason(state: ChangeState) -> str:
    if state.first_pass:
        return "first pass"
    if not state.code_changed:
        return "no code changes since prior Copilot review"
    return ""


def write_lessons_input(
    config: Config, inventories: Inventories
) -> list[dict[str, Any]]:
    eligible = new_findings(inventories.all_prior, inventories.current)
    marker = MARKER_LESSONS.format(review_id=config.review_id)
    if comment_exists(config.repo, config.pr, marker):
        log(f"lessons already posted for review {config.review_id}")
        eligible = []
    payload = [finding_payload(finding) for finding in eligible]
    lessons_path = config.out_dir / "lessons-input.json"
    lessons_path.write_text(json.dumps(payload, indent=2) + "\n", encoding="utf-8")
    log(f"wrote {lessons_path} with {len(payload)} eligible findings")
    return payload


def finding_payload(finding: Finding) -> dict[str, Any]:
    return {
        "path": finding.path,
        "line": finding.line,
        "body": finding.body,
        "source": finding.source,
        "comment_id": finding.comment_id,
    }


def write_outputs(config: Config, posted_recon: bool, payload: list[Any]) -> None:
    write_github_output("posted_recon", "true" if posted_recon else "false")
    write_github_output("has_lessons", "true" if payload else "false")
    write_github_output("lessons_count", str(len(payload)))
    write_github_output("review_id", str(config.review_id))
    write_github_output("pr_number", str(config.pr))


def main() -> int:
    config = load_config()
    if config is None:
        return 2
    config.out_dir.mkdir(parents=True, exist_ok=True)
    data = load_review_data(config)
    if data is None:
        return 1
    inventories = build_inventories(data)
    state = build_change_state(config, data)
    log_state(state, inventories)
    posted_recon = reconcile(config, data, inventories, state)
    payload = write_lessons_input(config, inventories)
    write_outputs(config, posted_recon, payload)
    return 0


if __name__ == "__main__":
    try:
        raise SystemExit(main())
    except subprocess.CalledProcessError as exc:
        log(exc.stderr or str(exc))
        raise SystemExit(exc.returncode or 1)
