package multiagentspec

import (
	"fmt"
	"io"
	"strings"
	"text/template"
)

// NarrativeSection holds prose content for narrative reports.
type NarrativeSection struct {
	// Problem describes the issue or context being addressed.
	Problem string `json:"problem,omitempty"`

	// Analysis contains the detailed findings.
	Analysis string `json:"analysis,omitempty"`

	// Recommendation describes the suggested action.
	Recommendation string `json:"recommendation,omitempty"`
}

// NarrativeRenderer renders TeamReport to Pandoc-friendly Markdown.
// Output is designed for conversion to PDF via:
//
//	pandoc report.md -o report.pdf --pdf-engine=xelatex
type NarrativeRenderer struct {
	w io.Writer
}

// NewNarrativeRenderer creates a new NarrativeRenderer writing to w.
func NewNarrativeRenderer(w io.Writer) *NarrativeRenderer {
	return &NarrativeRenderer{w: w}
}

// Render renders the report as Pandoc-friendly Markdown.
// No emojis are used - status is rendered as text (PASS, FAIL, WARNING, SKIP).
func (r *NarrativeRenderer) Render(report *TeamReport) error {
	report.SortByDAG()

	tmpl, err := template.New("narrative").Funcs(narrativeFuncs()).Parse(NarrativeTemplate)
	if err != nil {
		return fmt.Errorf("parsing template: %w", err)
	}
	return tmpl.Execute(r.w, report)
}

// narrativeFuncs returns the template function map for narrative rendering.
func narrativeFuncs() template.FuncMap {
	return template.FuncMap{
		"statusText":       statusText,
		"hasNarrative":     hasNarrative,
		"hasSummary":       hasSummary,
		"hasConclusion":    hasConclusion,
		"hasContentBlocks": hasContentBlocks,
		"renderBlockMD":    renderBlockMD,
		"renderBlocksMD":   renderBlocksMD,
		"indent":           indent,
	}
}

// statusText returns a text representation of status (no emojis).
func statusText(s Status) string {
	switch s {
	case StatusGo:
		return "PASS"
	case StatusWarn:
		return "WARNING"
	case StatusNoGo:
		return "FAIL"
	case StatusSkip:
		return "SKIP"
	default:
		return string(s)
	}
}

// hasNarrative returns true if the team has narrative content.
func hasNarrative(team TeamSection) bool {
	return team.Narrative != nil && (team.Narrative.Problem != "" ||
		team.Narrative.Analysis != "" ||
		team.Narrative.Recommendation != "")
}

// hasSummary returns true if the report has a summary.
func hasSummary(report *TeamReport) bool {
	return report.Summary != ""
}

// hasConclusion returns true if the report has a conclusion.
func hasConclusion(report *TeamReport) bool {
	return report.Conclusion != ""
}

// renderBlocksMD renders multiple content blocks as Markdown.
func renderBlocksMD(blocks []ContentBlock) string {
	var parts []string
	for _, block := range blocks {
		parts = append(parts, renderBlockMD(block))
	}
	return strings.Join(parts, "\n\n")
}

// renderBlockMD renders a single content block as Markdown.
func renderBlockMD(block ContentBlock) string {
	var sb strings.Builder

	if block.Title != "" {
		sb.WriteString("**")
		sb.WriteString(block.Title)
		sb.WriteString("**\n\n")
	}

	switch block.Type {
	case ContentBlockKVPairs:
		for _, pair := range block.Pairs {
			sb.WriteString("- **")
			sb.WriteString(pair.Key)
			sb.WriteString("**: ")
			sb.WriteString(pair.Value)
			sb.WriteString("\n")
		}
	case ContentBlockList:
		for _, item := range block.Items {
			sb.WriteString("- ")
			sb.WriteString(item.Text)
			sb.WriteString("\n")
		}
	case ContentBlockText:
		sb.WriteString(block.Content)
		sb.WriteString("\n")
	case ContentBlockTable:
		// Header row
		sb.WriteString("| ")
		sb.WriteString(strings.Join(block.Headers, " | "))
		sb.WriteString(" |\n")
		// Separator
		sep := make([]string, len(block.Headers))
		for i := range sep {
			sep[i] = "---"
		}
		sb.WriteString("| ")
		sb.WriteString(strings.Join(sep, " | "))
		sb.WriteString(" |\n")
		// Rows
		for _, row := range block.Rows {
			sb.WriteString("| ")
			sb.WriteString(strings.Join(row, " | "))
			sb.WriteString(" |\n")
		}
	case ContentBlockMetric:
		sb.WriteString("- **")
		sb.WriteString(block.Label)
		sb.WriteString("**: ")
		sb.WriteString(block.Value)
		if block.Target != "" {
			sb.WriteString(" (target: ")
			sb.WriteString(block.Target)
			sb.WriteString(")")
		}
		sb.WriteString(" — ")
		sb.WriteString(statusText(block.Status))
		sb.WriteString("\n")
	}

	return sb.String()
}

// indent adds prefix to each line of text.
func indent(prefix, text string) string {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		if line != "" {
			lines[i] = prefix + line
		}
	}
	return strings.Join(lines, "\n")
}

// NarrativeTemplate is the Pandoc-friendly Markdown template.
// Designed for PDF generation via: pandoc report.md -o report.pdf --pdf-engine=xelatex
const NarrativeTemplate = `---
title: "{{ .EffectiveTitle }}"
date: "{{ .GeneratedAt.Format "2006-01-02" }}"
---

# {{ .EffectiveTitle }}

**Project**: {{ .Project }}
**Version**: {{ .Version }}
**Phase**: {{ .Phase }}
**Overall Status**: {{ statusText .Status }}
{{- if hasSummary . }}

## Executive Summary

{{ .Summary }}
{{- end }}
{{- if .SummaryBlocks }}

## Overview

{{ renderBlocksMD .SummaryBlocks }}
{{- end }}

## Team Results
{{- range .Teams }}

### {{ .Name }}

**Status**: {{ statusText .Status }}
{{- if hasNarrative . }}
{{- if .Narrative.Problem }}

#### Problem

{{ .Narrative.Problem }}
{{- end }}
{{- if .Narrative.Analysis }}

#### Analysis

{{ .Narrative.Analysis }}
{{- end }}
{{- if .Narrative.Recommendation }}

#### Recommendation

{{ .Narrative.Recommendation }}
{{- end }}
{{- end }}
{{- if .Tasks }}

#### Tasks

| Task | Status | Detail |
| --- | --- | --- |
{{- range .Tasks }}
| {{ .ID }} | {{ statusText .Status }} | {{ .Detail }} |
{{- end }}
{{- end }}
{{- if hasContentBlocks . }}

#### Details

{{ renderBlocksMD .ContentBlocks }}
{{- end }}
{{- end }}
{{- if .FooterBlocks }}

## Action Items

{{ renderBlocksMD .FooterBlocks }}
{{- end }}
{{- if hasConclusion . }}

## Conclusion

{{ .Conclusion }}
{{- end }}
`
