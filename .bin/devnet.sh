#!/bin/sh

i=0
chainId=""

while true; do
  i=$((i + 1))
  out=$(curl -L localhost:5050/rpc \
  -H 'Content-Type: application/json' \
  -d'{
    "jsonrpc": "2.0",
    "method": "starknet_chainId",
    "params": [],
    "id": '${i}'
  }' \
  2>/dev/null)
  result=$?
  if [ $i -gt 10 -o $result -eq 0 ]; then
    chainId=$(echo $out | jq -r '.result')
    break
  fi
  echo "we will continue in a while, loop ${i}..."
  i=$((i + 1))
  sleep 3
done

if [ "$chainId" = "0x534e5f474f45524c49" ]; then
  echo "devnet is running with chainId $chainId..."
  exit 0
fi

echo "could not check devnet, chainId=$chainId; fail!!!"
exit 1
