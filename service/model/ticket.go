package model

type Ticket struct {
	Id                 int32
	MaintainerGitlabId int64
	RepositoryId       string
	HookTypes          map[GitHookType]interface{}
	ChatIds            []int64
}

type GitHookType string

const (
	HookTypePush               GitHookType = "PushEvents"
	HookTypeIssues             GitHookType = "IssuesEvents"
	HookTypeConfidentialIssues GitHookType = "ConfidentialIssuesEvents"
	HookTypeMergeRequests      GitHookType = "MergeRequestsEvents"
	HookTypeTagPush            GitHookType = "TagPushEvents"
	HookTypeNote               GitHookType = "NoteEvents"
	HookTypeJob                GitHookType = "JobEvents"
	HookTypePipeline           GitHookType = "PipelineEvents"
	HookTypeWikiPage           GitHookType = "WikiPageEvents"
)

type GitHookSubtype string

const (
	MRApproved GitHookSubtype = "approved"
	MRClose    GitHookSubtype = "close"
	MRMerge    GitHookSubtype = "merge"
	MROpen     GitHookSubtype = "open"
	MRReopen   GitHookSubtype = "reopen"
	MRUpdated  GitHookSubtype = "update"

	MRUnknown GitHookSubtype = ""
)

func (g GitHookType) GetSubs() []GitHookSubtype {
	if g == HookTypeMergeRequests {
		return []GitHookSubtype{MRApproved, MRClose, MRMerge, MROpen, MRReopen, MRUpdated}
	}
	return nil
}
