# 18 — Git และ release workflow

ข้อมูลจาก `git log` บน `origin/main` (20 commits ตั้งแต่ `2c50c4e` 2026-04-12) และ `gh pr list` (PR เดียว #1 merged 2026-07-26)
ไม่มี tag, ไม่มี release, ไม่มี branch protection (`gh api .../branches/main/protection` → 404) และไม่มี deployment: "release" ของ repo นี้คือ commit บน `main` ที่ CI เขียว

```mermaid
gitGraph
  commit id: "2c50c4e init"
  branch agent-bookmarketing-design-doc
  commit id: "ed3e465 design system"
  checkout main
  merge agent-bookmarketing-design-doc id: "PR 1"
  commit id: "25ca091 runbook stub"
  commit id: "b934a69 api-go compile"
  commit id: "ddd9fb8 web SEO contract"
  commit id: "5d754ca PRODUCT.md"
  commit id: "c197b0c license"
  commit id: "8937ac2 Go HTTP tests"
  commit id: "1ee7d36 red tests"
  commit id: "c7e232c Next 16"
  commit id: "76f10ed deps"
  commit id: "1095ae3 CI, red"
  commit id: "c796b31 CI green"
  commit id: "517ac6e next 16.4.2"
  commit id: "e182e01 sentry scrub"
  commit id: "92e0cf8 no demo creds"
  commit id: "8e1d0f7 deploy smoke"
  commit id: "60379ea smoke env"
  commit id: "a4061e5 KPI review"
```

หมายเหตุ: `git log` ปัจจุบันเป็นเส้นตรง (ไม่มี merge commit) จึงแสดง PR #1 ตาม `gh pr list` ไม่ใช่ตาม topology ของ git

```mermaid
flowchart TD
  A["Local clone or worktree"] --> B["Commit on main, or docs branch for packs"]
  B --> C["git push origin main"]
  C --> D["ci.yml on push main: web, api-python, api-go, security"]
  D -->|"red"| E["Fix forward with a new commit"]
  E --> C
  D -->|"green"| F["Release = commit on main"]
  F -.->|"not connected"| G["Vercel deployment_status, deploy-smoke.yml"]
  PR["pull_request event"] -.->|"also runs ci.yml, no PR in the last 90 days"| D
```

Sources: `git log --format="%h %ad %s"`, `gh pr list --state all`, `gh run list` (8 runs, all `push` on `main`), `gh api repos/bookchaowalit-website/marketing-campaign-platform/branches/main/protection`, `.github/workflows/ci.yml`, `.github/workflows/deploy-smoke.yml`, `PRODUCT.md` §Dirty-tree triage, `DIRTY-TRIAGE.md`

แนวปฏิบัติที่เห็นจาก history: conventional commits (`feat`, `fix`, `test`, `ci`, `docs`, `chore`), TDD "red" commit ก่อน fix (`1ee7d36` → `c7e232c`, `34a136f` → `e182e01`) ตาม ADR-0004 ของ architecture-docs และ commit แยก cluster ตาม `DIRTY-TRIAGE.md`

## Open questions

- ต้องการ branch protection / required checks บน `main` หรือไม่
- เมื่อมี deployment จะ tag release หรือใช้ commit SHA จาก `/api/health` อย่างเดียว
