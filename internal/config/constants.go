package config

const SourceFileExt = ".lang"

// SourceFileExtensions are all recognized source file extensions
var SourceFileExtensions = []string{".lang", ".funxy", ".fx"}

// Built-in trait and method names
const (
	IterTraitName  = "Iter"
	IterMethodName = "iter"
)

// Built-in function names
const (
	PrintFuncName    = "print"
	WriteFuncName    = "write"
	PanicFuncName    = "panic"
	DebugFuncName    = "debug"
	TraceFuncName    = "trace"
	LenFuncName      = "len"
	LenBytesFuncName = "lenBytes"
	TypeOfFuncName   = "typeOf"
	GetTypeFuncName  = "getType"
	DefaultFuncName  = "default"
	ShowFuncName     = "show"
	ReadFuncName     = "read"
	IdFuncName       = "id"
	ConstFuncName    = "const"
)

// Built-in type names
const (
	ListTypeName   = "List"
	MapTypeName    = "Map"
	BytesTypeName  = "Bytes"
	BitsTypeName   = "Bits"
	OptionTypeName = "Option"
	ResultTypeName = "Result"
	SomeCtorName   = "Some"
	ZeroCtorName   = "Zero"
	OkCtorName     = "Ok"
	FailCtorName   = "Fail"
)
