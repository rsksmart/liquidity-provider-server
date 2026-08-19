"""Render and post reconciliation / lessons artifacts."""

from __future__ import annotations

import json
import subprocess
from typing import Any

from .config import Config, log
from .finding import (
    MARKER_LESSONS,
    MARKER_RECONCILIATION,
    ChangeState,
    Classified,
    Finding,
    Inventories,
    ReviewData,
)
from .github import gh_paginated
from .process import classify, new_findings, reconciliation_skip_reason


def render_reconciliation(classified: Classified, review_id: int) -> str:
    sections: list[str] = [
        MARKER_RECONCILIATION.format(review_id=review_id),
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
    marker = MARKER_RECONCILIATION.format(review_id=config.review_id)
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
    body = render_reconciliation(classified, config.review_id)
    (config.out_dir / "reconciliation.md").write_text(body, encoding="utf-8")
    post_comment(config.repo, config.pr, body, config.dry_run)
    return True


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
