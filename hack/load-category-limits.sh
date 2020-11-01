#!/bin/sh -ex

limits_file=$1

[ -z $limits_file ] && echo "Must specify file" && exit 1

for line in $(cat $limits_file); do
  name=$(echo $line | cut -d ',' -f 1)
  amount=$(echo $line | cut -d ',' -f 2)
  if [[ -z $name ]]; then
    continue
  fi
  curl -X POST localhost:8000/v1/category-limits -d "{\"name\":\"$name\",\"limit\":"$amount"}"
done
