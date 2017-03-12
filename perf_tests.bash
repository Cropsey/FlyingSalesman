#!/bin/bash

if [ -z "$DONOTFETCH" ]; then

#	wget https://github.com/kiwicom/travelling-salesman/raw/master/real_data/data_5.txt.zip -O /tmp/data_5.txt.zip
#	wget https://github.com/kiwicom/travelling-salesman/raw/master/real_data/data_10.txt.zip -O /tmp/data_10.txt.zip
#	wget https://github.com/kiwicom/travelling-salesman/raw/master/real_data/data_15.txt.zip -O /tmp/data_15.txt.zip
#	wget https://github.com/kiwicom/travelling-salesman/raw/master/real_data/data_20.txt.zip -O /tmp/data_20.txt.zip
#	wget https://github.com/kiwicom/travelling-salesman/raw/master/real_data/data_30.txt.zip -O /tmp/data_30.txt.zip
#	wget https://github.com/kiwicom/travelling-salesman/raw/master/real_data/data_40.txt.zip -O /tmp/data_40.txt.zip
#	wget https://github.com/kiwicom/travelling-salesman/raw/master/real_data/data_50.txt.zip -O /tmp/data_50.txt.zip
#	wget https://github.com/kiwicom/travelling-salesman/raw/master/real_data/data_60.txt.zip -O /tmp/data_60.txt.zip
#	wget https://github.com/kiwicom/travelling-salesman/raw/master/real_data/data_70.txt.zip -O /tmp/data_70.txt.zip
#	wget https://github.com/kiwicom/travelling-salesman/raw/master/real_data/data_100.txt.zip -O /tmp/data_100.txt.zip
#	wget https://github.com/kiwicom/travelling-salesman/raw/master/real_data/data_200.txt.zip -O /tmp/data_200.txt.zip
#	wget https://github.com/kiwicom/travelling-salesman/raw/master/real_data/data_300.txt.zip -O /tmp/data_300.txt.zip
#	
#	unzip -o /tmp/data_5.txt.zip -d /tmp/
#	unzip -o /tmp/data_10.txt.zip -d /tmp/
#	unzip -o /tmp/data_15.txt.zip -d /tmp/
#	unzip -o /tmp/data_20.txt.zip -d /tmp/
#	unzip -o /tmp/data_30.txt.zip -d /tmp/
#	unzip -o /tmp/data_40.txt.zip -d /tmp/
#	unzip -o /tmp/data_50.txt.zip -d /tmp/
#	unzip -o /tmp/data_60.txt.zip -d /tmp/
#	unzip -o /tmp/data_70.txt.zip -d /tmp/
#	unzip -o /tmp/data_100.txt.zip -d /tmp/
#	unzip -o /tmp/data_200.txt.zip -d /tmp/
#	unzip -o /tmp/data_300.txt.zip -d /tmp/
    cp data/bottleneck_15.txt /tmp/data_bn_15.txt

fi


# RESULTS - reference with greedy DFS
# -------
#   /tmp/data_5.txt |  1950 |     Greedy |    4.9854ms
#  /tmp/data_10.txt |  5375 |     Greedy |   15.0398ms
#  /tmp/data_15.txt |  4344 |     Greedy |   763.039ms
#  /tmp/data_20.txt |  6864 |     Greedy | 13.1075727s
#  /tmp/data_30.txt |  8478 |     Greedy |   593.585ms
#  /tmp/data_40.txt |  9561 |     Greedy | 24.1656583s
#  /tmp/data_50.txt |  8886 |     Greedy |  261.6657ms
#  /tmp/data_60.txt | 11530 |     Greedy |   1.566712s
#  /tmp/data_70.txt | 15564 |     Greedy |  272.7369ms
# /tmp/data_100.txt | 17336 |     Greedy | 25.0097311s
# /tmp/data_200.txt | 41930 |     Greedy |  17.429643s
# /tmp/data_300.txt | 52060 |     Greedy |  9.7457127s
# -----------------------------------------------------
#            Total:  183878


RETVAL=0
declare -A results
declare -A info
declare -A reference=(  ["/tmp/data_5.txt"]=1950
			["/tmp/data_10.txt"]=5375
			["/tmp/data_15.txt"]=4344
			["/tmp/data_20.txt"]=6864
			["/tmp/data_30.txt"]=8478
			["/tmp/data_40.txt"]=9561
			["/tmp/data_50.txt"]=8886
			["/tmp/data_60.txt"]=11530
			["/tmp/data_70.txt"]=15564
			["/tmp/data_100.txt"]=17336
			["/tmp/data_200.txt"]=41930
			["/tmp/data_300.txt"]=52060
			["/tmp/data_bn_15.txt"]=22261
		)
reference_total=183878
for input in $(ls /tmp/data_*.txt | sort -n -t_ -k2); do
    echo "testing $input"
    cat "$input" | go run fspcmd/main.go -v > /tmp/out.txt 2> >(tee /tmp/errout.txt >&2)
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
#for k in $(echo "${!results[@]}" | sort -n -t_ -k2)
sum=0
for k in $(ls /tmp/data_*.txt | sort -n -t_ -k2)
do
	printf "%20s | %5d | %13s | %13s | %5d\n" $k ${results[$k]} ${info[$k]} $(( ${reference[$k]} - ${results[$k]} ))
	let sum+=${results[$k]}
done
printf "%68s\n" | tr ' ' -
printf "%20s %7d %31s %7d\n" "Total:" $sum "Improvement:" $(( $reference_total - $sum ))
exit $RETVAL
