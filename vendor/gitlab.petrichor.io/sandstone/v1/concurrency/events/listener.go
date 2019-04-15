package events

import (
	"sort"
)

type listener struct {
	binding  *binding
	triggers []string
	priority int
	all bool
}

func (l *listener) On(triggers ...string) *listener {
	l.triggers = append(l.triggers, triggers...)

	l.sort()

	return l
}

func (l *listener) All() *listener {
	l.all = true

	l.sort()

	return l
}

func (l *listener) SetPriority(p int) *listener {
	l.priority = p

	l.sort()

	return l
}

func (l *listener) sort() {
	sort.SliceStable(l.binding.d.ls, func(i, j int) bool {
		return l.binding.d.ls[i].priority > l.binding.d.ls[j].priority
	})
}
