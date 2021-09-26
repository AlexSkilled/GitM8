package processor

import (
	"encoding/json"
	"fmt"
	"gitlab-tg-bot/internal/emoji"
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

	messageText = mrMessagePreAppend[event.SubType] + messageText

	mergeRequestNameWithLink := "[" + pl.Name + "](" + pl.Link + ")"

	message.WriteStringNF(fmt.Sprintf(messageText, event.ProjectName, mergeRequestNameWithLink))

	message.WriteStringNF("🔀" + fmt.Sprintf(patterns[mergereq.Pattern_Branches], pl.SourceBranch, pl.TargetBranch))

	message.WriteStringNF(emoji.Man + fmt.Sprintf(patterns[mergereq.Pattern_OpenedBy], event.AuthorName))

	if len(pl.AssignedTo) != 0 && event.SubType != model.MRApproved {
		message.WriteStringNF("👁‍🗨" + fmt.Sprintf(patterns[mergereq.Pattern_AssignedTo], pl.AssignedTo))
	}

	writeMrPayload(event, message, patterns)

	return message.String()
}

func writeMrPayload(event model.GitEvent, message *utils.MessageBuilder, patterns map[string]string) {
	switch event.SubType {
	case model.MRApproved:
		message.WriteStringNF(emoji.Man+patterns[mergereq.Pattern_ApprovedBy], event.TriggeredByName)
	case model.MRClose:
		message.WriteStringNF(emoji.Man+patterns[mergereq.Pattern_ClosedBy], event.TriggeredByName)
	case model.MRMerge:
		message.WriteStringNF(emoji.Man+patterns[mergereq.Pattern_MergedBy], event.TriggeredByName)
	case model.MRUpdated:
		message.WriteStringNF(emoji.Man+patterns[mergereq.Pattern_UpdatedBy], event.TriggeredByName)
		message.WriteStringN(extractChanges(event.Payload[payload.Changes], patterns))

	case model.MRUnknown:
	}
}

func extractChanges(change []byte, patterns map[string]string) string {
	var update mergereq.Change
	err := json.Unmarshal(change, &update)
	if err != nil {
		panic(err) // TODO
	}

	switch update.Type {
	case mergereq.Rename:
		return fmt.Sprintf("🆕"+patterns[mergereq.Pattern_Rename], update.Old, update.New)
	case mergereq.Update:
		return fmt.Sprintf("🆕"+patterns[mergereq.Pattern_Update], "["+update.Old+"]("+update.New+")")
	case mergereq.ReAssignee:
		return fmt.Sprintf("🆕"+patterns[mergereq.Pattern_ReAssignee], update.Old, update.New)
	}
	return ""
}

var mrMessagePreAppend = map[model.GitHookSubtype]string{
	model.MRApproved: "✅",
	model.MRClose:    "🛑",
	model.MRMerge:    "🔀",
	model.MROpen:     "🆕" + emoji.GetEmoji(emoji.Loudspeaker),
	model.MRReopen:   "🔄",
	model.MRUpdated:  "⤴",
}
