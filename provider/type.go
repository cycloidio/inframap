package provider

// Type defines the type of the Provider
type Type int

//go:generate ./../bin/enumer -type=Type -transform=lower -output=type_string.go

// List of all the Providers supported
const (
	Raw Type = iota
	AWS
	FlexibleEngine
	OpenStack
	Google
)
