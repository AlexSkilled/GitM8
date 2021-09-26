package events

import (
	"encoding/json"
	"gitlab-tg-bot/service/model"
	"gitlab-tg-bot/service/payload"
	"gitlab-tg-bot/service/payload/mergereq"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

const MergeRequestHeader = "Merge Request Hook"

type MergeRequest struct {
	User struct {
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
		VisibilityLevel   int32       `json:"visibility_level"`
		PathWithNamespace string      `json:"path_with_namespace"`
		DefaultBranch     string      `json:"default_branch"`
		CiConfigPath      interface{} `json:"ci_config_path"`
		Homepage          string      `json:"homepage"`
		Url               string      `json:"url"`
		SshUrl            string      `json:"ssh_url"`
		HttpUrl           string      `json:"http_url"`
	} `json:"project"`
	ObjectAttributes struct {
		AssigneeId                int32       `json:"assignee_id"`
		AuthorId                  int         `json:"author_id"`
		CreatedAt                 string      `json:"created_at"`
		Description               string      `json:"description"`
		HeadPipelineId            interface{} `json:"head_pipeline_id"`
		MergeStatus               string      `json:"merge_status"`
		MergeUserId               interface{} `json:"merge_user_id"`
		MergeWhenPipelineSucceeds bool        `json:"merge_when_pipeline_succeeds"`
		MilestoneId               interface{} `json:"milestone_id"`
		SourceBranch              string      `json:"source_branch"`
		SourceProjectId           int32       `json:"source_project_id"`
		StateId                   int32       `json:"state_id"`
		TargetBranch              string      `json:"target_branch"`
		TargetProjectId           int32       `json:"target_project_id"`
		TimeEstimate              int32       `json:"time_estimate"`
		Title                     string      `json:"title"`
		UpdatedAt                 string      `json:"updated_at"`
		UpdatedById               interface{} `json:"updated_by_id"`
		Url                       string      `json:"url"`
		Source                    struct {
			Id                int32       `json:"id"`
			Name              string      `json:"name"`
			Description       string      `json:"description"`
			WebUrl            string      `json:"web_url"`
			AvatarUrl         interface{} `json:"avatar_url"`
			GitSshUrl         string      `json:"git_ssh_url"`
			GitHttpUrl        string      `json:"git_http_url"`
			Namespace         string      `json:"namespace"`
			VisibilityLevel   int32       `json:"visibility_level"`
			PathWithNamespace string      `json:"path_with_namespace"`
			DefaultBranch     string      `json:"default_branch"`
			CiConfigPath      interface{} `json:"ci_config_path"`
			Homepage          string      `json:"homepage"`
			Url               string      `json:"url"`
			SshUrl            string      `json:"ssh_url"`
			HttpUrl           string      `json:"http_url"`
		} `json:"source"`
		Target struct {
			Id                int32       `json:"id"`
			Name              string      `json:"name"`
			Description       string      `json:"description"`
			WebUrl            string      `json:"web_url"`
			AvatarUrl         interface{} `json:"avatar_url"`
			GitSshUrl         string      `json:"git_ssh_url"`
			GitHttpUrl        string      `json:"git_http_url"`
			Namespace         string      `json:"namespace"`
			VisibilityLevel   int32       `json:"visibility_level"`
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
		State  string  `json:"state"`
		Action string  `json:"action"`
		Oldrev *string `json:"oldrev"`
	} `json:"object_attributes"`
	Changes struct {
		Title *struct {
			Previous string `json:"previous"`
			Current  string `json:"current"`
		} `json:"title"`
		Assignees *struct {
			Previous []struct {
				Name      string `json:"name"`
				Username  string `json:"username"`
				AvatarUrl string `json:"avatar_url"`
				Email     string `json:"email"`
			} `json:"previous"`
			Current []struct {
				Name      string `json:"name"`
				Username  string `json:"username"`
				AvatarUrl string `json:"avatar_url"`
				Email     string `json:"email"`
			} `json:"current"`
		} `json:"assignees"`
	} `json:"changes"`
	Repository struct {
		Name        string `json:"name"`
		Url         string `json:"url"`
		Description string `json:"description"`
		Homepage    string `json:"homepage"`
	} `json:"repository"`
	Assignees []struct {
		Name      string `json:"name"`
		Username  string `json:"username"`
		AvatarUrl string `json:"avatar_url"`
		Email     string `json:"email"`
	} `json:"assignees"`
}

func (p *MergeRequest) ToModel() model.GitEvent {
	payloadMap := make(model.Payload)

	pl := mergereq.MergeRequest{
		Name:         p.ObjectAttributes.Title,
		Link:         p.ObjectAttributes.Url,
		SourceBranch: p.ObjectAttributes.SourceBranch,
		TargetBranch: p.ObjectAttributes.TargetBranch,
	}

	if p.Assignees != nil {
		assignedToArray := make([]string, len(p.Assignees))
		for i, item := range p.Assignees {
			assignedToArray[i] = item.Name
		}
		pl.AssignedTo = strings.Join(assignedToArray, ", ")
	}

	payloadMap[payload.Main], _ = json.Marshal(pl)

	update := p.getUpdate()
	if update != nil {
		payloadMap[payload.Changes] = update
	}

	return model.GitEvent{
		GitSource:   model.Gitlab,
		ProjectId:   strconv.Itoa(p.Project.Id),
		ProjectName: p.Project.Name,

		HookType: model.HookTypeMergeRequests,
		SubType:  convertSubType(p.ObjectAttributes.Action),

		TriggeredByName: p.User.Name,
		AuthorId:        strconv.Itoa(p.User.Id),

		Payload: payloadMap,
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

func (p *MergeRequest) getUpdate() []byte {
	var update *mergereq.Change

	if p.Changes.Assignees != nil {
		update = &mergereq.Change{
			New:  p.Changes.Assignees.Current[0].Name,
			Type: mergereq.ReAssignee,
		}
		if len(p.Changes.Assignees.Previous) != 0 {
			update.Old = p.Changes.Assignees.Previous[0].Name
		}
	} else if p.Changes.Title != nil {
		update = &mergereq.Change{
			Old:  p.Changes.Title.Previous,
			New:  p.Changes.Title.Current,
			Type: mergereq.Rename,
		}
	} else if p.ObjectAttributes.Oldrev != nil {
		update = &mergereq.Change{
			Old:  p.ObjectAttributes.LastCommit.Url,
			New:  p.ObjectAttributes.LastCommit.Title,
			Type: mergereq.Update,
		}
	}
	if update == nil {
		return nil
	}

	out, err := json.Marshal(update)
	if err != nil {
		logrus.Error(err)
		return nil
	}

	return out
}
