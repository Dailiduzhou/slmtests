#!/bin/bash
# p5 测试脚本:启动 507 预约系统,用 curl 验证预约、取消与使用中规则。
# 端口预检防残留服务;退出时连 go run 的子进程一并杀掉。
cd "$(dirname "$0")"
export PORT=${PORT:-18050}
B=http://127.0.0.1:$PORT

if lsof -tiTCP:$PORT -sTCP:LISTEN >/dev/null 2>&1; then
  echo "FAIL: 端口 $PORT 已被占用,请先释放或用 PORT=xxxx 指定其他端口"
  exit 1
fi

go run . 2>/dev/null &
pid=$!
trap 'pkill -P $pid 2>/dev/null; kill $pid 2>/dev/null' EXIT
ready=0
for i in $(seq 1 100); do
  if curl -s -o /dev/null "$B/api/reservations"; then ready=1; break; fi
  sleep 0.1
done
if [ "$ready" -eq 0 ]; then
  echo "FAIL: 服务未在端口 $PORT 启动"
  exit 1
fi

pass=0 fail=0
check() {
  local name="$1" want="$2" got="$3"
  if [ "$got" = "$want" ]; then
    echo "PASS: $name"
    pass=$((pass+1))
  else
    echo "FAIL: $name expected=[$want] got=[$got]"
    fail=$((fail+1))
  fi
}

# 预约创建:201,状态默认 booking,全新状态 id=1
code=$(curl -s -o /tmp/p5a -w '%{http_code}' -X POST "$B/api/reservations" -d '{"project":"蓝桥杯","slot":"2026-10-17 14:00"}')
check "预约创建返回 201" "201" "$code"
check "预约状态默认 booking" '"status":"booking"' "$(grep -o '"status":"[^"]*"' /tmp/p5a)"
check "全新状态首个 id=1" '"id":1' "$(grep -o '"id":[0-9]*' /tmp/p5a)"

# 同一时间段冲突:409 预约失败
code=$(curl -s -o /tmp/p5b -w '%{http_code}' -X POST "$B/api/reservations" -d '{"project":"传智杯","slot":"2026-10-17 14:00"}')
check "同时间段冲突返回 409" "409" "$code"
check "冲突提示预约失败" "yes" "$(grep -q '"error":"预约失败"' /tmp/p5b && echo yes)"

# 取消预约:200,状态改为 cancelled
code=$(curl -s -o /tmp/p5c -w '%{http_code}' -X POST "$B/api/reservations/1/cancel")
check "取消预约返回 200" "200" "$code"
check "取消后状态 cancelled" '"status":"cancelled"' "$(grep -o '"status":"[^"]*"' /tmp/p5c)"

# 重复取消:409 取消失败
code=$(curl -s -o /tmp/p5d -w '%{http_code}' -X POST "$B/api/reservations/1/cancel")
check "重复取消返回 409" "409" "$code"
check "取消失败提示" "yes" "$(grep -q '"error":"取消失败"' /tmp/p5d && echo yes)"

# 已在使用中的预约不能取消:新建 算法研讨 → use → cancel 应 409
curl -s -o /tmp/p5e -X POST "$B/api/reservations" -d '{"project":"算法研讨","slot":"2026-10-18 09:00"}'
id2=$(grep -o '"id":[0-9]*' /tmp/p5e | cut -d: -f2)
code=$(curl -s -o /tmp/p5f -w '%{http_code}' -X POST "$B/api/reservations/$id2/use")
check "标记使用中返回 200" "200" "$code"
check "use 后状态 used" '"status":"used"' "$(grep -o '"status":"[^"]*"' /tmp/p5f)"
code=$(curl -s -o /tmp/p5g -w '%{http_code}' -X POST "$B/api/reservations/$id2/cancel")
check "已在使用中取消返回 409" "409" "$code"
check "使用中取消提示取消失败" "yes" "$(grep -q '"error":"取消失败"' /tmp/p5g && echo yes)"

echo "p5: pass=$pass fail=$fail"
[ "$fail" -eq 0 ]
