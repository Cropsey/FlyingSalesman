#!/bin/bash

go test

for input in data/input*.txt; do
    output="${input/input/output}"
    echo -n "comparing $input $output - "
    d=`diff <(cat "$input" | go run cmd/main.go) "$output"`
    if [ "" == "$d" ]; then
        echo "ok"
    else
        echo "error"
    fi
done
