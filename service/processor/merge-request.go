package processor

import (
	"encoding/json"
	"fmt"
	"gitlab-tg-bot/service/model"
	"gitlab-tg-bot/service/payload"
	"gitlab-tg-bot/service/payload/mergereq"
	"gitlab-tg-bot/utils"
)

func ProcessMergeRequest(event model.GitEvent, messageText string, patterns map[string]string) (string, error) {
	var pl mergereq.MergeRequest

	err := json.Unmarshal(event.Payload[payload.Main], &pl)
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
		message.WriteStringNF(
			extractChanges(event.Payload[payload.Changes], patterns))

	case model.MRUnknown:
	}
	return message.String()
}

func extractChanges(change []byte, patterns map[string]string) string {
	var update mergereq.Change
	err := json.Unmarshal(change, &update)
	if err != nil {
		panic(err) // TODO
	}
	switch update.Type {
	case mergereq.Rename:
		return fmt.Sprintf(patterns[mergereq.Pattern_Rename], update.Old, update.New)
	case mergereq.Update:
		return fmt.Sprintf(patterns[mergereq.Pattern_Update], "["+update.Old+"]("+update.New+")")
	case mergereq.ReAssignee:
		return fmt.Sprintf(patterns[mergereq.Pattern_ReAssignee], update.Old, update.New)
	}
	return ""
}
