package graph

import (
	"fmt"
	"sort"

	"github.com/haleyrc/clubhouse"
)

const GraphHeader = `digraph %q {
    rankdir=LR;
    splines=polyline;
    overlap=false;
    node[shape="box",style="filled",fillcolor="white",penwidth="3"];
    label="%s";
`

func NewGraph(name string, workspace *clubhouse.Workspace) *Graph {
	g := Graph{
		Name:          name,
		Nodes:         []*Node{},
		Ranks:         []int{},
		Projects:      []string{},
		ProjectColors: make(map[string]string),
	}

	allRanks := make(map[int]bool)
	allProjects := make(map[string]bool)
	for _, story := range workspace.Stories {
		n := Node{
			ID:        story.ID,
			Label:     story.Name,
			Project:   story.Project.Name,
			Color:     story.Project.Color,
			FillColor: colorFor(story),
			Rank:      rankFor(story),
		}
		for _, blocker := range story.Blockers {
			n.Blockers = append(n.Blockers, blocker.ID)
		}
		g.Nodes = append(g.Nodes, &n)
		allRanks[n.Rank] = true
		allProjects[story.Project.Name] = true
		g.ProjectColors[story.Project.Name] = story.Project.Color
	}

	for rank := range allRanks {
		g.Ranks = append(g.Ranks, rank)
	}
	sort.Ints(g.Ranks)

	for project := range allProjects {
		g.Projects = append(g.Projects, project)
	}
	sort.Strings(g.Projects)

	return &g
}

type Graph struct {
	Name          string
	Nodes         []*Node
	Ranks         []int
	Projects      []string
	ProjectColors map[string]string
}

func (g Graph) String() string {
	cluster := -1
	w := NewWriter()

	header := fmt.Sprintf(GraphHeader, g.Name, g.Name)
	w.WriteString(0, header)
	for _, project := range g.Projects {
		cluster++
		w.WriteStringf(1, "subgraph cluster_%d {\n", cluster)
		w.WriteString(2, "rank=same;\n")
		w.WriteStringf(2, "label=%q;\n", project)
		w.WriteStringf(2, "labeljust=\"l\";\n")
		w.WriteStringf(2, "style=\"filled\";\n")
		w.WriteStringf(2, "fillcolor=%q;\n", g.ProjectColors[project])
		for rank := range g.Ranks {
			cluster++
			w.WriteStringf(3, "subgraph cluster_%d {\n", cluster)
			w.WriteString(4, "rank=same;\n")
			w.WriteStringf(4, "label=\"\";\n")
			w.WriteStringf(4, "style=invis;\n")
			for _, node := range g.Nodes {
				if node.Rank != rank || node.Project != project {
					continue
				}
				w.WriteLn(4, node)
			}
			w.WriteString(3, "}\n")
		}
		w.WriteString(1, "}\n")
	}
	for _, node := range g.Nodes {
		for _, blocker := range node.Blockers {
			w.WriteStringf(1, "node%d->node%d;\n", blocker, node.ID)
		}
	}
	w.WriteString(0, "}")

	return w.String()
}

func rankFor(story *clubhouse.Story) int {
	predRanks := []int{0} // Seed with a 0 so we always get a result
	for _, blocker := range story.Blockers {
		predRanks = append(predRanks, rankFor(blocker)+1)
	}
	return highestRank(predRanks)
}

func highestRank(ranks []int) int {
	if len(ranks) < 1 {
		return 0
	}
	highest := ranks[0]
	for _, rank := range ranks {
		if rank > highest {
			highest = rank
		}
	}
	return highest
}

func colorFor(s *clubhouse.Story) string {
	color := "white"
	if s.Blocked {
		color = "red"
	}
	if s.Completed {
		color = "green"
	}
	return color
}
