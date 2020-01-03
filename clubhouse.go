package clubhouse

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"
)

const BaseURL = "https://api.clubhouse.io/api/v3"

type Story struct {
	AppURL              string    `json:"app_url"`
	Archived            bool      `json:"archived"`
	Blocked             bool      `json:"blocked"`
	Blocker             bool      `json:"blocker"`
	CommentIDs          []int64   `json:"comment_ids"`
	Completed           bool      `json:"completed"`
	CompleteAt          time.Time `json:"completed_at"`
	CompletedAtOverride time.Time `json:"completed_at_override"`
	CreatedAt           time.Time `json:"created_at"`
	Deadline            time.Time `json:"deadline"`
	EntityType          string    `json:"entity_type"`
	EpicID              int64     `json:"epic_id"`
	Estimate            int64     `json:"estimate"`
	ExternalID          string    `json:"external_id"`
	ExternalTickets     []struct {
		EpicIDs     []int64 `json:"epic_ids"`
		ExternalID  string  `json:"external_id"`
		ExternalURL string  `json:"external_url"`
		ID          string  `json:"id"`
		StoryIDs    []int64 `json:"story_ids"`
	} `json:"external_tickets"`
	FileIDs         []int64  `json:"file_ids"`
	FollowerIDs     []string `json:"follower_ids"`
	GroupMentionIDs []string `json:"group_mention_ids"`
	ID              int64    `json:"id"`
	IterationID     int64    `json:"iteration_id"`
	Labels          []struct {
		Archived    bool   `json:"archived"`
		Color       string `json:"color"`
		CreatedAt   string `json:"created_at"`
		Description string `json:"description"`
		EntityType  string `json:"entity_type"`
		ExternalID  string `json:"external_id"`
		ID          int64  `json:"id"`
		Name        string `json:"name"`
		Stats       struct {
			NumEpics              int64 `json:"num_epics"`
			NumPointsCompleted    int64 `json:"num_points_completed"`
			NumPointsInProgress   int64 `json:"num_points_in_progress"`
			NumPointsTotal        int64 `json:"num_points_total"`
			NumStoriesCompleted   int64 `json:"num_stories_completed"`
			NumStoriesInProgress  int64 `json:"num_stories_in_progress"`
			NumStoriesTotal       int64 `json:"num_stories_total"`
			NumStoriesUnestimated int64 `json:"num_stories_unestimated"`
		} `json:"stats"`
		UpdatedAt time.Time `json:"updated_at"`
	} `json:"labels"`
	LinkedFileIDs        []int64   `json:"linked_file_ids"`
	MemberMentionIDs     []string  `json:"member_mention_ids"`
	MentionIDs           []string  `json:"mention_ids"`
	MovedAt              time.Time `json:"moved_at"`
	Name                 string    `json:"name"`
	OwnerIDs             []string  `json:"owner_ids"`
	Position             int64     `json:"position"`
	PreviousIterationIDs []int64   `json:"previous_iteration_ids"`
	ProjectID            int64     `json:"project_id"`
	RequestedByID        string    `json:"requested_by_id"`
	Started              bool      `json:"started"`
	StartedAt            time.Time `json:"started_at"`
	StartedAtOverride    time.Time `json:"started_at_override"`
	StoryLinks           []struct {
		CreatedAt  time.Time `json:"created_at"`
		EntityType string    `json:"entity_type"`
		ID         int64     `json:"id"`
		ObjectID   int64     `json:"object_id"`
		SubjectID  int64     `json:"subject_id"`
		Type       string    `json:"type"`
		UpdatedAt  time.Time `json:"updated_at"`
		Verb       string    `json:"verb"`
	} `json:"story_links"`
	StoryType       string    `json:"story_type"`
	TaskIDs         []int64   `json:"task_ids"`
	UpdatedAt       time.Time `json:"updated_at"`
	WorkflowStateID int64     `json:"workflow_state_id"`

	Project  *Project
	Blockers []*Story
}

func NewClient(token string) *Client {
	return &Client{
		c:     &http.Client{Timeout: 5 * time.Second},
		token: token,
	}
}

type Client struct {
	c     *http.Client
	token string
}

type Workspace struct {
	Stories []*Story
}

type GetWorkspaceParams struct {
	OnlyProjects []string
}

func (c Client) GetWorkspace(ctx context.Context, params GetWorkspaceParams) (*Workspace, error) {
	projects, err := c.ListProjects(ctx)
	if err != nil {
		return nil, fmt.Errorf("get workspace failed: %w", err)
	}

	projectIDs := make([]int64, 0)
	for _, project := range projects {
		if params.OnlyProjects != nil && !contains(params.OnlyProjects, project.Name) {
			continue
		}
		projectIDs = append(projectIDs, project.ID)
	}

	stories, err := c.SearchStories(ctx, SearchStoriesParams{
		ProjectIDs: projectIDs,
	})
	if err != nil {
		return nil, fmt.Errorf("get workspace failed: %w", err)
	}

	var workspace Workspace
	for _, story := range stories {
		story.Project = findProject(story.ProjectID, projects)
		for _, link := range story.StoryLinks {
			if link.Type != "object" {
				continue
			}
			story.Blockers = append(story.Blockers, findStory(link.SubjectID, stories))
		}
		workspace.Stories = append(workspace.Stories, story)
	}

	return &workspace, nil
}

func findProject(id int64, projects []*Project) *Project {
	for _, proj := range projects {
		if proj.ID == id {
			return proj
		}
	}
	return nil
}

func findStory(id int64, stories []*Story) *Story {
	for _, story := range stories {
		if story.ID == id {
			return story
		}
	}
	return nil
}

type SearchStoriesParams struct {
	ProjectIDs []int64 `json:"project_ids,omitempty"`
}

func (c Client) SearchStories(ctx context.Context, params SearchStoriesParams) ([]*Story, error) {
	resp, err := c.post("/stories/search", params)
	if err != nil {
		return nil, fmt.Errorf("search stories failed: %w", err)
	}

	var stories []*Story
	if err := json.NewDecoder(resp.Body).Decode(&stories); err != nil {
		return nil, fmt.Errorf("search stories failed: %w", err)
	}
	resp.Body.Close()

	return stories, nil
}

type Project struct {
	Abbreviation      string    `json:"abbreviation"`
	Archived          bool      `json:"archived"`
	Color             string    `json:"color"`
	CreatedAt         time.Time `json:"created_at"`
	DaysToThermometer int64     `json:"days_to_thermometer"`
	Description       string    `json:"description"`
	EntityType        string    `json:"entity_type"`
	ExternalID        string    `json:"external_id"`
	FollowerIDs       []string  `json:"follower_ids"`
	ID                int64     `json:"id"`
	IterationLength   int64     `json:"iteration_length"`
	Name              string    `json:"name"`
	ShowThermometer   bool      `json:"show_thermometer"`
	StartTime         time.Time `json:"start_time"`
	Stats             struct {
		NumPoints  int64 `json:"num_points"`
		NumStories int64 `json:"num_stories"`
	} `json:"stats"`
	TeamID    int64     `json:"team_id"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (c Client) ListProjects(ctx context.Context) ([]*Project, error) {
	resp, err := c.get("/projects")
	if err != nil {
		return nil, fmt.Errorf("list projects failed: %w", err)
	}

	var projects []*Project
	if err := json.NewDecoder(resp.Body).Decode(&projects); err != nil {
		return nil, fmt.Errorf("list projects failed: %w", err)
	}
	resp.Body.Close()

	return projects, nil
}

func (c Client) get(path string) (*http.Response, error) {
	return c.makeRequest(http.MethodGet, path, nil)
}

func (c Client) post(path string, body interface{}) (*http.Response, error) {
	var buff bytes.Buffer
	if err := json.NewEncoder(&buff).Encode(body); err != nil {
		return nil, fmt.Errorf("post failed: %w", err)
	}
	return c.makeRequest(http.MethodPost, path, &buff)
}

func (c Client) makeRequest(method, path string, body io.Reader) (*http.Response, error) {
	theURL := c.makeURL(path)

	req, err := http.NewRequest(method, theURL, body)
	if err != nil {
		return nil, fmt.Errorf("make request failed: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("make request failed: %w", err)
	}
	// dump(resp)

	return resp, nil
}

func (c Client) makeURL(path string) string {
	vals := url.Values{}
	vals.Set("token", c.token)
	return BaseURL + path + "?" + vals.Encode()
}

func dump(resp *http.Response) {
	b, _ := httputil.DumpResponse(resp, true)
	fmt.Fprintln(os.Stderr, string(b))
}

func contains(ss []string, s string) bool {
	for _, fs := range ss {
		if fs == s {
			return true
		}
	}
	return false
}
