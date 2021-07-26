package model

type GitHookType string

const (
	GitEventPush         GitHookType = "Push"
	GitEventTag          GitHookType = "Tag"
	GitEventIssue        GitHookType = "Issue"
	GitEventNote         GitHookType = "Note"
	GitEventMergeRequest GitHookType = "Merge Request"
	GitEventWiki         GitHookType = "Wiki Page"
	GitEventPipeline     GitHookType = "Pipeline"
	GitEventJob          GitHookType = "Job"
	GitEventDeployment   GitHookType = "Deployment"
	GitEventMember       GitHookType = "Member"
	GitEventSubgroup     GitHookType = "Subgroup"
	GitEventFeatureFlag  GitHookType = "Feature Flag"
	GitEventRelease      GitHookType = "Release"
)
