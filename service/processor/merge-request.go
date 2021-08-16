package processor

import (
	"encoding/json"
	"fmt"
	"gitlab-tg-bot/service/model"
	"gitlab-tg-bot/service/payload/mergereq"
	"strings"
)

func ProcessMergeRequest(event model.GitEvent, messageText string, patterns map[string]string) (string, error) {
	var pl mergereq.MergeRequest

	err := json.Unmarshal(event.ActualPayload, &pl)
	if err != nil {
		panic(err) // TODO
	}

	sb := &strings.Builder{}

	mergeRequestNameWithLink := "[" + pl.Name + "](" + pl.Link + ")"

	sb.WriteString("\n" + fmt.Sprintf(messageText, event.ProjectName, mergeRequestNameWithLink))

	sb.WriteString("\n" + fmt.Sprintf(patterns[mergereq.Pattern_Branches], pl.SourceBranch, pl.TargetBranch))

	sb.WriteString("\n" + fmt.Sprintf(patterns[mergereq.Pattern_OpenedBy], event.AuthorName))

	if len(pl.AssignedTo) != 0 {
		sb.WriteString("\n" + fmt.Sprintf(patterns[mergereq.Pattern_AssignedTo], pl.AssignedTo))
	}

	switch event.SubType {
	case model.MRApproved:
		sb.WriteString(
			fmt.Sprintf("\n"+patterns[mergereq.Pattern_ApprovedBy], event.TriggeredByName))
	case model.MRClose:
		sb.WriteString(
			fmt.Sprintf("\n"+patterns[mergereq.Pattern_ClosedBy], event.TriggeredByName))
	case model.MRMerge:
		sb.WriteString(
			fmt.Sprintf("\n"+patterns[mergereq.Pattern_MergedBy], event.TriggeredByName))
	case model.MRUpdated:
		sb.WriteString(
			fmt.Sprintf("\n"+patterns[mergereq.Pattern_UpdatedBy], event.TriggeredByName))
	case model.MRUnknown:
	}
	return sb.String(), nil
}
