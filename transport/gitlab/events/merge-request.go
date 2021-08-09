package events

import (
	"gitlab-tg-bot/service/model"
	"strconv"
	"time"
)

const MergeRequestHeader = "Merge Request Hook"

type MergeRequest struct {
	ObjectKind string `json:"object_kind"`
	EventType  string `json:"event_type"`
	User       struct {
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
		Homepage          string      `json:"homepage"`
		Url               string      `json:"url"`
		SshUrl            string      `json:"ssh_url"`
		HttpUrl           string      `json:"http_url"`
	} `json:"project"`
	ObjectAttributes struct {
		AssigneeId     interface{} `json:"assignee_id"`
		AuthorId       int         `json:"author_id"`
		CreatedAt      string      `json:"created_at"`
		Description    string      `json:"description"`
		HeadPipelineId interface{} `json:"head_pipeline_id"`
		Id             int         `json:"id"`
		Iid            int         `json:"iid"`
		LastEditedAt   interface{} `json:"last_edited_at"`
		LastEditedById interface{} `json:"last_edited_by_id"`
		MergeCommitSha interface{} `json:"merge_commit_sha"`
		MergeError     interface{} `json:"merge_error"`
		MergeParams    struct {
			ForceRemoveSourceBranch string `json:"force_remove_source_branch"`
		} `json:"merge_params"`
		MergeStatus               string      `json:"merge_status"`
		MergeUserId               interface{} `json:"merge_user_id"`
		MergeWhenPipelineSucceeds bool        `json:"merge_when_pipeline_succeeds"`
		MilestoneId               interface{} `json:"milestone_id"`
		SourceBranch              string      `json:"source_branch"`
		SourceProjectId           int         `json:"source_project_id"`
		StateId                   int         `json:"state_id"`
		TargetBranch              string      `json:"target_branch"`
		TargetProjectId           int         `json:"target_project_id"`
		TimeEstimate              int         `json:"time_estimate"`
		Title                     string      `json:"title"`
		UpdatedAt                 string      `json:"updated_at"`
		UpdatedById               interface{} `json:"updated_by_id"`
		Url                       string      `json:"url"`
		Source                    struct {
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
			Homepage          string      `json:"homepage"`
			Url               string      `json:"url"`
			SshUrl            string      `json:"ssh_url"`
			HttpUrl           string      `json:"http_url"`
		} `json:"source"`
		Target struct {
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
			Homepage          string      `json:"homepage"`
			Url               string      `json:"url"`
			SshUrl            string      `json:"ssh_url"`
			HttpUrl           string      `json:"http_url"`
		} `json:"target"`
		LastCommit struct {
			Id        string    `json:"id"`
			Message   string    `json:"message"`
			Title     string    `json:"title"`
			Timestamp time.Time `json:"timestamp"`
			Url       string    `json:"url"`
			Author    struct {
				Name  string `json:"name"`
				Email string `json:"email"`
			} `json:"author"`
		} `json:"last_commit"`
		WorkInProgress      bool          `json:"work_in_progress"`
		TotalTimeSpent      int           `json:"total_time_spent"`
		HumanTotalTimeSpent interface{}   `json:"human_total_time_spent"`
		HumanTimeEstimate   interface{}   `json:"human_time_estimate"`
		AssigneeIds         []interface{} `json:"assignee_ids"`
		State               string        `json:"state"`
		Action              string        `json:"action"`
	} `json:"object_attributes"`
	Labels  []interface{} `json:"labels"`
	Changes struct {
		AuthorId struct {
			Previous interface{} `json:"previous"`
			Current  int         `json:"current"`
		} `json:"author_id"`
		CreatedAt struct {
			Previous interface{} `json:"previous"`
			Current  string      `json:"current"`
		} `json:"created_at"`
		Description struct {
			Previous interface{} `json:"previous"`
			Current  string      `json:"current"`
		} `json:"description"`
		Id struct {
			Previous interface{} `json:"previous"`
			Current  int         `json:"current"`
		} `json:"id"`
		Iid struct {
			Previous interface{} `json:"previous"`
			Current  int         `json:"current"`
		} `json:"iid"`
		MergeParams struct {
			Previous struct {
			} `json:"previous"`
			Current struct {
				ForceRemoveSourceBranch string `json:"force_remove_source_branch"`
			} `json:"current"`
		} `json:"merge_params"`
		SourceBranch struct {
			Previous interface{} `json:"previous"`
			Current  string      `json:"current"`
		} `json:"source_branch"`
		SourceProjectId struct {
			Previous interface{} `json:"previous"`
			Current  int         `json:"current"`
		} `json:"source_project_id"`
		TargetBranch struct {
			Previous interface{} `json:"previous"`
			Current  string      `json:"current"`
		} `json:"target_branch"`
		TargetProjectId struct {
			Previous interface{} `json:"previous"`
			Current  int         `json:"current"`
		} `json:"target_project_id"`
		Title struct {
			Previous interface{} `json:"previous"`
			Current  string      `json:"current"`
		} `json:"title"`
		UpdatedAt struct {
			Previous interface{} `json:"previous"`
			Current  string      `json:"current"`
		} `json:"updated_at"`
	} `json:"changes"`
	Repository struct {
		Name        string `json:"name"`
		Url         string `json:"url"`
		Description string `json:"description"`
		Homepage    string `json:"homepage"`
	} `json:"repository"`
}

func (p *MergeRequest) ToModel() *model.GitEvent {
	return &model.GitEvent{
		GitSource:   model.Gitlab,
		ProjectId:   strconv.Itoa(p.Project.Id),
		ProjectName: p.Project.Name,
		HookType:    model.HookTypeMergeRequests,
		SubType:     convertSubType(p.ObjectAttributes.Action),
	}
}

// Маппим действия (action) хуков гитлаба в наши
func convertSubType(subType string) model.GitHookSubtype {
	switch subType {
	case "approved":
		return model.MRApproved
	case "close":
		return model.MRClose
	case "merge":
		return model.MRMerge
	case "open":
		return model.MROpen
	case "reopen":
		return model.MRReopen
	case "update":
		return model.MRUpdated
	}
	return model.MRUnknown
}
