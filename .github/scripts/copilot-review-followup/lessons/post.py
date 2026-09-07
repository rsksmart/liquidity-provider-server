"""Post the authored lessons comment to the PR."""

from __future__ import annotations

import json
import subprocess

from .config import PostConfig, log

LESSONS_HEADER = "## Suggested lessons for the Rootstock lessons DB"


def post_lessons(config: PostConfig) -> int:
    if not config.body_path.exists():
        log(f"missing lessons body at {config.body_path}")
        return 1

    body = config.body_path.read_text(encoding="utf-8").strip()
    # Strip accidental fences the model sometimes wraps around the whole comment.
    if body.startswith("```") and body.endswith("```"):
        body = body.strip("`")
        if body.startswith("markdown"):
            body = body[len("markdown") :].lstrip()

    marker = f"<!-- ccr-followup:lessons:{config.review_id} -->"
    if marker not in body:
        body = marker + "\n" + body

    if LESSONS_HEADER not in body:
        log("model output does not look like a lessons comment; refusing to post")
        log(body[:500])
        return 1

    if config.dry_run:
        log("DRY_RUN: would post lessons comment:")
        print(body)
        return 0

    payload = json.dumps({"body": body})
    proc = subprocess.run(
        [
            "gh",
            "api",
            f"repos/{config.repo}/issues/{config.pr}/comments",
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
        log(proc.stderr)
        return proc.returncode
    log("posted lessons comment")
    return 0
