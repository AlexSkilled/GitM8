package model

const (
	MRStateOpened = "opened"
	MRStateClosed = "closed"
	MRStateLocked = "locked"
	MRStateMerged = "merged"
)

type MergeRequest struct {
	ObjectKind string `json:"object_kind"`
	User       struct {
		ID        int    `json:"id"`
		Name      string `json:"name" faker:"name"`
		Username  string `json:"username"`
		AvatarURL string `json:"avatar_url"`
		Email     string `json:"email"`
	} `json:"user"`
	Project struct {
		ID                int         `json:"id" faker:"boundary_start=1, boundary_end=1000"`
		Name              string      `json:"name" faker:"word"`
		Description       string      `json:"description" faker:"sentence"`
		WebURL            string      `json:"web_url"`
		AvatarURL         interface{} `json:"avatar_url" faker:"-"`
		GitSSHURL         string      `json:"git_ssh_url"`
		GitHTTPURL        string      `json:"git_http_url"`
		Namespace         string      `json:"namespace"`
		VisibilityLevel   int         `json:"visibility_level"`
		PathWithNamespace string      `json:"path_with_namespace"`
		DefaultBranch     string      `json:"default_branch"`
		Homepage          string      `json:"homepage"`
		URL               string      `json:"url"`
		SSHURL            string      `json:"ssh_url"`
		HTTPURL           string      `json:"http_url"`
	} `json:"project"`
	Repository struct {
		Name        string `json:"name"`
		URL         string `json:"url"`
		Description string `json:"description"`
		Homepage    string `json:"homepage"`
	} `json:"repository"`
	ObjectAttributes struct {
		ID              int         `json:"id"`
		TargetBranch    string      `json:"target_branch" faker:"word"`
		SourceBranch    string      `json:"source_branch" faker:"word"`
		SourceProjectID int         `json:"source_project_id"`
		AuthorID        int         `json:"author_id"`
		AssigneeID      int         `json:"assignee_id"`
		Title           string      `json:"title" faker:"sentence"`
		CreatedAt       Datetime    `json:"created_at" faker:"datetime"`
		UpdatedAt       Datetime    `json:"updated_at" faker:"datetime"`
		MilestoneID     interface{} `json:"milestone_id" faker:"-"`
		State           string      `json:"state" faker:"oneof: opened, closed, locked, merged"`
		MergeStatus     string      `json:"merge_status"`
		TargetProjectID int         `json:"target_project_id"`
		Iid             int         `json:"iid"`
		Description     string      `json:"description"`
		Source          struct {
			Name              string      `json:"name"`
			Description       string      `json:"description"`
			WebURL            string      `json:"web_url"`
			AvatarURL         interface{} `json:"avatar_url" faker:"-"`
			GitSSHURL         string      `json:"git_ssh_url"`
			GitHTTPURL        string      `json:"git_http_url"`
			Namespace         string      `json:"namespace"`
			VisibilityLevel   int         `json:"visibility_level"`
			PathWithNamespace string      `json:"path_with_namespace"`
			DefaultBranch     string      `json:"default_branch"`
			Homepage          string      `json:"homepage"`
			URL               string      `json:"url"`
			SSHURL            string      `json:"ssh_url"`
			HTTPURL           string      `json:"http_url"`
		} `json:"source"`
		Target struct {
			Name              string      `json:"name"`
			Description       string      `json:"description"`
			WebURL            string      `json:"web_url"`
			AvatarURL         interface{} `json:"avatar_url" faker:"-"`
			GitSSHURL         string      `json:"git_ssh_url"`
			GitHTTPURL        string      `json:"git_http_url"`
			Namespace         string      `json:"namespace"`
			VisibilityLevel   int         `json:"visibility_level"`
			PathWithNamespace string      `json:"path_with_namespace"`
			DefaultBranch     string      `json:"default_branch"`
			Homepage          string      `json:"homepage"`
			URL               string      `json:"url"`
			SSHURL            string      `json:"ssh_url"`
			HTTPURL           string      `json:"http_url"`
		} `json:"target"`
		LastCommit struct {
			ID        string   `json:"id"`
			Message   string   `json:"message"`
			Timestamp Datetime `json:"timestamp" faker:"datetime"`
			URL       string   `json:"url"`
			Author    struct {
				Name  string `json:"name"`
				Email string `json:"email"`
			} `json:"author"`
		} `json:"last_commit"`
		WorkInProgress bool   `json:"work_in_progress"`
		URL            string `json:"url" faker:"url"`
		Action         string `json:"action"`
		Assignee       struct {
			Name      string `json:"name" faker:"name"`
			Username  string `json:"username"`
			AvatarURL string `json:"avatar_url"`
		} `json:"assignee"`
	} `json:"object_attributes"`
	Labels []struct {
		ID          int      `json:"id"`
		Title       string   `json:"title"`
		Color       string   `json:"color"`
		ProjectID   int      `json:"project_id"`
		CreatedAt   Datetime `json:"created_at" faker:"datetime"`
		UpdatedAt   Datetime `json:"updated_at" faker:"datetime"`
		Template    bool     `json:"template"`
		Description string   `json:"description"`
		Type        string   `json:"type"`
		GroupID     int      `json:"group_id"`
	} `json:"labels"`
	Changes struct {
		UpdatedByID struct {
			Previous interface{} `json:"previous" faker:"-"`
			Current  int         `json:"current"`
		} `json:"updated_by_id"`
		UpdatedAt struct {
			Previous string `json:"previous"`
			Current  string `json:"current"`
		} `json:"updated_at"`
		Labels struct {
			Previous []struct {
				ID          int      `json:"id"`
				Title       string   `json:"title"`
				Color       string   `json:"color"`
				ProjectID   int      `json:"project_id"`
				CreatedAt   Datetime `json:"created_at" faker:"datetime"`
				UpdatedAt   Datetime `json:"updated_at" faker:"datetime"`
				Template    bool     `json:"template"`
				Description string   `json:"description"`
				Type        string   `json:"type"`
				GroupID     int      `json:"group_id"`
			} `json:"previous"`
			Current []struct {
				ID          int      `json:"id"`
				Title       string   `json:"title"`
				Color       string   `json:"color"`
				ProjectID   int      `json:"project_id"`
				CreatedAt   Datetime `json:"created_at" faker:"datetime"`
				UpdatedAt   Datetime `json:"updated_at" faker:"datetime"`
				Template    bool     `json:"template"`
				Description string   `json:"description"`
				Type        string   `json:"type"`
				GroupID     int      `json:"group_id"`
			} `json:"current"`
		} `json:"labels"`
	} `json:"changes"`
}
