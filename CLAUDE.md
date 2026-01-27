# MSLS Project Instructions for Claude Code

## Project Overview
MSLS (Multi-School Learning System) - School ERP with SaaS + On-Premise deployment.

## Agent Autonomy Rules

### Autonomous Operation Mode: ENABLED

Agents operate autonomously WITHOUT user approval for:

| Action | Autonomous | Notes |
|--------|-----------|-------|
| Code within documented patterns | YES | Follow architecture.md |
| Create/modify files per story | YES | Match project structure |
| Run tests | YES | Always run after changes |
| Fix lint/type errors | YES | Auto-fix immediately |
| Implement API endpoints | YES | Follow REST patterns |
| Create components | YES | Follow Tailwind + custom patterns |
| Write unit tests | YES | Required for all code |
| Database migrations | YES | Follow naming conventions |
| Bug fixes | YES | Document in commit |
| Refactor within module | YES | Keep boundaries |

### Requires User Approval:

| Action | Approval Needed |
|--------|----------------|
| New architectural patterns | YES |
| New external dependencies | YES |
| Breaking API changes | YES |
| Schema changes affecting multiple modules | YES |
| Security-related changes | YES |
| Deployment configuration | YES |

## Agent Communication Protocol

Agents communicate through **shared artifacts**, not direct messages:

```
ARTIFACT HANDOFF FLOW:
──────────────────────────────────────────────────────────────
1. PM creates story → writes to: stories/{epic}/{story}.md
2. Dev reads story → implements → updates story status
3. TEA reads story → writes tests → updates story status
4. Dev reads test results → fixes issues → marks complete
──────────────────────────────────────────────────────────────
```

### Key Artifacts (Read Before Acting):

| Artifact | Location | Purpose |
|----------|----------|---------|
| Architecture | `_bmad-output/planning-artifacts/architecture.md` | ALL technical decisions |
| Project Context | `_bmad-output/planning-artifacts/project-context.md` | Agent implementation rules |
| Sprint Status | `_bmad-output/implementation-artifacts/sprint-status.yaml` | Current work state |
| Stories | `_bmad-output/implementation-artifacts/stories/` | Implementation specs |

### Workflow Sequence (No User Prompts Between Steps):

```
create-epics-and-stories → sprint-planning → dev-story → code-review → (repeat)
```

## Technology Stack (Do Not Deviate)

- **Backend**: Go 1.23+, Gin, GORM + sqlc, PostgreSQL 16
- **Frontend**: Angular 21, Tailwind CSS, Custom Components, Signals
- **Auth**: JWT (RS256), Argon2id, TOTP 2FA
- **Multi-tenancy**: PostgreSQL RLS with tenant_id

## Critical Rules

1. **ALWAYS read architecture.md before implementing**
2. **ALWAYS read the story file before starting work**
3. **NEVER ask user for approval on documented patterns**
4. **ALWAYS update story status after completing tasks**
5. **ALWAYS run tests before marking complete**
6. **NEVER deviate from documented project structure**

## Directory Structure

```
msls/
├── msls-backend/          # Go backend (create when implementing)
├── msls-frontend/         # Angular frontend (create when implementing)
├── _bmad-output/
│   ├── planning-artifacts/
│   │   ├── architecture.md
│   │   ├── project-context.md
│   │   └── school-erp-prd/
│   └── implementation-artifacts/
│       ├── sprint-status.yaml
│       ├── epics/
│       └── stories/
└── CLAUDE.md              # This file
```

## When Stuck

1. Re-read architecture.md for guidance
2. Check if similar pattern exists in codebase
3. Follow the documented pattern exactly
4. If truly blocked (missing info), then ask user
