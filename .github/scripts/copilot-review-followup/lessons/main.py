"""CLI for building and posting lessons comments.

Env:
  REPO, PR_NUMBER, REVIEW_ID
  SCHEMA_PATH          path to lesson-schema.md
  PROMPT_TEMPLATE_PATH path to lessons-prompt.md
  FINDINGS_PATH        path to lessons-input.json
  PROMPT_OUT           where to write the filled prompt
  LESSONS_BODY_PATH    path to model output (for post)
  DRY_RUN=1            print instead of posting
"""

from __future__ import annotations

import argparse

from .config import load_post_config, load_prompt_config
from .post import post_lessons
from .prompt import build_prompt


def run_build_prompt() -> int:
    config = load_prompt_config()
    if config is None:
        return 2
    return build_prompt(config)


def run_post_lessons() -> int:
    config = load_post_config()
    if config is None:
        return 2
    return post_lessons(config)


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("command", choices=["build-prompt", "post"])
    args = parser.parse_args()
    if args.command == "build-prompt":
        return run_build_prompt()
    return run_post_lessons()


if __name__ == "__main__":
    raise SystemExit(main())
