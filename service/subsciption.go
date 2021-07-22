package service

import (
	"gitlab-tg-bot/internal"
	"gitlab-tg-bot/internal/interfaces"
	"gitlab-tg-bot/internal/model"
)

const (
	HookTypePushEvents               = "PushEvents"
	HookTypeIssuesEvents             = "IssuesEvents"
	HookTypeConfidentialIssuesEvents = "ConfidentialIssuesEvents"
	HookTypeMergeRequestsEvents      = "MergeRequestsEvents"
	HookTypeTagPushEvents            = "TagPushEvents"
	HookTypeNoteEvents               = "NoteEvents"
	HookTypeJobEvents                = "JobEvents"
	HookTypePipelineEvents           = "PipelineEvents"
	HookTypeWikiPageEvents           = "WikiPageEvents"
)

type SubscriptionService struct {
	TicketProvider interfaces.TicketProvider

	GitlabApi interfaces.GitlabApiService
}

var _ interfaces.SubscriptionService = (*SubscriptionService)(nil)

func NewSubscription(conf internal.Configuration, provider interfaces.ProviderStorage,
	gitlabApi interfaces.GitlabApiService) *SubscriptionService {
	return &SubscriptionService{
		TicketProvider: provider.Ticket(),

		GitlabApi: gitlabApi,
	}
}

func (s *SubscriptionService) Subscribe(gitlab model.GitlabUser, tgUserId int64, hookInfo model.Hook) (int32, error) {
	ticket := model.Ticket{
		MaintainerGitlabId: gitlab.Id,
		RepositoryId:       hookInfo.RepoId,
		HookTypes: map[string]interface{}{
			HookTypePushEvents:               hookInfo.PushEvents,
			HookTypeIssuesEvents:             hookInfo.IssuesEvents,
			HookTypeConfidentialIssuesEvents: hookInfo.ConfidentialIssuesEvents,
			HookTypeMergeRequestsEvents:      hookInfo.MergeRequestsEvents,
			HookTypeTagPushEvents:            hookInfo.TagPushEvents,
			HookTypeNoteEvents:               hookInfo.NoteEvents,
			HookTypeJobEvents:                hookInfo.JobEvents,
			HookTypePipelineEvents:           hookInfo.PipelineEvents,
			HookTypeWikiPageEvents:           hookInfo.WikiPageEvents,
		},
	}
	ticketId, err := s.TicketProvider.AddTicket(ticket)
	if err != nil {
		return 0, err
	}
	err = s.TicketProvider.AddTicketToUser(tgUserId, ticketId)
	if err != nil {
		return 0, err
	}
	return ticketId, nil
}

func (s *SubscriptionService) GetRepositories(user model.GitlabUser) ([]model.Repository, error) {
	return s.GitlabApi.GetRepositories(user)
}
