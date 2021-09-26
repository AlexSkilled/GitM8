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

	message := &utils.NoticeMessage{}

	mergeRequestNameWithLink := "[" + pl.Name + "](" + pl.Link + ")"

	message.Header.WriteWithEmoji(mrMessageEmoji[event.SubType], messageText, event.ProjectName, mergeRequestNameWithLink)

	message.MainInfo.WriteWithEmoji(emoji.Branches, patterns[mergereq.Pattern_Branches], pl.SourceBranch, pl.TargetBranch)

	message.Author.WriteWithEmoji(emoji.Man, patterns[mergereq.Pattern_OpenedBy], event.AuthorName)

	if len(pl.AssignedTo) != 0 && event.SubType != model.MRApproved {
		message.AssignedToUserName = utils.NewMessageRawWithEmoji(emoji.EyeWatch, patterns[mergereq.Pattern_AssignedTo], pl.AssignedTo)
	}

	writeMrPayload(event, message, patterns)

	return message.PrintWithEmoji(), nil
}

func writeMrPayload(event model.GitEvent, message *utils.NoticeMessage, patterns map[string]string) {
	switch event.SubType {
	case model.MRApproved:
		message.TriggeredByUserName = utils.NewMessageRawWithEmoji(emoji.Man, patterns[mergereq.Pattern_ApprovedBy], event.TriggeredByName)
	case model.MRClose:
		message.TriggeredByUserName = utils.NewMessageRawWithEmoji(emoji.Man, patterns[mergereq.Pattern_ClosedBy], event.TriggeredByName)
	case model.MRMerge:
		message.TriggeredByUserName = utils.NewMessageRawWithEmoji(emoji.Man, patterns[mergereq.Pattern_MergedBy], event.TriggeredByName)
	case model.MRUpdated:
		message.TriggeredByUserName = utils.NewMessageRawWithEmoji(emoji.Man, patterns[mergereq.Pattern_UpdatedBy], event.TriggeredByName)
		message.SubInfo = utils.NewMessageRawWithEmoji(extractChanges(event.Payload[payload.Changes], patterns))

	case model.MRUnknown:
	}
}

func extractChanges(change []byte, patterns map[string]string) (emo string, mes string) {
	var update mergereq.Change
	err := json.Unmarshal(change, &update)
	if err != nil {
		panic(err) // TODO
	}

	switch update.Type {
	case mergereq.Rename:
		return emoji.New, fmt.Sprintf(patterns[mergereq.Pattern_Rename], update.Old, update.New)
	case mergereq.Update:
		return emoji.New, fmt.Sprintf(patterns[mergereq.Pattern_Update], "["+update.Old+"]("+update.New+")")
	case mergereq.ReAssignee:
		return emoji.New, fmt.Sprintf(patterns[mergereq.Pattern_ReAssignee], update.Old, update.New)
	}
	return "", ""
}

var mrMessageEmoji = map[model.GitHookSubtype]string{
	model.MRApproved: "✅",
	model.MRClose:    "🛑",
	model.MRMerge:    "🔀",
	model.MROpen:     "🆕" + emoji.GetEmoji(emoji.Loudspeaker),
	model.MRReopen:   "🔄",
	model.MRUpdated:  "⤴",
}
