-- В этом скрипте в базу льём словарь с шаблонами ру-сообщений.
-- всё пишется в одном инсёрте
INSERT INTO message_pattern (hook_type, lang, patterns, additional_patterns)
VALUES ('MergeRequestsEvents', 'ru_RU',
        '{
          "approved": "В репозитории «%v» одобрен запрос на слияние «%v»",
          "close": "В репозитории «%v» закрыт запрос на слияние «%v»",
          "merge": "В репозитории «%v» завершён запрос на слияние «%v»",
          "open": "В репозитории «%v» открыт запрос на слияние «%v»",
          "reopen": "В репозитории «%v» переоткрыт запрос на слияние «%v»",
          "update": "В репозитории «%v» обновлён запрос на слияние «%v»"
        }',
        '{
          "branches" :"Ветки %v -> %v",
          "opened_by":"Открыто пользователем: %v",
          "assigned_to": "Ответственный пользователь: %v",
          "approved_by":"Одобрено пользователем: %v",
          "closed_by": "Закрыто пользователем: %v",
          "merged_by": "Слияние произведено пользователем: %v",
          "updated_by": "Обновлено пользователем: %v",

          "rename": "Переименование запроса на слияние %v -> %v",
          "update": "Коммит: %v",
          "reAssignee": "Ответственный переназначен %v -> %v"
        }')
ON CONFLICT (hook_type, lang) DO UPDATE SET
patterns = excluded.patterns,
additional_patterns = excluded.additional_patterns;