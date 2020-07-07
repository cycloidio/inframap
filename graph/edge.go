package graph

import (
	"github.com/cycloidio/infraview/errcode"
)

// Edge defines the standard format of an Edge
type Edge struct {
	ID string

	Canonicals []string
	// mCanonicals is used to know which Canonicals
	// are already on the Canonicals slice so we
	// do not repeat them
	mCanonicals map[string]struct{}

	Source string
	Target string
}

// Replace replaces the src (on Target or Source) for rep values
func (e *Edge) Replace(src, rep string) error {
	if e.Source == src {
		e.Source = rep
		return nil
	} else if e.Target == src {
		e.Target = rep
		return nil
	}

	return errcode.ErrGraphEdgeReplace
}

// AddCanonicals adds the cans to the internal list, if
// one is repeated it'll be ignored
func (e *Edge) AddCanonicals(cans ...string) {
	if e.mCanonicals == nil {
		e.mCanonicals = make(map[string]struct{})
	}

	for _, c := range cans {
		if _, ok := e.mCanonicals[c]; !ok {
			e.mCanonicals[c] = struct{}{}
			e.Canonicals = append(e.Canonicals, c)
		}
	}
}
