# websnap documentation

This folder organizes the technical proposal so the project can be explained clearly both during implementation and in an interview setting.

---

## Documentation status

| Field | Status |
| --- | --- |
| Current phase | Bootstrap / `v0.1.0` in progress |
| Published release | None |
| Immediate target | `v0.1.0` |
| Language | English-first |
| i18n strategy | Planned for the CLI after the core is stable |

---

## Document map

- [`ARCHITECTURE.md`](ARCHITECTURE.md) — technical design, core decisions, and proposed Go structure
- [`FEATURES.md`](FEATURES.md) — versioned roadmap, release scope, and backlog

---

## Recommended reading order

1. `../README.md`
2. `ARCHITECTURE.md`
3. `FEATURES.md`

---

## What these docs should answer

- what problem `websnap` solves
- why Go was chosen
- how a terminal tool can capture a web page
- what has already been implemented in the bootstrap
- what belongs in V1 and what is intentionally deferred
- how the project can grow into GIF and localization without polluting the first implementation
