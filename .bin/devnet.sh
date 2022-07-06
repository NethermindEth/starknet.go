#!/bin/sh

i=0
version=""

while true; do
  i=$((i + 1))
  out=$(curl -L localhost:5050/rpc \
  -H 'Content-Type: application/json' \
  -d'{
    "jsonrpc": "2.0",
    "method": "starknet_protocolVersion",
    "params": [],
    "id": '${i}'
  }' \
  2>/dev/null)
  result=$?
  if [ $i -gt 10 -o $result -eq 0 ]; then
    version=$(echo $out | jq -r '.result')
    break
  fi
  echo "we will continue in a while, loop ${i}..."
  i=$((i + 1))
  sleep 3
done

if [ "$version" = "0x302e31352e30" ]; then
  echo "devnet is running with protocol $version..."
  exit 0
fi

echo "could not check devnet, version=$version; fail!!!"
exit 1
