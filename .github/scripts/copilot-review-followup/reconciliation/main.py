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

import subprocess

from .config import load_config, log, write_outputs
from .parse import build_change_state, build_inventories, load_review_data, log_state
from .post import reconcile, write_lessons_input


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
    posted_reconciliation = reconcile(config, data, inventories, state)
    payload = write_lessons_input(config, inventories)
    write_outputs(config, posted_reconciliation, payload)
    return 0


if __name__ == "__main__":
    try:
        raise SystemExit(main())
    except subprocess.CalledProcessError as exc:
        log(exc.stderr or str(exc))
        raise SystemExit(exc.returncode or 1)
