package processor

import (
	"encoding/json"
	"fmt"
	"gitlab-tg-bot/service/model"
	"gitlab-tg-bot/service/payload/mergereq"
	"gitlab-tg-bot/utils"
)

func ProcessMergeRequest(event model.GitEvent, messageText string, patterns map[string]string) (string, error) {
	var pl mergereq.MergeRequest

	err := json.Unmarshal(event.Payload, &pl)
	if err != nil {
		panic(err) // TODO
	}

	message := &utils.MessageBuilder{}

	mergeRequestNameWithLink := "[" + pl.Name + "](" + pl.Link + ")"

	message.WriteStringNF(fmt.Sprintf(messageText, event.ProjectName, mergeRequestNameWithLink))

	message.WriteStringNF(fmt.Sprintf(patterns[mergereq.Pattern_Branches], pl.SourceBranch, pl.TargetBranch))

	message.WriteStringNF(fmt.Sprintf(patterns[mergereq.Pattern_OpenedBy], event.AuthorName))

	if len(pl.AssignedTo) != 0 {
		message.WriteStringNF(fmt.Sprintf(patterns[mergereq.Pattern_AssignedTo], pl.AssignedTo))
	}

	switch event.SubType {
	case model.MRApproved:
		message.WriteStringNF(
			fmt.Sprintf(patterns[mergereq.Pattern_ApprovedBy], event.TriggeredByName))
	case model.MRClose:
		message.WriteStringNF(
			fmt.Sprintf(patterns[mergereq.Pattern_ClosedBy], event.TriggeredByName))
	case model.MRMerge:
		message.WriteStringNF(
			fmt.Sprintf(patterns[mergereq.Pattern_MergedBy], event.TriggeredByName))
	case model.MRUpdated:
		message.WriteStringNF(
			fmt.Sprintf(patterns[mergereq.Pattern_UpdatedBy], event.TriggeredByName))
	case model.MRUnknown:
	}
	return message.String()
}
