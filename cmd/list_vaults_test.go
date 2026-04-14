package cmd

import (
	"bytes"
	"testing"

	"github.com/Yakitrak/notesmd-cli/pkg/obsidian"
	"github.com/stretchr/testify/assert"
)

func TestFormatVaultsTable(t *testing.T) {
	t.Run("Aligns columns with varying name lengths", func(t *testing.T) {
		vaults := []obsidian.VaultInfo{
			{Name: "Notes", Path: "/home/user/Notes"},
			{Name: "LongVaultName", Path: "/home/user/LongVaultName"},
			{Name: "Work", Path: "/home/user/Work"},
		}

		var buf bytes.Buffer
		formatVaultsTable(&buf, vaults, "")
		output := buf.String()

		// All path columns should start at the same position
		lines := bytes.Split(bytes.TrimSpace([]byte(output)), []byte("\n"))
		assert.Len(t, lines, 3)

		// Each line should contain both name and path
		assert.Contains(t, output, "Notes")
		assert.Contains(t, output, "/home/user/Notes")
		assert.Contains(t, output, "LongVaultName")
		assert.Contains(t, output, "/home/user/LongVaultName")

		// Paths should be aligned — find the byte offset of each path
		// With tabwriter, the path column should start at the same position
		pathOffsets := make([]int, len(lines))
		for i, line := range lines {
			pathOffsets[i] = bytes.Index(line, []byte("/home"))
		}
		assert.Equal(t, pathOffsets[0], pathOffsets[1], "path columns should be aligned")
		assert.Equal(t, pathOffsets[1], pathOffsets[2], "path columns should be aligned")
	})

	t.Run("Single vault produces output", func(t *testing.T) {
		vaults := []obsidian.VaultInfo{
			{Name: "MyVault", Path: "/tmp/MyVault"},
		}

		var buf bytes.Buffer
		formatVaultsTable(&buf, vaults, "")
		output := buf.String()

		assert.Contains(t, output, "MyVault")
		assert.Contains(t, output, "/tmp/MyVault")
	})

	t.Run("Marks default vault", func(t *testing.T) {
		vaults := []obsidian.VaultInfo{
			{Name: "Notes", Path: "/home/user/Notes"},
			{Name: "Work", Path: "/home/user/Work"},
		}

		var buf bytes.Buffer
		formatVaultsTable(&buf, vaults, "Work")
		output := buf.String()

		assert.Contains(t, output, "Work")
		assert.Contains(t, output, "(default)")
		// Notes line should not have (default)
		lines := bytes.Split(bytes.TrimSpace([]byte(output)), []byte("\n"))
		assert.Len(t, lines, 2)
		assert.NotContains(t, string(lines[0]), "(default)")
		assert.Contains(t, string(lines[1]), "(default)")
	})

	t.Run("Empty vault list produces no output", func(t *testing.T) {
		var buf bytes.Buffer
		formatVaultsTable(&buf, []obsidian.VaultInfo{}, "")

		assert.Empty(t, buf.String())
	})
}
