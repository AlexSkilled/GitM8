package entity

type SubscriptionActionTypeEnum string

const (
	SubscriptionActionType_Pipeline SubscriptionActionTypeEnum = "pipeline"
)

type SubscriptionRepository struct {
	UserGitlabId int32

	Repository string
	ActionType SubscriptionActionTypeEnum
	Branch     string
}
