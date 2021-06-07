package internal

import (
	"encoding/json"
	"fmt"
	"gitlab-tg-bot/internal/model"
	"strings"
)

type MergeRequestProcessor struct {
}

func (m MergeRequestProcessor) Process(payload []byte) (msg string, skip bool, err error) {
	var request model.MergeRequest
	err = json.Unmarshal(payload, &request)
	if err != nil {
		return msg, true, err
	}
	if m.isNew(request) {
		return m.parseNew(request), false, nil
	}
	if m.isUpdated(request) {
		return m.parseUpdated(request), false, nil
	}
	if request.ObjectAttributes.State == model.MRStateClosed {
		return m.parseClosed(request), false, nil
	}
	if request.ObjectAttributes.State == model.MRStateLocked {
		return m.parseLocked(request), false, nil
	}
	if request.ObjectAttributes.State == model.MRStateMerged {
		return m.parseMerged(request), false, nil
	}
	return msg, true, nil
}

func (m MergeRequestProcessor) isNew(request model.MergeRequest) bool {
	return request.ObjectAttributes.State == model.MRStateOpened &&
		request.ObjectAttributes.CreatedAt == request.ObjectAttributes.UpdatedAt
}

func (m MergeRequestProcessor) isUpdated(request model.MergeRequest) bool {
	return request.ObjectAttributes.State == model.MRStateOpened &&
		request.ObjectAttributes.CreatedAt != request.ObjectAttributes.UpdatedAt
}

func (m MergeRequestProcessor) parseNew(request model.MergeRequest) string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%s *Опубликован Merge Request* № *%d*\n", getEmoji(Loudspeaker), request.ObjectAttributes.ID))
	sb.WriteString(fmt.Sprintf("%s Проект *%s*\n", getEmoji(WhiteLargeCircle), request.Project.Name))
	sb.WriteString(fmt.Sprintf("%s source *%s* в *%s*\n", getEmoji(WhiteLargeCircle), request.ObjectAttributes.SourceBranch, request.ObjectAttributes.TargetBranch))
	sb.WriteString(fmt.Sprintf("%s Публикатор *%s*\n", getEmoji(WhiteLargeCircle), request.User.Name))
	sb.WriteString(fmt.Sprintf("%s Ответственный *%s*\n", getEmoji(WhiteLargeCircle), request.ObjectAttributes.Assignee.Name))
	sb.WriteString(fmt.Sprintf("%s Заголовок *%s*\n", getEmoji(WhiteLargeCircle), request.ObjectAttributes.Title))
	sb.WriteString(fmt.Sprintf("%s Ссылка *%s*", getEmoji(WhiteLargeCircle), request.ObjectAttributes.URL))
	return sb.String()
}

func (m MergeRequestProcessor) parseUpdated(request model.MergeRequest) string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%s *Обновлен Merge Request* № *%d*\n", getEmoji(Loudspeaker), request.ObjectAttributes.ID))
	sb.WriteString(fmt.Sprintf("%s Проект *%s*\n", getEmoji(OrangeLargeCircle), request.Project.Name))
	sb.WriteString(fmt.Sprintf("%s source *%s* -> target *%s*\n", getEmoji(OrangeLargeCircle), request.ObjectAttributes.SourceBranch, request.ObjectAttributes.TargetBranch))
	sb.WriteString(fmt.Sprintf("%s Публикатор *%s*\n", getEmoji(OrangeLargeCircle), request.User.Name))
	sb.WriteString(fmt.Sprintf("%s Ответственный *%s*\n", getEmoji(OrangeLargeCircle), request.ObjectAttributes.Assignee.Name))
	sb.WriteString(fmt.Sprintf("%s Заголовок *%s*\n", getEmoji(OrangeLargeCircle), request.ObjectAttributes.Title))
	sb.WriteString(fmt.Sprintf("%s Ссылка *%s*", getEmoji(OrangeLargeCircle), request.ObjectAttributes.URL))
	return sb.String()
}

func (m MergeRequestProcessor) parseLocked(request model.MergeRequest) string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%s *Заблокирован Merge Request* № *%d*\n", getEmoji(Loudspeaker), request.ObjectAttributes.ID))
	sb.WriteString(fmt.Sprintf("%s Проект *%s*\n", getEmoji(BlueLargeCircle), request.Project.Name))
	sb.WriteString(fmt.Sprintf("%s source *%s* -> target *%s*\n", getEmoji(BlueLargeCircle), request.ObjectAttributes.SourceBranch, request.ObjectAttributes.TargetBranch))
	sb.WriteString(fmt.Sprintf("%s Публикатор *%s*\n", getEmoji(BlueLargeCircle), request.User.Name))
	sb.WriteString(fmt.Sprintf("%s Заголовок *%s*\n", getEmoji(BlueLargeCircle), request.ObjectAttributes.Title))
	sb.WriteString(fmt.Sprintf("%s Ссылка *%s*", getEmoji(BlueLargeCircle), request.ObjectAttributes.URL))
	return sb.String()
}

func (m MergeRequestProcessor) parseClosed(request model.MergeRequest) string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%s *Закрыт Merge Request* № *%d*\n", getEmoji(Loudspeaker), request.ObjectAttributes.ID))
	sb.WriteString(fmt.Sprintf("%s Проект *%s*\n", getEmoji(BlackLargeCircle), request.Project.Name))
	sb.WriteString(fmt.Sprintf("%s source *%s* -> target *%s*\n", getEmoji(BlackLargeCircle), request.ObjectAttributes.SourceBranch, request.ObjectAttributes.TargetBranch))
	sb.WriteString(fmt.Sprintf("%s Заголовок *%s*\n", getEmoji(BlackLargeCircle), request.ObjectAttributes.Title))
	sb.WriteString(fmt.Sprintf("%s Ссылка *%s*", getEmoji(BlackLargeCircle), request.ObjectAttributes.URL))
	return sb.String()
}

func (m MergeRequestProcessor) parseMerged(request model.MergeRequest) string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%s *Успешное слияние Merge Request* № *%d*\n", getEmoji(Loudspeaker), request.ObjectAttributes.ID))
	sb.WriteString(fmt.Sprintf("%s Проект *%s*\n", getEmoji(GreenLargeCircle), request.Project.Name))
	sb.WriteString(fmt.Sprintf("%s source *%s* -> target *%s*\n", getEmoji(GreenLargeCircle), request.ObjectAttributes.SourceBranch, request.ObjectAttributes.TargetBranch))
	sb.WriteString(fmt.Sprintf("%s Заголовок *%s*\n", getEmoji(GreenLargeCircle), request.ObjectAttributes.Title))
	sb.WriteString(fmt.Sprintf("%s Ссылка *%s*", getEmoji(GreenLargeCircle), request.ObjectAttributes.URL))
	return sb.String()
}

var _ Processor = (*MergeRequestProcessor)(nil)
