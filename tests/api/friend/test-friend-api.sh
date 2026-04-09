#!/bin/bash
#
# Friend Service HTTP API 测试脚本
# 用于测试好友管理相关的 HTTP 接口
#

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/../common.sh"

# 配置
GATEWAY_URL="${GATEWAY_URL:-http://localhost:8080}"
API_BASE="${GATEWAY_URL}/api/v1"

# 测试数据
TIMESTAMP=$(date +%s)
TEST_PHONE_1="138${TIMESTAMP:(-8)}"
TEST_PHONE_2="139${TIMESTAMP:(-8)}"
TEST_EMAIL_1="user1_${TIMESTAMP}@example.com"
TEST_EMAIL_2="user2_${TIMESTAMP}@example.com"
TEST_PASSWORD="Test@123456"
TEST_DEVICE_ID="test-device-${TIMESTAMP}"

# 全局变量
USER1_TOKEN=""
USER2_TOKEN=""
USER1_ID=""
USER2_ID=""
FRIEND_REQUEST_ID=""
USER2_CONVERSATION_ID=""
POST_UNBLOCK_REQUEST_ID=""

# ========================================
# 准备工作：创建测试用户
# ========================================

setup_test_users() {
    print_header "准备测试用户"

    # 注册用户1
    print_info "注册用户1: ${TEST_EMAIL_1}"
    local response1
    response1=$(register_test_user "${API_BASE}" "${TEST_EMAIL_1}" "${TEST_PASSWORD}" "测试用户1_${TIMESTAMP}" "${TEST_DEVICE_ID}_1" "iOS")
    if check_response "$response1"; then
        USER1_ID=$(extract_user_id "$response1")
        USER1_TOKEN=$(extract_access_token "$response1")
        print_success "用户1注册成功 (ID: ${USER1_ID})"
    else
        print_error "用户1注册失败"
        return 1
    fi

    sleep 1

    # 注册用户2
    print_info "注册用户2: ${TEST_EMAIL_2}"
    local response2
    response2=$(register_test_user "${API_BASE}" "${TEST_EMAIL_2}" "${TEST_PASSWORD}" "测试用户2_${TIMESTAMP}" "${TEST_DEVICE_ID}_2" "iOS")
    if check_response "$response2"; then
        USER2_ID=$(extract_user_id "$response2")
        USER2_TOKEN=$(extract_access_token "$response2")
        print_success "用户2注册成功 (ID: ${USER2_ID})"
    else
        print_error "用户2注册失败"
        return 1
    fi
}

# ========================================
# 测试用例
# ========================================

# 1. 发送好友申请
test_send_friend_request() {
    print_header "1. 发送好友申请"

    local data=$(cat <<EOF
{
    "userId": "${USER2_ID}",
    "message": "你好，我想加你为好友",
    "source": "search"
}
EOF
)

    print_info "用户1向用户2发送好友申请"

    local response=$(http_post "${API_BASE}/friends/requests" "$data" "$USER1_TOKEN")
    print_info "响应: $response"

    if check_response "$response"; then
        # 注意：protobuf 的 request_id 在 JSON 中是 requestId (驼峰)
        FRIEND_REQUEST_ID=$(echo "$response" | jq -r '.data.requestId // .data.request_id // empty')
        local auto_accepted=$(echo "$response" | jq -r '.data.autoAccepted // .data.auto_accepted // false')

        if [ -z "$FRIEND_REQUEST_ID" ] || [ "$FRIEND_REQUEST_ID" = "null" ]; then
            print_error "无法获取申请ID，响应数据: $(echo "$response" | jq -r '.data')"
            return 1
        fi

        print_success "发送好友申请成功"
        print_info "申请ID: ${FRIEND_REQUEST_ID}"
        print_info "自动接受: ${auto_accepted}"
        return 0
    else
        return 1
    fi
}

# 2. 获取收到的好友申请
test_get_received_requests() {
    print_header "2. 获取收到的好友申请"

    print_info "用户2获取收到的好友申请"

    local response=$(http_get "${API_BASE}/friends/requests?type=received" "$USER2_TOKEN")
    print_info "响应: $response"

    if check_response "$response"; then
        local total=$(echo "$response" | jq -r '.data.total // 0')
        print_success "获取好友申请成功"
        print_info "收到 ${total} 个好友申请"

        # 备用方案：如果之前没有获取到 FRIEND_REQUEST_ID，从列表中提取
        if [ -z "$FRIEND_REQUEST_ID" ] || [ "$FRIEND_REQUEST_ID" = "null" ]; then
            FRIEND_REQUEST_ID=$(echo "$response" | jq -r '.data.requests[0].id // empty')
            if [ -n "$FRIEND_REQUEST_ID" ] && [ "$FRIEND_REQUEST_ID" != "null" ]; then
                print_info "从申请列表中获取到申请ID: ${FRIEND_REQUEST_ID}"
            fi
        fi

        return 0
    else
        return 1
    fi
}

# 3. 获取发送的好友申请
test_get_sent_requests() {
    print_header "3. 获取发送的好友申请"

    print_info "用户1获取发送的好友申请"

    local response=$(http_get "${API_BASE}/friends/requests?type=sent" "$USER1_TOKEN")
    print_info "响应: $response"

    if check_response "$response"; then
        local total=$(echo "$response" | jq -r '.data.total // 0')
        print_success "获取发送的好友申请成功"
        print_info "发送了 ${total} 个好友申请"
        return 0
    else
        return 1
    fi
}

# 4. 接受好友申请
test_accept_friend_request() {
    print_header "4. 接受好友申请"

    # 检查 FRIEND_REQUEST_ID 是否有效
    if [ -z "$FRIEND_REQUEST_ID" ] || [ "$FRIEND_REQUEST_ID" = "null" ]; then
        print_error "申请ID无效，跳过此测试"
        print_info "提示：请确保前面的测试成功执行"
        return 1
    fi

    local data=$(cat <<EOF
{
    "action": "accept"
}
EOF
)

    print_info "用户2接受好友申请 (ID: ${FRIEND_REQUEST_ID})"

    local response=$(http_put "${API_BASE}/friends/requests/${FRIEND_REQUEST_ID}" "$data" "$USER2_TOKEN")
    print_info "响应: $response"

    if check_response "$response"; then
        print_success "接受好友申请成功"
        return 0
    else
        return 1
    fi
}

# 5. 获取好友列表
test_get_friend_list() {
    print_header "5. 获取好友列表"

    # 用户1获取好友列表
    print_info "用户1获取好友列表"
    local response1=$(http_get "${API_BASE}/friends" "$USER1_TOKEN")
    print_info "响应: $response1"

    if check_response "$response1"; then
        local total1=$(echo "$response1" | jq -r '.data.total // 0')
        print_success "用户1获取好友列表成功 (共 ${total1} 个好友)"
    else
        return 1
    fi

    # 用户2获取好友列表
    print_info "用户2获取好友列表"
    local response2=$(http_get "${API_BASE}/friends" "$USER2_TOKEN")
    print_info "响应: $response2"

    if check_response "$response2"; then
        local total2=$(echo "$response2" | jq -r '.data.total // 0')
        print_success "用户2获取好友列表成功 (共 ${total2} 个好友)"
        return 0
    else
        return 1
    fi
}

# 5.1 准备用户2到用户1的单聊会话ID
prepare_user2_single_conversation() {
    print_header "5.1 准备单聊会话ID"

    local max_retry=5
    local i=1
    while [ $i -le $max_retry ]; do
        local response=$(http_get "${API_BASE}/conversations?limit=100" "$USER2_TOKEN")
        print_info "第 ${i} 次查询会话列表"

        if check_response "$response"; then
            USER2_CONVERSATION_ID=$(echo "$response" | jq -r --arg uid "$USER1_ID" '
                .data.conversations[]? |
                select((.conversationType // .conversation_type) == "single" and (.targetId // .target_id) == $uid) |
                (.conversationId // .conversation_id)
            ' | head -n 1)

            if [ -n "$USER2_CONVERSATION_ID" ] && [ "$USER2_CONVERSATION_ID" != "null" ]; then
                print_success "获取会话ID成功: ${USER2_CONVERSATION_ID}"
                return 0
            fi
        fi

        sleep 1
        i=$((i + 1))
    done

    print_error "未找到用户2与用户1的单聊会话ID，后续消息拦截测试将失败"
    return 1
}

# 6. 更新好友备注
test_update_friend_remark() {
    print_header "6. 更新好友备注"

    local data=$(cat <<EOF
{
    "remark": "我的好朋友"
}
EOF
)

    print_info "用户1更新用户2的备注"

    local response=$(http_put "${API_BASE}/friends/${USER2_ID}/remark" "$data" "$USER1_TOKEN")
    print_info "响应: $response"

    if check_response "$response"; then
        print_success "更新好友备注成功"
        return 0
    else
        return 1
    fi
}

# 7. 增量同步好友列表
test_incremental_sync() {
    print_header "7. 增量同步好友列表"

    # 使用过去的时间戳（5分钟前）来测试增量同步
    # 这样可以捕获刚才创建的好友关系
    local last_time=$(($(date +%s) - 300))
    print_info "使用时间戳进行增量同步: ${last_time}"

    local response=$(http_get "${API_BASE}/friends?lastUpdateTime=${last_time}" "$USER1_TOKEN")
    print_info "响应: $response"

    if check_response "$response"; then
        local total=$(echo "$response" | jq -r '.data.total // 0')
        print_success "增量同步成功"
        print_info "更新了 ${total} 个好友"
        return 0
    else
        return 1
    fi
}

# 8. 添加到黑名单
test_add_to_blacklist() {
    print_header "8. 添加到黑名单"

    local data=$(cat <<EOF
{
    "userId": "${USER2_ID}"
}
EOF
)

    print_info "用户1将用户2添加到黑名单"

    local response=$(http_post "${API_BASE}/friends/blacklist" "$data" "$USER1_TOKEN")
    print_info "响应: $response"

    if check_response "$response"; then
        print_success "添加到黑名单成功"
        return 0
    else
        return 1
    fi
}

# 9. 获取黑名单
test_get_blacklist() {
    print_header "9. 获取黑名单"

    print_info "用户1获取黑名单"

    local response=$(http_get "${API_BASE}/friends/blacklist" "$USER1_TOKEN")
    print_info "响应: $response"

    if check_response "$response"; then
        local total=$(echo "$response" | jq -r '.data.total // 0')
        print_success "获取黑名单成功"
        print_info "黑名单中有 ${total} 个用户"
        return 0
    else
        return 1
    fi
}

# 10. 验证拉黑后自动解除好友关系
test_blacklist_auto_remove_friend() {
    print_header "10. 验证拉黑后自动解除好友关系"

    local response1=$(http_get "${API_BASE}/friends" "$USER1_TOKEN")
    print_info "用户1好友列表响应: $response1"
    if ! check_response "$response1"; then
        return 1
    fi

    local total1=$(echo "$response1" | jq -r '.data.total // 0')
    if [ "$total1" -ne 0 ]; then
        print_error "用户1好友列表应为空，实际 total=${total1}"
        return 1
    fi

    local response2=$(http_get "${API_BASE}/friends" "$USER2_TOKEN")
    print_info "用户2好友列表响应: $response2"
    if ! check_response "$response2"; then
        return 1
    fi

    local total2=$(echo "$response2" | jq -r '.data.total // 0')
    if [ "$total2" -ne 0 ]; then
        print_error "用户2好友列表应为空，实际 total=${total2}"
        return 1
    fi

    print_success "拉黑后已自动解除双方好友关系"
    return 0
}

# 11. 验证黑名单限制：无法发送消息
test_blacklist_blocks_message() {
    print_header "11. 验证黑名单限制：无法发送消息"

    if [ -z "$USER2_CONVERSATION_ID" ] || [ "$USER2_CONVERSATION_ID" = "null" ]; then
        print_error "缺少会话ID，无法执行消息拦截测试"
        return 1
    fi

    local data=$(cat <<EOF
{
    "conversation_id": "${USER2_CONVERSATION_ID}",
    "content_type": "text",
    "content": "{\"text\":\"blacklist block test\"}",
    "local_id": "local-blacklist-${TIMESTAMP}"
}
EOF
)

    local response=$(http_post "${API_BASE}/messages" "$data" "$USER2_TOKEN")
    print_info "响应: $response"

    if check_response_fail "$response" && check_fail_code "$response" "403"; then
        print_success "黑名单消息拦截生效（403）"
        return 0
    fi
    return 1
}

# 12. 验证黑名单限制：无法发起音视频通话
test_blacklist_blocks_call() {
    print_header "12. 验证黑名单限制：无法发起音视频通话"

    local data=$(cat <<EOF
{
    "calleeId": "${USER1_ID}",
    "callType": "audio"
}
EOF
)

    local response=$(http_post "${API_BASE}/calling/calls" "$data" "$USER2_TOKEN")
    print_info "响应: $response"

    if check_response_fail "$response" && check_fail_code "$response" "403"; then
        print_success "黑名单通话拦截生效（403）"
        return 0
    fi
    return 1
}

# 13. 验证黑名单限制：无法查看用户资料
test_blacklist_blocks_user_info() {
    print_header "13. 验证黑名单限制：无法查看用户资料"

    local response=$(http_get "${API_BASE}/users/${USER1_ID}" "$USER2_TOKEN")
    print_info "响应: $response"

    if check_response_fail "$response" && check_fail_code "$response" "403"; then
        print_success "黑名单资料访问限制生效（403）"
        return 0
    fi
    return 1
}

# 14. 从黑名单移除
test_remove_from_blacklist() {
    print_header "14. 从黑名单移除"

    print_info "用户1将用户2从黑名单移除"

    local response=$(http_delete "${API_BASE}/friends/blacklist/${USER2_ID}" "$USER1_TOKEN")
    print_info "响应: $response"

    if check_response "$response"; then
        print_success "从黑名单移除成功"
        return 0
    else
        return 1
    fi
}

# 15. 验证移除黑名单后仍非好友（不会自动恢复）
test_verify_not_friend_after_unblock() {
    print_header "15. 验证移除黑名单后仍非好友"

    local response1=$(http_get "${API_BASE}/friends" "$USER1_TOKEN")
    print_info "用户1好友列表响应: $response1"
    if ! check_response "$response1"; then
        return 1
    fi

    local total1=$(echo "$response1" | jq -r '.data.total // 0')
    if [ "$total1" -ne 0 ]; then
        print_error "用户1好友列表应为空，实际 total=${total1}"
        return 1
    fi

    local response2=$(http_get "${API_BASE}/friends" "$USER2_TOKEN")
    print_info "用户2好友列表响应: $response2"
    if ! check_response "$response2"; then
        return 1
    fi

    local total2=$(echo "$response2" | jq -r '.data.total // 0')
    if [ "$total2" -ne 0 ]; then
        print_error "用户2好友列表应为空，实际 total=${total2}"
        return 1
    fi

    print_success "移除黑名单后好友关系未自动恢复（符合预期）"
    return 0
}

# 16. 验证移除黑名单后可重新发起好友申请
test_send_friend_request_after_unblock() {
    print_header "16. 验证移除黑名单后可重新发起好友申请"

    local data=$(cat <<EOF
{
    "userId": "${USER1_ID}",
    "message": "解除拉黑后重新申请好友",
    "source": "search"
}
EOF
)

    local response=$(http_post "${API_BASE}/friends/requests" "$data" "$USER2_TOKEN")
    print_info "响应: $response"

    if ! check_response "$response"; then
        return 1
    fi

    POST_UNBLOCK_REQUEST_ID=$(echo "$response" | jq -r '.data.requestId // .data.request_id // empty')
    if [ -z "$POST_UNBLOCK_REQUEST_ID" ] || [ "$POST_UNBLOCK_REQUEST_ID" = "null" ]; then
        print_error "未获取到重新申请的 requestId"
        return 1
    fi

    print_success "移除黑名单后可重新发起好友申请"
    print_info "新的申请ID: ${POST_UNBLOCK_REQUEST_ID}"
    return 0
}

# ========================================
# 主函数
# ========================================

main() {
    echo -e "${GREEN}"
    echo "╔═══════════════════════════════════════════╗"
    echo "║   Friend Service API 测试脚本             ║"
    echo "╚═══════════════════════════════════════════╝"
    echo -e "${NC}"

    echo "测试环境: ${GATEWAY_URL}"
    echo "开始时间: $(date '+%Y-%m-%d %H:%M:%S')"
    echo ""

    # 检查依赖
    if ! command -v jq &> /dev/null; then
        print_error "需要安装 jq 工具: apt-get install jq 或 brew install jq"
        exit 1
    fi

    # 准备测试用户
    setup_test_users || exit 1

    # 执行测试
    local failed=0

    test_send_friend_request || ((failed++))
    sleep 1
    test_get_received_requests || ((failed++))
    test_get_sent_requests || ((failed++))
    test_accept_friend_request || ((failed++))
    sleep 1
    test_get_friend_list || ((failed++))
    prepare_user2_single_conversation || ((failed++))
    test_update_friend_remark || ((failed++))
    sleep 1
    test_incremental_sync || ((failed++))
    test_add_to_blacklist || ((failed++))
    test_get_blacklist || ((failed++))
    test_blacklist_auto_remove_friend || ((failed++))
    test_blacklist_blocks_message || ((failed++))
    test_blacklist_blocks_call || ((failed++))
    test_blacklist_blocks_user_info || ((failed++))
    test_remove_from_blacklist || ((failed++))
    test_verify_not_friend_after_unblock || ((failed++))
    test_send_friend_request_after_unblock || ((failed++))

    # 输出测试结果
    echo ""
    echo -e "${YELLOW}========================================${NC}"
    echo -e "${YELLOW}测试结果${NC}"
    echo -e "${YELLOW}========================================${NC}"
    echo "结束时间: $(date '+%Y-%m-%d %H:%M:%S')"

    if [ $failed -eq 0 ]; then
        echo -e "${GREEN}所有测试通过! ✓${NC}"
        exit 0
    else
        echo -e "${RED}失败测试数: ${failed} ✗${NC}"
        exit 1
    fi
}

# 运行主函数
main "$@"
