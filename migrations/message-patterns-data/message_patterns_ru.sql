-- В этом скрипте в базу льём словарь с шаблонами ру-сообщений.
-- всё пишется в одном инсёрте
INSERT INTO message_pattern (hook_type, lang, patterns, additional_patterns)
VALUES ('MergeRequestsEvents', 'ru_RU',
        '{
          "approved": "В репозитории «%s» одобрен запрос на слияние «%s»",
          "close": "В репозитории «%s» закрыт запрос на слияние «%s»",
          "merge": "В репозитории «%s» завёршён запрос на слияние «%s»",
          "open": "В репозитории «%s» открыт запрос на слияние «%s»",
          "reopen": "В репозитории «%s» переоткрыт запрос на слияние «%s»",
          "update": "В репозитории «%s» обновлён запрос на слияние «%s»"
        }',
        '{
          "branches" :"Ветки %s -> %s",
          "opened_by":"Открыто пользователем: %s",
          "assigned_to": "Ответственный пользователь: %s",
          "approved_by":"Одобрено пользователем: %s",
          "closed_by": "Закрыто пользователем: %s",
          "merged_by": "Слияние произведено пользователем: %s"
        }')
ON CONFLICT (hook_type, lang) DO UPDATE SET
patterns = excluded.patterns,
additional_patterns = excluded.additional_patterns;