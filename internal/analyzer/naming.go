package analyzer

import (
	"github.com/funvibe/funxy/internal/diagnostics"
	"github.com/funvibe/funxy/internal/token"
	"unicode"
)

// isValueName checks if a name follows value naming convention (starts with lowercase or _)
// Also allows operator methods like (==), (+), etc.
func isValueName(name string) bool {
	if len(name) == 0 {
		return false
	}
	first := rune(name[0])
	// Allow operators in parens: (==), (+), etc.
	if first == '(' {
		return true
	}
	return unicode.IsLower(first) || first == '_'
}

// isTypeName checks if a name follows type naming convention (starts with uppercase)
func isTypeName(name string) bool {
	if len(name) == 0 {
		return false
	}
	return unicode.IsUpper(rune(name[0]))
}

// checkValueName validates that a name follows value naming convention
// and appends an error if not
func checkValueName(name string, tok token.Token, errors *[]*diagnostics.DiagnosticError) bool {
	if !isValueName(name) {
		*errors = append(*errors, diagnostics.NewAnalyzerError(
			diagnostics.ErrA008, tok,
			"value '"+name+"' must start with lowercase letter or underscore",
		))
		return false
	}
	return true
}

// checkTypeName validates that a name follows type naming convention
// and appends an error if not
func checkTypeName(name string, tok token.Token, errors *[]*diagnostics.DiagnosticError) bool {
	if !isTypeName(name) {
		*errors = append(*errors, diagnostics.NewAnalyzerError(
			diagnostics.ErrA008, tok,
			"type '"+name+"' must start with uppercase letter",
		))
		return false
	}
	return true
}
