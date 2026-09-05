## 1.1.0 (Unreleased)

REFACTOR:
- Fix permanent drift on every plan (groups/templates null round-trip, id -> known-after-apply)
- Read resolves hosts by exact-name search (API has no single-host GET); numeric ID stored in state for Update/Delete
- Shared plan-to-request converter: Update no longer drops icon_id and other optional fields
- HTTP client: 30s timeout (was none - hangs forever), context propagation on every call
- Removed the forced 1s sleep per host operation; writes stay serialized via a mutex, reads are parallel (parallel applies now work)
- Delete tolerates 404 (out-of-band deletions no longer wedge state)
- Reload failure after create/update no longer orphans the resource in state
- Secrets redacted from debug logs (macro values, API key); correct log levels (no more ERROR for INFO)
- URL query params properly escaped and search payloads JSON-marshalled (values with spaces/quotes no longer break)
- Unit tests: client layer covered (auth header, timeouts, ctx cancel, lifecycle, 409, concurrency, escaping)

DOCS:
- Fixed authenticating guide (correct POST /login payload shape, verified against Centreon 24.10)

SECOND PASS (independent review):
- terraform import support for centreon_host (by host name)
- Read: optional int fields keep null-vs-0 round-trip semantics documented (unconditional writes would reintroduce drift against unset config)
- Read: inherit value 2 now maps to 0 instead of 1 (matches schema 0/1 contract)
- Macro read failures surface as warning diagnostics instead of silently nil-ing state
- Macro values never logged (names/count only); debug logs were leaking password macros
- Computed id uses UseStateForUnknown (no known-after-apply churn on updates)
- Response bodies read via LimitReader (16 MiB cap)
- Data sources guard against unconfigured client (no more nil-pointer panic)
- HostResponse now decodes meta (pagination totals)
- GNUmakefile: duplicate targets removed, stale install path dropped
- CI: fake "Acceptance Tests" job renamed; TF_ACC no longer set with zero acceptance tests
## 0.1.0 (Unreleased)

FEATURES:
