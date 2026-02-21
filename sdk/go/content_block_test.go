package multiagentspec

import (
	"bytes"
	"strings"
	"testing"
)

func TestContentBlockConstructors(t *testing.T) {
	t.Run("NewKVPairsBlock", func(t *testing.T) {
		block := NewKVPairsBlock("METADATA",
			KVPair{Key: "Name", Value: "test"},
			KVPair{Key: "Version", Value: "1.0.0", Icon: "📦"},
		)
		if block.Type != ContentBlockKVPairs {
			t.Errorf("expected type %s, got %s", ContentBlockKVPairs, block.Type)
		}
		if block.Title != "METADATA" {
			t.Errorf("expected title METADATA, got %s", block.Title)
		}
		if len(block.Pairs) != 2 {
			t.Errorf("expected 2 pairs, got %d", len(block.Pairs))
		}
	})

	t.Run("NewListBlock", func(t *testing.T) {
		block := NewListBlock("FINDINGS",
			ListItem{Text: "Issue 1", Icon: "🔴"},
			ListItem{Text: "Issue 2", Status: StatusWarn},
		)
		if block.Type != ContentBlockList {
			t.Errorf("expected type %s, got %s", ContentBlockList, block.Type)
		}
		if len(block.Items) != 2 {
			t.Errorf("expected 2 items, got %d", len(block.Items))
		}
	})

	t.Run("NewTextBlock", func(t *testing.T) {
		block := NewTextBlock("DESCRIPTION", "This is a test description.")
		if block.Type != ContentBlockText {
			t.Errorf("expected type %s, got %s", ContentBlockText, block.Type)
		}
		if block.Content != "This is a test description." {
			t.Errorf("unexpected content: %s", block.Content)
		}
	})

	t.Run("NewTableBlock", func(t *testing.T) {
		block := NewTableBlock("COMPARISON",
			[]string{"Feature", "v1", "v2"},
			[][]string{
				{"Auth", "Basic", "OAuth2"},
				{"Cache", "None", "Redis"},
			},
		)
		if block.Type != ContentBlockTable {
			t.Errorf("expected type %s, got %s", ContentBlockTable, block.Type)
		}
		if len(block.Headers) != 3 {
			t.Errorf("expected 3 headers, got %d", len(block.Headers))
		}
		if len(block.Rows) != 2 {
			t.Errorf("expected 2 rows, got %d", len(block.Rows))
		}
	})

	t.Run("NewMetricBlock", func(t *testing.T) {
		block := NewMetricBlock("Coverage", "85%", StatusGo, "80%")
		if block.Type != ContentBlockMetric {
			t.Errorf("expected type %s, got %s", ContentBlockMetric, block.Type)
		}
		if block.Label != "Coverage" {
			t.Errorf("expected label Coverage, got %s", block.Label)
		}
		if block.Status != StatusGo {
			t.Errorf("expected status GO, got %s", block.Status)
		}
		if block.Target != "80%" {
			t.Errorf("expected target 80%%, got %s", block.Target)
		}
	})

	t.Run("NewMetricBlock without target", func(t *testing.T) {
		block := NewMetricBlock("Score", "95", StatusGo, "")
		if block.Target != "" {
			t.Errorf("expected empty target, got %s", block.Target)
		}
	})
}

func TestListItemEffectiveIcon(t *testing.T) {
	tests := []struct {
		name     string
		item     ListItem
		expected string
	}{
		{
			name:     "explicit icon takes precedence",
			item:     ListItem{Text: "test", Icon: "⚠️", Status: StatusGo},
			expected: "⚠️",
		},
		{
			name:     "derives from status when no icon",
			item:     ListItem{Text: "test", Status: StatusWarn},
			expected: "🟡",
		},
		{
			name:     "empty when neither set",
			item:     ListItem{Text: "test"},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.item.EffectiveIcon()
			if got != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}

func TestRenderWithContentBlocks(t *testing.T) {
	report := &TeamReport{
		Title:   "TEST REPORT",
		Project: "test-project",
		Version: "1.0.0",
		Target:  "v1.0.0",
		Phase:   "ANALYSIS",
		SummaryBlocks: []ContentBlock{
			NewKVPairsBlock("",
				KVPair{Key: "Project", Value: "test-project"},
				KVPair{Key: "Version", Value: "1.0.0"},
			),
		},
		Teams: []TeamSection{
			{
				ID:     "security",
				Name:   "Security Analysis",
				Status: StatusWarn,
				Tasks: []TaskResult{
					{ID: "vuln-scan", Status: StatusWarn, Detail: "2 findings"},
				},
				ContentBlocks: []ContentBlock{
					NewListBlock("",
						ListItem{Text: "CVE-2024-1234 (HIGH)", Icon: "🔴"},
						ListItem{Text: "Outdated dependency (MEDIUM)", Icon: "🟡"},
					),
				},
			},
		},
		FooterBlocks: []ContentBlock{
			NewKVPairsBlock("ACTION ITEMS",
				KVPair{Icon: "🔴", Key: "1", Value: "Fix CVE-2024-1234"},
			),
		},
		Status: StatusWarn,
	}

	var buf bytes.Buffer
	renderer := NewRenderer(&buf)
	if err := renderer.Render(report); err != nil {
		t.Fatalf("render failed: %v", err)
	}

	output := buf.String()

	// Verify title
	if !strings.Contains(output, "TEST REPORT") {
		t.Error("expected custom title in output")
	}

	// Verify summary blocks rendered
	if !strings.Contains(output, "Project: test-project") {
		t.Error("expected summary kv_pairs in output")
	}

	// Verify team content blocks rendered
	if !strings.Contains(output, "CVE-2024-1234 (HIGH)") {
		t.Error("expected team content block in output")
	}

	// Verify footer blocks rendered
	if !strings.Contains(output, "ACTION ITEMS") {
		t.Error("expected footer block title in output")
	}
	if !strings.Contains(output, "Fix CVE-2024-1234") {
		t.Error("expected footer content in output")
	}
}

func TestRenderTable(t *testing.T) {
	lines := renderTable(
		[]string{"Name", "Status"},
		[][]string{
			{"auth", "GO"},
			{"cache", "WARN"},
		},
	)

	if len(lines) != 4 { // header + separator + 2 data rows
		t.Errorf("expected 4 lines, got %d", len(lines))
	}

	// Check header
	if !strings.Contains(lines[0], "Name") || !strings.Contains(lines[0], "Status") {
		t.Error("header row missing column names")
	}

	// Check separator contains table characters
	if !strings.Contains(lines[1], "─") {
		t.Error("separator row missing horizontal line")
	}

	// Check data rows
	if !strings.Contains(lines[2], "auth") {
		t.Error("first data row missing")
	}
}

func TestWrapText(t *testing.T) {
	content := "This is a long line that should be wrapped to fit within the box width properly"
	lines := wrapText(content, 40)

	for _, line := range lines {
		// Each line should be a paddedLine with visual width of boxWidth+2
		// Use visualLength to check (excludes border chars which are counted separately)
		// The line format is: "║ " + content + padding + "║"
		// Visual width should be boxWidth + 2 (for the two border chars)
		vLen := visualLength(line)
		if vLen > boxWidth+2 {
			t.Errorf("line visual length too long: %d chars (expected max %d)", vLen, boxWidth+2)
		}
	}

	if len(lines) < 2 {
		t.Error("expected text to wrap to multiple lines")
	}
}

func TestBackwardCompatibility(t *testing.T) {
	// Ensure reports without content blocks still render correctly
	report := &TeamReport{
		Project: "legacy-project",
		Version: "1.0.0",
		Target:  "v1.0.0",
		Phase:   "PHASE 1: VALIDATION",
		Teams: []TeamSection{
			{
				ID:   "pm-validation",
				Name: "pm",
				Tasks: []TaskResult{
					{ID: "requirements", Status: StatusGo, Detail: "All documented"},
				},
				Status: StatusGo,
			},
		},
		Status: StatusGo,
	}

	var buf bytes.Buffer
	renderer := NewRenderer(&buf)
	if err := renderer.Render(report); err != nil {
		t.Fatalf("render failed: %v", err)
	}

	output := buf.String()

	// Should use default title
	if !strings.Contains(output, "TEAM STATUS REPORT") {
		t.Error("expected default title for legacy report")
	}

	// Should show project/target in default format
	if !strings.Contains(output, "Project: legacy-project") {
		t.Error("expected project line in legacy format")
	}

	// Should render tasks
	if !strings.Contains(output, "requirements") {
		t.Error("expected task to render")
	}
}

func TestRenderMetric(t *testing.T) {
	t.Run("without target", func(t *testing.T) {
		line := renderMetric("Coverage", "85%", StatusGo, "")

		if !strings.Contains(line, "Coverage") {
			t.Error("expected label in output")
		}
		if !strings.Contains(line, "85%") {
			t.Error("expected value in output")
		}
		if !strings.Contains(line, "🟢") {
			t.Error("expected GO icon in output")
		}
	})

	t.Run("with target", func(t *testing.T) {
		line := renderMetric("Coverage", "85%", StatusGo, "80%")

		if !strings.Contains(line, "Coverage") {
			t.Error("expected label in output")
		}
		if !strings.Contains(line, "85%") {
			t.Error("expected value in output")
		}
		if !strings.Contains(line, "target: 80%") {
			t.Error("expected target in output")
		}
	})
}

func TestRenderKVPairs(t *testing.T) {
	pairs := []KVPair{
		{Key: "Name", Value: "test"},
		{Key: "Status", Value: "active", Icon: "✅"},
	}
	lines := renderKVPairs(pairs)

	if len(lines) != 2 {
		t.Errorf("expected 2 lines, got %d", len(lines))
	}

	if !strings.Contains(lines[0], "Name: test") {
		t.Error("first pair not rendered correctly")
	}
	if !strings.Contains(lines[1], "✅") || !strings.Contains(lines[1], "Status: active") {
		t.Error("second pair with icon not rendered correctly")
	}
}

func TestRenderList(t *testing.T) {
	items := []ListItem{
		{Text: "Item without icon"},
		{Text: "Item with icon", Icon: "•"},
		{Text: "Item with status", Status: StatusWarn},
	}
	lines := renderList(items)

	if len(lines) != 3 {
		t.Errorf("expected 3 lines, got %d", len(lines))
	}

	// Item without icon should have indentation
	if !strings.Contains(lines[0], "  Item without icon") {
		t.Error("item without icon should be indented")
	}

	// Item with explicit icon
	if !strings.Contains(lines[1], "• Item with icon") {
		t.Error("item with icon not rendered correctly")
	}

	// Item with status-derived icon
	if !strings.Contains(lines[2], "🟡") {
		t.Error("item should have WARN icon derived from status")
	}
}

func TestAgentResultToTeamSectionWithContentBlocks(t *testing.T) {
	result := AgentResult{
		AgentID:    "security",
		StepID:     "security-scan",
		AgentModel: "sonnet",
		Tasks: []TaskResult{
			{ID: "scan", Status: StatusWarn, Detail: "Issues found"},
		},
		ContentBlocks: []ContentBlock{
			NewListBlock("Findings",
				ListItem{Text: "CVE-2024-001", Icon: "🔴"},
			),
		},
		Status: StatusWarn,
	}

	section := result.ToTeamSection()

	if section.ID != "security-scan" {
		t.Errorf("expected ID security-scan, got %s", section.ID)
	}
	if len(section.ContentBlocks) != 1 {
		t.Errorf("expected 1 content block, got %d", len(section.ContentBlocks))
	}
	if section.ContentBlocks[0].Type != ContentBlockList {
		t.Errorf("expected list block, got %s", section.ContentBlocks[0].Type)
	}
}

func TestQuickRendererParity(t *testing.T) {
	// Test that QuickRenderer produces similar output to Renderer
	report := &TeamReport{
		Title:   "TEST REPORT",
		Project: "test-project",
		Version: "1.0.0",
		Target:  "v1.0.0",
		Phase:   "ANALYSIS",
		Tags: map[string]string{
			"customer":    "acme",
			"environment": "staging",
		},
		Teams: []TeamSection{
			{
				ID:      "security",
				Name:    "Security Analysis",
				Status:  StatusWarn,
				Verdict: "NEEDS_REVIEW",
				Tasks: []TaskResult{
					{ID: "vuln-scan", Status: StatusWarn, Severity: "high", Detail: "2 findings"},
				},
			},
		},
		Status: StatusWarn,
	}

	var stdBuf, quickBuf bytes.Buffer

	// Render with standard renderer
	stdRenderer := NewRenderer(&stdBuf)
	if err := stdRenderer.Render(report); err != nil {
		t.Fatalf("standard render failed: %v", err)
	}

	// Render with quick renderer
	quickRenderer := NewQuickRenderer(&quickBuf)
	if err := quickRenderer.Render(report); err != nil {
		t.Fatalf("quick render failed: %v", err)
	}

	stdOutput := stdBuf.String()
	quickOutput := quickBuf.String()

	// Both should contain key elements
	for _, expected := range []string{
		"TEST REPORT",
		"Security Analysis",
		"NEEDS_REVIEW",
		"vuln-scan",
		"high",
		"customer: acme",
		"environment: staging",
	} {
		if !strings.Contains(stdOutput, expected) {
			t.Errorf("standard output missing %q", expected)
		}
		if !strings.Contains(quickOutput, expected) {
			t.Errorf("quick output missing %q", expected)
		}
	}
}

// TestQuickRendererBoxStructure verifies the box format structure
func TestQuickRendererBoxStructure(t *testing.T) {
	report := &TeamReport{
		Title:   "BOX STRUCTURE TEST",
		Project: "test",
		Phase:   "TEST PHASE",
		Teams: []TeamSection{
			{
				ID:     "team1",
				Name:   "Test Team",
				Status: StatusGo,
				Tasks: []TaskResult{
					{ID: "task1", Status: StatusGo, Detail: "Done"},
				},
			},
		},
		Status: StatusGo,
	}

	var buf bytes.Buffer
	renderer := NewQuickRenderer(&buf)
	if err := renderer.Render(report); err != nil {
		t.Fatalf("render failed: %v", err)
	}

	output := buf.String()
	lines := strings.Split(output, "\n")

	// Filter out empty lines
	var nonEmptyLines []string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			nonEmptyLines = append(nonEmptyLines, line)
		}
	}

	if len(nonEmptyLines) < 3 {
		t.Fatalf("expected at least 3 lines, got %d", len(nonEmptyLines))
	}

	// Check header (first non-empty line should start with ╔)
	if !strings.HasPrefix(nonEmptyLines[0], "╔") {
		t.Errorf("first line should be header starting with ╔, got: %s", nonEmptyLines[0])
	}
	if !strings.HasSuffix(nonEmptyLines[0], "╗") {
		t.Errorf("header should end with ╗, got: %s", nonEmptyLines[0])
	}

	// Check footer (last non-empty line should start with ╚)
	lastLine := nonEmptyLines[len(nonEmptyLines)-1]
	if !strings.HasPrefix(lastLine, "╚") {
		t.Errorf("last line should be footer starting with ╚, got: %s", lastLine)
	}
	if !strings.HasSuffix(lastLine, "╝") {
		t.Errorf("footer should end with ╝, got: %s", lastLine)
	}

	// Check that content lines have proper borders
	for i, line := range nonEmptyLines {
		if i == 0 || i == len(nonEmptyLines)-1 {
			continue // Skip header and footer
		}
		if strings.HasPrefix(line, "╠") {
			// Separator line
			if !strings.HasSuffix(line, "╣") {
				t.Errorf("separator should end with ╣, got: %s", line)
			}
		} else if strings.HasPrefix(line, "║") {
			// Content line
			if !strings.HasSuffix(line, "║") {
				t.Errorf("content line should end with ║, got: %s", line)
			}
		}
	}
}

// TestQuickRendererTags verifies tags are rendered in sorted order
func TestQuickRendererTags(t *testing.T) {
	report := &TeamReport{
		Project: "test",
		Phase:   "TEST",
		Tags: map[string]string{
			"zebra":  "last",
			"alpha":  "first",
			"middle": "center",
		},
		Teams:  []TeamSection{},
		Status: StatusGo,
	}

	var buf bytes.Buffer
	renderer := NewQuickRenderer(&buf)
	if err := renderer.Render(report); err != nil {
		t.Fatalf("render failed: %v", err)
	}

	output := buf.String()

	// Check all tags are present
	if !strings.Contains(output, "alpha: first") {
		t.Error("missing alpha tag")
	}
	if !strings.Contains(output, "middle: center") {
		t.Error("missing middle tag")
	}
	if !strings.Contains(output, "zebra: last") {
		t.Error("missing zebra tag")
	}

	// Verify sort order: alpha should appear before middle, middle before zebra
	alphaIdx := strings.Index(output, "alpha: first")
	middleIdx := strings.Index(output, "middle: center")
	zebraIdx := strings.Index(output, "zebra: last")

	if alphaIdx > middleIdx {
		t.Error("alpha should appear before middle")
	}
	if middleIdx > zebraIdx {
		t.Error("middle should appear before zebra")
	}
}

// TestQuickRendererSeverity verifies severity is displayed correctly
func TestQuickRendererSeverity(t *testing.T) {
	report := &TeamReport{
		Project: "test",
		Phase:   "TEST",
		Teams: []TeamSection{
			{
				ID:     "security",
				Name:   "Security",
				Status: StatusNoGo,
				Tasks: []TaskResult{
					{ID: "critical-vuln", Status: StatusNoGo, Severity: "critical", Detail: "SQL injection"},
					{ID: "high-vuln", Status: StatusWarn, Severity: "high", Detail: "XSS vulnerability"},
					{ID: "no-severity", Status: StatusGo, Detail: "All checks passed"},
				},
			},
		},
		Status: StatusNoGo,
	}

	var buf bytes.Buffer
	renderer := NewQuickRenderer(&buf)
	if err := renderer.Render(report); err != nil {
		t.Fatalf("render failed: %v", err)
	}

	output := buf.String()

	// Check severity is shown in brackets
	if !strings.Contains(output, "[critical]") {
		t.Error("expected [critical] severity")
	}
	if !strings.Contains(output, "[high]") {
		t.Error("expected [high] severity")
	}

	// Task without severity should not have brackets
	// Find the line with "no-severity" and ensure no brackets
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "no-severity") && strings.Contains(line, "All checks passed") {
			if strings.Contains(line, "[") && strings.Contains(line, "]") {
				// Make sure brackets aren't from severity
				if !strings.Contains(line, "[critical]") && !strings.Contains(line, "[high]") {
					t.Error("task without severity should not have severity brackets")
				}
			}
		}
	}
}

// TestQuickRendererVerdict verifies verdict is displayed in team header
func TestQuickRendererVerdict(t *testing.T) {
	report := &TeamReport{
		Project: "test",
		Phase:   "TEST",
		Teams: []TeamSection{
			{
				ID:      "team-with-verdict",
				Name:    "Security Team",
				Status:  StatusNoGo,
				Verdict: "BLOCKED_BY_CRITICAL",
			},
			{
				ID:     "team-without-verdict",
				Name:   "QA Team",
				Status: StatusGo,
			},
		},
		Status: StatusNoGo,
	}

	var buf bytes.Buffer
	renderer := NewQuickRenderer(&buf)
	if err := renderer.Render(report); err != nil {
		t.Fatalf("render failed: %v", err)
	}

	output := buf.String()

	// Team with verdict should show it
	if !strings.Contains(output, "BLOCKED_BY_CRITICAL") {
		t.Error("expected verdict BLOCKED_BY_CRITICAL in output")
	}

	// Find Security Team line and verify verdict is on same line
	lines := strings.Split(output, "\n")
	foundSecurityTeam := false
	for _, line := range lines {
		if strings.Contains(line, "Security Team") {
			foundSecurityTeam = true
			if !strings.Contains(line, "BLOCKED_BY_CRITICAL") {
				t.Error("verdict should be on same line as team name")
			}
			if !strings.Contains(line, "NO-GO") {
				t.Error("status should be on same line as team name")
			}
		}
	}
	if !foundSecurityTeam {
		t.Error("Security Team not found in output")
	}
}

// TestQuickRendererContentBlocks verifies content blocks are rendered
func TestQuickRendererContentBlocks(t *testing.T) {
	report := &TeamReport{
		Project: "test",
		Phase:   "TEST",
		SummaryBlocks: []ContentBlock{
			NewKVPairsBlock("Summary",
				KVPair{Key: "Version", Value: "1.0.0"},
				KVPair{Key: "Build", Value: "12345"},
			),
		},
		Teams: []TeamSection{
			{
				ID:     "analysis",
				Name:   "Analysis",
				Status: StatusWarn,
				Tasks: []TaskResult{
					{ID: "scan", Status: StatusWarn, Detail: "Issues found"},
				},
				ContentBlocks: []ContentBlock{
					NewListBlock("Findings",
						ListItem{Text: "CVE-2024-001", Icon: "🔴"},
						ListItem{Text: "CVE-2024-002", Icon: "🟡"},
					),
				},
			},
		},
		FooterBlocks: []ContentBlock{
			NewTextBlock("Next Steps", "Review and fix all critical issues before release."),
		},
		Status: StatusWarn,
	}

	var buf bytes.Buffer
	renderer := NewQuickRenderer(&buf)
	if err := renderer.Render(report); err != nil {
		t.Fatalf("render failed: %v", err)
	}

	output := buf.String()

	// Check summary blocks
	if !strings.Contains(output, "Summary") {
		t.Error("expected Summary title")
	}
	if !strings.Contains(output, "Version: 1.0.0") {
		t.Error("expected Version KV pair")
	}
	if !strings.Contains(output, "Build: 12345") {
		t.Error("expected Build KV pair")
	}

	// Check team content blocks
	if !strings.Contains(output, "Findings") {
		t.Error("expected Findings title")
	}
	if !strings.Contains(output, "CVE-2024-001") {
		t.Error("expected CVE-2024-001 in findings")
	}
	if !strings.Contains(output, "🔴") {
		t.Error("expected red icon for critical CVE")
	}

	// Check footer blocks
	if !strings.Contains(output, "Next Steps") {
		t.Error("expected Next Steps title")
	}
	if !strings.Contains(output, "Review and fix") {
		t.Error("expected footer text content")
	}
}

// TestQuickRendererStatusIcons verifies correct status icons
func TestQuickRendererStatusIcons(t *testing.T) {
	report := &TeamReport{
		Project: "test",
		Phase:   "TEST",
		Teams: []TeamSection{
			{
				ID:     "team1",
				Name:   "Go Team",
				Status: StatusGo,
				Tasks: []TaskResult{
					{ID: "go-task", Status: StatusGo, Detail: "Passed"},
				},
			},
			{
				ID:     "team2",
				Name:   "Warn Team",
				Status: StatusWarn,
				Tasks: []TaskResult{
					{ID: "warn-task", Status: StatusWarn, Detail: "Warning"},
				},
			},
			{
				ID:     "team3",
				Name:   "NoGo Team",
				Status: StatusNoGo,
				Tasks: []TaskResult{
					{ID: "nogo-task", Status: StatusNoGo, Detail: "Failed"},
				},
			},
		},
		Status: StatusNoGo,
	}

	var buf bytes.Buffer
	renderer := NewQuickRenderer(&buf)
	if err := renderer.Render(report); err != nil {
		t.Fatalf("render failed: %v", err)
	}

	output := buf.String()

	// Check status icons
	if !strings.Contains(output, "🟢") {
		t.Error("expected green icon (🟢) for GO status")
	}
	if !strings.Contains(output, "🟡") {
		t.Error("expected yellow icon (🟡) for WARN status")
	}
	if !strings.Contains(output, "🔴") {
		t.Error("expected red icon (🔴) for NO-GO status")
	}
}

// TestQuickRendererTableBlock verifies table rendering
func TestQuickRendererTableBlock(t *testing.T) {
	report := &TeamReport{
		Project: "test",
		Phase:   "TEST",
		Teams: []TeamSection{
			{
				ID:     "comparison",
				Name:   "Comparison",
				Status: StatusGo,
				ContentBlocks: []ContentBlock{
					NewTableBlock("Feature Comparison",
						[]string{"Feature", "v1", "v2"},
						[][]string{
							{"Auth", "Basic", "OAuth2"},
							{"Cache", "None", "Redis"},
						},
					),
				},
			},
		},
		Status: StatusGo,
	}

	var buf bytes.Buffer
	renderer := NewQuickRenderer(&buf)
	if err := renderer.Render(report); err != nil {
		t.Fatalf("render failed: %v", err)
	}

	output := buf.String()

	// Check table headers
	if !strings.Contains(output, "Feature") {
		t.Error("expected Feature header")
	}
	if !strings.Contains(output, "v1") {
		t.Error("expected v1 header")
	}
	if !strings.Contains(output, "v2") {
		t.Error("expected v2 header")
	}

	// Check table data
	if !strings.Contains(output, "Auth") {
		t.Error("expected Auth row")
	}
	if !strings.Contains(output, "OAuth2") {
		t.Error("expected OAuth2 value")
	}
	if !strings.Contains(output, "Redis") {
		t.Error("expected Redis value")
	}

	// Check table separator characters
	if !strings.Contains(output, "│") {
		t.Error("expected column separator │")
	}
	if !strings.Contains(output, "─") {
		t.Error("expected row separator ─")
	}
}

// TestQuickRendererMetricBlock verifies metric rendering
func TestQuickRendererMetricBlock(t *testing.T) {
	report := &TeamReport{
		Project: "test",
		Phase:   "TEST",
		Teams: []TeamSection{
			{
				ID:     "metrics",
				Name:   "Metrics",
				Status: StatusGo,
				ContentBlocks: []ContentBlock{
					NewMetricBlock("Coverage", "85%", StatusGo, "80%"),
					NewMetricBlock("Performance", "120ms", StatusWarn, "100ms"),
				},
			},
		},
		Status: StatusGo,
	}

	var buf bytes.Buffer
	renderer := NewQuickRenderer(&buf)
	if err := renderer.Render(report); err != nil {
		t.Fatalf("render failed: %v", err)
	}

	output := buf.String()

	// Check metric with target
	if !strings.Contains(output, "Coverage") {
		t.Error("expected Coverage label")
	}
	if !strings.Contains(output, "85%") {
		t.Error("expected 85% value")
	}
	if !strings.Contains(output, "target: 80%") {
		t.Error("expected target: 80%")
	}

	// Check metric icons match status
	// Coverage is GO so should have green icon nearby
	// Performance is WARN so should have yellow icon nearby
}
