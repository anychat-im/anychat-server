-- 会话表
CREATE TABLE sessions (
    session_id      VARCHAR(100) PRIMARY KEY,
    session_type    VARCHAR(20)  NOT NULL,                -- single/group/system
    user_id         VARCHAR(100) NOT NULL,
    target_id       VARCHAR(100) NOT NULL,                -- 单聊为对方用户ID，群聊为群ID
    last_message_id VARCHAR(100),
    last_message_content TEXT,
    last_message_time    TIMESTAMPTZ,
    unread_count    INT          NOT NULL DEFAULT 0,
    is_pinned       BOOLEAN      NOT NULL DEFAULT FALSE,
    is_muted        BOOLEAN      NOT NULL DEFAULT FALSE,
    pin_time        TIMESTAMPTZ,
    burn_after_reading INT       NOT NULL DEFAULT 0,       -- 阅后即焚时长(秒),0表示关闭
    auto_delete_duration INT     NOT NULL DEFAULT 0,       -- 自动删除时长(秒),0表示未启用
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX uk_session_user_target ON sessions (user_id, session_type, target_id);
CREATE INDEX idx_sessions_user_id      ON sessions (user_id);
CREATE INDEX idx_sessions_updated_at   ON sessions (updated_at);

-- 消息发送幂等表
CREATE TABLE IF NOT EXISTS message_send_idempotency (
    id BIGSERIAL PRIMARY KEY,
    sender_id VARCHAR(36) NOT NULL,
    conversation_id VARCHAR(64) NOT NULL,
    local_id VARCHAR(128) NOT NULL,
    message_id VARCHAR(64) NOT NULL DEFAULT '',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uk_sender_conversation_local UNIQUE (sender_id, conversation_id, local_id)
);

CREATE INDEX idx_message_idempotency_message_id ON message_send_idempotency(message_id);
