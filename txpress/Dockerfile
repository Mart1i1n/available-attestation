FROM ubuntu:22.04

WORKDIR /root

COPY ./build/bin/txpress /usr/bin/txpress
COPY ./app.json /root/app.json
COPY ./accounts.json /root/accounts.json

ENTRYPOINT [ "txpress --start --log /root/press.log " ]