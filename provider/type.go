package provider

type Type int

//go:generate enumer -type=Type -transform=snake -output=type_string.go

const (
	Default Type = iota
	AWS
)
