# 会话管理设计

## 1. 概述

会话管理提供会话列表、置顶、免打扰、未读数、阅后即焚、自动删除消息等核心功能。

## 2. 功能列表

- [x] [获取会话列表](./session.md#41-获取会话列表)
- [x] [创建/更新会话](./session.md#42-创建更新会话)
- [x] [删除会话](./session.md#43-删除会话)
- [x] [会话置顶](./session.md#44-会话置顶)
- [x] [会话免打扰](./session.md#45-会话免打扰)
- [x] [未读数管理](./session.md#46-未读数管理)
- [x] [阅后即焚](./session.md#47-阅后即焚)
- [x] [自动删除消息](./auto_delete.md)

## 3. 数据模型

### 3.1 Session 表

```go
type Session struct {
    SessionID          string     // 会话ID
    SessionType        string     // 会话类型: single/group/system
    UserID             string     // 用户ID
    TargetID           string     // 目标ID(用户或群组)
    LastMessageID      string     // 最后一条消息ID
    LastMessageContent string     // 最后消息摘要
    LastMessageTime    *time.Time // 最后消息时间
    UnreadCount        int32      // 未读数
    IsPinned           bool       // 是否置顶
    IsMuted            bool       // 是否免打扰
    PinTime            *time.Time // 置顶时间
    BurnAfterReading   int32      // 阅后即焚时长(秒),0表示未启用
    AutoDeleteDuration int32      // 自动删除时长(秒),0表示未启用
    CreatedAt          time.Time
    UpdatedAt          time.Time
}
```

## 4. 业务流程

### 4.1 获取会话列表

详见 [自动删除消息](./auto_delete.md) 文档。

```mermaid
sequenceDiagram
    participant Client
    participant Gateway
    participant SessionService
    participant MessageService
    participant DB

    Client->>Gateway: GET /session/list<br/>Header: Authorization: Bearer {token}
    Gateway->>Gateway: 从JWT解析userId
    Gateway->>SessionService: gRPC GetSessions(userId, limit, updatedBefore)
    SessionService->>DB: 查询会话列表(按置顶+更新时间排序)
    SessionService->>MessageService: 获取最后消息详情
    SessionService-->>Gateway: 返回会话列表
    Gateway-->>Client: 200 OK
```

### 4.2 创建/更新会话

消息到达时由消息服务调用，更新会话的最后消息信息。

### 4.3 删除会话

```mermaid
sequenceDiagram
    participant Client
    participant Gateway
    participant SessionService
    participant DB
    participant NATS

    Client->>Gateway: DELETE /session/{sessionId}<br/>Header: Authorization: Bearer {token}
    Gateway->>Gateway: 从JWT解析userId
    Gateway->>SessionService: gRPC DeleteSession(userId, sessionId)
    SessionService->>DB: 删除会话
    DB-->>SessionService: 成功
    SessionService->>NATS: 发布会话删除事件
    SessionService-->>Gateway: 成功
    Gateway-->>Client: 200 OK
```

### 4.4 会话置顶

```mermaid
sequenceDiagram
    participant Client
    participant Gateway
    participant SessionService
    participant DB
    participant NATS

    Client->>Gateway: PUT /session/pin<br/>Header: Authorization: Bearer {token}<br/>Body: {session_id, pinned: true/false}
    Gateway->>Gateway: 从JWT解析userId
    Gateway->>SessionService: gRPC SetPinned(userId, sessionId, pinned)
    SessionService->>DB: 更新置顶状态
    DB-->>SessionService: 成功
    SessionService->>NATS: 发布置顶变更事件
    SessionService-->>Gateway: 成功
    Gateway-->>Client: 200 OK
```

### 4.5 会话免打扰

```mermaid
sequenceDiagram
    participant Client
    participant Gateway
    participant SessionService
    participant DB
    participant NATS

    Client->>Gateway: PUT /session/mute<br/>Header: Authorization: Bearer {token}<br/>Body: {session_id, muted: true/false}
    Gateway->>Gateway: 从JWT解析userId
    Gateway->>SessionService: gRPC SetMuted(userId, sessionId, muted)
    SessionService->>DB: 更新免打扰状态
    DB-->>SessionService: 成功
    SessionService->>NATS: 发布免打扰变更事件
    SessionService-->>Gateway: 成功
    Gateway-->>Client: 200 OK
```

### 4.6 未读数管理

- **增加未读数**：消息服务发送消息时调用
- **清除未读数**：用户查看会话时调用

### 4.7 阅后即焚

用户设置后，接收方阅读消息时触发删除。

详见 [自动删除消息](./auto_delete.md) 了解与自动删除的区别。

```mermaid
sequenceDiagram
    participant Client
    participant Gateway
    participant SessionService
    participant DB

    Client->>Gateway: PUT /session/burn<br/>Header: Authorization: Bearer {token}<br/>Body: {session_id, duration: 30}
    Gateway->>Gateway: 从JWT解析userId
    Gateway->>SessionService: gRPC SetBurnAfterReading(userId, sessionId, duration)
    SessionService->>DB: 更新阅后即焚时长
    DB-->>SessionService: 成功
    SessionService-->>Gateway: 成功
    Gateway-->>Client: 200 OK
```
> duration为0表示取消阅后即焚

## 5. API设计

### 5.1 获取会话列表

```protobuf
message GetSessionsRequest {
    string user_id = 1;
    int32 limit = 2;
    int64 updated_before = 3;
}

message GetSessionsResponse {
    repeated Session sessions = 1;
}
```

### 5.2 置顶/免打扰

```protobuf
message SetPinnedRequest {
    string user_id = 1;
    string session_id = 2;
    bool pinned = 3;
}

message SetMutedRequest {
    string user_id = 1;
    string session_id = 2;
    bool muted = 3;
}
```

### 5.3 阅后即焚

```protobuf
message SetBurnAfterReadingRequest {
    string user_id = 1;
    string session_id = 2;
    int32 duration = 3;  // 秒,0表示取消
}
```

### 5.4 自动删除消息

详见 [自动删除消息](./auto_delete.md) 文档。

## 6. 通知主题

- `notification.session.pin_updated.{user_id}` - 置顶状态变更
- `notification.session.mute_updated.{user_id}` - 免打扰状态变更
- `notification.session.deleted.{user_id}` - 会话删除
- `notification.session.unread_updated.{user_id}` - 未读数变更
- `notification.session.burn_updated.{user_id}` - 阅后即焚配置变更
- `notification.session.auto_delete_updated.{user_id}` - 自动删除配置变更