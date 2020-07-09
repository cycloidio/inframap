package printer

type Type int

//go:generate enumer -type=Type -transform=lower -output=type_string.go

const (
	DOT Type = iota
)
