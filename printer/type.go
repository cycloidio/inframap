package printer

// Type defines the type of the Printer iota
type Type int

//go:generate enumer -type=Type -transform=lower -output=type_string.go

// List of all Types
const (
	DOT Type = iota
)
