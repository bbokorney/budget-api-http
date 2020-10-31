#!/bin/sh -e

categories_file=$1

[ -z $categories_file ] && echo "Must specify file" && exit 1

for cat in $(cat $categories_file); do
  curl -X POST localhost:8000/v1/categories -d "{\"name\":\"$cat\"}"
done
