#!/bin/bash
# p3 测试脚本:每组示例独立运行(重新启动程序),逐行比较 JSON(字段顺序不作要求)。
cd "$(dirname "$0")"

check() {
  local name="$1" input="$2" expected="$3"
  local got
  got=$(go run . <<< "$input")
  if python3 - "$got" "$expected" <<'PY'
import json, sys
got = [json.loads(l) for l in sys.argv[1].splitlines() if l.strip()]
exp = [json.loads(l) for l in sys.argv[2].splitlines() if l.strip()]
sys.exit(0 if got == exp else 1)
PY
  then
    echo "PASS: $name"
  else
    echo "FAIL: $name"
    echo "--- got ---"
    echo "$got"
    echo "--- expected ---"
    echo "$expected"
  fi
}

# 示例一:收信、重试与关闭
check "示例一" '{"op":"MEOW","client":"lemon"}
{"op":"SEND","client":"lemon","seq":1,"data":"小鱼干"}
{"op":"SEND","client":"lemon","seq":1,"data":"小鱼干"}
{"op":"SEND","client":"lemon","seq":1,"data":"罐头"}
{"op":"SEND","client":"lemon","seq":3,"data":"晚安"}
{"op":"SEND","client":"lemon","seq":2,"data":"晚安"}
{"op":"PEEK","client":"lemon"}
{"op":"SEAL","client":"lemon"}
{"op":"PEEK","client":"lemon"}
{"op":"SEND","client":"lemon","seq":3,"data":"早安"}
{"op":"MEOW","client":"lemon"}
{"op":"SEAL","client":"lemon"}' '{"code":"PURR"}
{"code":"ACK","seq":1}
{"code":"ACK","seq":1}
{"code":"ERR_CONFLICT"}
{"code":"ERR_OUT_OF_ORDER"}
{"code":"ACK","seq":2}
{"code":"PURR","state":"AWAKE","count":2,"next_seq":3}
{"code":"PURR"}
{"code":"PURR","state":"ASLEEP","count":2,"next_seq":3}
{"code":"ERR_ASLEEP"}
{"code":"ERR_ALREADY_EXISTS"}
{"code":"PURR"}'

# 示例二:不同信箱互不影响
check "示例二" '{"op":"PEEK","client":"ghost"}
{"op":"MEOW","client":"lemon"}
{"op":"MEOW","client":"orange"}
{"op":"SEND","client":"lemon","seq":1,"data":"喵"}
{"op":"SEND","client":"orange","seq":1,"data":"呼噜"}
{"op":"SEAL","client":"lemon"}
{"op":"PEEK","client":"orange"}' '{"code":"ERR_NO_SESSION"}
{"code":"PURR"}
{"code":"PURR"}
{"code":"ACK","seq":1}
{"code":"ACK","seq":1}
{"code":"PURR"}
{"code":"PURR","state":"AWAKE","count":1,"next_seq":2}'
