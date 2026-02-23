package cmd

import (
	"errors"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

type stubVaultManager struct {
	defaultName    string
	defaultNameErr error
	path           string
	pathErr        error
	openType       string
	openTypeErr    error
}

func (s *stubVaultManager) DefaultName() (string, error) {
	if s.defaultNameErr != nil {
		return "", s.defaultNameErr
	}
	if s.defaultName == "" {
		return "vault", nil
	}
	return s.defaultName, nil
}

func (s *stubVaultManager) SetDefaultName(name string) error {
	s.defaultName = name
	return nil
}

func (s *stubVaultManager) Path() (string, error) {
	if s.pathErr != nil {
		return "", s.pathErr
	}
	if s.path == "" {
		return "path", nil
	}
	return s.path, nil
}

func (s *stubVaultManager) DefaultOpenType() (string, error) {
	if s.openTypeErr != nil {
		return "", s.openTypeErr
	}
	if s.openType == "" {
		return "obsidian", nil
	}
	return s.openType, nil
}

func newSearchContentOptionsTestCmd() *cobra.Command {
	c := &cobra.Command{Use: "test"}
	c.Flags().BoolP("editor", "e", false, "")
	c.Flags().Bool("no-interactive", false, "")
	c.Flags().String("format", "text", "")
	return c
}

func TestSearchContentCommandFlagsWired(t *testing.T) {
	assert.NotNil(t, searchContentCmd.Flags().Lookup("no-interactive"))
	assert.NotNil(t, searchContentCmd.Flags().Lookup("format"))
	assert.NotNil(t, searchContentCmd.Flags().Lookup("editor"))
	assert.NotNil(t, searchContentCmd.Flags().Lookup("vault"))

	assert.Equal(t, "text", searchContentCmd.Flags().Lookup("format").DefValue)
	assert.Contains(t, searchContentCmd.Aliases, "sc")
}

func TestBuildSearchContentOptionsParsesExplicitFlags(t *testing.T) {
	c := newSearchContentOptionsTestCmd()
	err := c.ParseFlags([]string{"--editor", "--no-interactive", "--format", "json"})
	assert.NoError(t, err)

	vault := &stubVaultManager{openType: "obsidian"}
	options, err := buildSearchContentOptions(c, vault, false)
	assert.NoError(t, err)
	assert.True(t, options.UseEditor)
	assert.True(t, options.EditorFlagExplicit)
	assert.True(t, options.NoInteractive)
	assert.Equal(t, "json", options.Format)
	assert.False(t, options.InteractiveTerminal)
	assert.NotNil(t, options.Output)
}

func TestBuildSearchContentOptionsRespectsDefaultOpenType(t *testing.T) {
	c := newSearchContentOptionsTestCmd()
	err := c.ParseFlags([]string{})
	assert.NoError(t, err)

	vault := &stubVaultManager{openType: "editor"}
	options, err := buildSearchContentOptions(c, vault, true)
	assert.NoError(t, err)
	assert.True(t, options.UseEditor)
	assert.False(t, options.EditorFlagExplicit)
	assert.False(t, options.NoInteractive)
	assert.Equal(t, "text", options.Format)
	assert.True(t, options.InteractiveTerminal)
}

func TestBuildSearchContentOptionsDefaultOpenTypeErrorFallsBack(t *testing.T) {
	c := newSearchContentOptionsTestCmd()
	err := c.ParseFlags([]string{})
	assert.NoError(t, err)

	vault := &stubVaultManager{openTypeErr: errors.New("config error")}
	options, err := buildSearchContentOptions(c, vault, true)
	assert.NoError(t, err)
	assert.False(t, options.UseEditor)
	assert.False(t, options.EditorFlagExplicit)
}
