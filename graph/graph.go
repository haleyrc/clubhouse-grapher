package graph

import (
	"fmt"
	"sort"
	"strings"

	"github.com/haleyrc/clubhouse"
)

const GraphHeader = `digraph %q {
	layout="dot";
	rankdir=LR;
	ranksep=2;
    // splines=ortho;
    // overlap=false;
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
		rank := rankFor(story)
		if story.Completed {
			rank = -100 + rank
		}
		n := Node{
			ID:        story.ID,
			Label:     story.Name,
			Project:   story.Project.Name,
			Color:     story.Project.Color,
			FillColor: colorFor(story),
			Rank:      rank,
			Completed: story.Completed,
		}
		for _, blocker := range story.Blockers {
			if blocker == nil {
				continue
			}
			n.Blockers = append(n.Blockers, blocker.ID)
		}
		for _, blocked := range story.Blocks {
			if blocked == nil {
				continue
			}
			n.Blocks = append(n.Blocks, blocked.ID)
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
	normed := normalize(g.Ranks)
	for _, node := range g.Nodes {
		node.Rank = normed[node.Rank]
	}

	for project := range allProjects {
		g.Projects = append(g.Projects, project)
	}
	sort.Strings(g.Projects)

	return &g
}

func normalize(vals []int) map[int]int {
	normed := make(map[int]int)
	for i, val := range vals {
		normed[val] = i
	}
	return normed
}

type Graph struct {
	Name          string
	Nodes         []*Node
	Ranks         []int
	Projects      []string
	ProjectColors map[string]string
}

func (g Graph) ToString() string {
	w := NewWriter()

	header := fmt.Sprintf(GraphHeader, g.Name, g.Name)
	w.WriteString(0, header)
	for _, project := range g.Projects {
		for _, node := range g.Nodes {
			if node.Project != project {
				continue
			}
			w.WriteLn(1, node)
		}
	}

	for _, node := range g.Nodes {
		if node.Blocks == nil || len(node.Blocks) == 0 {
			continue
		}
		var allBlocks []string
		for _, blocks := range node.Blocks {
			s := fmt.Sprintf("node%d", blocks)
			allBlocks = append(allBlocks, s)
		}
		if len(allBlocks) < 1 {
			continue
		}
		blocks := strings.Join(allBlocks, " ")
		w.WriteStringf(1, "node%d -> { %s };\n", node.ID, blocks)
	}

	nodesInRanks := make(map[int][]string)
	for _, node := range g.Nodes {
		nodesInRanks[node.Rank] = append(
			nodesInRanks[node.Rank],
			fmt.Sprintf("node%d", node.ID),
		)
	}
	for rank, nodes := range nodesInRanks {
		w.WriteStringf(1,
			"subgraph cluster_%d { style=invis; rank=same; %s }\n",
			rank,
			strings.Join(nodes, ", "),
		)
	}

	w.WriteString(0, "}")

	return w.String()
}

func rankFor(story *clubhouse.Story) int {
	if story == nil || story.Blockers == nil {
		return 0
	}
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
