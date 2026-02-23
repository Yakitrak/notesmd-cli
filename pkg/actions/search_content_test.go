package actions_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
	"testing"

	"github.com/Yakitrak/notesmd-cli/mocks"
	"github.com/Yakitrak/notesmd-cli/pkg/actions"
	"github.com/Yakitrak/notesmd-cli/pkg/obsidian"
	"github.com/stretchr/testify/assert"
)

// CustomMockNoteForSingleMatch returns exactly one match for editor testing
// and for checking single-result interactive behavior.
type CustomMockNoteForSingleMatch struct{}

func (m *CustomMockNoteForSingleMatch) Delete(string) error                        { return nil }
func (m *CustomMockNoteForSingleMatch) Move(string, string) error                  { return nil }
func (m *CustomMockNoteForSingleMatch) UpdateLinks(string, string, string) error   { return nil }
func (m *CustomMockNoteForSingleMatch) GetContents(string, string) (string, error) { return "", nil }
func (m *CustomMockNoteForSingleMatch) SetContents(string, string, string) error   { return nil }
func (m *CustomMockNoteForSingleMatch) GetNotesList(string) ([]string, error)      { return nil, nil }
func (m *CustomMockNoteForSingleMatch) SearchNotesWithSnippets(string, string) ([]obsidian.NoteMatch, error) {
	return []obsidian.NoteMatch{
		{FilePath: "test-note.md", LineNumber: 5, MatchLine: "test content"},
	}, nil
}
func (m *CustomMockNoteForSingleMatch) FindBacklinks(string, string) ([]obsidian.NoteMatch, error) {
	return nil, nil
}

type searchContentJSONMatch struct {
	File      string `json:"file"`
	Line      int    `json:"line"`
	Content   string `json:"content"`
	MatchType string `json:"match_type"`
}

func defaultOptions(output *bytes.Buffer) actions.SearchContentOptions {
	return actions.SearchContentOptions{
		Format:              "text",
		InteractiveTerminal: true,
		Output:              output,
	}
}

func captureStderr(t *testing.T, fn func()) string {
	t.Helper()
	oldStderr := os.Stderr
	r, w, err := os.Pipe()
	assert.NoError(t, err)
	os.Stderr = w

	fn()

	_ = w.Close()
	os.Stderr = oldStderr
	stderrBytes, readErr := io.ReadAll(r)
	assert.NoError(t, readErr)
	_ = r.Close()
	return string(stderrBytes)
}

func TestSearchNotesContent(t *testing.T) {
	t.Run("Backward compatible SearchNotesContent API still works", func(t *testing.T) {
		vault := mocks.MockVaultOperator{Name: "myVault"}
		uri := mocks.MockUriManager{}
		note := mocks.MockNoteManager{}
		fuzzyFinder := mocks.MockFuzzyFinder{SelectedIndex: 0}

		err := actions.SearchNotesContent(&vault, &note, &uri, &fuzzyFinder, "test", false)
		assert.NoError(t, err)
		assert.Equal(t, 1, uri.ExecuteCalls)
	})

	t.Run("Successful interactive content search with multiple matches", func(t *testing.T) {
		vault := mocks.MockVaultOperator{Name: "myVault"}
		uri := mocks.MockUriManager{}
		note := mocks.MockNoteManager{}
		fuzzyFinder := mocks.MockFuzzyFinder{SelectedIndex: 0}
		output := &bytes.Buffer{}

		err := actions.SearchNotesContentWithOptions(&vault, &note, &uri, &fuzzyFinder, "test", defaultOptions(output))
		assert.NoError(t, err)
		assert.Equal(t, 1, uri.ExecuteCalls)
	})

	t.Run("No matches found in interactive mode", func(t *testing.T) {
		vault := mocks.MockVaultOperator{Name: "myVault"}
		uri := mocks.MockUriManager{}
		note := mocks.MockNoteManager{NoMatches: true}
		fuzzyFinder := mocks.MockFuzzyFinder{}
		output := &bytes.Buffer{}

		err := actions.SearchNotesContentWithOptions(&vault, &note, &uri, &fuzzyFinder, "nonexistent", defaultOptions(output))
		assert.NoError(t, err)
		assert.Contains(t, output.String(), "No notes found containing 'nonexistent'")
	})

	t.Run("No matches found in non-interactive text mode prints message to stderr", func(t *testing.T) {
		vault := mocks.MockVaultOperator{Name: "myVault"}
		uri := mocks.MockUriManager{}
		note := mocks.MockNoteManager{NoMatches: true}
		fuzzyFinder := mocks.MockFuzzyFinder{}
		output := &bytes.Buffer{}

		options := defaultOptions(output)
		options.NoInteractive = true

		stderr := captureStderr(t, func() {
			err := actions.SearchNotesContentWithOptions(&vault, &note, &uri, &fuzzyFinder, "nonexistent", options)
			assert.NoError(t, err)
		})

		assert.Equal(t, "", output.String())
		assert.Contains(t, stderr, "No notes found containing 'nonexistent'")
	})

	t.Run("SearchNotesWithSnippets returns error", func(t *testing.T) {
		vault := mocks.MockVaultOperator{Name: "myVault"}
		uri := mocks.MockUriManager{}
		note := mocks.MockNoteManager{GetContentsError: errors.New("search failed")}
		fuzzyFinder := mocks.MockFuzzyFinder{}
		output := &bytes.Buffer{}

		err := actions.SearchNotesContentWithOptions(&vault, &note, &uri, &fuzzyFinder, "test", defaultOptions(output))
		assert.Error(t, err)
	})

	t.Run("vault.DefaultName returns error", func(t *testing.T) {
		vault := mocks.MockVaultOperator{DefaultNameErr: errors.New("vault name error")}
		uri := mocks.MockUriManager{}
		note := mocks.MockNoteManager{}
		fuzzyFinder := mocks.MockFuzzyFinder{}
		output := &bytes.Buffer{}

		err := actions.SearchNotesContentWithOptions(&vault, &note, &uri, &fuzzyFinder, "test", defaultOptions(output))
		assert.Error(t, err)
	})

	t.Run("vault.Path returns error", func(t *testing.T) {
		vault := mocks.MockVaultOperator{PathError: errors.New("vault path error")}
		uri := mocks.MockUriManager{}
		note := mocks.MockNoteManager{}
		fuzzyFinder := mocks.MockFuzzyFinder{}
		output := &bytes.Buffer{}

		err := actions.SearchNotesContentWithOptions(&vault, &note, &uri, &fuzzyFinder, "test", defaultOptions(output))
		assert.Error(t, err)
	})

	t.Run("fuzzy finder returns error", func(t *testing.T) {
		vault := mocks.MockVaultOperator{Name: "myVault"}
		uri := mocks.MockUriManager{}
		note := mocks.MockNoteManager{}
		fuzzyFinder := mocks.MockFuzzyFinder{FindErr: errors.New("fuzzy finder error")}
		output := &bytes.Buffer{}

		err := actions.SearchNotesContentWithOptions(&vault, &note, &uri, &fuzzyFinder, "test", defaultOptions(output))
		assert.Error(t, err)
	})

	t.Run("uri execution returns error", func(t *testing.T) {
		vault := mocks.MockVaultOperator{Name: "myVault"}
		uri := mocks.MockUriManager{ExecuteErr: errors.New("uri execution error")}
		note := mocks.MockNoteManager{}
		fuzzyFinder := mocks.MockFuzzyFinder{}
		output := &bytes.Buffer{}

		err := actions.SearchNotesContentWithOptions(&vault, &note, &uri, &fuzzyFinder, "test", defaultOptions(output))
		assert.Error(t, err)
	})

	t.Run("Successful content search with editor flag - single match", func(t *testing.T) {
		vault := mocks.MockVaultOperator{Name: "myVault"}
		uri := mocks.MockUriManager{}
		note := &CustomMockNoteForSingleMatch{}
		fuzzyFinder := mocks.MockFuzzyFinder{}
		output := &bytes.Buffer{}

		originalEditor := os.Getenv("EDITOR")
		defer os.Setenv("EDITOR", originalEditor)
		os.Setenv("EDITOR", "true")

		options := defaultOptions(output)
		options.UseEditor = true
		options.EditorFlagExplicit = true

		err := actions.SearchNotesContentWithOptions(&vault, note, &uri, &fuzzyFinder, "test", options)
		assert.NoError(t, err)
		assert.Equal(t, 0, uri.ExecuteCalls)
	})

	t.Run("Successful content search with editor flag - multiple matches", func(t *testing.T) {
		vault := mocks.MockVaultOperator{Name: "myVault"}
		uri := mocks.MockUriManager{}
		note := mocks.MockNoteManager{}
		fuzzyFinder := mocks.MockFuzzyFinder{SelectedIndex: 0}
		output := &bytes.Buffer{}

		originalEditor := os.Getenv("EDITOR")
		defer os.Setenv("EDITOR", originalEditor)
		os.Setenv("EDITOR", "true")

		options := defaultOptions(output)
		options.UseEditor = true
		options.EditorFlagExplicit = true

		err := actions.SearchNotesContentWithOptions(&vault, &note, &uri, &fuzzyFinder, "test", options)
		assert.NoError(t, err)
		assert.Equal(t, 0, uri.ExecuteCalls)
	})

	t.Run("Content search with editor flag fails when editor fails", func(t *testing.T) {
		vault := mocks.MockVaultOperator{Name: "myVault"}
		uri := mocks.MockUriManager{}
		note := &CustomMockNoteForSingleMatch{}
		fuzzyFinder := mocks.MockFuzzyFinder{}
		output := &bytes.Buffer{}

		originalEditor := os.Getenv("EDITOR")
		defer os.Setenv("EDITOR", originalEditor)
		os.Setenv("EDITOR", "false")

		options := defaultOptions(output)
		options.UseEditor = true
		options.EditorFlagExplicit = true

		err := actions.SearchNotesContentWithOptions(&vault, note, &uri, &fuzzyFinder, "test", options)
		assert.Error(t, err)
	})

	t.Run("No-interactive flag forces text output", func(t *testing.T) {
		vault := mocks.MockVaultOperator{Name: "myVault"}
		uri := mocks.MockUriManager{}
		note := mocks.MockNoteManager{}
		fuzzyFinder := mocks.MockFuzzyFinder{}
		output := &bytes.Buffer{}

		options := defaultOptions(output)
		options.NoInteractive = true

		err := actions.SearchNotesContentWithOptions(&vault, &note, &uri, &fuzzyFinder, "test", options)
		assert.NoError(t, err)
		assert.Equal(t, "note1.md:5: example match line\nnote2.md:10: another match\n", output.String())
		assert.Equal(t, 0, uri.ExecuteCalls)
	})

	t.Run("JSON format outputs structured data", func(t *testing.T) {
		vault := mocks.MockVaultOperator{Name: "myVault"}
		uri := mocks.MockUriManager{}
		note := mocks.MockNoteManager{}
		fuzzyFinder := mocks.MockFuzzyFinder{}
		output := &bytes.Buffer{}

		options := defaultOptions(output)
		options.Format = "json"

		err := actions.SearchNotesContentWithOptions(&vault, &note, &uri, &fuzzyFinder, "test", options)
		assert.NoError(t, err)
		assert.Equal(t, 0, uri.ExecuteCalls)

		var result []searchContentJSONMatch
		decodeErr := json.Unmarshal(output.Bytes(), &result)
		assert.NoError(t, decodeErr)
		assert.Len(t, result, 2)
		assert.Equal(t, "note1.md", result[0].File)
		assert.Equal(t, 5, result[0].Line)
		assert.Equal(t, "example match line", result[0].Content)
		assert.Equal(t, "content", result[0].MatchType)
	})

	t.Run("JSON format with no matches prints empty array", func(t *testing.T) {
		vault := mocks.MockVaultOperator{Name: "myVault"}
		uri := mocks.MockUriManager{}
		note := mocks.MockNoteManager{NoMatches: true}
		fuzzyFinder := mocks.MockFuzzyFinder{}
		output := &bytes.Buffer{}

		options := defaultOptions(output)
		options.Format = "json"

		err := actions.SearchNotesContentWithOptions(&vault, &note, &uri, &fuzzyFinder, "test", options)
		assert.NoError(t, err)
		assert.Equal(t, "[]\n", output.String())
	})

	t.Run("Non-interactive terminals auto-fallback to text output", func(t *testing.T) {
		vault := mocks.MockVaultOperator{Name: "myVault"}
		uri := mocks.MockUriManager{}
		note := mocks.MockNoteManager{}
		fuzzyFinder := mocks.MockFuzzyFinder{}
		output := &bytes.Buffer{}

		options := defaultOptions(output)
		options.InteractiveTerminal = false

		err := actions.SearchNotesContentWithOptions(&vault, &note, &uri, &fuzzyFinder, "test", options)
		assert.NoError(t, err)
		assert.Equal(t, 0, uri.ExecuteCalls)
		assert.Equal(t, "note1.md:5: example match line\nnote2.md:10: another match\n", output.String())
	})

	t.Run("Explicit editor flag returns error in non-interactive output mode", func(t *testing.T) {
		vault := mocks.MockVaultOperator{Name: "myVault"}
		uri := mocks.MockUriManager{}
		note := mocks.MockNoteManager{}
		fuzzyFinder := mocks.MockFuzzyFinder{}
		output := &bytes.Buffer{}

		options := defaultOptions(output)
		options.NoInteractive = true
		options.UseEditor = true
		options.EditorFlagExplicit = true

		err := actions.SearchNotesContentWithOptions(&vault, &note, &uri, &fuzzyFinder, "test", options)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "--editor cannot be used")
	})

	t.Run("Configured editor default is ignored in non-interactive terminals", func(t *testing.T) {
		vault := mocks.MockVaultOperator{Name: "myVault"}
		uri := mocks.MockUriManager{}
		note := mocks.MockNoteManager{}
		fuzzyFinder := mocks.MockFuzzyFinder{}
		output := &bytes.Buffer{}

		options := defaultOptions(output)
		options.InteractiveTerminal = false
		options.UseEditor = true
		options.EditorFlagExplicit = false

		err := actions.SearchNotesContentWithOptions(&vault, &note, &uri, &fuzzyFinder, "test", options)
		assert.NoError(t, err)
		assert.Equal(t, 0, uri.ExecuteCalls)
		assert.Equal(t, "note1.md:5: example match line\nnote2.md:10: another match\n", output.String())
	})

	t.Run("Invalid format returns error", func(t *testing.T) {
		vault := mocks.MockVaultOperator{Name: "myVault"}
		uri := mocks.MockUriManager{}
		note := mocks.MockNoteManager{}
		fuzzyFinder := mocks.MockFuzzyFinder{}
		output := &bytes.Buffer{}

		options := defaultOptions(output)
		options.Format = "yaml"

		err := actions.SearchNotesContentWithOptions(&vault, &note, &uri, &fuzzyFinder, "test", options)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid format")
	})
}
