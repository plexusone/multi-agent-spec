package multiagentspec

// ContentBlockType discriminates content block variants.
type ContentBlockType string

const (
	ContentBlockKVPairs ContentBlockType = "kv_pairs"
	ContentBlockList    ContentBlockType = "list"
	ContentBlockTable   ContentBlockType = "table"
	ContentBlockText    ContentBlockType = "text"
	ContentBlockMetric  ContentBlockType = "metric"
)

// ContentBlock represents rich content within a report section.
// The Type field determines which other fields are relevant:
//   - kv_pairs: uses Pairs
//   - list: uses Items
//   - table: uses Headers, Rows
//   - text: uses Content
//   - metric: uses Label, Value, Status
type ContentBlock struct {
	// Type discriminates the block variant.
	Type ContentBlockType `json:"type"`

	// Title is an optional section heading.
	Title string `json:"title,omitempty"`

	// Pairs holds key-value data (for kv_pairs type).
	Pairs []KVPair `json:"pairs,omitempty"`

	// Items holds list entries (for list type).
	Items []ListItem `json:"items,omitempty"`

	// Headers holds table column headers (for table type).
	Headers []string `json:"headers,omitempty"`

	// Rows holds table data rows (for table type).
	Rows [][]string `json:"rows,omitempty"`

	// Content holds multi-line text (for text type).
	Content string `json:"content,omitempty"`

	// Label is the metric name (for metric type).
	Label string `json:"label,omitempty"`

	// Value is the metric value (for metric type).
	Value string `json:"value,omitempty"`

	// Status is used for metric blocks to indicate health.
	Status Status `json:"status,omitempty"`

	// Target is the target value for metric blocks (e.g., "80%").
	Target string `json:"target,omitempty"`
}

// KVPair is a key-value pair with optional icon.
type KVPair struct {
	// Key is the label/identifier.
	Key string `json:"key"`

	// Value is the associated value.
	Value string `json:"value"`

	// Icon is an optional prefix icon (e.g., "🔴", "🟡").
	Icon string `json:"icon,omitempty"`
}

// ListItem is a list entry with optional icon and status.
type ListItem struct {
	// Text is the item content.
	Text string `json:"text"`

	// Icon is an optional prefix icon.
	Icon string `json:"icon,omitempty"`

	// Status allows automatic icon selection if Icon is empty.
	Status Status `json:"status,omitempty"`
}

// EffectiveIcon returns the Icon if set, otherwise derives from Status.
func (li ListItem) EffectiveIcon() string {
	if li.Icon != "" {
		return li.Icon
	}
	if li.Status != "" {
		return li.Status.Icon()
	}
	return ""
}

// NewKVPairsBlock creates a kv_pairs content block.
func NewKVPairsBlock(title string, pairs ...KVPair) ContentBlock {
	return ContentBlock{
		Type:  ContentBlockKVPairs,
		Title: title,
		Pairs: pairs,
	}
}

// NewListBlock creates a list content block.
func NewListBlock(title string, items ...ListItem) ContentBlock {
	return ContentBlock{
		Type:  ContentBlockList,
		Title: title,
		Items: items,
	}
}

// NewTextBlock creates a text content block.
func NewTextBlock(title, content string) ContentBlock {
	return ContentBlock{
		Type:    ContentBlockText,
		Title:   title,
		Content: content,
	}
}

// NewTableBlock creates a table content block.
func NewTableBlock(title string, headers []string, rows [][]string) ContentBlock {
	return ContentBlock{
		Type:    ContentBlockTable,
		Title:   title,
		Headers: headers,
		Rows:    rows,
	}
}

// NewMetricBlock creates a metric content block.
// Target is optional - pass empty string to omit.
func NewMetricBlock(label, value string, status Status, target string) ContentBlock {
	return ContentBlock{
		Type:   ContentBlockMetric,
		Label:  label,
		Value:  value,
		Status: status,
		Target: target,
	}
}
