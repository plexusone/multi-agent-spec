package multiagentspec

import (
	"fmt"
	"io"
	"strings"
	"text/template"
)

const (
	// boxWidth is the inner width of the box (between the border characters).
	boxWidth = 78
)

// Renderer renders TeamReport to various formats using text/template.
type Renderer struct {
	w io.Writer
}

// NewRenderer creates a new Renderer writing to w.
func NewRenderer(w io.Writer) *Renderer {
	return &Renderer{w: w}
}

// Render renders the report using the box template.
// It automatically sorts teams by DAG order before rendering.
func (r *Renderer) Render(report *TeamReport) error {
	report.SortByDAG()
	return r.renderBox(report)
}

// renderBox renders the report in the box format.
func (r *Renderer) renderBox(report *TeamReport) error {
	tmpl, err := template.New("report").Funcs(templateFuncs()).Parse(BoxTemplate)
	if err != nil {
		return fmt.Errorf("parsing template: %w", err)
	}
	return tmpl.Execute(r.w, report)
}

// QuickRenderer renders TeamReport using quicktemplate (compile-time type-safe).
// This is an alternative to Renderer that uses generated code instead of reflection.
type QuickRenderer struct {
	w io.Writer
}

// NewQuickRenderer creates a new QuickRenderer writing to w.
func NewQuickRenderer(w io.Writer) *QuickRenderer {
	return &QuickRenderer{w: w}
}

// Render renders the report using quicktemplate.
// It automatically sorts teams by DAG order before rendering.
func (r *QuickRenderer) Render(report *TeamReport) error {
	report.SortByDAG()
	WriteBoxReport(r.w, report)
	return nil
}

// templateFuncs returns the template function map.
func templateFuncs() template.FuncMap {
	return template.FuncMap{
		"header":           header,
		"separator":        separator,
		"footer":           footer,
		"teamHeader":       teamHeader,
		"taskLine":         taskLine,
		"centerLine":       centerLine,
		"paddedLine":       paddedLine,
		"finalMessage":     finalMessage,
		"renderBlock":      renderBlock,
		"renderBlocks":     renderBlocks,
		"hasContentBlocks": hasContentBlocks,
		"hasSummaryBlocks": hasSummaryBlocks,
		"hasFooterBlocks":  hasFooterBlocks,
		"hasTags":          hasTags,
		"renderTags":       renderTags,
	}
}

// header returns the top border of the box.
func header() string {
	return "╔" + strings.Repeat("═", boxWidth) + "╗"
}

// separator returns a separator line.
func separator() string {
	return "╠" + strings.Repeat("═", boxWidth) + "╣"
}

// footer returns the bottom border of the box.
func footer() string {
	return "╚" + strings.Repeat("═", boxWidth) + "╝"
}

// centerLine centers text within the box.
func centerLine(text string) string {
	visualLen := visualLength(text)
	padding := max(0, boxWidth-visualLen)
	left := padding / 2
	right := padding - left
	return "║" + strings.Repeat(" ", left) + text + strings.Repeat(" ", right) + "║"
}

// paddedLine left-aligns text with padding.
func paddedLine(text string) string {
	visualLen := visualLength(text)
	padding := max(0, boxWidth-visualLen-1)
	return "║ " + text + strings.Repeat(" ", padding) + "║"
}

// teamHeader formats a team header line with status icon and optional verdict.
func teamHeader(team TeamSection) string {
	icon := team.Status.Icon()
	var text string
	if team.Verdict != "" {
		text = fmt.Sprintf("%s %s — %s — %s", icon, team.Name, team.Status, team.Verdict)
	} else {
		text = fmt.Sprintf("%s %s — %s", icon, team.Name, team.Status)
	}
	return paddedLine(text)
}

// taskLine formats a single task result line with optional severity.
func taskLine(task TaskResult) string {
	id := task.ID
	if len(id) > 24 {
		id = id[:21] + "..."
	}

	icon := task.Status.Icon()
	statusText := string(task.Status)

	// Add severity in brackets if present
	if task.Severity != "" {
		statusText = fmt.Sprintf("%s [%s]", statusText, task.Severity)
	}

	detail := task.Detail
	maxDetail := boxWidth - 45 // Reduced to accommodate severity
	if len(detail) > maxDetail {
		detail = detail[:maxDetail-3] + "..."
	}

	line := fmt.Sprintf("  %-24s %s %-15s %s", id, icon, statusText, detail)
	return paddedLine(line)
}

// finalMessage formats the final status message line.
func finalMessage(report *TeamReport) string {
	return centerLine(report.FinalMessage())
}

// visualLength calculates the visual length of a string,
// accounting for emoji characters that take 2 columns.
func visualLength(s string) int {
	length := 0
	for _, r := range s {
		if r >= 0x1F300 && r <= 0x1FAFF {
			length += 2
		} else if r >= 0x2600 && r <= 0x27BF {
			length += 2
		} else {
			length++
		}
	}
	return length
}

// hasContentBlocks returns true if the team has content blocks.
func hasContentBlocks(team TeamSection) bool {
	return len(team.ContentBlocks) > 0
}

// hasSummaryBlocks returns true if the report has summary blocks.
func hasSummaryBlocks(report *TeamReport) bool {
	return len(report.SummaryBlocks) > 0
}

// hasFooterBlocks returns true if the report has footer blocks.
func hasFooterBlocks(report *TeamReport) bool {
	return len(report.FooterBlocks) > 0
}

// hasTags returns true if the report has tags.
func hasTags(report *TeamReport) bool {
	return len(report.Tags) > 0
}

// renderTags renders tags as key-value lines, sorted by key.
func renderTags(tags map[string]string) string {
	// Sort keys for deterministic output
	keys := make([]string, 0, len(tags))
	for k := range tags {
		keys = append(keys, k)
	}
	sortStrings(keys)

	var lines []string
	for _, k := range keys {
		lines = append(lines, paddedLine(fmt.Sprintf("  %s: %s", k, tags[k])))
	}
	return strings.Join(lines, "\n")
}

// renderBlocks renders multiple content blocks, returning joined lines.
func renderBlocks(blocks []ContentBlock) string {
	var lines []string
	for _, block := range blocks {
		lines = append(lines, renderBlock(block))
	}
	return strings.Join(lines, "\n")
}

// renderBlock renders a single content block to box-formatted lines.
func renderBlock(block ContentBlock) string {
	var lines []string

	// Add title if present
	if block.Title != "" {
		lines = append(lines, paddedLine(block.Title))
	}

	switch block.Type {
	case ContentBlockKVPairs:
		lines = append(lines, renderKVPairs(block.Pairs)...)
	case ContentBlockList:
		lines = append(lines, renderList(block.Items)...)
	case ContentBlockText:
		lines = append(lines, wrapText(block.Content, boxWidth-2)...)
	case ContentBlockTable:
		lines = append(lines, renderTable(block.Headers, block.Rows)...)
	case ContentBlockMetric:
		lines = append(lines, renderMetric(block.Label, block.Value, block.Status, block.Target))
	}

	return strings.Join(lines, "\n")
}

// renderKVPairs renders key-value pairs.
func renderKVPairs(pairs []KVPair) []string {
	var lines []string
	for _, pair := range pairs {
		var text string
		if pair.Icon != "" {
			text = fmt.Sprintf("%s %s: %s", pair.Icon, pair.Key, pair.Value)
		} else {
			text = fmt.Sprintf("%s: %s", pair.Key, pair.Value)
		}
		lines = append(lines, paddedLine(text))
	}
	return lines
}

// renderList renders list items.
func renderList(items []ListItem) []string {
	var lines []string
	for _, item := range items {
		icon := item.EffectiveIcon()
		var text string
		if icon != "" {
			text = fmt.Sprintf("%s %s", icon, item.Text)
		} else {
			text = fmt.Sprintf("  %s", item.Text)
		}
		lines = append(lines, paddedLine(text))
	}
	return lines
}

// wrapText wraps text to fit within maxWidth, returning padded lines.
func wrapText(content string, maxWidth int) []string {
	var lines []string
	words := strings.Fields(content)
	if len(words) == 0 {
		return lines
	}

	currentLine := ""
	for _, word := range words {
		if currentLine == "" {
			currentLine = word
		} else if len(currentLine)+1+len(word) <= maxWidth {
			currentLine += " " + word
		} else {
			lines = append(lines, paddedLine(currentLine))
			currentLine = word
		}
	}
	if currentLine != "" {
		lines = append(lines, paddedLine(currentLine))
	}
	return lines
}

// renderTable renders a simple table.
func renderTable(headers []string, rows [][]string) []string {
	var lines []string

	// Calculate column widths
	colWidths := make([]int, len(headers))
	for i, h := range headers {
		colWidths[i] = len(h)
	}
	for _, row := range rows {
		for i, cell := range row {
			if i < len(colWidths) && len(cell) > colWidths[i] {
				colWidths[i] = len(cell)
			}
		}
	}

	// Render header row
	headerParts := make([]string, len(headers))
	for i, h := range headers {
		headerParts[i] = fmt.Sprintf("%-*s", colWidths[i], h)
	}
	lines = append(lines, paddedLine(strings.Join(headerParts, " │ ")))

	// Render separator
	sepParts := make([]string, len(headers))
	for i := range headers {
		sepParts[i] = strings.Repeat("─", colWidths[i])
	}
	lines = append(lines, paddedLine(strings.Join(sepParts, "─┼─")))

	// Render data rows
	for _, row := range rows {
		rowParts := make([]string, len(headers))
		for i := range headers {
			cell := ""
			if i < len(row) {
				cell = row[i]
			}
			rowParts[i] = fmt.Sprintf("%-*s", colWidths[i], cell)
		}
		lines = append(lines, paddedLine(strings.Join(rowParts, " │ ")))
	}

	return lines
}

// renderMetric renders a single metric with status icon and optional target.
func renderMetric(label, value string, status Status, target string) string {
	icon := status.Icon()
	var text string
	if target != "" {
		text = fmt.Sprintf("%s %s: %s (target: %s)", icon, label, value, target)
	} else {
		text = fmt.Sprintf("%s %s: %s", icon, label, value)
	}
	return paddedLine(text)
}

// BoxTemplate is the text/template for the box format report.
// This is the reference implementation for rendering TeamReport to text.
// Each task is rendered on its own line with status indicator.
// Content blocks are rendered after tasks within each team section.
const BoxTemplate = `{{ header }}
{{ centerLine .EffectiveTitle }}
{{ separator }}
{{- if hasSummaryBlocks . }}
{{ renderBlocks .SummaryBlocks }}
{{ separator }}
{{- else }}
{{ paddedLine (printf "Project: %s" .Project) }}
{{ paddedLine (printf "Target:  %s" .Target) }}
{{- if hasTags . }}
{{ paddedLine "Tags:" }}
{{ renderTags .Tags }}
{{- end }}
{{ separator }}
{{- end }}
{{ paddedLine .Phase }}
{{- range .Teams }}
{{ separator }}
{{ teamHeader . }}
{{- range .Tasks }}
{{ taskLine . }}
{{- end }}
{{- if hasContentBlocks . }}
{{ renderBlocks .ContentBlocks }}
{{- end }}
{{- end }}
{{- if hasFooterBlocks . }}
{{ separator }}
{{ renderBlocks .FooterBlocks }}
{{- end }}
{{ separator }}
{{ finalMessage . }}
{{ footer }}
`
