#!/bin/bash

if [ -z "$DONOTFETCH" ]; then

	echo -en 'travis_fold:start:Fetch-data\r'
	echo Fetch testing data from kiwi repository
	wget https://github.com/kiwicom/travelling-salesman/raw/master/real_data/data_5.txt.zip -O /tmp/data_5.txt.zip
	wget https://github.com/kiwicom/travelling-salesman/raw/master/real_data/data_10.txt.zip -O /tmp/data_10.txt.zip
	wget https://github.com/kiwicom/travelling-salesman/raw/master/real_data/data_15.txt.zip -O /tmp/data_15.txt.zip
	wget https://github.com/kiwicom/travelling-salesman/raw/master/real_data/data_20.txt.zip -O /tmp/data_20.txt.zip
	wget https://github.com/kiwicom/travelling-salesman/raw/master/real_data/data_30.txt.zip -O /tmp/data_30.txt.zip
	wget https://github.com/kiwicom/travelling-salesman/raw/master/real_data/data_40.txt.zip -O /tmp/data_40.txt.zip
	wget https://github.com/kiwicom/travelling-salesman/raw/master/real_data/data_50.txt.zip -O /tmp/data_50.txt.zip
	wget https://github.com/kiwicom/travelling-salesman/raw/master/real_data/data_60.txt.zip -O /tmp/data_60.txt.zip
	wget https://github.com/kiwicom/travelling-salesman/raw/master/real_data/data_70.txt.zip -O /tmp/data_70.txt.zip
	wget https://github.com/kiwicom/travelling-salesman/raw/master/real_data/data_100.txt.zip -O /tmp/data_100.txt.zip
	wget https://github.com/kiwicom/travelling-salesman/raw/master/real_data/data_200.txt.zip -O /tmp/data_200.txt.zip
	wget https://github.com/kiwicom/travelling-salesman/raw/master/real_data/data_300.txt.zip -O /tmp/data_300.txt.zip
	
	unzip -o /tmp/data_5.txt.zip -d /tmp/
	unzip -o /tmp/data_10.txt.zip -d /tmp/
	unzip -o /tmp/data_15.txt.zip -d /tmp/
	unzip -o /tmp/data_20.txt.zip -d /tmp/
	unzip -o /tmp/data_30.txt.zip -d /tmp/
	unzip -o /tmp/data_40.txt.zip -d /tmp/
	unzip -o /tmp/data_50.txt.zip -d /tmp/
	unzip -o /tmp/data_60.txt.zip -d /tmp/
	unzip -o /tmp/data_70.txt.zip -d /tmp/
	unzip -o /tmp/data_100.txt.zip -d /tmp/
	unzip -o /tmp/data_200.txt.zip -d /tmp/
	unzip -o /tmp/data_300.txt.zip -d /tmp/
	#cp data/bottleneck_15.txt /tmp/data_bn_15.txt
	echo -en 'travis_fold:end:Fetch-data\r'

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
#reference_total=206139 # with bn_15

declare -A best_reference=(  ["/tmp/data_5.txt"]=1950
			     ["/tmp/data_10.txt"]=5375
			     ["/tmp/data_15.txt"]=4281
			     ["/tmp/data_20.txt"]=6053
			     ["/tmp/data_30.txt"]=7629
			     ["/tmp/data_40.txt"]=7751
			     ["/tmp/data_50.txt"]=7235
			     ["/tmp/data_60.txt"]=9180
			     ["/tmp/data_70.txt"]=12358
			     ["/tmp/data_100.txt"]=15609
			     ["/tmp/data_200.txt"]=28338
			     ["/tmp/data_300.txt"]=37957
			     ["/tmp/data_bn_15.txt"]=22261
		)
best_reference_total=143716

go build && go build fspcmd/main.go
for input in $(ls /tmp/data_*.txt | sort -n -t_ -k2); do
    echo -en "travis_fold:start:${input##*/}\r"
    echo "testing $input"
    #cat "$input" | go run fspcmd/main.go -v > /tmp/out.txt 2> >(tee /tmp/errout.txt >&2)
    cat "$input" | ./main -v -t 30 > /tmp/out.txt 2> >(tee /tmp/errout.txt >&2)
    if [ $? -eq 0 ]; then
	    results[$input]=$(head -n1 /tmp/out.txt)
	    info[$input]=$(grep "New best" /tmp/errout.txt | tail -1 | cut -f6,11 -d" ")
    else
        echo "error: run time error"
        RETVAL=1
    fi
    echo -en "travis_fold:end:${input##*/}\r"
done
echo
echo "RESULTS"
echo "-------"
printf "%20s | %5s | %23s | %13s | %6s | %6s | %15s\n" "input" "price" "engine" "time" "d(ref)" "d(bst)" "score"
printf "%106s\n" | tr ' ' -
#for k in $(echo "${!results[@]}" | sort -n -t_ -k2)
sum=0
total_points=0
for k in $(ls /tmp/data_*.txt | sort -n -t_ -k2)
do
	size=$(echo ${k} | sed -e s/[^0-9]//g)
	[ $k = "/tmp/data_bn_15.txt" ] && size=0
	max_points=$(echo - | awk "{print log(${size})/log(2)}")
	earned_points=$(echo - | awk "{print ((${best_reference[$k]} / ${results[$k]}) * ${max_points})}")
	printf "%20s | %5d | %23s | %13s | %6d | %6d | %6.5f/%6.5f \n" \
		$k ${results[$k]} ${info[$k]} \
		$(( ${reference[$k]} - ${results[$k]} )) \
		$(( ${best_reference[$k]} - ${results[$k]} )) \
		${earned_points} ${max_points}
	let sum+=${results[$k]}
	total_points=$(echo - | awk "{print ${total_points} + ${earned_points}}")
done
printf "%106s\n" | tr ' ' -
reference_improvement=$(( $reference_total - $sum))
best_improvement=$(( $best_reference_total - $sum))
printf "%20s %7d %41s %8d %8d\n" "Total:" $sum "Improvement:" $reference_improvement $best_improvement
printf "%71s (%5.1f%%) (%5.1f%%)\n" " "   $( echo - | awk "{ print $reference_improvement / $reference_total * 100 }" ) \
					$( echo - | awk "{ print $best_improvement / $best_reference_total * 100 }" )


max_total_points=64.29805444314002 # sum([log2(s) for s in data_set_sizes])
printf "%30s  %6.5f/%6.5f\n" "Score:" ${total_points} ${max_total_points}


exit $RETVAL
