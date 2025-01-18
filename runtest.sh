#!/bin/bash
casetype=${1:-"1"}
caseduration=${2:-"9000"}

basedir=$(pwd)
casedir="${basedir}/case"
export BASEDIR="$basedir/"

if [ -f ".env.sh" ]; then
    source .env.sh
    echo "validator num is $VAL_NUM"
else
    echo ".env.sh is not found, please execute ./build.sh first"
    exit -1
fi


updategenesis() {
	docker run -it --rm -v "${basedir}/config:/root/config" --name generate --entrypoint /usr/bin/prysmctl tscel/ethnettools:0627 \
		testnet \
		generate-genesis \
		--fork=deneb \
		--num-validators=${VAL_NUM}\
		--genesis-time-delay=15 \
		--output-ssz=/root/config/genesis.ssz \
		--chain-config-file=/root/config/config.yml \
		--geth-genesis-json-in=/root/config/genesis.json \
		--geth-genesis-json-out=/root/config/genesis.json
}

testcase1() {
	subdir="blockcost"
	targetdir="${casedir}/${subdir}"
	resultdir="${basedir}/results/${subdir}"
	# if resultdir exist, delete it.
	if [ -d $resultdir ]; then
		rm -rf $resultdir
	fi
	mkdir -p $resultdir

	echo "Running testcase $subdir"
	echo "first test with vanilla version"
	updategenesis
	docker compose -f $targetdir/docker-compose-normal.yml up -d 
	echo "wait $caseduration seconds" && sleep $caseduration
	docker compose -f $targetdir/docker-compose-normal.yml down
	echo "result collect"
	sudo mv data $resultdir/data-normal
	cd $resultdir/data-normal
	sudo find . -name GetBeaconBlock.csv | xargs cat > /tmp/_b.csv
	sort -t "," -k 1n,1 /tmp/_b.csv > /tmp/normal_getblockcost.csv
	awk -F, '{sum+=$2}END{print "Generate block cost avg=", sum/NR}' /tmp/normal_getblockcost.csv
	sudo find . -name VerifyAttest.csv | xargs cat > /tmp/_b.csv
	sort -t "," -k 1n,1 /tmp/_b.csv > /tmp/normal_verifyatt.csv
	awk -F, '{sum+=$2}END{print "Verify attestation cost avg=", sum/NR}' /tmp/normal_verifyatt.csv
	sudo find . -name VerifyBeaconBlock.csv | xargs cat > /tmp/_b.csv
	sort -t "," -k 1n,1 /tmp/_b.csv > /tmp/normal_verifyblk.csv
	awk -F, '{sum+=$2}END{print "Verify block cost avg=", sum/NR}' /tmp/normal_verifyblk.csv
	sudo find . -name GetAttest.csv | xargs cat > /tmp/_b.csv
	sort -t "," -k 1n,1 /tmp/_b.csv > /tmp/normal_getatt.csv
	awk -F, '{sum+=$2}END{print "Generate attestation cost avg=", sum/NR}' /tmp/normal_getatt.csv
	cd $basedir

	echo "second test with modified version"
	updategenesis
	docker compose -f $targetdir/docker-compose-reorg.yml up -d
	echo "wait $caseduration seconds" && sleep $caseduration
	docker compose -f $targetdir/docker-compose-reorg.yml down
	echo "result collect"
	sudo mv data $resultdir/data-reorg
	cd $resultdir/data-reorg
	sudo find . -name GetBeaconBlock.csv | xargs cat > /tmp/_b.csv
	sort -t "," -k 1n,1 /tmp/_b.csv > /tmp/normal_getblockcost.csv
	awk -F, '{sum+=$2}END{print "Generate block cost avg=", sum/NR}' /tmp/normal_getblockcost.csv
	sudo find . -name VerifyAttest.csv | xargs cat > /tmp/_b.csv
	sort -t "," -k 1n,1 /tmp/_b.csv > /tmp/normal_verifyatt.csv
	awk -F, '{sum+=$2}END{print "Verify attestation cost avg=", sum/NR}' /tmp/normal_verifyatt.csv
	sudo find . -name VerifyBeaconBlock.csv | xargs cat > /tmp/_b.csv
	sort -t "," -k 1n,1 /tmp/_b.csv > /tmp/normal_verifyblk.csv
	awk -F, '{sum+=$2}END{print "Verify block cost avg=", sum/NR}' /tmp/normal_verifyblk.csv
	sudo find . -name GetAttest.csv | xargs cat > /tmp/_b.csv
	sort -t "," -k 1n,1 /tmp/_b.csv > /tmp/normal_getatt.csv
	awk -F, '{sum+=$2}END{print "Generate attestation cost avg=", sum/NR}' /tmp/normal_getatt.csv
	cd $basedir
	echo "test done and result in $resultdir"
}

testcase2() {
	subdir="attack-exante"
	targetdir="${casedir}/${subdir}"
	resultdir="${basedir}/results/${subdir}"
	# if resultdir exist, delete it.
	if [ -d $resultdir ]; then
		rm -rf $resultdir
	fi
	mkdir -p $resultdir

	epochsToWait=24

	echo "Running testcase $subdir"
	echo "first test with vanilla version"
	updategenesis
	docker compose -f $targetdir/docker-compose-normal.yml up -d 
	echo "wait $epochsToWait epochs" && sleep $(($epochsToWait * 12 * 32))
	docker compose -f $targetdir/docker-compose-normal.yml down
	sudo mv data $resultdir/data-normal

	echo "second test with modified version"
	updategenesis
	docker compose -f $targetdir/docker-compose-reorg.yml up -d
	echo "wait $epochsToWait epochs" && sleep $(($epochsToWait * 12 * 32))
	docker compose -f $targetdir/docker-compose-reorg.yml down
	sudo mv data $resultdir/data-reorg

	echo "test done and result in $resultdir"
}

testcase3() {
	subdir="tps-normal"
	targetdir="${casedir}/${subdir}"
	resultdir="${basedir}/results/${subdir}"
	# if resultdir exist, delete it.
	if [ -d $resultdir ]; then
		rm -rf $resultdir
	fi
	mkdir -p $resultdir

	echo "Running testcase $subdir"
	echo "first test with vanilla version"
	updategenesis
	docker compose -f $targetdir/docker-compose-normal.yml up -d 
	echo "wait $caseduration seconds" && sleep $caseduration
	docker compose -f $targetdir/docker-compose-normal.yml down
	sudo mv data $resultdir/data-normal

	echo "second test with modified version"
	updategenesis
	docker compose -f $targetdir/docker-compose-reorg.yml up -d 
	echo "wait $caseduration seconds" && sleep $caseduration
	docker compose -f $targetdir/docker-compose-reorg.yml down
	sudo mv data $resultdir/data-reorg
	echo "test done and result in $resultdir"
}

testcase4() {
	subdir="attack-sandwich"
	targetdir="${casedir}/${subdir}"
	resultdir="${basedir}/results/${subdir}"
	# if resultdir exist, delete it.
	if [ -d $resultdir ]; then
		rm -rf $resultdir
	fi
	mkdir -p $resultdir

	epochsToWait=24

	echo "Running testcase $subdir"
	echo "first test with vanilla version"
	updategenesis
	docker compose -f $targetdir/docker-compose-normal.yml up -d 
	echo "wait $epochsToWait epochs" && sleep $(($epochsToWait * 12 * 32))
	docker compose -f $targetdir/docker-compose-normal.yml down
	sudo mv data $resultdir/data-normal

	echo "second test with modified version"
	updategenesis
	docker compose -f $targetdir/docker-compose-reorg.yml up -d
	echo "wait $epochsToWait epochs" && sleep $(($epochsToWait * 12 * 32))
	docker compose -f $targetdir/docker-compose-reorg.yml down
	sudo mv data $resultdir/data-reorg

	echo "test done and result in $resultdir"
}

testcase5() {
	subdir="attack-unrealized"
	targetdir="${casedir}/${subdir}"
	resultdir="${basedir}/results/${subdir}"
	# if resultdir exist, delete it.
	if [ -d $resultdir ]; then
		rm -rf $resultdir
	fi
	mkdir -p $resultdir

	epochsToWait=24

	echo "Running testcase $subdir"
	echo "first test with vanilla version"
	updategenesis
	docker compose -f $targetdir/docker-compose-normal.yml up -d 
	echo "wait $epochsToWait epochs" && sleep $(($epochsToWait * 12 * 32))
	docker compose -f $targetdir/docker-compose-normal.yml down
	sudo mv data $resultdir/data-normal

	echo "second test with modified version"
	updategenesis
	docker compose -f $targetdir/docker-compose-reorg.yml up -d
	echo "wait $epochsToWait epochs" && sleep $(($epochsToWait * 12 * 32))
	docker compose -f $targetdir/docker-compose-reorg.yml down
	sudo mv data $resultdir/data-reorg

	echo "test done and result in $resultdir"
}

testcase6() {
	subdir="attack-withholding"
	targetdir="${casedir}/${subdir}"
	resultdir="${basedir}/results/${subdir}"
	# if resultdir exist, delete it.
	if [ -d $resultdir ]; then
		rm -rf $resultdir
	fi
	mkdir -p $resultdir

	epochsToWait=24

	echo "Running testcase $subdir"
	echo "first test with vanilla version"
	updategenesis
	docker compose -f $targetdir/docker-compose-normal.yml up -d 
	echo "wait $epochsToWait epochs" && sleep $(($epochsToWait * 12 * 32))
	docker compose -f $targetdir/docker-compose-normal.yml down
	sudo mv data $resultdir/data-normal

	echo "second test with modified version"
	updategenesis
	docker compose -f $targetdir/docker-compose-reorg.yml up -d
	echo "wait $epochsToWait epochs" && sleep $(($epochsToWait * 12 * 32))
	docker compose -f $targetdir/docker-compose-reorg.yml down
	sudo mv data $resultdir/data-reorg

	echo "test done and result in $resultdir"
}

testcase7() {
	subdir="attack-staircase"
	targetdir="${casedir}/${subdir}"
	resultdir="${basedir}/results/${subdir}"
	# if resultdir exist, delete it.
	if [ -d $resultdir ]; then
		rm -rf $resultdir
	fi
	mkdir -p $resultdir

	epochsToWait=24

	echo "Running testcase $subdir"
	echo "first test with vanilla version"
	updategenesis
	docker compose -f $targetdir/docker-compose-normal.yml up -d 
	echo "wait $epochsToWait epochs" && sleep $(($epochsToWait * 12 * 32))
	docker compose -f $targetdir/docker-compose-normal.yml down
	sudo mv data $resultdir/data-normal

	echo "second test with modified version"
	updategenesis
	docker compose -f $targetdir/docker-compose-reorg.yml up -d
	echo "wait $epochsToWait epochs" && sleep $(($epochsToWait * 12 * 32))
	docker compose -f $targetdir/docker-compose-reorg.yml down
	sudo mv data $resultdir/data-reorg

	echo "test done and result in $resultdir"
}

echo "casetype is $casetype"
case $casetype in
	1)
		testcase2
		;;
	2)
		testcase4
		;;
	3)
		testcase5
		;;
	4)
		testcase6
		;;
	5)
		testcase7
		;;
	6)
		testcase3
		;;
	7)
		testcase1
		;;
	*)
		echo "Invalid case type"
		;;
esac
