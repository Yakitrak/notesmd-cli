package rules_test

import (
	"testing"

	"github.com/Yakitrak/notesmd-cli/pkg/validate"
	"github.com/Yakitrak/notesmd-cli/pkg/validate/rules"
	"github.com/stretchr/testify/assert"
)

func TestCheckStructure_HeadingHierarchy(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    int
	}{
		{
			"valid hierarchy H1 H2 H3",
			"# Title\n## Section\n### Subsection\n",
			0,
		},
		{
			"skipped level H1 to H3",
			"# Title\n### Subsection\n",
			1,
		},
		{
			"skipped level H2 to H4",
			"## Section\n#### Deep\n",
			1,
		},
		{
			"multiple skips",
			"# Title\n### Skip1\n##### Skip2\n",
			2,
		},
		{
			"no headings",
			"Just some text\nMore text\n",
			0,
		},
		{
			"heading without space is not a heading",
			"#notaheading\n## Real heading\n",
			0,
		},
		{
			"frontmatter is ignored",
			"---\ntitle: Test\n---\n# Title\n## Section\n",
			0,
		},
		{
			"going down then back up is valid",
			"# Title\n## Section\n### Sub\n## Another Section\n",
			0,
		},
		{
			"H2 then H3 is valid even without H1",
			"## Section\n### Subsection\n",
			0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := rules.CheckStructure(tt.content, "test.md", defaultConfig())
			headingResults := filterByRule(results, "structure.heading-hierarchy")
			assert.Len(t, headingResults, tt.want)
			for _, r := range headingResults {
				assert.Equal(t, validate.SeverityWarning, r.Severity)
				assert.Greater(t, r.Line, 0)
			}
		})
	}
}

func TestCheckStructure_HeadingHierarchyLineNumbers(t *testing.T) {
	t.Run("reports correct line number", func(t *testing.T) {
		content := "# Title\nSome text\n### Skipped H3\n"
		results := rules.CheckStructure(content, "test.md", defaultConfig())
		assert.Len(t, results, 1)
		assert.Equal(t, 3, results[0].Line)
	})
}

func TestCheckStructure_SeverityFilter(t *testing.T) {
	t.Run("error-only severity skips heading warnings", func(t *testing.T) {
		content := "# Title\n### Skipped\n"
		config := validate.Config{
			ExcludeDirs: validate.DefaultExcludeDirs(),
			MinSeverity: validate.SeverityError,
		}
		results := rules.CheckStructure(content, "test.md", config)
		assert.Empty(t, results)
	})
}
