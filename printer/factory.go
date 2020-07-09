package printer

import (
	"fmt"

	"github.com/cycloidio/infraview/errcode"
	"github.com/cycloidio/infraview/printer/dot"
)

var (
	printers = map[Type]Printer{
		DOT: dot.Dot{},
	}
)

func Get(t string) (Printer, error) {
	ty, err := TypeString(t)
	if err != nil {
		return nil, err
	}

	p, ok := printers[ty]
	if !ok {
		return nil, fmt.Errorf("no printer defined for %s: %w", ty, errcode.ErrPrinterNotFound)
	}
	return p, nil
}
