# commandcode-catalog

An unofficial, machine-generated snapshot of Command Code's per-model monthly
allowances for the **Go**, **GOAT**, and **Pro** plans.

Command Code publishes these numbers as "boosted usage" deals on its pricing
page. This repository turns that page into a small, stable JSON file other tools
can consume. It is refreshed automatically every week.

- Source: <https://commandcode.ai/docs/resources/pricing-limits>
- File: [`catalog.json`](./catalog.json)
- Raw URL: `https://raw.githubusercontent.com/peach0x33a/commandcode-catalog/main/catalog.json`

> Not affiliated with Command Code. The numbers belong to them and can change at
> any time; the upstream credit numbers remain authoritative. This snapshot is
> provided as-is, with no warranty.

## Format

`catalog.json` is a single JSON document with one model per line, so a change to
one model is a one-line diff.

```json
{"schema":1,"source":"https://commandcode.ai/docs/resources/pricing-limits","captured":"2026-11-16","models":[
{"id":"claude-fable-5","name":"Claude Fable 5","category":"premium","go":6,"goat":20,"pro":20},
{"id":"tencent/hy4-preview","name":"Tencent Hy4 Preview","category":"opensource","go":6,"goat":20,"pro":30}
]}
```

| Field | Type | Meaning |
| --- | --- | --- |
| `schema` | int | Format version. Currently `1`. |
| `source` | string | Page the numbers were captured from. |
| `captured` | string | `YYYY-MM-DD` when the data last changed (not when it was last checked). |
| `models[].id` | string | Catalog id, e.g. `deepseek-v4.1-flash` or `tencent/hy4-preview`. |
| `models[].name` | string | Display name. May carry a `(latest)` suffix. |
| `models[].category` | string | `opensource` or `premium`. |
| `models[].go` | number | Monthly USD allowance on the Go plan. |
| `models[].goat` | number | Monthly USD allowance on the GOAT plan. |
| `models[].pro` | number | Monthly USD allowance on the Pro plan. |

Consumers should reject a file whose `schema` is newer than they understand,
rather than guessing at unknown fields. Today `pro` equals `goat + 10` for
open-source models and equals `goat` for premium models, but that is an upstream
convention, not part of this contract: read the numbers.

## How it is built

`cmd/gen-allowances` fetches the pricing page's Next.js RSC payload, finds the
`models` array, and keeps `id`, `name`, `category`, and `planAllowanceUsd`. That
payload is an internal detail of the page, so parsing lives here and not in the
consumers. The command refuses to overwrite the catalog on a suspicious result
(fewer than 50 models, duplicate ids, missing fields).

```sh
go run ./cmd/gen-allowances
go test ./...
```

The generator is idempotent: when the upstream catalog has not changed it leaves
the file and `captured` alone. The [refresh workflow](.github/workflows/refresh.yml)
runs it weekly and commits only when the file actually changes.

## Consumers

- [`cmdcode2api`](https://github.com/peach0x33a/cmdcode2api) vendors a pinned
  copy of `catalog.json` for its boosted-allowance tracker, so its builds stay
  offline.
