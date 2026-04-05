# User Service (用户服务)

## 1. 服务概述

**职责**: 用户资料管理、个人设置、二维码、推送Token管理

**核心功能**:
- 用户资料管理（获取、修改、搜索）
- 个人设置（验证开关、通知设置、隐私设置）
- 二维码功能（生成、刷新、扫码）
- 推送Token管理
- 用户状态（在线状态、最后活跃时间）

## 2. 文档导航

| 功能 | 文档 | 说明 |
|------|------|------|
| 用户资料 | [profile.md](profile.md) | 获取/修改/搜索用户资料 |
| 个人设置 | [settings.md](settings.md) | 用户设置管理 |
| 二维码 | [qrcode.md](qrcode.md) | 二维码生成与扫码 |
| 推送Token | [push-token.md](push-token.md) | 推送Token管理 |

## 3. 数据模型

- **UserProfile**: 用户详细资料
- **UserSettings**: 用户个人设置
- **UserQRCode**: 用户二维码记录
- **UserPushToken**: 推送Token

## 4. 推送通知

- `notification.user.profile_updated.{user_id}` - 用户资料更新通知
- `notification.user.friend_profile_changed.{user_id}` - 好友资料变更通知
- `notification.user.status_changed.{user_id}` - 在线状态变更通知

## 5. 依赖服务

- **Auth Service**: 用户认证
- **File Service**: 头像上传
- **Redis**: 在线状态、二维码缓存
- **PostgreSQL**: 用户资料持久化

---

返回: [后端总体设计](../backend-design.md)
