#!/bin/bash
casetype=${1:-"1"}
caseduration=${2:-"9000"}

basedir=$(pwd)
casedir="${basedir}/case"
export BASEDIR="$basedir/"


PYTHON=$(which python3)

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
	awk -F, '{sum+=$2}END{print "Vanilla version block generation cost avg =", sum/NR, "ms"}' /tmp/normal_getblockcost.csv > $reportfile
	sudo find . -name VerifyAttest.csv | xargs cat > /tmp/_b.csv
	sort -t "," -k 1n,1 /tmp/_b.csv > /tmp/normal_verifyatt.csv
	awk -F, '{sum+=$2}END{print "Vanilla version attestation verification cost avg =", sum/NR, "us" }' /tmp/normal_verifyatt.csv >> $reportfile
	sudo find . -name VerifyBeaconBlock.csv | xargs cat > /tmp/_b.csv
	sort -t "," -k 1n,1 /tmp/_b.csv > /tmp/normal_verifyblk.csv
	awk -F, '{sum+=$2}END{print "Vanilla version block verification cost avg =", sum/NR, "ms" }' /tmp/normal_verifyblk.csv >> $reportfile
	sudo find . -name GetAttest.csv | xargs cat > /tmp/_b.csv
	sort -t "," -k 1n,1 /tmp/_b.csv > /tmp/normal_getatt.csv
	awk -F, '{sum+=$2}END{print "Vanilla version attestation generation cost avg =", sum/NR, "ms" }' /tmp/normal_getatt.csv >> $reportfile
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
	awk -F, '{sum+=$2}END{print "Modified version block generation cost avg=", sum/NR, "ms" }' /tmp/normal_getblockcost.csv >> $reportfile
	sudo find . -name VerifyAttest.csv | xargs cat > /tmp/_b.csv
	sort -t "," -k 1n,1 /tmp/_b.csv > /tmp/normal_verifyatt.csv
	awk -F, '{sum+=$2}END{print "Modified version attestation verification cost avg=", sum/NR, "us" }' /tmp/normal_verifyatt.csv >> $reportfile
	sudo find . -name VerifyBeaconBlock.csv | xargs cat > /tmp/_b.csv
	sort -t "," -k 1n,1 /tmp/_b.csv > /tmp/normal_verifyblk.csv
	awk -F, '{sum+=$2}END{print "Modified version block verification cost avg=", sum/NR, "ms" }' /tmp/normal_verifyblk.csv >> $reportfile
	sudo find . -name GetAttest.csv | xargs cat > /tmp/_b.csv
	sort -t "," -k 1n,1 /tmp/_b.csv > /tmp/normal_getatt.csv
	awk -F, '{sum+=$2}END{print "Modified version attestation generation cost avg=", sum/NR, "ms" }' /tmp/normal_getatt.csv >> $reportfile
	cd $basedir
  echo "test done and all data in $resultdir, report as bellow"
	# if PYTHON is not empty, then run the following script
	if [ -n "$PYTHON" ]; then
    $PYTHON ./collect/CollectLatency.py $reportfile
  else
    cat $reportfile
  fi
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

  # if PYTHON is not empty, then run the following script
  if [ -n "$PYTHON" ]; then
    $PYTHON ./collect/CollectTps.py $reportfile
    echo "test done and all data in $resultdir, report in $basedir/tps.png"
  else
    echo "test done and all data in $resultdir, report as bellow"
    cat $reportfile
  fi
  echo ""
  echo ""
}

testReorgs() {
  reorgs_result_dir="${basedir}/results/reorgtest"
  reportfile="${reorgs_result_dir}/report.txt"
  # define a array to store all test cases
  testcases=("attack-exante" "attack-sandwich" "attack-unrealized" "attack-withholding" "attack-staircase")
  # loop all test cases
  for testcase in ${testcases[@]}; do
    targetdir="${casedir}/${testcase}"
    resultdir="${reorgs_result_dir}/${testcase}"

    # if resultdir exist, delete it.
    if [ -d $resultdir ]; then
      rm -rf $resultdir
    fi
    mkdir -p $resultdir

    echo "Running testcase $testcase"
    echo "first test with vanilla version"
    updategenesis
    docker compose -f $targetdir/docker-compose-normal.yml up -d
    echo "wait $caseduration seconds" && sleep $caseduration
    docker compose -f $targetdir/docker-compose-normal.yml down
    sudo mv data $resultdir/data-normal
    echo "$testcase Vanilla version reorg event info: " >> $reportfile
    grep "reorg event" $resultdir/data-normal/attacker-1/d.log >> $reportfile

    echo "second test with modified version"
    updategenesis
    docker compose -f $targetdir/docker-compose-reorg.yml up -d
    echo "wait $caseduration seconds" && sleep $caseduration
    docker compose -f $targetdir/docker-compose-reorg.yml down
    sudo mv data $resultdir/data-reorg
    echo "$testcase Modified version reorg event info: " >> $reportfile
    grep "reorg event" $resultdir/data-reorg/attacker-1/d.log >> $reportfile
    echo "$testcase test done and all data in $resultdir"
  done

  if [ -n "$PYTHON" ]; then
    $PYTHON ./collect/CollectReorg.py $reorgs_result_dir
    echo "test done and all data in $resultdir, report in $basedir/reorgs.png"
  else
    echo "all test done and all data in $reorgs_result_dir, report as bellow"
    cat $reportfile
  fi
  echo ""
  echo ""

}

testOneReorg() {
  testcase=${1}
  targetdir="${casedir}/${testcase}"
  resultdir="${basedir}/results/${testcase}"
  reportfile="${resultdir}/report.txt"
  # if resultdir exist, delete it.
  if [ -d $resultdir ]; then
    rm -rf $resultdir
  fi
  mkdir -p $resultdir

  echo "Running testcase ${testcase}"
  echo "first test with vanilla version"
  updategenesis
  docker compose -f $targetdir/docker-compose-normal.yml up -d
  echo "wait $caseduration seconds" && sleep $caseduration
  docker compose -f $targetdir/docker-compose-normal.yml down
  sudo mv data $resultdir/data-normal
  echo "$testcase Vanilla version reorg event info: " >> $reportfile
  grep "reorg event" $resultdir/data-normal/attacker-1/d.log >> $reportfile

  echo "second test with modified version"
  updategenesis
  docker compose -f $targetdir/docker-compose-reorg.yml up -d
  echo "wait $caseduration seconds" && sleep $caseduration
  docker compose -f $targetdir/docker-compose-reorg.yml down
  sudo mv data $resultdir/data-reorg
  echo "$testcase Modified version reorg event info: " >> $reportfile
  grep "reorg event" $resultdir/data-reorg/attacker-1/d.log >> $reportfile

  echo "test done and all data in $resultdir, report as bellow"
  cat $reportfile
  echo ""
  echo ""
}

testReorg1() {
  testOneReorg "attack-exante"
}


testReorg2() {
  testOneReorg "attack-sandwich"
}

testReorg3() {
  testOneReorg "attack-unrealized"
}

testReorg4() {
  testOneReorg "attack-withholding"
}

testReorg5() {
  testOneReorg "attack-staircase"
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
    testReorgs
    ;;
	"latency")
		testLatency
		;;
	*)
		echo "Invalid case type"
		;;
esac
