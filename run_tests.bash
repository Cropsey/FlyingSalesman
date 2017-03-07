#!/bin/bash

go test

RETVAL=0
for input in data/input*.txt; do
    output="${input/input/output}"
    echo -n "comparing $input $output - "
    cat "$input" | go run fspcmd/main.go > out.txt
    if [ $? -eq 0 ]; then
        d=`diff out.txt "$output"`
        if [ "" == "$d" ]; then
            echo "ok"
        else
            echo "error: bad result"
	    RETVAL=1
        fi
    else
        echo "error: run time error"
        RETVAL=1
    fi
done
exit $RETVAL
