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
		--num-validators=16384\
		--genesis-time-delay=15 \
		--output-ssz=/root/config/genesis.ssz \
		--chain-config-file=/root/config/config.yml \
		--geth-genesis-json-in=/root/config/genesis.json \
		--geth-genesis-json-out=/root/config/genesis.json
}

testLatency() {
	subdir="blockcost"
	targetdir="${casedir}/${subdir}"
	resultdir="${basedir}/results/${subdir}"
	reportfile="${resultdir}/report.txt"
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
	cd $resultdir/data-normal
	sudo find . -name GetBeaconBlock.csv | xargs cat > /tmp/_b.csv
	sort -t "," -k 1n,1 /tmp/_b.csv > /tmp/normal_getblockcost.csv
	awk -F, '{sum+=$2}END{print "Vanilla version block generate cost avg =", sum/NR, "ms"}' /tmp/normal_getblockcost.csv > $reportfile
	sudo find . -name VerifyAttest.csv | xargs cat > /tmp/_b.csv
	sort -t "," -k 1n,1 /tmp/_b.csv > /tmp/normal_verifyatt.csv
	awk -F, '{sum+=$2}END{print "Vanilla version attestation verify cost avg =", sum/NR, "us" }' /tmp/normal_verifyatt.csv >> $reportfile
	sudo find . -name VerifyBeaconBlock.csv | xargs cat > /tmp/_b.csv
	sort -t "," -k 1n,1 /tmp/_b.csv > /tmp/normal_verifyblk.csv
	awk -F, '{sum+=$2}END{print "Vanilla version block verify cost avg =", sum/NR, "ms" }' /tmp/normal_verifyblk.csv >> $reportfile
	sudo find . -name GetAttest.csv | xargs cat > /tmp/_b.csv
	sort -t "," -k 1n,1 /tmp/_b.csv > /tmp/normal_getatt.csv
	awk -F, '{sum+=$2}END{print "Vanilla version attestation generate cost avg =", sum/NR, "ms" }' /tmp/normal_getatt.csv >> $reportfile
	cd $basedir

	echo "second test with modified version"
	updategenesis
	docker compose -f $targetdir/docker-compose-reorg.yml up -d
	echo "wait $caseduration seconds" && sleep $caseduration
	docker compose -f $targetdir/docker-compose-reorg.yml down
	sudo mv data $resultdir/data-reorg
	cd $resultdir/data-reorg
	sudo find . -name GetBeaconBlock.csv | xargs cat > /tmp/_b.csv
	sort -t "," -k 1n,1 /tmp/_b.csv > /tmp/normal_getblockcost.csv
	awk -F, '{sum+=$2}END{print "Modified version block generate cost avg=", sum/NR, "ms" }' /tmp/normal_getblockcost.csv >> $reportfile
	sudo find . -name VerifyAttest.csv | xargs cat > /tmp/_b.csv
	sort -t "," -k 1n,1 /tmp/_b.csv > /tmp/normal_verifyatt.csv
	awk -F, '{sum+=$2}END{print "Modified version attestation verify cost avg=", sum/NR, "us" }' /tmp/normal_verifyatt.csv >> $reportfile
	sudo find . -name VerifyBeaconBlock.csv | xargs cat > /tmp/_b.csv
	sort -t "," -k 1n,1 /tmp/_b.csv > /tmp/normal_verifyblk.csv
	awk -F, '{sum+=$2}END{print "Modified version block verify cost avg=", sum/NR, "ms" }' /tmp/normal_verifyblk.csv >> $reportfile
	sudo find . -name GetAttest.csv | xargs cat > /tmp/_b.csv
	sort -t "," -k 1n,1 /tmp/_b.csv > /tmp/normal_getatt.csv
	awk -F, '{sum+=$2}END{print "Modified version attestation generate cost avg=", sum/NR, "ms" }' /tmp/normal_getatt.csv >> $reportfile
	cd $basedir
	echo "test done and all data in $resultdir, report as bellow"
	cat $reportfile
	echo ""
	echo ""

}

testReorg1() {
	subdir="attack-exante"
	targetdir="${casedir}/${subdir}"
	resultdir="${basedir}/results/${subdir}"
	reportfile="${resultdir}/report.txt"
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
	echo "Vanilla version reorg event info: " >> $reportfile
	grep "reorg event" $resultdir/data-normal/attacker-1/d.log >> $reportfile

	echo "second test with modified version"
	updategenesis
	docker compose -f $targetdir/docker-compose-reorg.yml up -d
	echo "wait $caseduration seconds" && sleep $caseduration
	docker compose -f $targetdir/docker-compose-reorg.yml down
	sudo mv data $resultdir/data-reorg
	echo "Modified version reorg event info: " >> $reportfile
  grep "reorg event" $resultdir/data-reorg/attacker-1/d.log >> $reportfile

	echo "test done and all data in $resultdir, report as bellow"
  cat $reportfile
  echo ""
  echo ""

}

testTps() {
	subdir="tps-normal"
	targetdir="${casedir}/${subdir}"
	resultdir="${basedir}/results/${subdir}"
	reportfile="${resultdir}/report.txt"
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
	echo "Vanilla version tps detail info: " > $reportfile
	grep "test history" $resultdir/data-normal/txpress/press.log >> $reportfile

	echo "second test with modified version"
	updategenesis
	docker compose -f $targetdir/docker-compose-reorg.yml up -d 
	echo "wait $caseduration seconds" && sleep $caseduration
	docker compose -f $targetdir/docker-compose-reorg.yml down
	sudo mv data $resultdir/data-reorg
	echo "Modified version tps detail info: " >> $reportfile
  grep "test history" $resultdir/data-reorg/txpress/press.log >> $reportfile
	echo "test done and all data in $resultdir, report as bellow"
  cat $reportfile
  echo ""
  echo ""
}

testReorg2() {
	subdir="attack-sandwich"
	targetdir="${casedir}/${subdir}"
	resultdir="${basedir}/results/${subdir}"
	reportfile="${resultdir}/report.txt"
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
	echo "Vanilla version reorg event info: " >> $reportfile
  grep "reorg event" $resultdir/data-normal/attacker-1/d.log >> $reportfile

	echo "second test with modified version"
	updategenesis
	docker compose -f $targetdir/docker-compose-reorg.yml up -d
	echo "wait $caseduration seconds" && sleep $caseduration
	docker compose -f $targetdir/docker-compose-reorg.yml down
	sudo mv data $resultdir/data-reorg

	echo "Modified version reorg event info: " >> $reportfile
  grep "reorg event" $resultdir/data-reorg/attacker-1/d.log >> $reportfile

  echo "test done and all data in $resultdir, report as bellow"
  cat $reportfile
  echo ""
  echo ""
}

testReorg3() {
	subdir="attack-unrealized"
	targetdir="${casedir}/${subdir}"
	resultdir="${basedir}/results/${subdir}"
	reportfile="${resultdir}/report.txt"
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
	echo "Vanilla version reorg event info: " >> $reportfile
  grep "reorg event" $resultdir/data-normal/attacker-1/d.log >> $reportfile

	echo "second test with modified version"
	updategenesis
	docker compose -f $targetdir/docker-compose-reorg.yml up -d
	echo "wait $caseduration seconds" && sleep $caseduration
	docker compose -f $targetdir/docker-compose-reorg.yml down
	sudo mv data $resultdir/data-reorg

	echo "Modified version reorg event info: " >> $reportfile
  grep "reorg event" $resultdir/data-reorg/attacker-1/d.log >> $reportfile

  echo "test done and all data in $resultdir, report as bellow"
  cat $reportfile
  echo ""
  echo ""
}

testReorg4() {
	subdir="attack-withholding"
	targetdir="${casedir}/${subdir}"
	resultdir="${basedir}/results/${subdir}"
	reportfile="${resultdir}/report.txt"
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
	echo "Vanilla version reorg event info: " >> $reportfile
  grep "reorg event" $resultdir/data-normal/attacker-1/d.log >> $reportfile

	echo "second test with modified version"
	updategenesis
	docker compose -f $targetdir/docker-compose-reorg.yml up -d
	echo "wait $caseduration seconds" && sleep $caseduration
	docker compose -f $targetdir/docker-compose-reorg.yml down
	sudo mv data $resultdir/data-reorg

	echo "Modified version reorg event info: " >> $reportfile
  grep "reorg event" $resultdir/data-reorg/attacker-1/d.log >> $reportfile

  echo "test done and all data in $resultdir, report as bellow"
  cat $reportfile
  echo ""
  echo ""
}

testReorg5() {
	subdir="attack-staircase"
	targetdir="${casedir}/${subdir}"
	resultdir="${basedir}/results/${subdir}"
	reportfile="${resultdir}/report.txt"
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
	grep "reorg event" $resultdir/data-normal/att
	echo "Vanilla version reorg event info: " >> $reportfile
  grep "reorg event" $resultdir/data-normal/attacker-1/d.log >> $reportfile

	echo "second test with modified version"
	updategenesis
	docker compose -f $targetdir/docker-compose-reorg.yml up -d
	echo "wait $caseduration seconds" && sleep $caseduration
	docker compose -f $targetdir/docker-compose-reorg.yml down
	sudo mv data $resultdir/data-reorg

	echo "Modified version reorg event info: " >> $reportfile
  grep "reorg event" $resultdir/data-reorg/attacker-1/d.log >> $reportfile

  echo "test done and all data in $resultdir, report as bellow"
  cat $reportfile
  echo ""
  echo ""
}

echo "casetype is $casetype"
case $casetype in
	1)
		testReorg1
		;;
	2)
		testReorg2
		;;
	3)
		testReorg3
		;;
	4)
		testReorg4
		;;
	5)
		testReorg5
		;;
	"tps")
		testTps
		;;
  "reorg")
    testcase3
    ;;
	"latency")
		testLatency
		;;
	*)
		echo "Invalid case type"
		;;
esac
