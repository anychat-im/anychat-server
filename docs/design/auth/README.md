# Auth Service (认证服务)

## 1. 服务概述

**职责**: 用户注册、登录、Token管理、多端登录策略

**核心功能**:
- 用户注册（手机号/邮箱，需验证码）
- 用户登录（账号密码/验证码）
- Token管理（JWT: AccessToken + RefreshToken）
- 多端登录策略与设备管理

## 2. 文档导航

| 功能 | 文档 | 说明 |
|------|------|------|
| 用户注册 | [register.md](register.md) | 手机号/邮箱注册 |
| 用户登录 | [login.md](login.md) | 登录方式与流程 |
| Token管理 | [token.md](token.md) | JWT令牌管理 |
| 会话管理 | [session.md](session.md) | 用户会话 |
| 设备管理 | [device.md](device.md) | 设备登录记录 |
| 密码管理 | [password.md](password.md) | 修改/重置密码 |
| 验证码 | [verification-code.md](verification-code.md) | 验证码发送与验证 |

## 3. 数据模型

- **User**: 用户基本信息
- **UserDevice**: 设备登录记录
- **UserSession**: 用户会话信息

## 4. 推送通知

- `notification.auth.force_logout.{user_id}` - 多端互踢通知
- `notification.auth.unusual_login.{user_id}` - 异常登录提醒
- `notification.auth.password_changed.{user_id}` - 密码修改通知

## 5. 依赖服务

- **ZITADEL**: 身份认证
- **Redis**: Token缓存、在线状态
- **PostgreSQL**: 设备登录记录
- **NATS**: 强制下线消息推送

---

返回: [后端总体设计](../backend-design.md)
