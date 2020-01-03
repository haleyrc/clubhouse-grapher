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
}

func (n Node) ToString() string {
	return fmt.Sprintf(
		"node%d[label=\"%s\",fillcolor=\"%s\",color=\"%s\"];",
		n.ID,
		n.Label,
		n.FillColor,
		n.Color,
	)
}
