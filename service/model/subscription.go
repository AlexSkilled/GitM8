package model

type SubscriptionActionTypeEnum string

const (
	SubscriptionActionType_Pipeline SubscriptionActionTypeEnum = "pipeline"
)

type Subscription struct {
	User User

	Repository string
	ActionType SubscriptionActionTypeEnum
	Branch     string
}
