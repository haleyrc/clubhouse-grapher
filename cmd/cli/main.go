package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/haleyrc/clubhouse"
	"github.com/haleyrc/clubhouse/graph"
)

func main() {
	ctx := context.Background()

	token := os.Getenv("CLUBHOUSE_API_TOKEN")
	if token == "" {
		panic(errors.New("token is required"))
	}

	projects := []string{"Backend", "Frontend", "Back Office"}
	if len(os.Args) > 1 {
		projects = strings.Split(os.Args[1], ",")
	}

	c := clubhouse.NewClient(token)
	workspace, err := c.GetWorkspace(ctx, clubhouse.GetWorkspaceParams{
		OnlyProjects: projects,
	})
	if err != nil {
		panic(err)
	}

	g := graph.NewGraph("All Projects", workspace)
	fmt.Println(g.ToString())
}
