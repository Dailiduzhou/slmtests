#!/bin/bash
# p1 测试脚本:go run . < 输入 > 实际输出,与期望对比。
cd "$(dirname "$0")"
pass=0 fail=0
check() {
  local in="$1" want="$2"
  local got
  got=$(go run . <<< "$in")
  if [ "$got" = "$want" ]; then
    echo "PASS: input=[$in] output=[$got]"
    pass=$((pass+1))
  else
    echo "FAIL: input=[$in] expected=[$want] got=[$got]"
    fail=$((fail+1))
  fi
}
# 示例 1: 1/3 位于编号 5,排数 5/2=2
check "1 3" "2"
# 示例 2: 2/4 约分为 1/2,编号 1,排数 0
check "2 4" "0"
# 示例 3: 大数,int64 内精确
check "987654321 987654322" "975461057789971040"
# 边界: 1/1 编号 0
check "1 1" "0"
echo "p1: pass=$pass fail=$fail"
[ "$fail" -eq 0 ]
