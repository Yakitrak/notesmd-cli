package actions

import (
	"fmt"
	"github.com/Yakitrak/notesmd-cli/pkg/obsidian"
	"path/filepath"
)

func SearchNotes(vault obsidian.VaultManager, note obsidian.NoteManager, uri obsidian.UriManager, fuzzyFinder obsidian.FuzzyFinderManager, useEditor bool) error {
	vaultName, err := vault.DefaultName()
	if err != nil {
		return err
	}

	vaultPath, err := vault.Path()
	if err != nil {
		return err
	}

	notes, err := note.GetNotesList(vaultPath)
	if err != nil {
		return err
	}

	index, err := fuzzyFinder.Find(notes, func(i int) string {
		return notes[i]
	})

	if err != nil {
		return err
	}

	if useEditor {
		fmt.Printf("Opening note: %s\n", notes[index])
		filePath := filepath.Join(vaultPath, notes[index])
		return obsidian.OpenInEditor(filePath)
	}

	obsidianUri := uri.Construct(ObsOpenUrl, map[string]string{
		"file":  notes[index],
		"vault": vaultName,
	})

	err = uri.Execute(obsidianUri)
	if err != nil {
		return err
	}

	return nil
}
