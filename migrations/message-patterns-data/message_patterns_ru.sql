-- В этом скрипте в базу льём словарь с шаблонами ру-сообщений.
-- всё пишется в одном инсёрте
INSERT INTO message_pattern (hook_type, lang, patterns, additional_patterns)
VALUES ('MergeRequestsEvents', 'ru_RU',
        '{
          "approved": "В репозитории «%s» одобрен запрос на слияние «%s».\nВетки %s -> %s.\nОткрыто пользователем:%s",
          "close": "В репозитории «%s» закрыт запрос на слияние «%s».\nВетки %s -> %s.\nОткрыто пользователем:%s",
          "merge": "В репозитории «%s» завёршён запрос на слияние «%s».\nВетки %s -> %s.\nОткрыто пользователем:%s",
          "open": "В репозитории «%s» открыт запрос на слияние «%s».\nВетки %s -> %s.\nОткрыто пользователем:%s",
          "reopen": "В репозитории «%s» переоткрыт запрос на слияние «%s».\nВетки %s -> %s.\nОткрыто пользователем:%s",
          "update": "В репозитории «%s» обновлён запрос на слияние «%s».\nВетки %s -> %s.\nОткрыто пользователем:%s"
        }',
        '{
          "assigned_to": "\nОтветственный пользователь: %s",
          "approved_by":"\nОдобрено пользователем: %s",
          "closed_by": "\nЗакрыто пользователем: %s",
          "merged_by": "\nСлияние произведено пользователем: %s",
          "updated_by": "\nОбновлено пользователем: %s",
          "mr_link": "\nСсылка на запрос на слияние: %s",
          "pipeline_link": "\nСсылка на пайп: %s"
        }')
ON CONFLICT (hook_type, lang) DO UPDATE SET patterns = excluded.patterns;