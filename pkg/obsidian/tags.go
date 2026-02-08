package obsidian

import (
	"regexp"
	"strings"

	"github.com/Yakitrak/obsidian-cli/pkg/frontmatter"
)

// TagLocation indicates where a tag was found
type TagLocation string

const (
	TagFrontmatter TagLocation = "frontmatter"
	TagInline      TagLocation = "inline"
)

// TagMatch represents a single tag occurrence
type TagMatch struct {
	Tag        string
	Location   TagLocation
	LineNumber int
	MatchLine  string
}

// TagSearcher provides tag extraction and search functionality
type TagSearcher struct{}

// tagRegex matches inline #tags but excludes headers and hex codes
// Matches: #tag, #multi-word-tag, #abc123, #programming/philosophy
// Excludes: #123 (numeric headers), #fff (hex codes), ## (markdown)
var tagRegex = regexp.MustCompile(`(?:^|\s)#([a-zA-Z][\w/-]*)`)

// headerRegex detects markdown headers (lines starting with #)
var headerRegex = regexp.MustCompile(`^#{1,6}\s`)

// ExtractTags extracts all tags from note content
func (t *TagSearcher) ExtractTags(content string) []TagMatch {
	var matches []TagMatch

	// Track if we're in a code block
	inCodeBlock := false
	codeBlockDelimiter := ""

	lines := strings.Split(content, "\n")
	inFrontmatter := false
	frontmatterEnded := false

	for i, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		// Handle frontmatter delimiters
		if trimmedLine == "---" && !frontmatterEnded {
			if !inFrontmatter {
				inFrontmatter = true
				continue
			} else {
				inFrontmatter = false
				frontmatterEnded = true
				continue
			}
		}

		// Extract frontmatter tags
		if inFrontmatter && !frontmatterEnded {
			// We handle frontmatter tags separately after parsing
			continue
		}

		// Handle code blocks (fenced)
		if strings.HasPrefix(trimmedLine, "```") || strings.HasPrefix(trimmedLine, "~~~") {
			if !inCodeBlock {
				inCodeBlock = true
				codeBlockDelimiter = trimmedLine[:3]
			} else if strings.HasPrefix(trimmedLine, codeBlockDelimiter) {
				inCodeBlock = false
				codeBlockDelimiter = ""
			}
			continue
		}

		// Skip code block content
		if inCodeBlock {
			continue
		}

		// Extract inline tags from this line
		lineTags := t.extractInlineTagsFromLine(line, i+1)
		matches = append(matches, lineTags...)
	}

	// Parse frontmatter for tags
	if frontmatter.HasFrontmatter(content) {
		fm, _, err := frontmatter.Parse(content)
		if err == nil && fm != nil {
			if tags, ok := fm["tags"]; ok {
				frontmatterTags := t.parseFrontmatterTags(tags)
				for _, tag := range frontmatterTags {
					matches = append(matches, TagMatch{
						Tag:        tag,
						Location:   TagFrontmatter,
						LineNumber: 0,
						MatchLine:  "(frontmatter tag)",
					})
				}
			}
		}
	}

	return matches
}

// extractInlineTagsFromLine extracts #tags from a single line
func (t *TagSearcher) extractInlineTagsFromLine(line string, lineNum int) []TagMatch {
	var matches []TagMatch

	// Skip if line is a markdown header
	if headerRegex.MatchString(strings.TrimSpace(line)) {
		return matches
	}

	// Find all inline tags
	foundTags := tagRegex.FindAllStringSubmatchIndex(line, -1)
	for _, match := range foundTags {
		if len(match) >= 4 {
			// match[0] and match[1] are the full match including the #
			// match[2] and match[3] are the captured group (tag without #)
			tagName := line[match[2]:match[3]]

			// Skip purely numeric tags (hex codes)
			if t.isHexCode(tagName) {
				continue
			}

			// Get context around the tag
			contextStart := match[0]
			contextEnd := match[1]

			// Expand context to full word/phrase if needed
			lineMatches := TagMatch{
				Tag:        tagName,
				Location:   TagInline,
				LineNumber: lineNum,
				MatchLine:  strings.TrimSpace(line[contextStart:contextEnd]),
			}
			matches = append(matches, lineMatches)
		}
	}

	return matches
}

// parseFrontmatterTags extracts tag strings from various frontmatter formats
func (t *TagSearcher) parseFrontmatterTags(tags interface{}) []string {
	var result []string

	switch v := tags.(type) {
	case []interface{}:
		// Format: tags: [a, b, c] or tags:
		//           - a
		//           - b
		for _, tag := range v {
			if str, ok := tag.(string); ok {
				result = append(result, str)
			}
		}
	case string:
		// Format: tags: single-tag
		result = append(result, v)
	case []string:
		// Direct string slice
		result = append(result, v...)
	}

	return result
}

// isHexCode checks if a string looks like a CSS hex color code
// 
// Rules:
// - 3 chars: all must be hex digits (fff, abc, 123) -> hex code
// - 6 chars: must be ALL hex letters (a-f) with NO digits -> hex code like "aabbcc"
// - 6 chars with digits mixed in (like "abc123") -> likely a tag, not hex
func (t *TagSearcher) isHexCode(s string) bool {
	// 3-char hex codes (CSS shorthand colors)
	if len(s) == 3 {
		for _, r := range s {
			if !((r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')) {
				return false
			}
		}
		return true
	}
	
	// 6-char hex codes - only if pure hex letters (no digits)
	// "aabbcc" is a hex code, "abc123" is a tag
	if len(s) == 6 {
		hasOnlyHexLetters := true
		for _, r := range s {
			// Must be hex digit
			if !((r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')) {
				return false
			}
			// If it's a digit (0-9), then it's not pure hex letters
			if r >= '0' && r <= '9' {
				hasOnlyHexLetters = false
			}
		}
		// Only hex code if no digits present (pure letters a-f)
		return hasOnlyHexLetters
	}
	
	return false
}

// HasTag checks if content contains a specific tag
func (t *TagSearcher) HasTag(content string, tag string) bool {
	tags := t.ExtractTags(content)
	tagLower := strings.ToLower(tag)
	for _, t := range tags {
		if strings.ToLower(t.Tag) == tagLower {
			return true
		}
	}
	return false
}

// FindTagLines finds all lines containing a specific tag
func (t *TagSearcher) FindTagLines(content string, tag string) []TagMatch {
	var matches []TagMatch
	tagLower := strings.ToLower(tag)

	allTags := t.ExtractTags(content)
	for _, tm := range allTags {
		if strings.ToLower(tm.Tag) == tagLower {
			matches = append(matches, tm)
		}
	}

	return matches
}

// TagContentMatch represents a tag with surrounding context
type TagContentMatch struct {
	FilePath     string
	TagLineNum   int
	ContextStart int
	ContextEnd   int
	ContextLines []string
	Tag          string
}

// FindTagBlocks finds tag occurrences with surrounding paragraph/block context
func (t *TagSearcher) FindTagBlocks(content string, tag string, contextLines int, includeSubtags bool) []TagContentMatch {
	var matches []TagContentMatch
	tagLower := strings.ToLower(tag)

	lines := strings.Split(content, "\n")

	// Find all lines with the tag
	tagLineNums := make(map[int]string) // line number -> matched tag
	for i, line := range lines {
		lineTags := t.extractInlineTagsFromLine(line, i+1)
		for _, lt := range lineTags {
			ltTagLower := strings.ToLower(lt.Tag)
			isMatch := ltTagLower == tagLower
			if includeSubtags && !isMatch {
				isMatch = strings.HasPrefix(ltTagLower, tagLower+"/")
			}
			if isMatch {
				tagLineNums[i] = lt.Tag
				break
			}
		}
	}

	// Build blocks around each tag line
	for lineNum, matchedTag := range tagLineNums {
		start := lineNum - contextLines
		if start < 0 {
			start = 0
		}
		end := lineNum + contextLines + 1
		if end > len(lines) {
			end = len(lines)
		}

		var context []string
		for i := start; i < end && i < len(lines); i++ {
			context = append(context, lines[i])
		}

		matches = append(matches, TagContentMatch{
			TagLineNum:   lineNum + 1,
			ContextStart: start + 1,
			ContextEnd:   end,
			ContextLines: context,
			Tag:          matchedTag,
		})
	}

	return matches
}
