package processor

import (
	"fmt"
	"gitlab-tg-bot/service/model"
	"strings"
)

func ProcessMergeRequest(event model.GitEvent, messageText string, additional map[string]string) string {
	sb := &strings.Builder{}

	mergeRequestNameWithLink := "[" + event.Payload[model.PKHeader] + "](" + event.Payload[model.PKLink] + ")"

	//utils.AppendWithPattern(sb, messageText, event.ProjectName, mergeRequestNameWithLink)

	sb.WriteString("\n" + fmt.Sprintf(messageText, event.ProjectName, mergeRequestNameWithLink))

	sb.WriteString("\n" + fmt.Sprintf(additional[model.MRPattern_Branches],

		event.Payload[model.PKSourceBranch],
		event.Payload[model.PKTargetBranch]))

	sb.WriteString("\n" + fmt.Sprintf(additional[model.MRPattern_OpenedBy], event.Payload[model.PKCreatedByUser]))

	if assignedToUser, ok := event.Payload[model.PKmrAssignedToUser]; ok {
		sb.WriteString(fmt.Sprintf(additional[model.MRPattern_AssignedTo], assignedToUser))
	}

	switch event.SubType {
	case model.MRApproved:
		sb.WriteString(
			fmt.Sprintf(additional[model.MRPattern_ApprovedBy],
				event.Payload[model.PKTriggeredByUser]))
	case model.MRClose:
		sb.WriteString(
			fmt.Sprintf(additional[model.MRPattern_ClosedBy],
				event.Payload[model.PKmrClosedBy]))
	case model.MRMerge:
		sb.WriteString(
			fmt.Sprintf(additional[model.MRPattern_MergedBy],
				event.Payload[model.PKTriggeredByUser]))
	case model.MRUpdated:
		sb.WriteString(
			fmt.Sprintf(additional[model.MRPattern_UpdatedBy],
				event.Payload[model.PKmrUpdatedBy]))
	case model.MRUnknown:

	}

	//sb.WriteString("\n[link](https://google.com)")
	return sb.String()
}
