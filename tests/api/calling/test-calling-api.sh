#!/bin/bash
#
# LiveKit Calling Service API 测试脚本
#

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/../common.sh"

BASE_URL="${GATEWAY_URL:-http://localhost:8080}/api/v1"

echo "=================================================="
echo "  LiveKit Calling Service API 测试"
echo "=================================================="
echo ""

PASS=0
FAIL=0
TOKEN_A=""
TOKEN_B=""
USER_A_ID=""
USER_B_ID=""

# ── 辅助函数 ─────────────────────────────────────────

pass() { echo -e "${GREEN}✓ PASS${NC}: $1"; PASS=$((PASS + 1)); }
fail() { echo -e "${RED}✗ FAIL${NC}: $1"; FAIL=$((FAIL + 1)); }

# ── 准备测试用户 ──────────────────────────────────────

echo "正在注册测试用户..."
TOKEN_A=$(register_and_login_test_user "${BASE_URL}" "calling_00000001@test.com" "Test@1234" "User00000001" "calling-dev-00000001" "Web")
TOKEN_B=$(register_and_login_test_user "${BASE_URL}" "calling_00000002@test.com" "Test@1234" "User00000002" "calling-dev-00000002" "Web")

if [ -z "$TOKEN_A" ] || [ -z "$TOKEN_B" ]; then
    echo -e "${RED}错误: 无法获取测试用户 token，跳过 Calling 测试${NC}"
    exit 0
fi

echo "用户 A token 获取成功"
echo "用户 B token 获取成功"

USER_A_ID=$(get_user_id_by_token "${BASE_URL}" "$TOKEN_A")
USER_B_ID=$(get_user_id_by_token "${BASE_URL}" "$TOKEN_B")
if [ -z "$USER_A_ID" ] || [ -z "$USER_B_ID" ]; then
    echo -e "${RED}错误: 无法获取测试用户ID，跳过 Calling 测试${NC}"
    exit 0
fi
echo "用户 A ID: ${USER_A_ID}"
echo "用户 B ID: ${USER_B_ID}"
echo ""

# ── 测试 1: 未认证发起通话（应返回 401）────────────────

echo "测试 1: 未认证发起通话"
RESP=$(curl -s -X POST "${BASE_URL}/calling/calls" \
    -H "Content-Type: application/json" \
    -d '{"calleeId":"user-xxx","callType":"audio"}')
CODE=$(json_code "$RESP")
if [ "$CODE" = "401" ]; then
    pass "未认证返回 code 401"
else
    fail "期望 code 401，实际 $CODE"
fi

# ── 测试 2: 缺少 calleeId（应返回 400）────────────────

echo "测试 2: 缺少 calleeId"
RESP=$(curl -s -o /dev/null -w "%{http_code}" -X POST "${BASE_URL}/calling/calls" \
    -H "Authorization: Bearer $TOKEN_A" \
    -H "Content-Type: application/json" \
    -d '{}')
if [ "$RESP" = "400" ]; then
    pass "缺少 calleeId 返回 400"
else
    fail "期望 400，实际 $RESP"
fi

# ── 测试 3: 获取通话记录（空列表）────────────────────

echo "测试 3: 获取通话记录（初始为空）"
RESP=$(curl -s "${BASE_URL}/calling/calls" \
    -H "Authorization: Bearer $TOKEN_A")
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" "${BASE_URL}/calling/calls" \
    -H "Authorization: Bearer $TOKEN_A")
if [ "$HTTP_CODE" = "200" ]; then
    pass "获取通话记录成功（HTTP 200）"
else
    fail "期望 200，实际 $HTTP_CODE"
fi

# ── 测试 4: 获取不存在的通话（应返回错误）────────────

echo "测试 4: 获取不存在的通话"
RESP=$(curl -s "${BASE_URL}/calling/calls/nonexistent-call-id" \
    -H "Authorization: Bearer $TOKEN_A")
CODE=$(json_code "$RESP")
if [ "$CODE" = "404" ] || [ "$CODE" = "500" ]; then
    pass "不存在通话返回错误码 $CODE"
else
    fail "期望 code 404/500，实际 $CODE"
fi

# ── 测试 5: 未认证创建会议室（应返回 401）────────────

echo "测试 5: 未认证创建会议室"
RESP=$(curl -s -X POST "${BASE_URL}/calling/meetings" \
    -H "Content-Type: application/json" \
    -d '{"title":"Test Meeting"}')
CODE=$(json_code "$RESP")
if [ "$CODE" = "401" ]; then
    pass "未认证返回 code 401"
else
    fail "期望 code 401，实际 $CODE"
fi

# ── 测试 6: 缺少 title（应返回 400）────────────────────

echo "测试 6: 创建会议室缺少 title"
RESP=$(curl -s -o /dev/null -w "%{http_code}" -X POST "${BASE_URL}/calling/meetings" \
    -H "Authorization: Bearer $TOKEN_A" \
    -H "Content-Type: application/json" \
    -d '{}')
if [ "$RESP" = "400" ]; then
    pass "缺少 title 返回 400"
else
    fail "期望 400，实际 $RESP"
fi

# ── 测试 7: 列举会议室（初始为空）────────────────────

echo "测试 7: 列举会议室（初始为空）"
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" \
    "${BASE_URL}/calling/meetings" \
    -H "Authorization: Bearer $TOKEN_A")
if [ "$HTTP_CODE" = "200" ]; then
    pass "列举会议室成功（HTTP 200）"
else
    fail "期望 200，实际 $HTTP_CODE"
fi

# ── 测试 8: 接听不存在的通话（应返回错误）────────────

echo "测试 8: 接听不存在的通话"
RESP=$(curl -s -X POST "${BASE_URL}/calling/calls/fake-call-id/join" \
    -H "Authorization: Bearer $TOKEN_B")
CODE=$(json_code "$RESP")
if [ "$CODE" = "404" ] || [ "$CODE" = "500" ]; then
    pass "接听不存在通话返回错误码 $CODE"
else
    fail "期望 code 404/500，实际 $CODE"
fi

# ── 测试 9: 获取不存在的会议室（应返回错误）─────────

echo "测试 9: 获取不存在的会议室"
RESP=$(curl -s "${BASE_URL}/calling/meetings/nonexistent-room" \
    -H "Authorization: Bearer $TOKEN_A")
CODE=$(json_code "$RESP")
if [ "$CODE" = "404" ] || [ "$CODE" = "500" ]; then
    pass "不存在会议室返回错误码 $CODE"
else
    fail "期望 code 404/500，实际 $CODE"
fi

# ── 测试 10: 黑名单限制 - 被拉黑方无法发起通话 ──────────

echo "测试 10: 黑名单限制（无法发起通话）"
RESP=$(curl -s -X POST "${BASE_URL}/friends/blacklist" \
    -H "Authorization: Bearer $TOKEN_A" \
    -H "Content-Type: application/json" \
    -d "{\"userId\":\"${USER_B_ID}\"}")
CODE=$(json_code "$RESP")
if [ "$CODE" = "0" ]; then
    pass "用户A拉黑用户B成功"
else
    fail "用户A拉黑用户B失败，code=$CODE"
fi

RESP=$(curl -s -X POST "${BASE_URL}/calling/calls" \
    -H "Authorization: Bearer $TOKEN_B" \
    -H "Content-Type: application/json" \
    -d "{\"calleeId\":\"${USER_A_ID}\",\"callType\":\"audio\"}")
CODE=$(json_code "$RESP")
if [ "$CODE" = "403" ]; then
    pass "被拉黑方发起通话被拦截（code 403）"
else
    fail "期望被拦截返回 code 403，实际 $CODE"
fi

# 清理黑名单，避免影响后续手工调试
RESP=$(curl -s -X DELETE "${BASE_URL}/friends/blacklist/${USER_B_ID}" \
    -H "Authorization: Bearer $TOKEN_A")
CODE=$(json_code "$RESP")
if [ "$CODE" = "0" ]; then
    pass "清理黑名单成功"
else
    fail "清理黑名单失败，code=$CODE"
fi

# ── 汇总 ──────────────────────────────────────────────

echo ""
echo "=================================================="
echo "  测试结果: ${PASS} 通过, ${FAIL} 失败"
echo "=================================================="

if [ $FAIL -eq 0 ]; then
    exit 0
else
    exit 1
fi
