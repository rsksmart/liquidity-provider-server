"""Environment config and GitHub Actions outputs."""

from __future__ import annotations

import os
import sys
from dataclasses import dataclass
from pathlib import Path
from typing import Any


@dataclass(frozen=True)
class Config:
    repo: str
    pr: int
    review_id: int
    out_dir: Path
    dry_run: bool
    github_output: Path | None


def log(msg: str) -> None:
    print(msg, file=sys.stderr)


def load_config() -> Config | None:
    repo = os.environ.get("REPO") or os.environ.get("GITHUB_REPOSITORY")
    if not repo:
        log("REPO / GITHUB_REPOSITORY is required")
        return None
    github_output = os.environ.get("GITHUB_OUTPUT")
    return Config(
        repo=repo,
        pr=int(os.environ["PR_NUMBER"]),
        review_id=int(os.environ["REVIEW_ID"]),
        out_dir=Path(os.environ.get("OUT_DIR") or "."),
        dry_run=os.environ.get("DRY_RUN") == "1",
        github_output=Path(github_output) if github_output else None,
    )


def write_github_output(config: Config, name: str, value: str) -> None:
    if config.github_output is None:
        return
    with open(config.github_output, "a", encoding="utf-8") as fh:
        if "\n" in value:
            fh.write(f"{name}<<EOF\n{value}\nEOF\n")
        else:
            fh.write(f"{name}={value}\n")


def write_outputs(
    config: Config, posted_reconciliation: bool, payload: list[Any]
) -> None:
    write_github_output(
        config,
        "posted_reconciliation",
        "true" if posted_reconciliation else "false",
    )
    write_github_output(config, "has_lessons", "true" if payload else "false")
    write_github_output(config, "lessons_count", str(len(payload)))
    write_github_output(config, "review_id", str(config.review_id))
    write_github_output(config, "pr_number", str(config.pr))
