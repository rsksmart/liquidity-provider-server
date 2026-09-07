"""Fill the lessons prompt template from schema + findings."""

from __future__ import annotations

import json

from .config import PromptConfig, log


def build_prompt(config: PromptConfig) -> int:
    schema = config.schema_path.read_text(encoding="utf-8")
    template = config.template_path.read_text(encoding="utf-8")
    findings = config.findings_path.read_text(encoding="utf-8")

    data = json.loads(findings)
    if not data:
        log("no findings; skipping prompt build")
        return 0

    filled = (
        template.replace("<<<SCHEMA>>>", schema.strip())
        .replace("<<<REVIEW_ID>>>", config.review_id)
        .replace("<<<FINDINGS>>>", findings.strip())
    )
    config.prompt_out.write_text(filled, encoding="utf-8")
    log(f"wrote prompt to {config.prompt_out}")
    return 0
