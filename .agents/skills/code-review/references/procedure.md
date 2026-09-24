# Full review procedure notes

This file keeps the longer original wording for Steps 3–6. Follow [SKILL.md](../SKILL.md) first.
Read this when you need the original detail.

## Inventory sources

Posted comments come from `get_review_comments`. Suppressed findings live in each prior Copilot
review `body` under "Suppressed comments" (sometimes nested under "Review details", sometimes in a
`<summary>Suppressed comments</summary>` block). The follow-up workflow inventories suppressed
findings for reconciliation and never uses them as sources for lesson cards.

## Classification evidence

Record per item: the file and line, a one-sentence restatement of the original point, the
classification, and the evidence. For **Unverified**, record what blocked the determination.

When you cannot read the code behind one finding, classify it **Unverified** and name what failed.
A tool failure never silences a finding and never counts as evidence of a fix.

An empty **Addressed** section is a perfectly acceptable outcome. A finding that was never fixed
and never mentioned is not.

## Incremental diff fallbacks

If the review carries no `commit_id`, use `submitted_at` as the cutoff. Fall back to the full pull
request diff only when neither git nor the API can produce a change set, and say plainly that you
did. When the pull request adds a file, the full diff contains every line of it; the inventory
check is what stops old findings from looking new.

Do not reopen reconciliation in Step 4. The Step 3 classifications stand as recorded.

## Reporting

Silence is never a classification. Every finding you did not verify as fixed reaches the reader as
either a **Not addressed** comment or an **Unverified** question until the code shows it was fixed
or a human declines it.

If a decline rests on a factual mistake, say so once, in one sentence, next to that entry. Do not
re-post the inline comment and do not reclassify the finding away from **Declined**.

If the same point has come up more than twice, say how many times. For an **Unverified** item,
state the history but give no verdict.
