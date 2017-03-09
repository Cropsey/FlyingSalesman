#!/bin/bash

go test

wget https://github.com/kiwicom/travelling-salesman/raw/master/real_data/data_5.txt.zip -O /tmp/data_5.txt.zip
wget https://github.com/kiwicom/travelling-salesman/raw/master/real_data/data_10.txt.zip -O /tmp/data_10.txt.zip
wget https://github.com/kiwicom/travelling-salesman/raw/master/real_data/data_15.txt.zip -O /tmp/data_15.txt.zip
wget https://github.com/kiwicom/travelling-salesman/raw/master/real_data/data_20.txt.zip -O /tmp/data_20.txt.zip
wget https://github.com/kiwicom/travelling-salesman/raw/master/real_data/data_20.txt.zip -O /tmp/data_30.txt.zip
wget https://github.com/kiwicom/travelling-salesman/raw/master/real_data/data_20.txt.zip -O /tmp/data_50.txt.zip
wget https://github.com/kiwicom/travelling-salesman/raw/master/real_data/data_20.txt.zip -O /tmp/data_100.txt.zip
wget https://github.com/kiwicom/travelling-salesman/raw/master/real_data/data_20.txt.zip -O /tmp/data_200.txt.zip
wget https://github.com/kiwicom/travelling-salesman/raw/master/real_data/data_20.txt.zip -O /tmp/data_300.txt.zip

unzip -o /tmp/data_5.txt.zip -d /tmp/
unzip -o /tmp/data_10.txt.zip -d /tmp/
unzip -o /tmp/data_15.txt.zip -d /tmp/
unzip -o /tmp/data_20.txt.zip -d /tmp/
unzip -o /tmp/data_30.txt.zip -d /tmp/
unzip -o /tmp/data_50.txt.zip -d /tmp/
unzip -o /tmp/data_100.txt.zip -d /tmp/
unzip -o /tmp/data_200.txt.zip -d /tmp/
unzip -o /tmp/data_300.txt.zip -d /tmp/


RETVAL=0
declare -A results
for input in /tmp/data_*.txt; do
    output="${input/data/output}"
    echo -n "testing $input $output - "
    cat "$input" | go run fspcmd/main.go -v > /tmp/out.txt
    if [ $? -eq 0 ]; then
            results[$input]=`head -n1 /tmp/out.txt`
    else
        echo "error: run time error"
        RETVAL=1
    fi
done
echo
echo "RESULTS"
echo "-------"
for k in "${!results[@]}"
do
	echo "$k : ${results[$k]}"
done
exit $RETVAL
