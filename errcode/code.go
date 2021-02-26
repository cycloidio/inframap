package errcode

import "errors"

// List of all the errors
var (
	ErrGraphEdgeReplace              = errors.New("graph edge replace source not found")
	ErrGraphRequiredEdgeID           = errors.New("graph edge ID is required")
	ErrGraphRequiredEdgeTarget       = errors.New("graph edge target is required")
	ErrGraphRequiredEdgeSource       = errors.New("graph edge source is required")
	ErrGraphNotFoundEdgeTarget       = errors.New("graph edge target not found")
	ErrGraphNotFoundEdgeSource       = errors.New("graph edge source not found")
	ErrGraphAlreadyExistsEdge        = errors.New("graph edge already exists")
	ErrGraphAlreadyExistsEdgeID      = errors.New("graph edge ID already exists")
	ErrGraphRequiredNodeCanonical    = errors.New("graph node canonical is required")
	ErrGraphRequiredNodeID           = errors.New("graph node ID is required")
	ErrGraphAlreadyExistsNode        = errors.New("graph node already exists")
	ErrGraphAlreadyExistsNodeID      = errors.New("graph node ID already exists")
	ErrGraphNotFoundNode             = errors.New("graph node not found")
	ErrGraphRequiredEdgeBetweenNodes = errors.New("graph requires edge between nodes")

	ErrProviderNotFoundResource   = errors.New("provider resource not found")
	ErrProviderNotFoundDataSource = errors.New("provider data source not found")
	ErrProviderNotFound           = errors.New("provider not found")

	ErrInvalidTFStateFile    = errors.New("invalid Terraform State file")
	ErrInvalidTFStateVersion = errors.New("invalid Terraform State version, we only support version 3 and 4")

	ErrPrinterNotFound = errors.New("printer not found")

	ErrGenerateFromJSON = errors.New("we do not support JSON HCL")
)
