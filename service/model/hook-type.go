package model

type XGitEvent string

const (
	GitEventPush             XGitEvent = "Push"
	GitEventPushTag          XGitEvent = "Tag Push"
	GitEventPushIssue        XGitEvent = "Issue"
	GitEventPushNote         XGitEvent = "Note"
	GitEventPushMergeRequest XGitEvent = "Merge Request"
	GitEventPushWiki         XGitEvent = "Wiki Page"
	GitEventPushPipeline     XGitEvent = "Pipeline"
	GitEventPushJob          XGitEvent = "Job"
	GitEventPushDeployment   XGitEvent = "Deployment"
	GitEventPushMember       XGitEvent = "Member"
	GitEventPushSubgroup     XGitEvent = "Subgroup"
	GitEventPushFeatureFlag  XGitEvent = "Feature Flag"
	GitEventPushRelease      XGitEvent = "Release"
)
