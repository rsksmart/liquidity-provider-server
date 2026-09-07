"""Thin wrappers around the `gh` CLI for GitHub API access."""

from __future__ import annotations

import json
import subprocess
from typing import Any


def gh_json(args: list[str]) -> Any:
    cmd = ["gh", "api", *args]
    result = subprocess.run(cmd, check=True, capture_output=True, text=True)
    if not result.stdout.strip():
        return None
    return json.loads(result.stdout)


def gh_paginated(path: str) -> list[Any]:
    page_size = 100
    items: list[Any] = []
    page = 1
    batch = gh_list_page(path, page, page_size)
    while len(batch) == page_size:
        items.extend(batch)
        page += 1
        batch = gh_list_page(path, page, page_size)
    items.extend(batch)
    return items


def gh_list_page(path: str, page: int, page_size: int) -> list[Any]:
    batch = gh_json([f"{path}?per_page={page_size}&page={page}"])
    if not batch:
        return []
    if not isinstance(batch, list):
        raise RuntimeError(f"expected list from {path}, got {type(batch)}")
    return batch
