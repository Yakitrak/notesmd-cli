package obsidian

import (
	"strings"
	"testing"
)

func TestTagSearcher_ExtractTags_Inline(t *testing.T) {
	searcher := &TagSearcher{}

	tests := []struct {
		name     string
		content  string
		expected []TagMatch
	}{
		{
			name:    "simple inline tag",
			content: "This is a note with #daily tag",
			expected: []TagMatch{
				{Tag: "daily", Location: TagInline, LineNumber: 1, MatchLine: "#daily"},
			},
		},
		{
			name:    "multiple inline tags on same line",
			content: "Meeting notes #work #urgent #follow-up",
			expected: []TagMatch{
				{Tag: "work", Location: TagInline, LineNumber: 1},
				{Tag: "urgent", Location: TagInline, LineNumber: 1},
				{Tag: "follow-up", Location: TagInline, LineNumber: 1},
			},
		},
		{
			name:    "multi-word tag with hyphens",
			content: "Project #alpha-release is ready",
			expected: []TagMatch{
				{Tag: "alpha-release", Location: TagInline, LineNumber: 1},
			},
		},
		{
			name:    "tags on multiple lines",
			content: "Line one with #first\nLine two with #second",
			expected: []TagMatch{
				{Tag: "first", Location: TagInline, LineNumber: 1},
				{Tag: "second", Location: TagInline, LineNumber: 2},
			},
		},
		{
			name:    "no tags",
			content: "This note has no tags at all",
			expected: []TagMatch{},
		},
		{
			name:    "tag at start of line",
			content: "#morning routine starts here",
			expected: []TagMatch{
				{Tag: "morning", Location: TagInline, LineNumber: 1},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := searcher.ExtractTags(tt.content)

			if len(result) != len(tt.expected) {
				t.Errorf("expected %d tags, got %d", len(tt.expected), len(result))
			}

			for i, expected := range tt.expected {
				if i >= len(result) {
					break
				}
				if result[i].Tag != expected.Tag {
					t.Errorf("expected tag %s, got %s", expected.Tag, result[i].Tag)
				}
				if result[i].Location != expected.Location {
					t.Errorf("expected location %s, got %s", expected.Location, result[i].Location)
				}
				if result[i].LineNumber != expected.LineNumber {
					t.Errorf("expected line %d, got %d", expected.LineNumber, result[i].LineNumber)
				}
			}
		})
	}
}

func TestTagSearcher_ExtractTags_Frontmatter(t *testing.T) {
	searcher := &TagSearcher{}

	tests := []struct {
		name     string
		content  string
		expected []string
	}{
		{
			name: "frontmatter with array",
			content: `---
tags: [daily, review, meeting]
---
This is the content`,
			expected: []string{"daily", "review", "meeting"},
		},
		{
			name: "frontmatter with list",
			content: `---
tags:
  - project
  - planning
  - q1-2024
---
Content here`,
			expected: []string{"project", "planning", "q1-2024"},
		},
		{
			name: "frontmatter with single string",
			content: `---
tags: single-tag
---
Content`,
			expected: []string{"single-tag"},
		},
		{
			name: "frontmatter with empty tags",
			content: `---
tags: []
---
Content`,
			expected: []string{},
		},
		{
			name: "no frontmatter",
			content: `This note has no frontmatter
Just plain content`,
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := searcher.ExtractTags(tt.content)

			// Filter only frontmatter tags
			var frontmatterTags []string
			for _, match := range result {
				if match.Location == TagFrontmatter {
					frontmatterTags = append(frontmatterTags, match.Tag)
				}
			}

			if len(frontmatterTags) != len(tt.expected) {
				t.Errorf("expected %d frontmatter tags, got %d: %v", len(tt.expected), len(frontmatterTags), frontmatterTags)
			}

			for i, expected := range tt.expected {
				if i >= len(frontmatterTags) {
					break
				}
				if frontmatterTags[i] != expected {
					t.Errorf("expected tag %s, got %s", expected, frontmatterTags[i])
				}
			}
		})
	}
}

func TestTagSearcher_ExtractTags_Combined(t *testing.T) {
	searcher := &TagSearcher{}

	content := `---
tags: [daily, work]
---
This is my #morning routine
Need to check #emails and #calendar`

	result := searcher.ExtractTags(content)

	// Should have 2 frontmatter + 3 inline = 5 total
	if len(result) != 5 {
		t.Errorf("expected 5 tags total, got %d", len(result))
	}

	// Check frontmatter tags
	var frontmatterCount, inlineCount int
	for _, match := range result {
		if match.Location == TagFrontmatter {
			frontmatterCount++
		} else if match.Location == TagInline {
			inlineCount++
		}
	}

	if frontmatterCount != 2 {
		t.Errorf("expected 2 frontmatter tags, got %d", frontmatterCount)
	}
	if inlineCount != 3 {
		t.Errorf("expected 3 inline tags, got %d", inlineCount)
	}
}

func TestTagSearcher_ExtractTags_Exclusions(t *testing.T) {
	searcher := &TagSearcher{}

	tests := []struct {
		name        string
		content     string
		shouldFind  []string
		shouldNotFind []string
	}{
		{
			name:    "markdown headers excluded",
			content: "# Header 1\n## Header 2\nContent with #real-tag",
			shouldFind: []string{"real-tag"},
			shouldNotFind: []string{"Header", "1", "2"},
		},
		{
			name:    "hex codes excluded",
			content: "Color #fff and #aabbcc and #123\nBut #design tag is real",
			shouldFind: []string{"design"},
			shouldNotFind: []string{"fff", "aabbcc", "123"},
		},
		{
			name:    "code blocks excluded",
			content: "```\n#not-a-tag in code\n```\nBut #outside is real",
			shouldFind: []string{"outside"},
			shouldNotFind: []string{"not-a-tag"},
		},
		{
			name:    "numeric-only excluded",
			content: "#123abc is valid, #abc123 is valid",
			shouldFind: []string{"abc123"},
			shouldNotFind: []string{"123abc"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := searcher.ExtractTags(tt.content)
			foundTags := make(map[string]bool)
			for _, match := range result {
				foundTags[match.Tag] = true
			}

			for _, tag := range tt.shouldFind {
				if !foundTags[tag] {
					t.Errorf("should find tag %s but didn't", tag)
				}
			}

			for _, tag := range tt.shouldNotFind {
				if foundTags[tag] {
					t.Errorf("should NOT find tag %s but did", tag)
				}
			}
		})
	}
}

func TestTagSearcher_HasTag(t *testing.T) {
	searcher := &TagSearcher{}

	content := "Note with #daily and #work tags"

	if !searcher.HasTag(content, "daily") {
		t.Error("expected to find 'daily' tag")
	}
	if !searcher.HasTag(content, "work") {
		t.Error("expected to find 'work' tag")
	}
	if searcher.HasTag(content, "nonexistent") {
		t.Error("should not find 'nonexistent' tag")
	}

	// Case insensitive
	if !searcher.HasTag(content, "DAILY") {
		t.Error("expected case-insensitive match for 'DAILY'")
	}
}

func TestTagSearcher_FindTagLines(t *testing.T) {
	searcher := &TagSearcher{}

	content := `Line 1 has #first tag
Line 2 is empty
Line 3 has #second tag
Line 4 also has #first tag again`

	result := searcher.FindTagLines(content, "first")

	if len(result) != 2 {
		t.Errorf("expected 2 matches for 'first', got %d", len(result))
	}

	// Check line numbers
	if result[0].LineNumber != 1 {
		t.Errorf("expected first match on line 1, got %d", result[0].LineNumber)
	}
	if result[1].LineNumber != 4 {
		t.Errorf("expected second match on line 4, got %d", result[1].LineNumber)
	}
}

func TestTagSearcher_FindTagBlocks(t *testing.T) {
	searcher := &TagSearcher{}

	// Note: Order of results may vary due to map iteration in FindTagBlocks.
	// We check for presence of both matches rather than their order.
	content := `Paragraph one
This has the #target tag here
Paragraph three

Another section
With another #target mention
End of content`

	result := searcher.FindTagBlocks(content, "target", 1, false)

	if len(result) != 2 {
		t.Errorf("expected 2 blocks for 'target', got %d", len(result))
	}

	// Check that we found both lines (order may vary due to map iteration)
	foundLine2 := false
	foundLine6 := false
	for _, r := range result {
		if r.TagLineNum == 2 {
			foundLine2 = true
			// With contextLines=1, block should have lines 1, 2, 3
			if r.ContextStart != 1 {
				t.Errorf("expected context start at 1, got %d", r.ContextStart)
			}
			if r.ContextEnd != 3 {
				t.Errorf("expected context end at 3, got %d", r.ContextEnd)
			}
			// Check content includes surrounding lines
			contextStr := strings.Join(r.ContextLines, "\n")
			if !strings.Contains(contextStr, "Paragraph one") {
				t.Error("context should include 'Paragraph one'")
			}
			if !strings.Contains(contextStr, "#target") {
				t.Error("context should include '#target'")
			}
		}
		if r.TagLineNum == 6 {
			foundLine6 = true
		}
	}

	if !foundLine2 {
		t.Error("should have found tag on line 2")
	}
	if !foundLine6 {
		t.Error("should have found tag on line 6")
	}
}

func TestTagSearcher_FindTagBlocks_Subtags(t *testing.T) {
	searcher := &TagSearcher{}

	content := `Paragraph one
This has the #parent/child tag here
Paragraph three

Another section
With #parent/grandchild/nested
End of content`

	// Without subtags - searching "parent" should not match "parent/child"
	result := searcher.FindTagBlocks(content, "parent", 1, false)
	if len(result) != 0 {
		t.Errorf("expected 0 blocks without subtags, got %d", len(result))
	}

	// With subtags - searching "parent" should match "parent/child" and "parent/grandchild/nested"
	result = searcher.FindTagBlocks(content, "parent", 1, true)
	if len(result) != 2 {
		t.Errorf("expected 2 blocks with subtags, got %d", len(result))
	}

	// Check that both subtags were captured (order may vary)
	foundChild := false
	foundGrandchild := false
	for _, r := range result {
		if r.Tag == "parent/child" {
			foundChild = true
		}
		if r.Tag == "parent/grandchild/nested" {
			foundGrandchild = true
		}
	}
	if !foundChild {
		t.Error("expected to find 'parent/child'")
	}
	if !foundGrandchild {
		t.Error("expected to find 'parent/grandchild/nested'")
	}
}

func TestTagSearcher_parseFrontmatterTags(t *testing.T) {
	searcher := &TagSearcher{}

	tests := []struct {
		name     string
		input    interface{}
		expected []string
	}{
		{
			name:     "string slice",
			input:    []string{"a", "b", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "interface slice",
			input:    []interface{}{"x", "y", "z"},
			expected: []string{"x", "y", "z"},
		},
		{
			name:     "single string",
			input:    "single",
			expected: []string{"single"},
		},
		{
			name:     "nil",
			input:    nil,
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := searcher.parseFrontmatterTags(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("expected %d tags, got %d", len(tt.expected), len(result))
			}
			for i, expected := range tt.expected {
				if i >= len(result) {
					break
				}
				if result[i] != expected {
					t.Errorf("expected %s, got %s", expected, result[i])
				}
			}
		})
	}
}

func TestTagSearcher_isHexCode(t *testing.T) {
	searcher := &TagSearcher{}

	tests := []struct {
		input    string
		expected bool
	}{
		{"fff", true},
		{"FFF", true},
		{"abc", true},
		{"ABC", true},
		{"123", true},
		{"aabbcc", true},
		{"AAbbCC", true},
		{"daily", false},
		{"work", false},
		{"project-alpha", false},
		{"ff", false},     // too short
		{"ggg", false},    // invalid hex
		{"fffff", false},  // wrong length
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := searcher.isHexCode(tt.input)
			if result != tt.expected {
				t.Errorf("isHexCode(%s) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}
