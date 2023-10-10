#!/bin/sh

i=0
version=""
export TEST_ACCOUNT_ADDRESS=0x2a5338be59c7d134754add779b427db76adb734148c075efcc072a7e79858a7


while true; do
  i=$((i + 1))
  curl --fail localhost:5050/is_alive 2>/dev/null 2>&1
  result=$?
  if [ $result -eq 0 ]; then
    sleep 5
    curl --fail -H 'Content-Type: application/json' -XPOST http://localhost:5050/mint \
      -d '{ "address": "'${ACCOUNT_ADDRESS}'", "amount": 1000000000000000}'
    exit 0
  fi
  if [ $i -gt 10 ]; then
    break
  fi
  echo "we will continue in a while, loop ${i}..."
  i=$((i + 1))
  sleep 3
done

echo "could not check devnet is_alive; fail!!!"
exit 1

