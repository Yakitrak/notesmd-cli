package validate

// Config holds validation configuration.
type Config struct {
	// ExcludeDirs are directory names to skip during vault traversal.
	ExcludeDirs []string

	// MinSeverity is the minimum severity level to report.
	MinSeverity Severity
}

// DefaultConfig returns the default validation configuration.
func DefaultConfig() Config {
	return Config{
		ExcludeDirs: DefaultExcludeDirs(),
		MinSeverity: SeverityWarning,
	}
}

// DefaultExcludeDirs returns the default directories to exclude from validation.
func DefaultExcludeDirs() []string {
	return []string{
		".obsidian",
		".trash",
		".git",
		"node_modules",
	}
}

// ShouldReport returns true if the given severity meets the minimum threshold.
func (c *Config) ShouldReport(s Severity) bool {
	return SeverityOrder(s) <= SeverityOrder(c.MinSeverity)
}
