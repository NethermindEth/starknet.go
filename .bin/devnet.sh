#!/bin/sh

i=0
version=""

while true; do
  i=$((i + 1))
  curl --fail localhost:5050/is_alive 2>/dev/null 2>&1
  result=$?
  if [ $result -eq 0 ]; then
    exit 0
  fi
  if [ $i -gt 10  ]; then
    break
  fi
  echo "we will continue in a while, loop ${i}..."
  i=$((i + 1))
  sleep 3
done

echo "could not check devnet is_alive; fail!!!"
exit 1

