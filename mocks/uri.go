package mocks

type MockUriManager struct {
	ConstructedURI string
	LastBase       string
	LastParams     map[string]string
	ExecuteErr     error
	ExecuteCalls   int
}

func (m *MockUriManager) Construct(base string, params map[string]string) string {
	m.LastBase = base
	m.LastParams = params
	return m.ConstructedURI
}

func (m *MockUriManager) Execute(uri string) error {
	m.ExecuteCalls++
	return m.ExecuteErr
}
