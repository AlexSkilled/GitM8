package events

import "time"

const PipelineHeader = "Pipeline Hook"

type Pipeline struct {
	ObjectKind       string `json:"object_kind"`
	ObjectAttributes struct {
		Id             int           `json:"id"`
		Ref            string        `json:"ref"`
		Tag            bool          `json:"tag"`
		Sha            string        `json:"sha"`
		BeforeSha      string        `json:"before_sha"`
		Source         string        `json:"source"`
		Status         string        `json:"status"`
		DetailedStatus string        `json:"detailed_status"`
		Stages         []string      `json:"stages"`
		CreatedAt      time.Time     `json:"created_at"`
		FinishedAt     time.Time     `json:"finished_at"`
		Duration       int           `json:"duration"`
		Variables      []interface{} `json:"variables"`
	} `json:"object_attributes"`
	MergeRequest interface{} `json:"merge_request"`
	User         struct {
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
		CiConfigPath      interface{} `json:"ci_config_path"`
	} `json:"project"`
	Commit struct {
		Id        string    `json:"id"`
		Message   string    `json:"message"`
		Title     string    `json:"title"`
		Timestamp time.Time `json:"timestamp"`
		Url       string    `json:"url"`
		Author    struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"author"`
	} `json:"commit"`
	Builds []struct {
		Id           int       `json:"id"`
		Stage        string    `json:"stage"`
		Name         string    `json:"name"`
		Status       string    `json:"status"`
		CreatedAt    time.Time `json:"created_at"`
		StartedAt    time.Time `json:"started_at"`
		FinishedAt   time.Time `json:"finished_at"`
		When         string    `json:"when"`
		Manual       bool      `json:"manual"`
		AllowFailure bool      `json:"allow_failure"`
		User         struct {
			Name      string `json:"name"`
			Username  string `json:"username"`
			AvatarUrl string `json:"avatar_url"`
			Email     string `json:"email"`
		} `json:"user"`
		Runner struct {
			Id          int    `json:"id"`
			Description string `json:"description"`
			Active      bool   `json:"active"`
			IsShared    bool   `json:"is_shared"`
		} `json:"runner"`
		ArtifactsFile struct {
			Filename interface{} `json:"filename"`
			Size     interface{} `json:"size"`
		} `json:"artifacts_file"`
	} `json:"builds"`
}
