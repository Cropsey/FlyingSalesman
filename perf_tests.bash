#!/bin/bash

go test

wget https://github.com/kiwicom/travelling-salesman/raw/master/real_data/data_5.txt.zip -O /tmp/data_5.txt.zip
wget https://github.com/kiwicom/travelling-salesman/raw/master/real_data/data_10.txt.zip -O /tmp/data_10.txt.zip
wget https://github.com/kiwicom/travelling-salesman/raw/master/real_data/data_15.txt.zip -O /tmp/data_15.txt.zip
wget https://github.com/kiwicom/travelling-salesman/raw/master/real_data/data_20.txt.zip -O /tmp/data_20.txt.zip
wget https://github.com/kiwicom/travelling-salesman/raw/master/real_data/data_30.txt.zip -O /tmp/data_30.txt.zip
wget https://github.com/kiwicom/travelling-salesman/raw/master/real_data/data_50.txt.zip -O /tmp/data_50.txt.zip
wget https://github.com/kiwicom/travelling-salesman/raw/master/real_data/data_100.txt.zip -O /tmp/data_100.txt.zip
wget https://github.com/kiwicom/travelling-salesman/raw/master/real_data/data_200.txt.zip -O /tmp/data_200.txt.zip
wget https://github.com/kiwicom/travelling-salesman/raw/master/real_data/data_300.txt.zip -O /tmp/data_300.txt.zip

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
declare -A info
for input in $(ls /tmp/data_*.txt | sort -n -t_ -k2); do
    echo "testing $input"
    cat "$input" | go run fspcmd/main.go -v > >(tee /tmp/out.txt) 2> >(tee /tmp/errout.txt >&2)
    if [ $? -eq 0 ]; then
	    results[$input]=$(head -n1 /tmp/out.txt)
	    info[$input]=$(grep "New best" /tmp/errout.txt | tail -1 | cut -f6,11 -d" ")
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
	printf "%16s | %5d | %10s | %10s\n" $k ${results[$k]} ${info[$k]}
done
exit $RETVAL
