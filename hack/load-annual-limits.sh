#!/bin/sh -ex

limits_file=$1

[ -z $limits_file ] && echo "Must specify file" && exit 1

for line in $(cat $limits_file); do
  if [[ -z $line ]]; then
    continue
  fi
  curl -X POST localhost:8000/v1/annual-limits -d "{\"amount\":"$line"}"
done
