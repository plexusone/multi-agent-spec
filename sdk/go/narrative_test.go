package multiagentspec

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestStatusText(t *testing.T) {
	tests := []struct {
		status   Status
		expected string
	}{
		{StatusGo, "PASS"},
		{StatusWarn, "WARNING"},
		{StatusNoGo, "FAIL"},
		{StatusSkip, "SKIP"},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			got := statusText(tt.status)
			if got != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}

func TestRenderNarrative(t *testing.T) {
	report := &TeamReport{
		Title:       "Security Analysis Report",
		Project:     "test-project",
		Version:     "1.0.0",
		Target:      "v1.0.0",
		Phase:       "SECURITY REVIEW",
		Summary:     "This report summarizes the security analysis findings.",
		Conclusion:  "The project is ready for release with minor fixes.",
		GeneratedAt: time.Date(2026, 2, 12, 0, 0, 0, 0, time.UTC),
		Teams: []TeamSection{
			{
				ID:     "security",
				Name:   "Security Analysis",
				Status: StatusWarn,
				Tasks: []TaskResult{
					{ID: "vuln-scan", Status: StatusWarn, Detail: "2 findings"},
				},
				Narrative: &NarrativeSection{
					Problem:        "The application uses outdated dependencies with known vulnerabilities.",
					Analysis:       "CVE-2024-1234 affects the authentication module. CVE-2024-5678 is in a logging library.",
					Recommendation: "Upgrade dependencies to latest versions before release.",
				},
				ContentBlocks: []ContentBlock{
					NewListBlock("Vulnerabilities",
						ListItem{Text: "CVE-2024-1234 - Authentication bypass (HIGH)"},
						ListItem{Text: "CVE-2024-5678 - Log injection (MEDIUM)"},
					),
				},
			},
		},
		FooterBlocks: []ContentBlock{
			NewKVPairsBlock("Required Actions",
				KVPair{Key: "1", Value: "Upgrade auth library to v2.0.0"},
				KVPair{Key: "2", Value: "Upgrade logging library to v1.5.0"},
			),
		},
		Status: StatusWarn,
	}

	var buf bytes.Buffer
	renderer := NewNarrativeRenderer(&buf)
	if err := renderer.Render(report); err != nil {
		t.Fatalf("render failed: %v", err)
	}

	output := buf.String()

	// Check YAML frontmatter
	if !strings.Contains(output, "title: \"Security Analysis Report\"") {
		t.Error("expected title in frontmatter")
	}
	if !strings.Contains(output, "date: \"2026-02-12\"") {
		t.Error("expected date in frontmatter")
	}

	// Check no emojis
	if strings.Contains(output, "🟢") || strings.Contains(output, "🟡") || strings.Contains(output, "🔴") {
		t.Error("narrative should not contain emojis")
	}

	// Check status text
	if !strings.Contains(output, "**Overall Status**: WARNING") {
		t.Error("expected overall status as text")
	}

	// Check summary
	if !strings.Contains(output, "## Executive Summary") {
		t.Error("expected executive summary section")
	}
	if !strings.Contains(output, "This report summarizes") {
		t.Error("expected summary content")
	}

	// Check narrative sections
	if !strings.Contains(output, "#### Problem") {
		t.Error("expected problem section")
	}
	if !strings.Contains(output, "outdated dependencies") {
		t.Error("expected problem content")
	}
	if !strings.Contains(output, "#### Analysis") {
		t.Error("expected analysis section")
	}
	if !strings.Contains(output, "#### Recommendation") {
		t.Error("expected recommendation section")
	}

	// Check tasks table (now includes Severity column)
	if !strings.Contains(output, "| vuln-scan | WARNING |  | 2 findings |") {
		t.Error("expected task in table")
	}

	// Check conclusion
	if !strings.Contains(output, "## Conclusion") {
		t.Error("expected conclusion section")
	}
	if !strings.Contains(output, "ready for release") {
		t.Error("expected conclusion content")
	}

	// Check footer blocks rendered
	if !strings.Contains(output, "## Action Items") {
		t.Error("expected action items section")
	}
}

func TestRenderBlockMD(t *testing.T) {
	t.Run("kv_pairs", func(t *testing.T) {
		block := NewKVPairsBlock("Metadata",
			KVPair{Key: "Author", Value: "John"},
			KVPair{Key: "Version", Value: "1.0"},
		)
		output := renderBlockMD(block)

		if !strings.Contains(output, "**Metadata**") {
			t.Error("expected title")
		}
		if !strings.Contains(output, "- **Author**: John") {
			t.Error("expected kv pair")
		}
	})

	t.Run("list", func(t *testing.T) {
		block := NewListBlock("Items",
			ListItem{Text: "First item"},
			ListItem{Text: "Second item"},
		)
		output := renderBlockMD(block)

		if !strings.Contains(output, "- First item") {
			t.Error("expected list item")
		}
	})

	t.Run("table", func(t *testing.T) {
		block := NewTableBlock("Comparison",
			[]string{"Name", "Value"},
			[][]string{
				{"foo", "bar"},
				{"baz", "qux"},
			},
		)
		output := renderBlockMD(block)

		if !strings.Contains(output, "| Name | Value |") {
			t.Error("expected table header")
		}
		if !strings.Contains(output, "| --- | --- |") {
			t.Error("expected table separator")
		}
		if !strings.Contains(output, "| foo | bar |") {
			t.Error("expected table row")
		}
	})

	t.Run("metric with target", func(t *testing.T) {
		block := NewMetricBlock("Coverage", "85%", StatusGo, "80%")
		output := renderBlockMD(block)

		if !strings.Contains(output, "**Coverage**: 85%") {
			t.Error("expected metric label and value")
		}
		if !strings.Contains(output, "(target: 80%)") {
			t.Error("expected target")
		}
		if !strings.Contains(output, "PASS") {
			t.Error("expected status text")
		}
	})

	t.Run("text", func(t *testing.T) {
		block := NewTextBlock("Description", "This is a paragraph of text.")
		output := renderBlockMD(block)

		if !strings.Contains(output, "**Description**") {
			t.Error("expected title")
		}
		if !strings.Contains(output, "This is a paragraph of text.") {
			t.Error("expected text content")
		}
	})
}

func TestNarrativeWithoutOptionalFields(t *testing.T) {
	// Minimal report without summary, conclusion, or narrative sections
	report := &TeamReport{
		Project:     "minimal-project",
		Version:     "1.0.0",
		Phase:       "TEST",
		GeneratedAt: time.Now(),
		Teams: []TeamSection{
			{
				ID:     "test",
				Name:   "Test Team",
				Status: StatusGo,
				Tasks: []TaskResult{
					{ID: "task1", Status: StatusGo, Detail: "Passed"},
				},
			},
		},
		Status: StatusGo,
	}

	var buf bytes.Buffer
	renderer := NewNarrativeRenderer(&buf)
	if err := renderer.Render(report); err != nil {
		t.Fatalf("render failed: %v", err)
	}

	output := buf.String()

	// Should not contain optional sections
	if strings.Contains(output, "## Executive Summary") {
		t.Error("should not have summary section when empty")
	}
	if strings.Contains(output, "## Conclusion") {
		t.Error("should not have conclusion section when empty")
	}
	if strings.Contains(output, "#### Problem") {
		t.Error("should not have problem section when no narrative")
	}

	// Should still have core content
	if !strings.Contains(output, "## Team Results") {
		t.Error("expected team results section")
	}
	if !strings.Contains(output, "### Test Team") {
		t.Error("expected team name")
	}
}

func TestQuickNarrativeRendererParity(t *testing.T) {
	// Test that QuickNarrativeRenderer produces similar output to NarrativeRenderer
	report := &TeamReport{
		Title:       "Security Analysis Report",
		Project:     "test-project",
		Version:     "1.0.0",
		Phase:       "SECURITY REVIEW",
		Summary:     "This report summarizes the security analysis findings.",
		GeneratedAt: time.Date(2026, 2, 12, 0, 0, 0, 0, time.UTC),
		Tags: map[string]string{
			"customer": "acme",
		},
		Teams: []TeamSection{
			{
				ID:      "security",
				Name:    "Security Analysis",
				Status:  StatusWarn,
				Verdict: "NEEDS_ATTENTION",
				Tasks: []TaskResult{
					{ID: "vuln-scan", Status: StatusWarn, Severity: "high", Detail: "2 findings"},
				},
			},
		},
		Status: StatusWarn,
	}

	var stdBuf, quickBuf bytes.Buffer

	// Render with standard renderer
	stdRenderer := NewNarrativeRenderer(&stdBuf)
	if err := stdRenderer.Render(report); err != nil {
		t.Fatalf("standard render failed: %v", err)
	}

	// Render with quick renderer
	quickRenderer := NewQuickNarrativeRenderer(&quickBuf)
	if err := quickRenderer.Render(report); err != nil {
		t.Fatalf("quick render failed: %v", err)
	}

	stdOutput := stdBuf.String()
	quickOutput := quickBuf.String()

	// Both should contain key elements
	for _, expected := range []string{
		"Security Analysis Report",
		"Security Analysis",
		"NEEDS_ATTENTION",
		"vuln-scan",
		"WARNING",
		"high",
		"customer",
		"acme",
	} {
		if !strings.Contains(stdOutput, expected) {
			t.Errorf("standard output missing %q", expected)
		}
		if !strings.Contains(quickOutput, expected) {
			t.Errorf("quick output missing %q", expected)
		}
	}
}
