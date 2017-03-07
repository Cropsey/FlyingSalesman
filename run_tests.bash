#!/bin/bash

go test

for input in data/input*.txt; do
    output="${input/input/output}"
    echo -n "comparing $input $output - "
    cat "$input" | go run fspcmd/main.go > out.txt
    d=`diff out.txt "$output"`
    if [ "" == "$d" ]; then
        echo "ok"
    else
        echo "error"
        exit
    fi
done
