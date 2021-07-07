package processors

import "time"

type MergeRequest struct {
	ObjectKind string `json:"object_kind"`
	User       struct {
		Id        int    `json:"id"`
		Name      string `json:"name"`
		Username  string `json:"username"`
		AvatarUrl string `json:"avatar_url"`
		Email     string `json:"email"`
	} `json:"user"`
	Project struct {
		Id                int         `json:"id"`
		Name              string      `json:"name"`
		Description       string      `json:"description"`
		WebUrl            string      `json:"web_url"`
		AvatarUrl         interface{} `json:"avatar_url"`
		GitSshUrl         string      `json:"git_ssh_url"`
		GitHttpUrl        string      `json:"git_http_url"`
		Namespace         string      `json:"namespace"`
		VisibilityLevel   int         `json:"visibility_level"`
		PathWithNamespace string      `json:"path_with_namespace"`
		DefaultBranch     string      `json:"default_branch"`
		Homepage          string      `json:"homepage"`
		Url               string      `json:"url"`
		SshUrl            string      `json:"ssh_url"`
		HttpUrl           string      `json:"http_url"`
	} `json:"project"`
	Repository struct {
		Name        string `json:"name"`
		Url         string `json:"url"`
		Description string `json:"description"`
		Homepage    string `json:"homepage"`
	} `json:"repository"`
	ObjectAttributes struct {
		Id              int         `json:"id"`
		TargetBranch    string      `json:"target_branch"`
		SourceBranch    string      `json:"source_branch"`
		SourceProjectId int         `json:"source_project_id"`
		AuthorId        int         `json:"author_id"`
		AssigneeId      int         `json:"assignee_id"`
		Title           string      `json:"title"`
		CreatedAt       time.Time   `json:"created_at"`
		UpdatedAt       time.Time   `json:"updated_at"`
		MilestoneId     interface{} `json:"milestone_id"`
		State           string      `json:"state"`
		MergeStatus     string      `json:"merge_status"`
		TargetProjectId int         `json:"target_project_id"`
		Iid             int         `json:"iid"`
		Description     string      `json:"description"`
		Source          struct {
			Name              string      `json:"name"`
			Description       string      `json:"description"`
			WebUrl            string      `json:"web_url"`
			AvatarUrl         interface{} `json:"avatar_url"`
			GitSshUrl         string      `json:"git_ssh_url"`
			GitHttpUrl        string      `json:"git_http_url"`
			Namespace         string      `json:"namespace"`
			VisibilityLevel   int         `json:"visibility_level"`
			PathWithNamespace string      `json:"path_with_namespace"`
			DefaultBranch     string      `json:"default_branch"`
			Homepage          string      `json:"homepage"`
			Url               string      `json:"url"`
			SshUrl            string      `json:"ssh_url"`
			HttpUrl           string      `json:"http_url"`
		} `json:"source"`
		Target struct {
			Name              string      `json:"name"`
			Description       string      `json:"description"`
			WebUrl            string      `json:"web_url"`
			AvatarUrl         interface{} `json:"avatar_url"`
			GitSshUrl         string      `json:"git_ssh_url"`
			GitHttpUrl        string      `json:"git_http_url"`
			Namespace         string      `json:"namespace"`
			VisibilityLevel   int         `json:"visibility_level"`
			PathWithNamespace string      `json:"path_with_namespace"`
			DefaultBranch     string      `json:"default_branch"`
			Homepage          string      `json:"homepage"`
			Url               string      `json:"url"`
			SshUrl            string      `json:"ssh_url"`
			HttpUrl           string      `json:"http_url"`
		} `json:"target"`
		LastCommit struct {
			Id        string    `json:"id"`
			Message   string    `json:"message"`
			Timestamp time.Time `json:"timestamp"`
			Url       string    `json:"url"`
			Author    struct {
				Name  string `json:"name"`
				Email string `json:"email"`
			} `json:"author"`
		} `json:"last_commit"`
		WorkInProgress bool   `json:"work_in_progress"`
		Url            string `json:"url"`
		Action         string `json:"action"`
		Assignee       struct {
			Name      string `json:"name"`
			Username  string `json:"username"`
			AvatarUrl string `json:"avatar_url"`
		} `json:"assignee"`
	} `json:"object_attributes"`
	Labels []struct {
		Id          int       `json:"id"`
		Title       string    `json:"title"`
		Color       string    `json:"color"`
		ProjectId   int       `json:"project_id"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
		Template    bool      `json:"template"`
		Description string    `json:"description"`
		Type        string    `json:"type"`
		GroupId     int       `json:"group_id"`
	} `json:"labels"`
	Changes struct {
		UpdatedById struct {
			Previous interface{} `json:"previous"`
			Current  int         `json:"current"`
		} `json:"updated_by_id"`
		UpdatedAt struct {
			Previous string `json:"previous"`
			Current  string `json:"current"`
		} `json:"updated_at"`
		Labels struct {
			Previous []struct {
				Id          int       `json:"id"`
				Title       string    `json:"title"`
				Color       string    `json:"color"`
				ProjectId   int       `json:"project_id"`
				CreatedAt   time.Time `json:"created_at"`
				UpdatedAt   time.Time `json:"updated_at"`
				Template    bool      `json:"template"`
				Description string    `json:"description"`
				Type        string    `json:"type"`
				GroupId     int       `json:"group_id"`
			} `json:"previous"`
			Current []struct {
				Id          int       `json:"id"`
				Title       string    `json:"title"`
				Color       string    `json:"color"`
				ProjectId   int       `json:"project_id"`
				CreatedAt   time.Time `json:"created_at"`
				UpdatedAt   time.Time `json:"updated_at"`
				Template    bool      `json:"template"`
				Description string    `json:"description"`
				Type        string    `json:"type"`
				GroupId     int       `json:"group_id"`
			} `json:"current"`
		} `json:"labels"`
	} `json:"changes"`
}

func (m *MergeRequest) Process([]byte) (string, error) {
	return "", nil
}
