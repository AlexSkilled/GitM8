package processor

import (
	"fmt"
	"gitlab-tg-bot/service/model"
	"strings"
)

func ProcessMergeRequest(event model.GitEvent, messageText string, additional map[string]string) string {
	sb := strings.Builder{}

	mergeRequestNameWithLink := "[" + event.Payload[model.PKHeader] + "](" + event.Payload[model.PKLink] + ")"

	sb.WriteString(
		fmt.Sprintf(messageText,
			event.ProjectName,
			mergeRequestNameWithLink,
			event.Payload[model.PKSourceBranch],
			event.Payload[model.PKTargetBranch],
			event.Payload[model.PKCreatedByUser]))

	if assignedToUser, ok := event.Payload[model.PKmrAssignedToUser]; ok {
		sb.WriteString(fmt.Sprintf(additional[model.MRSubInfoKey_AssignedTo], assignedToUser))
	}

	switch event.SubType {
	case model.MRApproved:
		sb.WriteString(
			fmt.Sprintf(additional[model.MRSubInfoKey_ApprovedBy],
				event.Payload[model.PKTriggeredByUser]))
	case model.MRClose:
		sb.WriteString(
			fmt.Sprintf(additional[model.MRSubInfoKey_ClosedBy],
				event.Payload[model.PKmrClosedBy]))
	case model.MRMerge:
		sb.WriteString(
			fmt.Sprintf(additional[model.MRSubInfoKey_MergedBy],
				event.Payload[model.PKTriggeredByUser]))
	case model.MRUpdated:
		sb.WriteString(
			fmt.Sprintf(additional[model.MRSubInfoKey_UpdatedBy],
				event.Payload[model.PKmrUpdatedBy]))
	case model.MRUnknown:

	}

	//sb.WriteString("\n[link](https://google.com)")
	return sb.String()
}
