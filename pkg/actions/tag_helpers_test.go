package actions

import (
	"testing"

	"github.com/Yakitrak/obsidian-cli/pkg/obsidian"
)

func TestMatchesTag(t *testing.T) {
	tests := []struct {
		name           string
		tag            string
		searchTag      string
		includeSubtags bool
		want           bool
	}{
		{"exact match", "programming", "programming", false, true},
		{"case insensitive exact", "Programming", "programming", false, true},
		{"no match", "python", "programming", false, false},
		{"subtag without flag", "programming/philosophy", "programming", false, false},
		{"subtag with flag", "programming/philosophy", "programming", true, true},
		{"deeply nested subtag", "a/b/c/d", "a", true, true},
		{"not a subtag - different prefix", "abc/def", "ab", true, false},
		{"partial match not subtag", "programmer", "programming", true, false},
		{"empty tag", "", "programming", false, false},
		{"empty search", "programming", "", false, false},
		{"both empty", "", "", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := matchesTag(tt.tag, tt.searchTag, tt.includeSubtags)
			if got != tt.want {
				t.Errorf("matchesTag(%q, %q, %v) = %v, want %v",
					tt.tag, tt.searchTag, tt.includeSubtags, got, tt.want)
			}
		})
	}
}

func TestFilterByLocation(t *testing.T) {
	tests := []struct {
		name     string
		location string
		tagLoc   obsidian.TagLocation
		want     bool
	}{
		{"all accepts frontmatter", "all", obsidian.TagFrontmatter, true},
		{"all accepts inline", "all", obsidian.TagInline, true},
		{"frontmatter filter matches", "frontmatter", obsidian.TagFrontmatter, true},
		{"frontmatter filter rejects inline", "frontmatter", obsidian.TagInline, false},
		{"inline filter matches", "inline", obsidian.TagInline, true},
		{"inline filter rejects frontmatter", "inline", obsidian.TagFrontmatter, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm := obsidian.TagMatch{Location: tt.tagLoc}
			got := filterByLocation(tm, tt.location)
			if got != tt.want {
				t.Errorf("filterByLocation(%v, %q) = %v, want %v",
					tt.tagLoc, tt.location, got, tt.want)
			}
		})
	}
}

func TestFormatLocationMessage(t *testing.T) {
	tests := []struct {
		name           string
		location       string
		includeSubtags bool
		want           string
	}{
		{"no filters", "all", false, ""},
		{"empty location", "", false, ""},
		{"subtags only", "all", true, " (including subtags)"},
		{"location only", "frontmatter", false, " (in frontmatter)"},
		{"both filters", "inline", true, " (in inline, including subtags)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatLocationMessage(tt.location, tt.includeSubtags)
			if got != tt.want {
				t.Errorf("formatLocationMessage(%q, %v) = %q, want %q",
					tt.location, tt.includeSubtags, got, tt.want)
			}
		})
	}
}