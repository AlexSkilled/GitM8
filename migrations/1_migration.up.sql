CREATE TABLE IF NOT EXISTS tg_user
(
    id          INTEGER GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    name        TEXT,

    tg_username TEXT,
    locale      TEXT,

    urn         TEXT
);

CREATE TABLE IF NOT EXISTS git_user
(
    id         INTEGER GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    name       TEXT,
    user_id    INTEGER REFERENCES tg_user (id),
    email      TEXT,
    token      TEXT,
    domain     TEXT,
    git_source TEXT
);

ALTER TABLE IF EXISTS git_user
    DROP CONSTRAINT IF EXISTS gitlab_PKs;

ALTER TABLE IF EXISTS git_user
    ADD CONSTRAINT gitlab_PKs
        UNIQUE (user_id, token, domain);


CREATE TABLE IF NOT EXISTS ticket
(
    id                   INTEGER GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    maintainer_gitlab_id INTEGER REFERENCES git_user (id),
    repository_id        TEXT,
    hook_types           JSONB,
    name                 TEXT
);

CREATE TABLE IF NOT EXISTS tickets_chat_id
(
    chat_id      INTEGER,
    ticket_id    INTEGER,
    is_active    BOOLEAN DEFAULT TRUE,
    is_notifying BOOLEAN DEFAULT TRUE,

    UNIQUE (chat_id, ticket_id)
);

CREATE TYPE message_type_pattern_enum AS ENUM (
    'PushEvents',
    'IssuesEvents',
    'ConfidentialIssuesEvents',
    'MergeRequestsEvents',
    'TagPushEvents',
    'NoteEvents',
    'JobEvents',
    'PipelineEvents',
    'WikiPageEvents');

CREATE TABLE message_pattern
(
    hook_type           message_type_pattern_enum,
    lang                TEXT,
    patterns            JSONB,
    additional_patterns JSONB,

    UNIQUE (hook_type, lang)
);

CREATE TABLE IF NOT EXISTS webhook (
    id INTEGER GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    git_source TEXT,
    uri TEXT,
    UNIQUE (uri)
);

CREATE TABLE IF NOT EXISTS webhook_to_user (
    webhook_id INTEGER REFERENCES webhook (id),
    user_id    INTEGER REFERENCES tg_user (id)
);