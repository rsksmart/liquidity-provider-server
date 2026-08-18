#!/usr/bin/env python3
"""Build the Copilot CLI prompt and optionally post the lessons comment.

Env:
  REPO, PR_NUMBER, REVIEW_ID
  SCHEMA_PATH          path to lesson-schema.md
  PROMPT_TEMPLATE_PATH path to lessons-prompt.md
  FINDINGS_PATH        path to lessons-input.json
  PROMPT_OUT           where to write the filled prompt
  LESSONS_BODY_PATH    path to model output (for --post)
  DRY_RUN=1            print instead of posting
"""

from __future__ import annotations

import argparse
import json
import os
import subprocess
import sys
from pathlib import Path


LESSONS_HEADER = "## Suggested lessons for the Rootstock lessons DB"


def log(msg: str) -> None:
    print(msg, file=sys.stderr)


def build_prompt() -> int:
    schema = Path(os.environ["SCHEMA_PATH"]).read_text(encoding="utf-8")
    template = Path(os.environ["PROMPT_TEMPLATE_PATH"]).read_text(encoding="utf-8")
    findings = Path(os.environ["FINDINGS_PATH"]).read_text(encoding="utf-8")
    review_id = os.environ["REVIEW_ID"]
    out = Path(os.environ["PROMPT_OUT"])

    data = json.loads(findings)
    if not data:
        log("no findings; skipping prompt build")
        return 0

    filled = (
        template.replace("<<<SCHEMA>>>", schema.strip())
        .replace("<<<REVIEW_ID>>>", review_id)
        .replace("<<<FINDINGS>>>", findings.strip())
    )
    out.write_text(filled, encoding="utf-8")
    log(f"wrote prompt to {out}")
    return 0


def post_lessons() -> int:
    repo = os.environ.get("REPO") or os.environ.get("GITHUB_REPOSITORY")
    pr = os.environ["PR_NUMBER"]
    review_id = os.environ["REVIEW_ID"]
    body_path = Path(os.environ["LESSONS_BODY_PATH"])
    dry_run = os.environ.get("DRY_RUN") == "1"

    if not repo:
        log("REPO required")
        return 2
    if not body_path.exists():
        log(f"missing lessons body at {body_path}")
        return 1

    body = body_path.read_text(encoding="utf-8").strip()
    # Strip accidental fences the model sometimes wraps around the whole comment.
    if body.startswith("```") and body.endswith("```"):
        body = body.strip("`")
        if body.startswith("markdown"):
            body = body[len("markdown") :].lstrip()

    marker = f"<!-- ccr-followup:lessons:{review_id} -->"
    if marker not in body:
        body = marker + "\n" + body

    if LESSONS_HEADER not in body:
        log("model output does not look like a lessons comment; refusing to post")
        log(body[:500])
        return 1

    if dry_run:
        log("DRY_RUN: would post lessons comment:")
        print(body)
        return 0

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
        log(proc.stderr)
        return proc.returncode
    log("posted lessons comment")
    return 0


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("command", choices=["build-prompt", "post"])
    args = parser.parse_args()
    if args.command == "build-prompt":
        return build_prompt()
    return post_lessons()


if __name__ == "__main__":
    raise SystemExit(main())
