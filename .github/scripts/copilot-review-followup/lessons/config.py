"""Environment config for lessons commands."""

from __future__ import annotations

import os
import sys
from dataclasses import dataclass
from pathlib import Path


@dataclass(frozen=True)
class PromptConfig:
    schema_path: Path
    template_path: Path
    findings_path: Path
    review_id: str
    prompt_out: Path


@dataclass(frozen=True)
class PostConfig:
    repo: str
    pr: int
    review_id: str
    body_path: Path
    dry_run: bool


def log(msg: str) -> None:
    print(msg, file=sys.stderr)


def load_prompt_config() -> PromptConfig | None:
    try:
        return PromptConfig(
            schema_path=Path(os.environ["SCHEMA_PATH"]),
            template_path=Path(os.environ["PROMPT_TEMPLATE_PATH"]),
            findings_path=Path(os.environ["FINDINGS_PATH"]),
            review_id=os.environ["REVIEW_ID"],
            prompt_out=Path(os.environ["PROMPT_OUT"]),
        )
    except KeyError as exc:
        log(f"missing required env: {exc.args[0]}")
        return None


def load_post_config() -> PostConfig | None:
    repo = os.environ.get("REPO") or os.environ.get("GITHUB_REPOSITORY")
    if not repo:
        log("REPO / GITHUB_REPOSITORY is required")
        return None
    try:
        return PostConfig(
            repo=repo,
            pr=int(os.environ["PR_NUMBER"]),
            review_id=os.environ["REVIEW_ID"],
            body_path=Path(os.environ["LESSONS_BODY_PATH"]),
            dry_run=os.environ.get("DRY_RUN") == "1",
        )
    except KeyError as exc:
        log(f"missing required env: {exc.args[0]}")
        return None
