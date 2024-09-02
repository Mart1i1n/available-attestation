#!/bin/sh
initialize() {
	if [ ! -d "/root/gethdata/geth" ]; then
		echo "genesis init"
		base64 -d /root/config/base64.json > /tmp/secret.json
		/usr/bin/geth --datadir /root/gethdata init /root/config/genesis.json
		/usr/bin/geth --datadir /root/gethdata account import /tmp/secret.json <<EOM
123
123
EOM
		echo "123" > /root/config/password.txt
	fi
}

initialize 

/usr/bin/geth --datadir /root/gethdata \
	--http \
	--http.api=eth,net,web3,admin,txpool,debug \
	--http.addr=0.0.0.0 \
	--http.corsdomain=* \
	--ws \
	--ws.api=eth,net,web3 \
	--ws.addr=0.0.0.0 \
	--ws.origins=* \
	--authrpc.vhosts=* \
	--authrpc.addr=0.0.0.0 \
	--authrpc.port=8551 \
	--authrpc.jwtsecret=/root/config/jwtsecret \
	--datadir=gethdata \
	--bootnodes "enode://34c9a172990c85f80df96fe601743fbed3b3a00032b96dc44d56196f04a389f680b9b087d2214265fb4a9d6d1283199f33660965a307798f51e59342bbba881f@172.99.1.200:0?discport=30301" \
	--allow-insecure-unlock \
	--unlock=0x123463a4b065722e99115d6c222f267d9cabb524 \
	--password=/root/config/password.txt \
	--syncmode=full > /root/gethdata/d.log 2>&1
