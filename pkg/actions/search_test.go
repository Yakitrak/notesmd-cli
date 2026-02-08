package actions_test

import (
	"errors"
	"testing"

	"github.com/Yakitrak/obsidian-cli/pkg/actions"
	"github.com/Yakitrak/obsidian-cli/pkg/obsidian"
)

// Simple mock implementations for testing
type testVault struct {
	name string
	path string
	err  error
}

func (v *testVault) DefaultName() (string, error) {
	if v.err != nil {
		return "", v.err
	}
	return v.name, nil
}

func (v *testVault) SetDefaultName(name string) error {
	v.name = name
	return nil
}

func (v *testVault) Path() (string, error) {
	if v.err != nil {
		return "", v.err
	}
	return v.path, nil
}

type testNote struct {
	notes    []string
	contents map[string]string
}

func (n *testNote) Delete(string) error                                   { return nil }
func (n *testNote) Move(string, string) error                            { return nil }
func (n *testNote) UpdateLinks(string, string, string) error             { return nil }
func (n *testNote) GetContents(_, notePath string) (string, error)       { return n.contents[notePath], nil }
func (n *testNote) SetContents(string, string, string) error             { return nil }
func (n *testNote) GetNotesList(string) ([]string, error)                { return n.notes, nil }
func (n *testNote) SearchNotesWithSnippets(string, string) ([]obsidian.NoteMatch, error) { return nil, nil }
func (n *testNote) FindBacklinks(string, string) ([]obsidian.NoteMatch, error)          { return nil, nil }

type testUri struct{}

func (u *testUri) Construct(urlType string, params map[string]string) string { return "" }
func (u *testUri) Execute(uri string) error                                  { return nil }

func TestSearch_UnknownFormat(t *testing.T) {
	vault := &testVault{name: "test", path: "/vault"}
	note := &testNote{}
	uri := &testUri{}

	err := actions.Search(vault, note, uri, "test", "invalid", "all", false, 0, false)
	if err == nil {
		t.Error("expected error for unknown format")
	}
	if err.Error() != "unknown format: invalid (use: names, locations, lines, blocks)" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestSearch_VaultError(t *testing.T) {
	vault := &testVault{err: errors.New("vault error")}
	note := &testNote{}
	uri := &testUri{}

	err := actions.Search(vault, note, uri, "test", "names", "all", false, 0, false)
	if err == nil {
		t.Error("expected error when vault fails")
	}
}

func TestSearch_EmptyVault(t *testing.T) {
	vault := &testVault{name: "test", path: "/vault"}
	note := &testNote{notes: []string{}}
	uri := &testUri{}

	err := actions.Search(vault, note, uri, "test", "names", "all", false, 0, false)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSearch_NamesBasic(t *testing.T) {
	vault := &testVault{name: "test", path: "/vault"}
	note := &testNote{
		notes: []string{"note1.md", "note2.md", "note3.md"},
		contents: map[string]string{
			"note1.md": "#programming content",
			"note2.md": "no tags",
			"note3.md": "#programming/philosophy subtag",
		},
	}
	uri := &testUri{}

	// Without subtags - only note1
	err := actions.Search(vault, note, uri, "programming", "names", "all", false, 0, false)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// With subtags - note1 and note3
	err = actions.Search(vault, note, uri, "programming", "names", "all", true, 0, false)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSearch_NamesWithLocationFilter(t *testing.T) {
	vault := &testVault{name: "test", path: "/vault"}
	note := &testNote{
		notes: []string{"note1.md", "note2.md"},
		contents: map[string]string{
			"note1.md": "---\ntags: [programming]\n---\ncontent",
			"note2.md": "#programming inline tag",
		},
	}
	uri := &testUri{}

	// Frontmatter filter
	err := actions.Search(vault, note, uri, "programming", "names", "frontmatter", false, 0, false)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Inline filter
	err = actions.Search(vault, note, uri, "programming", "names", "inline", false, 0, false)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}