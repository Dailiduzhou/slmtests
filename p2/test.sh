#!/bin/bash
# p2 测试脚本:go run . < 输入 > 实际输出,与期望对比。
cd "$(dirname "$0")"
run() {
  local in="$1" want="$2"
  local got
  got=$(go run . <<< "$in")
  if [ "$got" = "$want" ]; then
    echo "PASS: input=[$in] output=[$got]"
  else
    echo "FAIL: input=[$in] expected=[$want] got=[$got]"
  fi
}
# 示例 1: M=30 N=21 -> 7
run "30 21" "7"
# 示例 2: M=10000 N=1 -> 0 (起点即终点)
run "10000 1" "0"
# 边界: M=1 N=1 -> 0
run "1 1" "0"
# 边界: M=2 N=2 -> 1 (走廊 1->2)
run "2 2" "1"
