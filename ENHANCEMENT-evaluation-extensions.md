# Enhancement Request: Domain-Specific Evaluation Extensions

## Context

Using multi-agent-spec `TeamReport` for a batch evaluation pipeline that analyzes 80 custom extension JARs across 5 dimensions (TRMS migration, JDK 17 compatibility, PM use case, security audit, performance). Each dimension maps to a `TeamSection`, each finding maps to a `TaskResult`.

The base schema works well — we're 100% compatible with zero custom fields. But three additions would make evaluation/audit use cases first-class without requiring downstream schema extensions.

## Current Workarounds

| Need | Current Workaround | Problem |
|---|---|---|
| Finding severity (critical/high/medium/low) | `metadata.severity` | Untyped, not renderable, not aggregatable |
| Domain verdict (e.g., "BLOCKED_PENDING_ENHANCEMENT") | Stuffed into `detail` string or table cell | Lost in rendering, can't filter/aggregate |
| Report tags for aggregation | Not possible | Can't group/filter across reports |

## Proposed Changes

### 1. `severity` on TaskResult

```json
{
  "id": "sql-injection",
  "status": "NO-GO",
  "severity": "critical",
  "detail": "String concatenation in WHERE clause"
}
```

**Schema addition** to `TaskResult`:
```json
"severity": {
  "type": "string",
  "enum": ["critical", "high", "medium", "low", "info"],
  "description": "Finding severity/impact level. Orthogonal to status (status = pass/fail, severity = impact)."
}
```

**Go SDK addition** to `TaskResult`:
```go
Severity string `json:"severity,omitempty"`
```

**Rendering**: Box format could show `🔴 NO-GO [critical]`, narrative could add a Severity column to the task table.

**Rationale**: Status answers "did it pass?" Severity answers "how bad is it?" Every audit/evaluation use case needs both. Currently forced into untyped `metadata` map which renderers can't access.

### 2. `verdict` on TeamSection

```json
{
  "id": "trms-migration",
  "name": "TRMS Compatibility",
  "status": "NO-GO",
  "verdict": "BLOCKED_PENDING_ENHANCEMENT",
  ...
}
```

**Schema addition** to `TeamSection`:
```json
"verdict": {
  "type": "string",
  "description": "Domain-specific verdict label. Richer than status (4 values). Status is machine-readable GO/NO-GO; verdict is the human-readable domain assessment."
}
```

**Go SDK addition** to `TeamSection`:
```go
Verdict string `json:"verdict,omitempty"`
```

**Rendering**: Box format shows verdict after status: `🔴 NO-GO — BLOCKED_PENDING_ENHANCEMENT`. Narrative shows as bold text under status.

**Rationale**: The 4-value Status enum is great for machine logic (is it GO or not?) but evaluation domains need richer labels. Security audits have PASSED/FAILED/CONDITIONAL. Migration checks have COMPATIBLE/NEEDS_WORK/BLOCKED. Forcing these into `detail` strings loses them for aggregation.

### 3. `tags` on TeamReport

```json
{
  "project": "aap/AS400.jar",
  "tags": {
    "customer": "aap",
    "environment": "aap-prod-aks",
    "use_case": "provisioning",
    "target_system": "as400",
    "source_type": "full"
  },
  ...
}
```

**Schema addition** to `TeamReport`:
```json
"tags": {
  "type": "object",
  "additionalProperties": { "type": "string" },
  "description": "Key-value tags for filtering and aggregation across reports."
}
```

**Go SDK addition** to `TeamReport`:
```go
Tags map[string]string `json:"tags,omitempty"`
```

**Rendering**: Box format shows tags as a kv_pairs block after the header. Narrative shows as metadata in the YAML frontmatter.

**Rationale**: When you have 80+ reports, you need to aggregate: "how many JARs are TRMS-blocked?", "which customers have security criticals?", "group by use_case". Tags enable this without parsing `summary_blocks` or `project` strings.

## Impact Assessment

| Change | Schema | Go SDK | Renderer | Breaking? |
|---|---|---|---|---|
| `severity` on TaskResult | +5 lines | +1 field | +icon in taskLine | No (optional field) |
| `verdict` on TeamSection | +4 lines | +1 field | +text in teamHeader | No (optional field) |
| `tags` on TeamReport | +4 lines | +1 field | +block in header | No (optional field) |

All three are optional fields — zero breaking changes. Existing reports validate unchanged.

## Example: Before and After

### Before (current workaround)
```json
{
  "id": "sql-injection",
  "status": "NO-GO",
  "detail": "SQL injection — string concatenation in WHERE clause",
  "metadata": { "severity": "critical", "count": 2 }
}
```

### After (with severity)
```json
{
  "id": "sql-injection",
  "status": "NO-GO",
  "severity": "critical",
  "detail": "SQL injection — string concatenation in WHERE clause",
  "metadata": { "count": 2 }
}
```

Severity moves from untyped metadata to a typed, renderable, aggregatable field.

## Use Cases Beyond JAR Evaluation

- **Security audits**: severity is standard (CVSS maps to critical/high/medium/low)
- **Compliance checks**: verdict = "COMPLIANT" / "NON_COMPLIANT" / "PARTIAL"
- **Code reviews**: severity on findings, verdict on review sections
- **Infrastructure audits**: tags for region, account, service
- **PRD evaluations**: verdict = "APPROVE" / "REVISE" / "REJECT" (already in llm-evaluation.schema.json as a concept)
