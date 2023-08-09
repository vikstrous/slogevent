package slogevent

import "log/slog"

// copied from https://github.com/jba/slog/blob/main/withsupport/withsupport.go

// groupOrAttrs holds either a group name or a list of slog.Attrs.
type groupOrAttrs struct {
	Group string      // group name if non-empty
	Attrs []slog.Attr // attrs if non-empty
	Next  *groupOrAttrs
}

// WithGroup returns a GroupOrAttrs that includes the given group.
func (g *groupOrAttrs) WithGroup(name string) *groupOrAttrs {
	if name == "" {
		return g
	}
	return &groupOrAttrs{
		Group: name,
		Next:  g,
	}
}

// WithAttrs returns a GroupOrAttrs that includes the given attrs.
func (g *groupOrAttrs) WithAttrs(attrs []slog.Attr) *groupOrAttrs {
	if len(attrs) == 0 {
		return g
	}
	return &groupOrAttrs{
		Attrs: attrs,
		Next:  g,
	}
}

// Apply calls f on each Attr in g. The first argument to f is the list
// of groups that precede the Attr.
// Apply returns the complete list of groups.
func (g *groupOrAttrs) Apply(f func(groups []string, a slog.Attr)) []string {
	var groups []string

	var rec func(*groupOrAttrs)
	rec = func(g *groupOrAttrs) {
		if g == nil {
			return
		}
		rec(g.Next)
		if g.Group != "" {
			groups = append(groups, g.Group)
		} else {
			for _, a := range g.Attrs {
				f(groups, a)
			}
		}
	}
	rec(g)

	return groups
}

// Collect returns a slice of the GroupOrAttrs in reverse order.
func (g *groupOrAttrs) Collect() []*groupOrAttrs {
	n := 0
	for ga := g; ga != nil; ga = ga.Next {
		n++
	}
	res := make([]*groupOrAttrs, n)
	i := 0
	for ga := g; ga != nil; ga = ga.Next {
		res[len(res)-i-1] = ga
		i++
	}
	return res
}
