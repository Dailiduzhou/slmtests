#!/bin/bash
# p4 测试脚本:运行 Go 求解器,验证网站展示 🍥 与密钥(挑战成功)。
cd "$(dirname "$0")"
got=$(go run . 2>&1)
echo "$got"
if echo "$got" | grep -q "🍥" && echo "$got" | grep -q "密钥"; then
  echo "PASS: 挑战成功,已找到 🍥 和密钥"
else
  echo "FAIL: 未找到 🍥 或密钥"
  exit 1
fi
