package provider

type Type int

//go:generate enumer -type=Type -transform=lower -output=type_string.go

const (
	Default Type = iota
	AWS
	FlexibleEngine
	OpenStack
)
