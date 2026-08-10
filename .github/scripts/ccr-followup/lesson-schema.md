# Lesson card schema

Use this schema for every suggested lesson. Do not invent extra top-level fields.

```yaml
id: short-stable-id
title: Human-readable title
tags:
  - language-or-domain
  - bug-class
severity: high | medium | low
frequency: 1
trigger_patterns:
  - code pattern or keyword
problem: |
  What goes wrong.
bad: |
  Example bad pattern.
good: |
  Example safer pattern.
why: |
  Why this matters.
```

## Field rules

- `id`: lowercase kebab-case, stable, unique within its category folder.
- `title`: one sentence a human can scan in an index table.
- `tags`: short labels used for retrieval (language, framework, subsystem, bug class).
- `severity`: `high`, `medium`, or `low`.
- `frequency`: start at `1` for a first observation.
- `trigger_patterns`: concrete tokens or patterns that should surface this lesson later.
- `problem` / `bad` / `good` / `why`: use block scalars (`|`). Keep examples minimal and real.

## Install path

Put the install path in the lessons comment; do not paste this reference file
into the pull request.

Save each card as `skills/lessons/<category>/<id>.yaml`, then add a row to
`skills/lessons/README.md` under the matching category table.
Common categories: `frontend`, `go`, `reliability`, `security`, `testing`,
`process`, `observability`, `sql`, `rpc`.

## Example

```yaml
id: nil-callfee-panic-in-aggregation
title: Nil CallFee must be skipped before AsBigInt in volume aggregation
tags:
  - go
  - reliability
  - nil
severity: high
frequency: 1
trigger_patterns:
  - CallFee
  - AsBigInt
  - volumeByDay
  - TotalPeginVolume
problem: |
  Aggregating quote volumes by adding CallFee.AsBigInt() without a nil check
  panics on incomplete records. Totals that sum only Value then disagree with
  per-day buckets that included CallFee.
bad: |
  volumeByDay[day].Add(volumeByDay[day], pair.Quote.Value.AsBigInt())
  volumeByDay[day].Add(volumeByDay[day], pair.Quote.CallFee.AsBigInt())
  total.Add(total, pair.Quote.Value.AsBigInt())
good: |
  for _, amount := range []*entities.Wei{pair.Quote.Value, pair.Quote.CallFee} {
      if amount != nil {
          volumeByDay[day].Add(volumeByDay[day], amount.AsBigInt())
          total.Add(total, amount.AsBigInt())
      }
  }
why: |
  CallFee is a pointer and older or incomplete quotes leave it nil. The same
  amounts must feed both the daily bucket and the running total so the report
  stays internally consistent.
```
