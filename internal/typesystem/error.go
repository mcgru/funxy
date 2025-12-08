package typesystem

import "fmt"

// SymbolNotFoundError indicates a symbol was not found in scope.
type SymbolNotFoundError struct {
	Name string
}

func (e *SymbolNotFoundError) Error() string {
	return fmt.Sprintf("symbol not found: %s", e.Name)
}

func NewSymbolNotFoundError(name string) error {
	return &SymbolNotFoundError{Name: name}
}

// UnificationError represents a type unification failure.
type UnificationError struct {
	Expected Type
	Actual   Type
	Context  string // e.g., "record field 'x'", "function parameter 1"
	Cause    error  // nested error for detailed messages
}

func (e *UnificationError) Error() string {
	if e.Context != "" {
		if e.Cause != nil {
			return fmt.Sprintf("%s: %v", e.Context, e.Cause)
		}
		return fmt.Sprintf("%s: cannot unify %s with %s", e.Context, e.Expected, e.Actual)
	}
	if e.Cause != nil {
		return e.Cause.Error()
	}
	return fmt.Sprintf("cannot unify %s with %s", e.Expected, e.Actual)
}

func (e *UnificationError) Unwrap() error {
	return e.Cause
}

// Unification error constructors

func errUnify(t1, t2 Type) error {
	return &UnificationError{Expected: t1, Actual: t2}
}

func errUnifyMsg(t1, t2 Type, msg string) error {
	return &UnificationError{Expected: t1, Actual: t2, Context: msg}
}

func errUnifyContext(context string, cause error) error {
	return &UnificationError{Context: context, Cause: cause}
}

func errMismatch(msg string) error {
	return &UnificationError{Context: msg}
}
