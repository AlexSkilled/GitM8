package provider

import (
	"errors"
	"fmt"
	"gitlab-tg-bot/data/entity"
	"gitlab-tg-bot/internal/interfaces"
	"gitlab-tg-bot/service/model"

	"github.com/go-pg/pg/v9"
)

type MessagePatternProvider struct {
	*pg.DB
}

var _ interfaces.MessagePatternProvider = (*MessagePatternProvider)(nil)

func NewMessagePatternProvider(conn *pg.DB) *MessagePatternProvider {
	return &MessagePatternProvider{conn}
}

func (m *MessagePatternProvider) GetMessage(lang string, hookType model.GitHookType, subType model.GitHookSubtype) (string, map[string]string, error) {
	var pattern entity.MessagePattern
	_, err := m.Query(&pattern, `
				SELECT 
					  patterns->> ?,
					  additional_patterns
 			  	FROM  message_pattern
				WHERE lang      = ?
				AND   hook_type = ?`, subType, lang, hookType)
	if err != nil {
		return "", nil, errors.New(fmt.Sprintf("ошибка при попопытке извлечь шаблон сообщения для %v с подтипом %v на %v языке", hookType, subType, lang))
	}
	return pattern.Message, pattern.AdditionalPatterns, nil
}
