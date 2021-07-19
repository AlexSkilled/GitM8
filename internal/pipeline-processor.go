package internal

import (
	"encoding/json"
	"fmt"
	"gitlab-tg-bot/internal/model"
	"strings"
)

type PipelineProcessor struct {
}

func (p PipelineProcessor) Process(payload []byte) (msg string, skip bool, err error) {
	var request model.Pipeline
	err = json.Unmarshal(payload, &request)
	if err != nil {
		return msg, true, err
	}
	if request.ObjectAttributes.Status == model.PLStatusCanceled {
		return p.parseCanceled(request), false, err
	}
	if request.ObjectAttributes.Status == model.PLStatusSkipped {
		return p.parseSkipped(request), false, err
	}
	if request.ObjectAttributes.Status == model.PLStatusFailed {
		return p.parseFailed(request), false, err
	}
	if request.ObjectAttributes.Status == model.PLStatusSuccess {
		return p.parseSuccess(request), false, err
	}
	return msg, true, nil
}

func (p PipelineProcessor) parseCanceled(request model.Pipeline) string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%s *Отменен Pipeline* № *%d*\n", GetEmoji(Loudspeaker), request.ObjectAttributes.ID))
	sb.WriteString(fmt.Sprintf("%s Проект *%s*\n", GetEmoji(BlackLargeSquare), request.Project.Name))
	sb.WriteString(fmt.Sprintf("%s Ссылка на связанный Merge Request *%s*\n", GetEmoji(BlackLargeSquare), request.MergeRequest.URL))
	sb.WriteString(fmt.Sprintf("%s Инициатор *%s*", GetEmoji(BlackLargeSquare), request.User.Name))
	return sb.String()
}

func (p PipelineProcessor) parseSkipped(request model.Pipeline) string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%s *Пропущен Pipeline* № *%d*\n", GetEmoji(Loudspeaker), request.ObjectAttributes.ID))
	sb.WriteString(fmt.Sprintf("%s Проект *%s*\n", GetEmoji(YellowLargeSquare), request.Project.Name))
	sb.WriteString(fmt.Sprintf("%s Ссылка на связанный Merge Request *%s*\n", GetEmoji(YellowLargeSquare), request.MergeRequest.URL))
	sb.WriteString(fmt.Sprintf("%s Инициатор *%s*", GetEmoji(YellowLargeSquare), request.User.Name))
	return sb.String()
}

func (p PipelineProcessor) parseFailed(request model.Pipeline) string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%s *Завершился ошибкой Pipeline* № *%d*\n", GetEmoji(Loudspeaker), request.ObjectAttributes.ID))
	sb.WriteString(fmt.Sprintf("%s Проект *%s*\n", GetEmoji(RedLargeSquare), request.Project.Name))
	sb.WriteString(fmt.Sprintf("%s Ссылка на связанный Merge Request *%s*\n", GetEmoji(RedLargeSquare), request.MergeRequest.URL))
	sb.WriteString(fmt.Sprintf("%s Инициатор *%s*\n", GetEmoji(RedLargeSquare), request.User.Name))
	sb.WriteString(fmt.Sprintf("%s Время исполнения *%d* сек.", GetEmoji(RedLargeSquare), request.ObjectAttributes.Duration))
	return sb.String()
}

func (p PipelineProcessor) parseSuccess(request model.Pipeline) string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%s *Успешно завершен Pipeline* № *%d*\n", GetEmoji(Loudspeaker), request.ObjectAttributes.ID))
	sb.WriteString(fmt.Sprintf("%s Проект *%s*\n", GetEmoji(GreenLargeSquare), request.Project.Name))
	sb.WriteString(fmt.Sprintf("%s Запрос на слияние *%s*\n", GetEmoji(GreenLargeSquare), request.MergeRequest.URL))
	sb.WriteString(fmt.Sprintf("%s Инициатор *%s*\n", GetEmoji(GreenLargeSquare), request.User.Name))
	sb.WriteString(fmt.Sprintf("%s Время исполнения *%d* сек.", GetEmoji(GreenLargeSquare), request.ObjectAttributes.Duration))
	return sb.String()
}

var _ Processor = (*PipelineProcessor)(nil)
