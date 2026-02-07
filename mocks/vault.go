package mocks

import "github.com/Yakitrak/obsidian-cli/pkg/obsidian"

type MockVaultOperator struct {
	DefaultNameErr error
	PathError      error
	Name           string
	ID             string
}

func (m *MockVaultOperator) DefaultName() (string, error) {
	if m.DefaultNameErr != nil {
		return "", m.DefaultNameErr
	}
	return m.Name, nil
}

func (m *MockVaultOperator) SetDefaultName(_ string) error {
	if m.DefaultNameErr != nil {
		return m.DefaultNameErr
	}
	return nil
}

func (m *MockVaultOperator) Path() (string, error) {
	if m.PathError != nil {
		return "", m.PathError
	}
	return "path", nil
}

func (m *MockVaultOperator) GetVaultInfo() (*obsidian.VaultInfo, error) {
	if m.DefaultNameErr != nil {
		return nil, m.DefaultNameErr
	}
	return &obsidian.VaultInfo{
		ID:   m.ID,
		Name: m.Name,
		Path: "path",
	}, nil
}
