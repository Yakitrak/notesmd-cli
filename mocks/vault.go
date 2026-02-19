package mocks

type MockVaultOperator struct {
	DefaultNameErr   error
	PathError        error
	Name             string
	PathValue        string
	OpenType         string
	OpenTypeErr      error
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
	if m.PathValue != "" {
		return m.PathValue, nil
	}
	return "path", nil
}

func (m *MockVaultOperator) DefaultOpenType() (string, error) {
	if m.OpenTypeErr != nil {
		return "", m.OpenTypeErr
	}
	if m.OpenType != "" {
		return m.OpenType, nil
	}
	return "obsidian", nil
}
