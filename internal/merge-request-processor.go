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
	if m.isApproved(request) {
		return m.parseApproved(request), false, nil
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
		request.ObjectAttributes.Action == model.MRActionOpen
}

func (m MergeRequestProcessor) isUpdated(request model.MergeRequest) bool {
	return request.ObjectAttributes.State == model.MRStateOpened &&
		request.ObjectAttributes.Action == model.MRActionUpdate
}

func (m MergeRequestProcessor) isApproved(request model.MergeRequest) bool {
	return request.ObjectAttributes.State == model.MRStateOpened &&
		request.ObjectAttributes.Action == model.MRActionApproved
}

func (m MergeRequestProcessor) parseNew(request model.MergeRequest) string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%s *Опубликован Merge Request* № *%d*\n", GetEmoji(Loudspeaker), request.ObjectAttributes.Iid))
	sb.WriteString(fmt.Sprintf("%s Проект *%s*\n", GetEmoji(WhiteLargeCircle), request.Project.Name))
	sb.WriteString(fmt.Sprintf("%s source *%s* -> target *%s*\n", GetEmoji(WhiteLargeCircle), request.ObjectAttributes.SourceBranch, request.ObjectAttributes.TargetBranch))
	sb.WriteString(fmt.Sprintf("%s Публикатор *%s*\n", GetEmoji(WhiteLargeCircle), request.User.Name))
	sb.WriteString(fmt.Sprintf("%s Ответственный *%s*\n", GetEmoji(WhiteLargeCircle), m.getAssigneeName(request)))
	sb.WriteString(fmt.Sprintf("%s Заголовок *%s*\n", GetEmoji(WhiteLargeCircle), request.ObjectAttributes.Title))
	sb.WriteString(fmt.Sprintf("%s Ссылка *%s*", GetEmoji(WhiteLargeCircle), request.ObjectAttributes.URL))
	return sb.String()
}

func (m MergeRequestProcessor) parseUpdated(request model.MergeRequest) string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%s *Обновлен Merge Request* № *%d*\n", GetEmoji(Loudspeaker), request.ObjectAttributes.Iid))
	sb.WriteString(fmt.Sprintf("%s Проект *%s*\n", GetEmoji(OrangeLargeCircle), request.Project.Name))
	sb.WriteString(fmt.Sprintf("%s source *%s* -> target *%s*\n", GetEmoji(OrangeLargeCircle), request.ObjectAttributes.SourceBranch, request.ObjectAttributes.TargetBranch))
	sb.WriteString(fmt.Sprintf("%s Публикатор *%s*\n", GetEmoji(OrangeLargeCircle), request.User.Name))
	sb.WriteString(fmt.Sprintf("%s Ответственный *%s*\n", GetEmoji(OrangeLargeCircle), m.getAssigneeName(request)))
	sb.WriteString(fmt.Sprintf("%s Заголовок *%s*\n", GetEmoji(OrangeLargeCircle), request.ObjectAttributes.Title))
	sb.WriteString(fmt.Sprintf("%s Ссылка *%s*", GetEmoji(OrangeLargeCircle), request.ObjectAttributes.URL))
	return sb.String()
}

func (m MergeRequestProcessor) parseLocked(request model.MergeRequest) string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%s *Заблокирован Merge Request* № *%d*\n", GetEmoji(Loudspeaker), request.ObjectAttributes.Iid))
	sb.WriteString(fmt.Sprintf("%s Проект *%s*\n", GetEmoji(BlueLargeCircle), request.Project.Name))
	sb.WriteString(fmt.Sprintf("%s source *%s* -> target *%s*\n", GetEmoji(BlueLargeCircle), request.ObjectAttributes.SourceBranch, request.ObjectAttributes.TargetBranch))
	sb.WriteString(fmt.Sprintf("%s Публикатор *%s*\n", GetEmoji(BlueLargeCircle), request.User.Name))
	sb.WriteString(fmt.Sprintf("%s Заголовок *%s*\n", GetEmoji(BlueLargeCircle), request.ObjectAttributes.Title))
	sb.WriteString(fmt.Sprintf("%s Ссылка *%s*", GetEmoji(BlueLargeCircle), request.ObjectAttributes.URL))
	return sb.String()
}

func (m MergeRequestProcessor) parseClosed(request model.MergeRequest) string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%s *Закрыт Merge Request* № *%d*\n", GetEmoji(Loudspeaker), request.ObjectAttributes.Iid))
	sb.WriteString(fmt.Sprintf("%s Проект *%s*\n", GetEmoji(BlackLargeCircle), request.Project.Name))
	sb.WriteString(fmt.Sprintf("%s source *%s* -> target *%s*\n", GetEmoji(BlackLargeCircle), request.ObjectAttributes.SourceBranch, request.ObjectAttributes.TargetBranch))
	sb.WriteString(fmt.Sprintf("%s Заголовок *%s*\n", GetEmoji(BlackLargeCircle), request.ObjectAttributes.Title))
	sb.WriteString(fmt.Sprintf("%s Ссылка *%s*", GetEmoji(BlackLargeCircle), request.ObjectAttributes.URL))
	return sb.String()
}

func (m MergeRequestProcessor) parseMerged(request model.MergeRequest) string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%s *Успешное слияние Merge Request* № *%d*\n", GetEmoji(Loudspeaker), request.ObjectAttributes.Iid))
	sb.WriteString(fmt.Sprintf("%s Проект *%s*\n", GetEmoji(GreenLargeCircle), request.Project.Name))
	sb.WriteString(fmt.Sprintf("%s source *%s* -> target *%s*\n", GetEmoji(GreenLargeCircle), request.ObjectAttributes.SourceBranch, request.ObjectAttributes.TargetBranch))
	sb.WriteString(fmt.Sprintf("%s Заголовок *%s*\n", GetEmoji(GreenLargeCircle), request.ObjectAttributes.Title))
	sb.WriteString(fmt.Sprintf("%s Ссылка *%s*", GetEmoji(GreenLargeCircle), request.ObjectAttributes.URL))
	return sb.String()
}

func (m MergeRequestProcessor) parseApproved(request model.MergeRequest) string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%s *Одобрен Merge Request* № *%d*\n", GetEmoji(Loudspeaker), request.ObjectAttributes.Iid))
	sb.WriteString(fmt.Sprintf("%s Проект *%s*\n", GetEmoji(WhiteCheckMark), request.Project.Name))
	sb.WriteString(fmt.Sprintf("%s source *%s* -> target *%s*\n", GetEmoji(WhiteCheckMark), request.ObjectAttributes.SourceBranch, request.ObjectAttributes.TargetBranch))
	sb.WriteString(fmt.Sprintf("%s Публикатор *%s*\n", GetEmoji(WhiteCheckMark), request.User.Name))
	sb.WriteString(fmt.Sprintf("%s Ответственный *%s*\n", GetEmoji(WhiteCheckMark), m.getAssigneeName(request)))
	sb.WriteString(fmt.Sprintf("%s Заголовок *%s*\n", GetEmoji(WhiteCheckMark), request.ObjectAttributes.Title))
	sb.WriteString(fmt.Sprintf("%s Ссылка *%s*", GetEmoji(WhiteCheckMark), request.ObjectAttributes.URL))
	return sb.String()
}

func (m MergeRequestProcessor) getAssigneeName(request model.MergeRequest) string {
	if len(request.Assignees) == 0 {
		return request.ObjectAttributes.Assignee.Name
	} else {
		return request.Assignees[0].Name
	}
}

var _ Processor = (*MergeRequestProcessor)(nil)
