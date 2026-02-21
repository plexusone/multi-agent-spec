# mas CLI

Command-line interface for multi-agent-spec.

## Installation

```bash
go install github.com/agentplexus/multi-agent-spec/cmd/mas@latest
```

## Commands

### render

Render TeamReport JSON to terminal or markdown format.

```bash
mas render <file> [flags]
```

**Flags:**

| Flag | Default | Description |
|------|---------|-------------|
| `--format`, `-f` | `box` | Output format: `box` or `narrative` |
| `--output`, `-o` | stdout | Output file path |

**Examples:**

```bash
# Render to terminal (box format)
mas render report.json

# Render to markdown (narrative format)
mas render report.json --format=narrative

# Save to file
mas render report.json --format=narrative -o report.md
```

### version

Print version information.

```bash
mas version
```

## Output Formats

### Box Format

Terminal-friendly format with Unicode box drawing:

```
╔══════════════════════════════════════════════════════════════════════════════╗
║                         RELEASE VALIDATION REPORT                            ║
╠══════════════════════════════════════════════════════════════════════════════╣
║  Project: my-app                                                             ║
║  Version: v1.2.0                                                             ║
║  Phase:   PHASE 1: REVIEW                                                    ║
╠══════════════════════════════════════════════════════════════════════════════╣
║  🔴 security — NO-GO                                                         ║
║     ├── 🟢 GO   hardcoded-secrets                                            ║
║     └── 🔴 NO-GO sql-injection [critical]                                    ║
╠──────────────────────────────────────────────────────────────────────────────╣
║  🟢 qa — GO                                                                  ║
║     └── 🟢 GO   test-coverage                                                ║
╠══════════════════════════════════════════════════════════════════════════════╣
║                        🛑 TEAM: NO-GO for v1.2.0 🛑                          ║
╚══════════════════════════════════════════════════════════════════════════════╝
```

### Narrative Format

Markdown format for documentation:

```markdown
# Release Validation Report

**Project:** my-app
**Version:** v1.2.0
**Phase:** PHASE 1: REVIEW

## Security Audit

**Status:** 🔴 NO-GO
**Verdict:** BLOCKED_SECURITY_ISSUES

| Task | Status | Severity | Detail |
|------|--------|----------|--------|
| hardcoded-secrets | 🟢 GO | - | No hardcoded secrets found |
| sql-injection | 🔴 NO-GO | critical | SQL injection in UserRepository |

## QA Validation

**Status:** 🟢 GO

| Task | Status | Detail |
|------|--------|--------|
| test-coverage | 🟢 GO | Coverage: 87% |

---

🛑 **TEAM: NO-GO for v1.2.0** 🛑
```

## Shell Completion

Generate shell completion scripts:

```bash
# Bash
mas completion bash > /etc/bash_completion.d/mas

# Zsh
mas completion zsh > "${fpath[1]}/_mas"

# Fish
mas completion fish > ~/.config/fish/completions/mas.fish
```

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Error (invalid input, file not found, etc.) |

## See Also

- [Report Schema](../schemas/report.md) — TeamReport format reference
- [Go SDK](../sdk/go.md) — Programmatic rendering
