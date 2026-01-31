package mocks

type MockUriManager struct {
	ConstructedURI string
	LastBase       string
	LastParams     map[string]string
	ExecuteErr     error
}

func (m *MockUriManager) Construct(base string, params map[string]string) string {
	m.LastBase = base
	m.LastParams = params
	return m.ConstructedURI
}

func (m *MockUriManager) Execute(uri string) error {
	return m.ExecuteErr
}
