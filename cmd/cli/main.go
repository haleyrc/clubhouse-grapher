package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/haleyrc/clubhouse"
	"github.com/haleyrc/clubhouse/graph"
)

func main() {
	ctx := context.Background()

	token := os.Getenv("CLUBHOUSE_API_TOKEN")
	if token == "" {
		panic(errors.New("token is required"))
	}

	c := clubhouse.NewClient(token)
	workspace, err := c.GetWorkspace(ctx, clubhouse.GetWorkspaceParams{
		OnlyProjects: []string{"Backend", "Frontend", "Back Office"},
	})
	if err != nil {
		panic(err)
	}

	g := graph.NewGraph("All Projects", workspace)
	fmt.Println(g)

	// for _, sg := range g.Subgraphs {
	// 	fn := fmt.Sprintf("%s.dot", sg.Name)
	// 	f, err := os.Create(fn)
	// 	if err != nil {
	// 		log.Printf("error creating subgraph: %s: %v\n", fn, err)
	// 		continue
	// 	}
	// 	ng := graph.AsGraph(sg)
	// 	if ng.Name == "Frontend" {
	// 		for _, rank := range ng.Nodes {
	// 			for _, node := range rank {
	// 				fmt.Fprintf(os.Stderr, "%-3d: %s\n", node.Rank, node.Story.Name)
	// 			}
	// 		}
	// 	}
	// 	fmt.Fprintf(f, ng.ToString())
	// 	f.Close()
	// }
}
