#!/bin/bash
casetype=${1:-"1"}
caseduration=${2:-"9000"}

basedir=$(pwd)
casedir="${basedir}/case"
export BASEDIR="$basedir/"


updategenesis() {
	docker run -it --rm -v "${basedir}/config:/root/config" --name generate --entrypoint /usr/bin/prysmctl tscel/ethnettools:0627 \
		testnet \
		generate-genesis \
		--fork=deneb \
		--num-validators=256 \
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
	echo "first test with normal version"
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

	echo "second test with reorg-fix version"
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
	subdir="attack-reorg"
	targetdir="${casedir}/${subdir}"
	resultdir="${basedir}/results/${subdir}"
	# if resultdir exist, delete it.
	if [ -d $resultdir ]; then
		rm -rf $resultdir
	fi
	mkdir -p $resultdir

	epochsToWait=20

	echo "Running testcase $subdir"
	echo "first test with normal version"
	updategenesis
	docker compose -f $targetdir/docker-compose-normal.yml up -d 
	echo "wait $epochsToWait epochs" && sleep $(($epochsToWait * 12 * 32))
	docker compose -f $targetdir/docker-compose-normal.yml down
	sudo mv data $resultdir/data-normal

	echo "second test with reorg-fix version"
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
	echo "first test with normal version"
	updategenesis
	docker compose -f $targetdir/docker-compose-normal.yml up -d 
	echo "wait $caseduration seconds" && sleep $caseduration
	docker compose -f $targetdir/docker-compose-normal.yml down
	sudo mv data $resultdir/data-normal

	echo "second test with reorg version"
	updategenesis
	docker compose -f $targetdir/docker-compose-reorg.yml up -d 
	echo "wait $caseduration seconds" && sleep $caseduration
	docker compose -f $targetdir/docker-compose-reorg.yml down
	sudo mv data $resultdir/data-reorg
	echo "test done and result in $resultdir"
}

echo "casetype is $casetype"
case $casetype in
	1)
		testcase1
		;;
	2)
		testcase2
		;;
	3)
		testcase3
		;;
	*)
		echo "Invalid case type"
		;;
esac
