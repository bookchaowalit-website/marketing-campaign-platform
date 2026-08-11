# BookMarketing — Dirty Tree Triage

**Snapshot:** 2026-08-11  
**Working tree:** ~39 paths (do not `git add -A` without cluster review)

## Counts by top-level

| Path | Approx. files | Cluster |
|---|---:|---|
| `api-go/` | 26 | A — backend remediation |
| `web/` | 10 | B — frontend |
| `packages/` | 1 | D — protos Makefile |
| `apps/` (deleted file) | 1 | E — dead tree |
| `README.md` | 1 | C — docs |
| `PRODUCT.md` (new) | 1 | C — docs |

## Suggested commit order (when owner approves)

1. **api-go + migrations 016 + protos Makefile + apps removal**  
   Message theme: “Make Go backend compile and migrations apply; remove dead apps layout.”
2. **web ESLint, api-client, SEO dashboard, package bumps**  
   Message theme: “Frontend lint/build and SEO response contract.”
3. **PRODUCT.md + README polish**  
   Message theme: “Interview case study and claim hygiene.”

## Explicit non-goals for this triage

- No `git reset` / clean of dirty work  
- No GitHub push  
- No inventing multi-tenant auth without product decision  
- No Next 16 major upgrade in this pass  
