## 1.1.0 (Unreleased)

REFACTOR:
- Fix permanent drift on every plan (groups/templates null round-trip, id -> known-after-apply)
- ID-based host Read/Update/Delete (name fallback for legacy states)
- Shared plan-to-request converter: Update no longer drops icon_id and other optional fields
- HTTP client: 60s timeout (was none - hangs forever), context propagation on every call
- Removed global mutex + forced 1s sleep per host operation (parallel applies now work)
- Delete tolerates 404 (out-of-band deletions no longer wedge state)
- Reload failure after create/update no longer orphans the resource in state
- Secrets redacted from debug logs (macro values, API key); correct log levels (no more ERROR for INFO)
- URL query params properly escaped (search values with spaces/quotes no longer break)
- Unit tests: client layer covered (auth header, timeouts, ctx cancel, lifecycle, 409, concurrency, escaping)

DOCS:
- Fixed authenticating guide (correct POST /login payload shape, verified against Centreon 24.10)
## 0.1.0 (Unreleased)

FEATURES:
