package factory

import (
	"fmt"

	"github.com/cycloidio/inframap/errcode"
	"github.com/cycloidio/inframap/printer"
	"github.com/cycloidio/inframap/printer/dot"
)

var (
	printers = map[printer.Type]printer.Printer{
		printer.DOT: dot.Dot{},
	}
)

// Get returns the specific Printer for t
func Get(t string) (printer.Printer, error) {
	ty, err := printer.TypeString(t)
	if err != nil {
		return nil, err
	}

	p, ok := printers[ty]
	if !ok {
		return nil, fmt.Errorf("no printer defined for %s: %w", ty, errcode.ErrPrinterNotFound)
	}
	return p, nil
}
