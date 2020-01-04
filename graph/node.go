package graph

import (
	"fmt"
)

type Node struct {
	ID        int64
	Label     string
	Project   string
	Color     string
	FillColor string
	Rank      int
	Blockers  []int64
	Blocks    []int64
	Completed bool
}

func (n Node) ToString() string {
	return fmt.Sprintf(
		"node%d[label=\"%s (%d)\",fillcolor=%q,color=%q];",
		n.ID,
		n.Label,
		n.Rank,
		n.FillColor,
		n.Color,
	)
}
