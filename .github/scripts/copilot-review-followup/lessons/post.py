"""Post the authored lessons comment to the PR."""

from __future__ import annotations

import json
import re
import subprocess

from .config import PostConfig, log

LESSONS_HEADER = "## Suggested lessons for the Rootstock lessons DB"
LESSONS_INSTALL = (
    "Copy each block into `skills/lessons/<category>/<id>.yaml`, then add a row to\n"
    "`skills/lessons/README.md`."
)
NO_LESSONS = "No eligible new findings for lesson cards on this pass."
MAX_LESSONS = 5
REQUIRED_CARD_FIELDS = ("title", "tags", "triggers", "problem", "bad", "good", "why")
DETAIL_RE = re.compile(
    r"<details>\n"
    r"<summary>(?P<category>[a-z0-9][a-z0-9-]*) / "
    r"(?P<summary_id>[a-z0-9][a-z0-9-]*) "
    r"\((?P<summary_severity>high|medium|low)\)</summary>\n\n"
    r"```yaml\n(?P<card>.*?)\n```\n\n"
    r"</details>",
    re.DOTALL,
)
CARD_RE = re.compile(
    r"id: (?P<id>[a-z0-9][a-z0-9-]*)\n"
    r"title: (?P<title>[^\n]+)\n"
    r"tags:\n(?P<tags>(?:  - [^\n]+\n)+)"
    r"severity: (?P<severity>high|medium|low)\n"
    r"frequency: 1\n"
    r"trigger_patterns:\n(?P<triggers>(?:  - [^\n]+\n)+)"
    r"problem: \|\n(?P<problem>(?:  [^\n]*(?:\n|$))+?)"
    r"bad: \|\n(?P<bad>(?:  [^\n]*(?:\n|$))+?)"
    r"good: \|\n(?P<good>(?:  [^\n]*(?:\n|$))+?)"
    r"why: \|\n(?P<why>(?:  [^\n]*(?:\n|$))+)",
)


def post_lessons(config: PostConfig) -> int:
    if not config.body_path.exists():
        log(f"missing lessons body at {config.body_path}")
        return 1

    body = config.body_path.read_text(encoding="utf-8").strip()
    marker = f"<!-- ccr-followup:lessons:{config.review_id} -->"
    error = validate_lessons_comment(body, marker)
    if error:
        log(f"invalid lessons comment: {error}; refusing to post")
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


def validate_lessons_comment(body: str, marker: str) -> str:
    content = lessons_content_after_header(body, marker)
    if content is None:
        return "missing or incorrect review marker/header"
    if content == NO_LESSONS:
        return ""
    blocks = lesson_blocks_after_install(content)
    if blocks is None:
        return "missing or incorrect install instructions"
    if not blocks:
        return "no lesson cards"
    matches, parse_error = parse_lesson_blocks(blocks)
    if parse_error:
        return parse_error
    return validate_lesson_cards(matches)


def lessons_content_after_header(body: str, marker: str) -> str | None:
    prefix = f"{marker}\n{LESSONS_HEADER}\n\n"
    if not body.startswith(prefix):
        return None
    return body[len(prefix) :]


def lesson_blocks_after_install(content: str) -> str | None:
    install_prefix = LESSONS_INSTALL + "\n\n"
    if not content.startswith(install_prefix):
        return None
    return content[len(install_prefix) :]


def parse_lesson_blocks(
    blocks: str,
) -> tuple[list[re.Match[str]], str]:
    matches: list[re.Match[str]] = []
    remaining = blocks
    while remaining:
        match = lesson_block_at(remaining)
        if match is None:
            return [], "extra prose or malformed details block"
        matches.append(match)
        remaining = remaining[match.end() :]
        if remaining:
            if not remaining.startswith("\n\n"):
                return [], "lesson blocks must be separated by one blank line"
            remaining = remaining[2:]
    if len(matches) > MAX_LESSONS:
        return [], "more than five lesson cards"
    return matches, ""


def lesson_block_at(blocks: str) -> re.Match[str] | None:
    return DETAIL_RE.match(blocks)


def validate_lesson_cards(matches: list[re.Match[str]]) -> str:
    seen_ids: set[str] = set()
    for match in matches:
        error = validate_lesson_card(match, seen_ids)
        if error:
            return error
    return ""


def validate_lesson_card(match: re.Match[str], seen_ids: set[str]) -> str:
    card = CARD_RE.fullmatch(match.group("card"))
    if card is None:
        return "card does not match the required YAML field structure"
    lesson_id = card.group("id")
    if lesson_id in seen_ids:
        return f"duplicate lesson id {lesson_id}"
    seen_ids.add(lesson_id)
    if lesson_id != match.group("summary_id"):
        return f"summary id does not match card id {lesson_id}"
    if card.group("severity") != match.group("summary_severity"):
        return f"summary severity does not match card {lesson_id}"
    return empty_card_field(card, lesson_id)


def empty_card_field(card: re.Match[str], lesson_id: str) -> str:
    for field in REQUIRED_CARD_FIELDS:
        if not card.group(field).strip():
            return f"empty {field} in card {lesson_id}"
    return ""
