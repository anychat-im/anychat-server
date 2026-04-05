# Group Service (群组服务)

## 1. 服务概述

**职责**: 群组创建、成员管理、群设置

**核心功能**:
- 群组创建与解散
- 成员管理（邀请、申请、移除、退出）
- 群主与管理员管理
- 群组设置（名称、头像、公告、权限）
- 群成员设置（群昵称、免打扰、置顶）
- 群组信息同步

## 2. 文档导航

| 功能 | 文档 | 说明 |
|------|------|------|
| 群组管理 | [group.md](group.md) | 群组创建、解散、信息 |
| 群成员管理 | [member.md](member.md) | 成员邀请、移除、角色管理 |
| 群设置 | [settings.md](settings.md) | 群名称、公告、权限设置 |

## 3. 数据模型

- **Group**: 群组基本信息
- **GroupMember**: 群成员关系
- **GroupSettings**: 群组设置
- **GroupAdmin**: 群管理员
- **GroupMute**: 群禁言记录
- **GroupJoinRequest**: 入群申请

## 4. 推送通知

- `notification.group.invited.{user_id}` - 群组邀请通知
- `notification.group.member_joined.{group_id}` - 新成员加入通知
- `notification.group.member_left.{group_id}` - 成员退出/被移除通知
- `notification.group.info_updated.{group_id}` - 群组信息更新通知
- `notification.group.role_changed.{group_id}` - 成员角色变更通知
- `notification.group.muted.{group_id}` - 群组禁言通知
- `notification.group.disbanded.{group_id}` - 群组解散通知

## 5. 依赖服务

- **User Service**: 用户信息
- **Message Service**: 群系统消息
- **File Service**: 群头像
- **NATS**: 群变更事件
- **Redis**: 群信息缓存、成员缓存
- **PostgreSQL**: 群数据持久化

---

返回: [后端总体设计](../backend-design.md)
