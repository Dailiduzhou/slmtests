#!/bin/bash
B=http://117.72.42.214:8080
JWT=$(curl -sS -m 15 $B/api/ticket | grep -o 'eyJ[A-Za-z0-9_.-]*' | head -1)
echo "JWT=$JWT"
for s in 1 3 5 7 9 2 4 6 8 10; do
  echo "--- POST seq=$s ---"
  curl -sS -m 15 -X POST -H "Authorization: Bearer $JWT" -H "Content-Type: application/json" -d "{\"seq\": $s}" $B/api/challenge
  echo
done
