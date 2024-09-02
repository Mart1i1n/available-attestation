#!/bin/sh
chmod +x /usr/bin/bootnode
/usr/bin/bootnode -nodekey /root/boot.key -addr 0.0.0.0:30301  > /root/bootnode/d.log 2>&1
