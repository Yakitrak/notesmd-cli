package validate

import "encoding/json"

// Severity represents the severity level of a validation finding.
type Severity string

const (
	SeverityError   Severity = "error"
	SeverityWarning Severity = "warning"
	SeverityInfo    Severity = "info"
)

// SeverityOrder returns the numeric order for severity (lower = more severe).
func SeverityOrder(s Severity) int {
	switch s {
	case SeverityError:
		return 0
	case SeverityWarning:
		return 1
	case SeverityInfo:
		return 2
	default:
		return 3
	}
}

// ValidationResult represents a single validation finding.
type ValidationResult struct {
	File       string   `json:"file"`
	Severity   Severity `json:"severity"`
	Rule       string   `json:"rule"`
	Message    string   `json:"message"`
	Suggestion string   `json:"suggestion,omitempty"`
	Line       int      `json:"line,omitempty"`
}

// SeverityCounts holds the count of findings by severity.
type SeverityCounts struct {
	Errors   int `json:"errors"`
	Warnings int `json:"warnings"`
	Infos    int `json:"infos"`
}

// ValidationSummary is the top-level result of a vault validation.
type ValidationSummary struct {
	Vault      string             `json:"vault"`
	TotalNotes int                `json:"totalNotes"`
	Results    []ValidationResult `json:"results"`
	Summary    SeverityCounts     `json:"summary"`
}

// HasErrors returns true if there are any error-severity findings.
func (s *ValidationSummary) HasErrors() bool {
	return s.Summary.Errors > 0
}

// ToJSON marshals the summary to indented JSON.
func (s *ValidationSummary) ToJSON() (string, error) {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ComputeSummary recalculates the severity counts from the results.
func (s *ValidationSummary) ComputeSummary() {
	s.Summary = SeverityCounts{}
	for _, r := range s.Results {
		switch r.Severity {
		case SeverityError:
			s.Summary.Errors++
		case SeverityWarning:
			s.Summary.Warnings++
		case SeverityInfo:
			s.Summary.Infos++
		}
	}
}
