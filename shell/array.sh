#!/bin/bash

ary=(a b c d)

## 获取数组长度
## ${#ary[*]}

## 自增遍历
for (( i = 0; i < ${#ary[*]} ;i++ ));
do
 echo "${ary[i]}"
done


## in
for i in ${ary[*]};
do
  echo "$i"
done


## while
i=0
while [ $i \< "${#ary[*]}" ]; do
  echo "${ary[$i]}"
  (( i++ ))
done
set -e
if ! sh sub.sh; then
  echo "succ"
else
  echo "fail"
fi

sh sub.sh

echo "end"
