package app

import "context"

type Commands []Command

type Command struct {
	Keyword     string
	UsageFlags  string
	Usage       string
	Description string
	Fn          func(ctx context.Context) (err error)
}

func (c Commands) Len() int {
	return len(c)
}

func (c Commands) Less(i, j int) bool {
	return c[i].Keyword < c[j].Keyword
}

func (c Commands) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
