//go:build tools
// +build tools

package tools

import (
	_ "github.com/dmarkham/enumer"
	_ "github.com/markbates/pkger/cmd/pkger"
	_ "golang.org/x/lint/golint"
	_ "golang.org/x/tools/cmd/goimports"
)
