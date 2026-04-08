# AnyChat Server 功能列表

> 用于跟踪开发进度的功能清单

## 模块说明

| 模块 | 描述 |
|------|------|
| Auth | 认证授权模块 |
| User | 用户信息管理模块 |
| Friend | 好友关系模块 |
| Group | 群组管理模块 |
| File | 文件上传下载模块 |
| Conversation | 会话消息模块 |
| Sync | 数据同步模块 |
| Calling | 音视频通话模块 |
| Version | 版本管理模块 |
| WebSocket | WebSocket接入模块 |

---

## 1. Auth 模块

| HTTP接口 | 功能说明 | 状态 |
|----------|----------|------|
| POST /api/v1/auth/send-code | 发送验证码 | ✅ |
| POST /api/v1/auth/register | 用户注册 | ✅ |
| POST /api/v1/auth/login | 用户登录 | ✅ |
| POST /api/v1/auth/refresh | 刷新Token | ✅ |
| POST /api/v1/auth/logout | 用户登出 | ✅ |
| POST /api/v1/auth/password/change | 修改密码 | ✅ |
| POST /api/v1/auth/password/reset | 重置密码（忘记密码） | ✅ |

---

## 2. User 模块

| HTTP接口 | 功能说明 | 状态 |
|----------|----------|------|
| GET /api/v1/users/me | 获取个人资料 | ✅ |
| PUT /api/v1/users/me | 更新个人资料 | ✅ |
| POST /api/v1/users/me/phone/bind | 绑定手机号 | ✅ |
| POST /api/v1/users/me/phone/change | 更换手机号 | ✅ |
| POST /api/v1/users/me/email/bind | 绑定邮箱 | ✅ |
| POST /api/v1/users/me/email/change | 更换邮箱 | ✅ |
| GET /api/v1/users/:userId | 获取指定用户信息 | ✅ |
| GET /api/v1/users/search | 搜索用户 | ✅ |
| GET /api/v1/users/me/settings | 获取用户设置 | ✅ |
| PUT /api/v1/users/me/settings | 更新用户设置 | ✅ |
| POST /api/v1/users/me/qrcode/refresh | 刷新二维码 | ✅ |
| GET /api/v1/users/qrcode | 通过二维码获取用户 | ✅ |
| POST /api/v1/users/me/push-token | 更新推送Token | ✅ |

---

## 3. Friend 模块

| HTTP接口 | 功能说明 | 状态 |
|----------|----------|------|
| GET /api/v1/friends | 获取好友列表 | ✅ |
| GET /api/v1/friends/requests | 获取好友申请列表 | ✅ |
| POST /api/v1/friends/requests | 发送好友申请 | ✅ |
| PUT /api/v1/friends/requests/:id | 处理好友申请 | ✅ |
| DELETE /api/v1/friends/:id | 删除好友 | ✅ |
| PUT /api/v1/friends/:id/remark | 修改好友备注 | ✅ |
| GET /api/v1/friends/blacklist | 获取黑名单 | ✅ |
| POST /api/v1/friends/blacklist | 添加到黑名单 | ✅ |
| DELETE /api/v1/friends/blacklist/:id | 从黑名单移除 | ✅ |

---

## 4. Group 模块

| HTTP接口 | 功能说明 | 状态 |
|----------|----------|------|
| POST /api/v1/groups | 创建群组 | ✅ |
| GET /api/v1/groups | 获取我的群组列表 | ✅ |
| GET /api/v1/groups/:id | 获取群组信息 | ✅ |
| PUT /api/v1/groups/:id | 更新群组信息 | ✅ |
| DELETE /api/v1/groups/:id | 解散群组 | ✅ |
| GET /api/v1/groups/:id/members | 获取群成员列表 | ✅ |
| POST /api/v1/groups/:id/members | 邀请成员入群 | ✅ |
| DELETE /api/v1/groups/:id/members/:userId | 移除群成员 | ✅ |
| PUT /api/v1/groups/:id/members/:userId/role | 更新成员角色 | ✅ |
| PUT /api/v1/groups/:id/nickname | 修改群昵称 | ✅ |
| POST /api/v1/groups/:id/quit | 退出群组 | ✅ |
| POST /api/v1/groups/:id/transfer | 转让群主 | ✅ |
| POST /api/v1/groups/:id/join | 申请入群 | ✅ |
| GET /api/v1/groups/:id/requests | 获取入群申请列表 | ✅ |
| PUT /api/v1/groups/:id/requests/:requestId | 处理入群申请 | ✅ |

---

## 5. File 模块

| HTTP接口 | 功能说明 | 状态 |
|----------|----------|------|
| POST /api/v1/files/upload-token | 生成上传Token | ✅ |
| POST /api/v1/files/:fileId/complete | 完成文件上传 | ✅ |
| GET /api/v1/files/:fileId/download | 生成下载URL | ✅ |
| GET /api/v1/files/:fileId | 获取文件信息 | ✅ |
| DELETE /api/v1/files/:fileId | 删除文件 | ✅ |
| GET /api/v1/files | 获取文件列表 | ✅ |

---

## 6. Conversation 模块

| HTTP接口 | 功能说明 | 状态 |
|----------|----------|------|
| GET /api/v1/conversations | 获取会话列表 | ✅ |
| GET /api/v1/conversations/unread/total | 获取总未读数 | ✅ |
| GET /api/v1/conversations/:conversationId | 获取会话详情 | ✅ |
| DELETE /api/v1/conversations/:conversationId | 删除会话 | ✅ |
| PUT /api/v1/conversations/:conversationId/pin | 置顶会话 | ✅ |
| PUT /api/v1/conversations/:conversationId/mute | 静音会话 | ✅ |
| POST /api/v1/conversations/:conversationId/read | 标记已读 | ✅ |

---

## 7. Sync 模块

| HTTP接口 | 功能说明 | 状态 |
|----------|----------|------|
| POST /api/v1/sync | 全量同步 | ✅ |
| POST /api/v1/sync/messages | 消息同步 | ✅ |

---

## 8. Calling 模块

### 一对一通话

| HTTP接口 | 功能说明 | 状态 |
|----------|----------|------|
| POST /api/v1/calling/calls | 发起呼叫 | ✅ |
| GET /api/v1/calling/calls | 获取通话记录 | ✅ |
| GET /api/v1/calling/calls/:callId | 获取通话会话 | ✅ |
| POST /api/v1/calling/calls/:callId/join | 加入通话 | ✅ |
| POST /api/v1/calling/calls/:callId/reject | 拒绝通话 | ✅ |
| POST /api/v1/calling/calls/:callId/end | 结束通话 | ✅ |

### 会议室

| HTTP接口 | 功能说明 | 状态 |
|----------|----------|------|
| POST /api/v1/calling/meetings | 创建会议室 | ✅ |
| GET /api/v1/calling/meetings | 获取会议室列表 | ✅ |
| GET /api/v1/calling/meetings/:roomId | 获取会议室详情 | ✅ |
| POST /api/v1/calling/meetings/:roomId/join | 加入会议室 | ✅ |
| POST /api/v1/calling/meetings/:roomId/end | 结束会议室 | ✅ |

---

## 9. Version 模块

| HTTP接口 | 功能说明 | 状态 |
|----------|----------|------|
| GET /api/v1/versions/check | 检查版本更新 | ✅ |
| GET /api/v1/versions/latest | 获取最新版本 | ✅ |
| GET /api/v1/versions/list | 获取版本列表 | ✅ |
| POST /api/v1/versions/report | 上报版本信息 | ✅ |

---

## 10. WebSocket 模块

| HTTP接口 | 功能说明 | 状态 |
|----------|----------|------|
| GET /api/v1/ws | WebSocket接入点 | ✅ |

---

## 11. 其他

| HTTP接口 | 功能说明 | 状态 |
|----------|----------|------|
| GET /health | 健康检查 | ✅ |

---

## 统计汇总

| 模块 | 功能数量 |
|------|----------|
| Auth | 6 |
| User | 13 |
| Friend | 9 |
| Group | 14 |
| File | 6 |
| Conversation | 7 |
| Sync | 2 |
| Calling | 11 |
| Version | 4 |
| WebSocket | 1 |
| 其他 | 1 |
| **总计** | **74** |