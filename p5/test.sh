#!/bin/bash
# p5 测试脚本:启动 507 预约系统,用 curl 验证预约与取消预约。
cd "$(dirname "$0")"
export PORT=${PORT:-18050}
go run . 2>/dev/null &
pid=$!
trap 'kill $pid 2>/dev/null' EXIT
for i in $(seq 1 50); do
  curl -s -o /dev/null "http://127.0.0.1:$PORT/api/reservations" && break
  sleep 0.1
done
B=http://127.0.0.1:$PORT
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

code=$(curl -s -o /tmp/p5a -w '%{http_code}' -X POST "$B/api/reservations" -d '{"project":"蓝桥杯","slot":"2026-10-17 14:00"}')
check "预约创建返回 201" "201" "$code"
check "预约状态默认 booking" '"status":"booking"' "$(grep -o '"status":"[^"]*"' /tmp/p5a)"

code=$(curl -s -o /tmp/p5b -w '%{http_code}' -X POST "$B/api/reservations" -d '{"project":"传智杯","slot":"2026-10-17 14:00"}')
check "同时间段冲突返回 409" "409" "$code"
check "冲突提示预约失败" "yes" "$(grep -q '"error":"预约失败"' /tmp/p5b && echo yes)"

code=$(curl -s -o /tmp/p5c -w '%{http_code}' -X POST "$B/api/reservations/1/cancel")
check "取消预约返回 200" "200" "$code"
check "取消后状态 cancelled" '"status":"cancelled"' "$(grep -o '"status":"[^"]*"' /tmp/p5c)"

code=$(curl -s -o /tmp/p5d -w '%{http_code}' -X POST "$B/api/reservations/1/cancel")
check "重复取消返回 409" "409" "$code"
check "取消失败提示" "yes" "$(grep -q '"error":"取消失败"' /tmp/p5d && echo yes)"

curl -s -o /dev/null -X POST "$B/api/reservations" -d '{"project":"算法研讨","slot":"2026-10-18 09:00"}'
curl -s -o /dev/null -X POST "$B/api/reservations/3/use"
code=$(curl -s -o /tmp/p5e -w '%{http_code}' -X POST "$B/api/reservations/3/cancel")
check "已在使用中取消返回 409" "409" "$code"

echo "p5: pass=$pass fail=$fail"
[ "$fail" -eq 0 ]
